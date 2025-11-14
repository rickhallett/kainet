import { createClient } from '@libsql/client';
import { encrypt, decrypt, deriveKey } from './crypto.js';

// Embedded credentials
const DB_URL = 'libsql://bt-phone-home-rickhallett.aws-eu-west-1.turso.io';
const AUTH_TOKEN = 'eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NjMwOTg5MjEsImlkIjoiNjc4YjAxZDQtYzExZi00OTMyLTk2MzktZDUxODNlZTVmMTI3IiwicmlkIjoiNzc0ZDk3NmItNTMzYS00NzE1LTlkZmItY2RmYjU1N2M2ZmRjIn0.4hZBJRmkoSMutMqpGVYhmIWmo2-5lTKle9p5QkSyF6B-OHGS1xpdoFi_wEZVn53qcHEaUEHAzg8rPiy0ZsjFCQ';

let db;
let encKey;
let username;
let roomName;
let lastID = 0;
let pollInterval;

const output = document.getElementById('output');
const input = document.getElementById('input');
const prompt = document.getElementById('prompt');

// Terminal functions
function print(text, className = '') {
  const line = document.createElement('div');
  line.className = `line ${className}`;
  line.textContent = text;
  output.appendChild(line);
  scrollToBottom();
}

function printHTML(html) {
  const line = document.createElement('div');
  line.className = 'line';
  line.innerHTML = html;
  output.appendChild(line);
  scrollToBottom();
}

function scrollToBottom() {
  output.scrollTop = output.scrollHeight;
}

function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// Boot sequence
async function printBanner() {
  // Check if mobile
  const isMobile = window.innerWidth <= 768;

  if (isMobile) {
    // Simple mobile banner
    print('╔════════════════════════════╗', 'green');
    print('║                            ║', 'green');
    print('║        K A I N E T         ║', 'green');
    print('║                            ║', 'green');
    print('║   CLASSIFIED TERMINAL      ║', 'green');
    print('║   UNAUTHORIZED ACCESS      ║', 'green');
    print('║   PROHIBITED               ║', 'green');
    print('║                            ║', 'green');
    print('╚════════════════════════════╝', 'green');
  } else {
    // Full desktop banner
    const banner = `
╔══════════════════════════════════════════════════════════╗
║                                                          ║
║   ██╗  ██╗ █████╗ ██╗███╗   ██╗███████╗████████╗        ║
║   ██║ ██╔╝██╔══██╗██║████╗  ██║██╔════╝╚══██╔══╝        ║
║   █████╔╝ ███████║██║██╔██╗ ██║█████╗     ██║           ║
║   ██╔═██╗ ██╔══██║██║██║╚██╗██║██╔══╝     ██║           ║
║   ██║  ██╗██║  ██║██║██║ ╚████║███████╗   ██║           ║
║   ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝   ╚═╝           ║
║                                                          ║
║           CLASSIFIED COMMUNICATIONS TERMINAL             ║
║              UNAUTHORIZED ACCESS PROHIBITED              ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝`;

    print(banner, 'banner green');
  }
}

async function bootSequence() {
  print('SYSTEM INITIALIZING...', 'cyan');
  await sleep(200);

  const steps = [
    'loading security protocols',
    'verifying operator credentials',
    `connecting to room: ${roomName}`,
    'establishing quantum-resistant handshake',
    'activating end-to-end encryption'
  ];

  for (const step of steps) {
    printHTML(`<span class="cyan">&gt;&gt;&gt;</span> ${step}<span class="green"> OK</span>`);
    await sleep(150);
  }

  print('');
}

async function initDB() {
  printHTML(`<span class="cyan">&gt;&gt;&gt;</span> establishing secure connection...<span class="green"> CONNECTED</span>`);

  db = createClient({
    url: DB_URL,
    authToken: AUTH_TOKEN
  });

  printHTML(`<span class="cyan">&gt;&gt;&gt;</span> initializing secure channel...<span class="green"> OK</span>`);

  await db.execute(`
    CREATE TABLE IF NOT EXISTS messages (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      room_name TEXT NOT NULL,
      username TEXT NOT NULL,
      message TEXT NOT NULL,
      timestamp INTEGER NOT NULL
    )
  `);

  await db.execute(`
    CREATE INDEX IF NOT EXISTS idx_room_id ON messages(room_name, id)
  `);

  printHTML(`<span class="cyan">&gt;&gt;&gt;</span> generating AES-256 encryption keys...<span class="green"> READY</span>`);
  encKey = await deriveKey(AUTH_TOKEN);

  print('');
  printHTML(`
<span class="green box-line">╔════════════════════════════════════════╗</span>
<span class="green box-line">║  SECURE CHANNEL ACTIVE                 ║</span>
<span class="green box-line">║  OPERATOR: ${username.toUpperCase().padEnd(28)}║</span>
<span class="green box-line">║  ROOM: ${roomName.toUpperCase().padEnd(32)}║</span>
<span class="green box-line">║  ENCRYPTION: AES-256-GCM               ║</span>
<span class="green box-line">╚════════════════════════════════════════╝</span>
  `);
  print('');
  print('commands: /burn (wipe history)', 'dim');
  print('');
}

async function displayHistory(limit = 20) {
  const result = await db.execute({
    sql: 'SELECT username, message, timestamp FROM messages WHERE room_name = ? ORDER BY id DESC LIMIT ?',
    args: [roomName, limit]
  });

  const messages = [];
  for (const row of result.rows) {
    try {
      const decrypted = await decrypt(encKey, row.message);
      messages.push({
        username: row.username,
        message: decrypted,
        timestamp: row.timestamp
      });
    } catch (e) {
      messages.push({
        username: row.username,
        message: '[DECRYPT ERROR]',
        timestamp: row.timestamp
      });
    }
  }

  // Display in chronological order
  for (let i = messages.length - 1; i >= 0; i--) {
    const m = messages[i];
    const time = new Date(m.timestamp * 1000).toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    });
    printHTML(`<span class="dim">[${time}]</span> <span class="cyan">${m.username}:</span> ${m.message}`);
  }
}

async function getLastMessageID() {
  const result = await db.execute({
    sql: 'SELECT COALESCE(MAX(id), 0) as max_id FROM messages WHERE room_name = ?',
    args: [roomName]
  });
  return result.rows[0].max_id;
}

async function pollNewMessages() {
  const result = await db.execute({
    sql: 'SELECT id, username, message, timestamp FROM messages WHERE room_name = ? AND id > ? ORDER BY id ASC',
    args: [roomName, lastID]
  });

  for (const row of result.rows) {
    try {
      const decrypted = await decrypt(encKey, row.message);
      const time = new Date(row.timestamp * 1000).toLocaleTimeString('en-US', {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
      });
      printHTML(`<span class="dim">[${time}]</span> <span class="cyan">${row.username}:</span> ${decrypted}`);

      if (row.id > lastID) {
        lastID = row.id;
      }
    } catch (e) {
      console.error('Decrypt error:', e);
    }
  }
}

async function sendMessage(message) {
  const encrypted = await encrypt(encKey, message);
  const timestamp = Math.floor(Date.now() / 1000);

  const result = await db.execute({
    sql: 'INSERT INTO messages (room_name, username, message, timestamp) VALUES (?, ?, ?, ?)',
    args: [roomName, username, encrypted, timestamp]
  });

  const time = new Date().toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  });
  printHTML(`<span class="dim">[${time}]</span> <span class="yellow">${username}:</span> ${message}`);

  // Update lastID to skip our own message in polling
  if (result.lastInsertRowid > lastID) {
    lastID = Number(result.lastInsertRowid);
  }
}

async function handleCommand(cmd) {
  if (cmd === '/burn') {
    print('WARNING: purging all encrypted transmissions for this room...', 'yellow');
    await sleep(500);
    await db.execute({
      sql: 'DELETE FROM messages WHERE room_name = ?',
      args: [roomName]
    });
    print('PURGED', 'green');
    print('all message history for this room has been permanently erased', 'dim');
    lastID = 0;
  } else {
    print(`COMMAND FAILED: unknown command: ${cmd}`, 'red');
  }
}

// Main initialization
async function main() {
  const welcomeScreen = document.getElementById('welcome-screen');
  const terminal = document.getElementById('terminal');
  const loginForm = document.getElementById('login-form');
  const usernameInput = document.getElementById('username-input');
  const roomInput = document.getElementById('room-input');

  // Focus username input on load
  usernameInput.focus();

  loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    username = usernameInput.value.trim();
    roomName = roomInput.value.trim();

    if (!username || !roomName) {
      return;
    }

    // Hide welcome screen and show terminal
    welcomeScreen.style.display = 'none';
    terminal.style.display = 'flex';

    await printBanner();
    await bootSequence();
    await initDB();
    await displayHistory(20);

    lastID = await getLastMessageID();

    // Start polling
    pollInterval = setInterval(pollNewMessages, 1000);

    // Set up input
    prompt.textContent = '>';
    input.focus();

    input.addEventListener('keydown', async (e) => {
      if (e.key === 'Enter') {
        const message = input.value.trim();
        input.value = '';

        if (!message) return;

        if (message.startsWith('/')) {
          await handleCommand(message);
        } else {
          await sendMessage(message);
        }
      }
    });
  });
}

// Start the app
main().catch(err => {
  console.error('CRITICAL ERROR:', err);
});
