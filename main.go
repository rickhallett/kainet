package main

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// embedded credentials - set these before building
const (
	embeddedDBURL      = "libsql://bt-phone-home-rickhallett.aws-eu-west-1.turso.io"
	embeddedAuthToken  = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NjMwOTg5MjEsImlkIjoiNjc4YjAxZDQtYzExZi00OTMyLTk2MzktZDUxODNlZTVmMTI3IiwicmlkIjoiNzc0ZDk3NmItNTMzYS00NzE1LTlkZmItY2RmYjU1N2M2ZmRjIn0.4hZBJRmkoSMutMqpGVYhmIWmo2-5lTKle9p5QkSyF6B-OHGS1xpdoFi_wEZVn53qcHEaUEHAzg8rPiy0ZsjFCQ"
)

var (
	green  = color.New(color.FgGreen, color.Bold)
	cyan   = color.New(color.FgCyan)
	red    = color.New(color.FgRed, color.Bold)
	yellow = color.New(color.FgYellow)
	dim    = color.New(color.Faint)
)

func main() {
	if len(os.Args) < 3 {
		red.Println("ERROR: username and room name required")
		fmt.Println("usage: kainet <username> <room-name>")
		os.Exit(1)
	}

	username := os.Args[1]
	roomName := os.Args[2]

	// check embedded credentials first, fall back to env vars
	dbURL := embeddedDBURL
	authToken := embeddedAuthToken

	if dbURL == "" {
		dbURL = os.Getenv("DB_URL")
	}
	if authToken == "" {
		authToken = os.Getenv("AUTH_TOKEN")
	}

	if dbURL == "" || authToken == "" {
		red.Println("SECURITY ERROR: credentials not found")
		fmt.Println("set embeddedDBURL/embeddedAuthToken in main.go or use DB_URL/AUTH_TOKEN environment variables")
		os.Exit(1)
	}

	// Boot sequence
	printBanner()
	bootSequence(username, roomName)

	// Connect to Turso
	cyan.Print(">>> ")
	fmt.Print("establishing secure connection...")
	connStr := dbURL + "?authToken=" + authToken
	db, err := sql.Open("libsql", connStr)
	if err != nil {
		red.Printf("\nCONNECTION FAILED: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	green.Println(" CONNECTED")

	// Initialize database schema
	cyan.Print(">>> ")
	fmt.Print("initializing secure channel...")
	if err := initDB(db); err != nil {
		red.Printf("\nINITIALIZATION FAILED: %v\n", err)
		os.Exit(1)
	}
	green.Println(" OK")

	// Derive encryption key from auth token
	cyan.Print(">>> ")
	fmt.Print("generating AES-256 encryption keys...")
	encKey := deriveKey(authToken)
	green.Println(" READY")

	fmt.Println()
	green.Println("╔════════════════════════════════════════╗")
	green.Printf("║  SECURE CHANNEL ACTIVE                 ║\n")
	green.Printf("║  OPERATOR: %-28s║\n", strings.ToUpper(username))
	green.Printf("║  ROOM: %-32s║\n", strings.ToUpper(roomName))
	green.Printf("║  ENCRYPTION: AES-256-GCM               ║\n")
	green.Println("╚════════════════════════════════════════╗")
	fmt.Println()
	dim.Println("commands: /burn (wipe history)")
	fmt.Println()

	// Display recent history
	if err := displayHistory(db, encKey, roomName, 20); err != nil {
		fmt.Printf("error loading history: %v\n", err)
	}

	// Track last seen message ID
	lastID := getLastMessageID(db, roomName)

	// Start polling for new messages in background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				newLastID, err := pollNewMessages(db, encKey, roomName, lastID)
				if err != nil {
					fmt.Printf("\nerror polling messages: %v\n", err)
				} else if newLastID > lastID {
					lastID = newLastID
				}
			}
		}
	}()

	// Read user input and send messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())
		if message == "" {
			continue
		}

		// Handle commands
		if strings.HasPrefix(message, "/") {
			if err := handleCommand(db, message, &lastID, username, roomName); err != nil {
				// Clear input line
				fmt.Print("\033[A\033[2K\r")
				red.Printf("COMMAND FAILED: %v\n", err)
			}
			continue
		}

		newID, err := sendMessage(db, encKey, username, roomName, message)
		if err != nil {
			red.Printf("TRANSMISSION FAILED: %v\n", err)
		} else {
			// Clear the input line and replace with formatted message
			// \033[A moves cursor up one line
			// \033[2K clears the entire line
			// \r returns to start of line
			fmt.Print("\033[A\033[2K\r")
			t := time.Now().Format("15:04")
			dim.Print("[" + t + "] ")
			yellow.Print(username + ": ")
			fmt.Println(message)
			// Update lastID so polling doesn't re-display this message
			if newID > lastID {
				lastID = newID
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error reading input: %v\n", err)
	}
}

func initDB(db *sql.DB) error {
	// Check if table exists and has old schema
	var tableExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) > 0
		FROM sqlite_master
		WHERE type='table' AND name='messages'
	`).Scan(&tableExists)
	if err != nil {
		return err
	}

	if tableExists {
		// Check if room_name column exists
		var hasRoomName bool
		err = db.QueryRow(`
			SELECT COUNT(*) > 0
			FROM pragma_table_info('messages')
			WHERE name='room_name'
		`).Scan(&hasRoomName)
		if err != nil {
			return err
		}

		// Migrate old schema by adding room_name column
		if !hasRoomName {
			_, err = db.Exec(`ALTER TABLE messages ADD COLUMN room_name TEXT DEFAULT 'default'`)
			if err != nil {
				return fmt.Errorf("migration failed: %v", err)
			}
		}
	} else {
		// Create new table with current schema
		_, err = db.Exec(`
			CREATE TABLE messages (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				room_name TEXT NOT NULL,
				username TEXT NOT NULL,
				message TEXT NOT NULL,
				timestamp INTEGER NOT NULL
			)
		`)
		if err != nil {
			return err
		}
	}

	// Create index for efficient room-based queries
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_room_id ON messages(room_name, id)
	`)
	return err
}

func printBanner() {
	green.Println(`
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
╚══════════════════════════════════════════════════════════╝
	`)
}

func bootSequence(username, roomName string) {
	cyan.Println("SYSTEM INITIALIZING...")
	time.Sleep(200 * time.Millisecond)

	steps := []string{
		"loading security protocols",
		"verifying operator credentials",
		fmt.Sprintf("connecting to room: %s", roomName),
		"establishing quantum-resistant handshake",
		"activating end-to-end encryption",
	}

	for _, step := range steps {
		cyan.Print(">>> ")
		fmt.Print(step)
		animate(step)
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println()
}

func animate(step string) {
	frames := []string{".", "..", "..."}
	for i := 0; i < 3; i++ {
		for _, frame := range frames {
			fmt.Printf("\r>>> %s%s   ", step, frame)
			time.Sleep(50 * time.Millisecond)
		}
	}
	fmt.Print("\r")
	cyan.Print(">>> ")
	fmt.Print(step)
	green.Println(" OK")
}

func displayHistory(db *sql.DB, encKey []byte, roomName string, limit int) error {
	rows, err := db.Query(`
		SELECT username, message, timestamp
		FROM messages
		WHERE room_name = ?
		ORDER BY id DESC
		LIMIT ?
	`, roomName, limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	var messages []struct {
		username  string
		message   string
		timestamp int64
	}

	for rows.Next() {
		var m struct {
			username  string
			message   string
			timestamp int64
		}
		var encryptedMsg string
		if err := rows.Scan(&m.username, &encryptedMsg, &m.timestamp); err != nil {
			return err
		}

		// Decrypt message
		decrypted, err := decrypt(encKey, encryptedMsg)
		if err != nil {
			// If decryption fails, show error indicator
			m.message = "[decrypt error]"
		} else {
			m.message = decrypted
		}

		messages = append(messages, m)
	}

	// Display in chronological order
	for i := len(messages) - 1; i >= 0; i-- {
		m := messages[i]
		t := time.Unix(m.timestamp, 0).Format("15:04")
		dim.Print("[" + t + "] ")
		cyan.Print(m.username + ": ")
		fmt.Println(m.message)
	}

	return rows.Err()
}

func getLastMessageID(db *sql.DB, roomName string) int64 {
	var id int64
	db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM messages WHERE room_name = ?", roomName).Scan(&id)
	return id
}

func pollNewMessages(db *sql.DB, encKey []byte, roomName string, lastID int64) (int64, error) {
	rows, err := db.Query(`
		SELECT id, username, message, timestamp
		FROM messages
		WHERE room_name = ? AND id > ?
		ORDER BY id ASC
	`, roomName, lastID)
	if err != nil {
		return lastID, err
	}
	defer rows.Close()

	newLastID := lastID
	for rows.Next() {
		var (
			id           int64
			username     string
			encryptedMsg string
			timestamp    int64
		)
		if err := rows.Scan(&id, &username, &encryptedMsg, &timestamp); err != nil {
			return newLastID, err
		}

		// Decrypt message
		message, err := decrypt(encKey, encryptedMsg)
		if err != nil {
			message = "[DECRYPT ERROR]"
		}

		t := time.Unix(timestamp, 0).Format("15:04")
		dim.Print("[" + t + "] ")
		cyan.Print(username + ": ")
		fmt.Println(message)

		if id > newLastID {
			newLastID = id
		}
	}

	return newLastID, rows.Err()
}

func handleCommand(db *sql.DB, cmd string, lastID *int64, username, roomName string) error {
	// Clear input line
	fmt.Print("\033[A\033[2K\r")

	switch cmd {
	case "/burn":
		yellow.Print("WARNING: purging all encrypted transmissions for this room...")
		time.Sleep(500 * time.Millisecond)
		_, err := db.Exec("DELETE FROM messages WHERE room_name = ?", roomName)
		if err != nil {
			return fmt.Errorf("purge failed: %v", err)
		}
		green.Println(" PURGED")
		dim.Println("all message history for this room has been permanently erased")
		*lastID = 0
		return nil
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

func deriveKey(authToken string) []byte {
	hash := sha256.Sum256([]byte(authToken))
	return hash[:]
}

func encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func sendMessage(db *sql.DB, encKey []byte, username, roomName, message string) (int64, error) {
	// Encrypt the message before storing
	encryptedMsg, err := encrypt(encKey, message)
	if err != nil {
		return 0, fmt.Errorf("encryption failed: %v", err)
	}

	result, err := db.Exec(`
		INSERT INTO messages (room_name, username, message, timestamp)
		VALUES (?, ?, ?, ?)
	`, roomName, username, encryptedMsg, time.Now().Unix())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
