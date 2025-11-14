import { deriveKey, encrypt, decrypt, AUTH_TOKEN } from './crypto.js';

const DB_URL = "https://bt-phone-home-rickhallett.aws-eu-west-1.turso.io";
const POLL_INTERVAL = 1000; // 1 second

class KaiNetClient {
    constructor() {
        this.username = null;
        this.roomName = null;
        this.key = null;
        this.lastCheck = new Date().toISOString();
        this.pollTimer = null;
        this.isConnected = false;

        this.initializeEventListeners();
    }

    initializeEventListeners() {
        // Login screen
        document.getElementById('connect-btn').addEventListener('click', () => this.connect());
        document.getElementById('username').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') document.getElementById('roomname').focus();
        });
        document.getElementById('roomname').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.connect();
        });

        // Chat screen
        document.getElementById('message-input').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') this.sendMessage();
        });
        document.getElementById('burn-btn').addEventListener('click', () => this.burnRoom());
        document.getElementById('disconnect-btn').addEventListener('click', () => this.disconnect());
    }

    async connect() {
        const usernameInput = document.getElementById('username');
        const roomnameInput = document.getElementById('roomname');

        this.username = usernameInput.value.trim();
        this.roomName = roomnameInput.value.trim();

        if (!this.username || !this.roomName) {
            alert('Please enter both username and room name');
            return;
        }

        try {
            // Derive encryption key
            this.key = await deriveKey(AUTH_TOKEN);

            // Initialize database schema
            await this.initDB();

            // Show chat screen
            document.getElementById('login-screen').classList.remove('active');
            document.getElementById('chat-screen').classList.add('active');

            // Update header
            document.getElementById('current-username').textContent = this.username;
            document.getElementById('current-room').textContent = this.roomName;

            // Show boot sequence
            await this.showBootSequence();

            // Start polling for messages
            this.isConnected = true;
            this.startPolling();

            // Focus message input
            document.getElementById('message-input').focus();
        } catch (error) {
            console.error('Connection failed:', error);
            alert('Failed to connect: ' + error.message);
        }
    }

    async showBootSequence() {
        const bootSeq = document.getElementById('boot-sequence');
        const lines = bootSeq.querySelectorAll('.boot-line');

        bootSeq.style.display = 'block';

        for (let i = 0; i < lines.length; i++) {
            await new Promise(resolve => setTimeout(resolve, 200));
            lines[i].style.opacity = '1';
        }

        await new Promise(resolve => setTimeout(resolve, 1000));
        bootSeq.style.display = 'none';
    }

    async initDB() {
        const schema = `
            CREATE TABLE IF NOT EXISTS messages (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                room_name TEXT NOT NULL,
                username TEXT NOT NULL,
                message TEXT NOT NULL,
                timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
            );
            CREATE INDEX IF NOT EXISTS idx_room_timestamp ON messages(room_name, timestamp);
        `;

        await this.executeSql(schema);
    }

    async executeSql(sql, params = []) {
        const response = await fetch(DB_URL, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${AUTH_TOKEN}`,
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                statements: [{ q: sql, params }]
            })
        });

        if (!response.ok) {
            throw new Error(`Database error: ${response.statusText}`);
        }

        const data = await response.json();
        return data;
    }

    async sendMessage() {
        const input = document.getElementById('message-input');
        const message = input.value.trim();

        if (!message) return;

        try {
            const encrypted = await encrypt(message, this.key);

            await this.executeSql(
                'INSERT INTO messages (room_name, username, message) VALUES (?, ?, ?)',
                [this.roomName, this.username, encrypted]
            );

            input.value = '';
        } catch (error) {
            console.error('Failed to send message:', error);
            this.displaySystemMessage('ERROR: Failed to send message', 'error');
        }
    }

    async burnRoom() {
        if (!confirm('Are you sure you want to wipe all messages in this room?')) {
            return;
        }

        try {
            await this.executeSql(
                'DELETE FROM messages WHERE room_name = ?',
                [this.roomName]
            );

            // Clear messages display
            document.getElementById('messages').innerHTML = '';
            this.displaySystemMessage('BURN: Room history wiped', 'burn');
        } catch (error) {
            console.error('Failed to burn room:', error);
            this.displaySystemMessage('ERROR: Failed to burn room', 'error');
        }
    }

    startPolling() {
        this.pollTimer = setInterval(() => this.pollMessages(), POLL_INTERVAL);
    }

    stopPolling() {
        if (this.pollTimer) {
            clearInterval(this.pollTimer);
            this.pollTimer = null;
        }
    }

    async pollMessages() {
        if (!this.isConnected) return;

        try {
            const result = await this.executeSql(
                `SELECT id, room_name, username, message, timestamp
                 FROM messages
                 WHERE room_name = ? AND timestamp > ?
                 ORDER BY timestamp ASC`,
                [this.roomName, this.lastCheck]
            );

            const rows = result[0]?.results?.rows || [];

            for (const row of rows) {
                const decrypted = await decrypt(row[3], this.key);
                if (decrypted) {
                    this.displayMessage({
                        username: row[2],
                        message: decrypted,
                        timestamp: row[4]
                    });
                }
            }

            if (rows.length > 0) {
                this.lastCheck = new Date().toISOString();
            }
        } catch (error) {
            console.error('Failed to poll messages:', error);
        }
    }

    displayMessage(msg) {
        const messagesDiv = document.getElementById('messages');
        const messageEl = document.createElement('div');
        messageEl.className = 'message';

        const timestamp = new Date(msg.timestamp).toLocaleTimeString('en-US', {
            hour12: false,
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });

        const isOwnMessage = msg.username === this.username;
        const usernameClass = isOwnMessage ? 'username-self' : 'username-other';

        messageEl.innerHTML = `
            <span class="timestamp">${timestamp}</span>
            <span class="${usernameClass}">${msg.username}</span>:
            <span class="message-text">${this.escapeHtml(msg.message)}</span>
        `;

        messagesDiv.appendChild(messageEl);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
    }

    displaySystemMessage(message, type = 'info') {
        const messagesDiv = document.getElementById('messages');
        const messageEl = document.createElement('div');
        messageEl.className = `message system-message system-${type}`;
        messageEl.textContent = `[${message}]`;
        messagesDiv.appendChild(messageEl);
        messagesDiv.scrollTop = messagesDiv.scrollHeight;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    disconnect() {
        this.stopPolling();
        this.isConnected = false;

        document.getElementById('chat-screen').classList.remove('active');
        document.getElementById('login-screen').classList.add('active');

        // Clear messages
        document.getElementById('messages').innerHTML = '';

        // Reset inputs
        document.getElementById('username').value = '';
        document.getElementById('roomname').value = '';
        document.getElementById('message-input').value = '';

        // Focus username input
        document.getElementById('username').focus();
    }
}

// Initialize client when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new KaiNetClient();
});
