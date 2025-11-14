package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

const (
	dbURL      = "libsql://bt-phone-home-rickhallett.aws-eu-west-1.turso.io?authToken=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NjMwOTg5MjEsImlkIjoiNjc4YjAxZDQtYzExZi00OTMyLTk2MzktZDUxODNlZTVmMTI3IiwicmlkIjoiNzc0ZDk3NmItNTMzYS00NzE1LTlkZmItY2RmYjU1N2M2ZmRjIn0.4hZBJRmkoSMutMqpGVYhmIWmo2-5lTKle9p5QkSyF6B-OHGS1xpdoFi_wEZVn53qcHEaUEHAzg8rPiy0ZsjFCQ"
	authToken  = "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NjMwOTg5MjEsImlkIjoiNjc4YjAxZDQtYzExZi00OTMyLTk2MzktZDUxODNlZTVmMTI3IiwicmlkIjoiNzc0ZDk3NmItNTMzYS00NzE1LTlkZmItY2RmYjU1N2M2ZmRjIn0.4hZBJRmkoSMutMqpGVYhmIWmo2-5lTKle9p5QkSyF6B-OHGS1xpdoFi_wEZVn53qcHEaUEHAzg8rPiy0ZsjFCQ"
	pollInterval = 1 * time.Second
)

var (
	green      = color.New(color.FgGreen).SprintFunc()
	greenBold  = color.New(color.FgGreen, color.Bold).SprintFunc()
	redBold    = color.New(color.FgRed, color.Bold).SprintFunc()
	yellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	cyanBold   = color.New(color.FgCyan, color.Bold).SprintFunc()
)

type Message struct {
	ID        int64
	RoomName  string
	Username  string
	Message   string
	Timestamp time.Time
}

func printBanner() {
	banner := `
    ██╗  ██╗ █████╗ ██╗███╗   ██╗███████╗████████╗
    ██║ ██╔╝██╔══██╗██║████╗  ██║██╔════╝╚══██╔══╝
    █████╔╝ ███████║██║██╔██╗ ██║█████╗     ██║
    ██╔═██╗ ██╔══██║██║██║╚██╗██║██╔══╝     ██║
    ██║  ██╗██║  ██║██║██║ ╚████║███████╗   ██║
    ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝   ╚═╝
    `
	fmt.Println(greenBold(banner))
}

func bootSequence() {
	sequences := []string{
		"[SYSTEM] Initializing secure connection...",
		"[CRYPTO] Loading AES-256-GCM encryption...",
		"[DATABASE] Connecting to Turso libSQL...",
		"[NETWORK] Establishing encrypted channel...",
		"[READY] kaiNET operational",
	}

	for _, seq := range sequences {
		fmt.Println(green(seq))
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println()
}

func deriveKey(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}

func encrypt(plaintext string, key []byte) (string, error) {
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

func decrypt(ciphertext string, key []byte) (string, error) {
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

	nonce, cipherData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func initDB(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_name TEXT NOT NULL,
		username TEXT NOT NULL,
		message TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_room_timestamp ON messages(room_name, timestamp);
	`
	_, err := db.Exec(schema)
	return err
}

func sendMessage(db *sql.DB, roomName, username, message string, key []byte) error {
	encrypted, err := encrypt(message, key)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO messages (room_name, username, message) VALUES (?, ?, ?)",
		roomName, username, encrypted,
	)
	return err
}

func getMessages(db *sql.DB, roomName string, since time.Time, key []byte) ([]Message, error) {
	rows, err := db.Query(
		`SELECT id, room_name, username, message, timestamp
		 FROM messages
		 WHERE room_name = ? AND timestamp > ?
		 ORDER BY timestamp ASC`,
		roomName, since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		var encryptedMsg string
		var timestampStr string

		err := rows.Scan(&m.ID, &m.RoomName, &m.Username, &encryptedMsg, &timestampStr)
		if err != nil {
			continue
		}

		decrypted, err := decrypt(encryptedMsg, key)
		if err != nil {
			continue
		}

		m.Message = decrypted
		m.Timestamp, _ = time.Parse("2006-01-02 15:04:05", timestampStr)
		messages = append(messages, m)
	}

	return messages, nil
}

func burnRoom(db *sql.DB, roomName string) error {
	_, err := db.Exec("DELETE FROM messages WHERE room_name = ?", roomName)
	return err
}

func displayMessage(msg Message, currentUser string) {
	timestamp := msg.Timestamp.Format("15:04:05")
	if msg.Username == currentUser {
		fmt.Printf("%s %s: %s\n",
			green(timestamp),
			cyanBold(msg.Username),
			msg.Message,
		)
	} else {
		fmt.Printf("%s %s: %s\n",
			green(timestamp),
			yellowBold(msg.Username),
			msg.Message,
		)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <username> <room-name>\n", os.Args[0])
		os.Exit(1)
	}

	username := os.Args[1]
	roomName := os.Args[2]

	printBanner()
	bootSequence()

	fmt.Printf("%s Connecting as %s to room %s\n\n",
		green("[INFO]"),
		cyanBold(username),
		yellowBold(roomName),
	)

	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		fmt.Printf("%s Failed to connect to database: %v\n", redBold("[ERROR]"), err)
		os.Exit(1)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		fmt.Printf("%s Failed to initialize database: %v\n", redBold("[ERROR]"), err)
		os.Exit(1)
	}

	key := deriveKey(authToken)
	lastCheck := time.Now()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	inputChan := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			inputChan <- strings.TrimSpace(text)
		}
	}()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	fmt.Println(green("Connected. Type your message and press Enter. Use /burn to wipe room history."))
	fmt.Println(green("Press Ctrl+C to exit."))
	fmt.Println()

	for {
		select {
		case <-sigChan:
			fmt.Println("\n" + green("[SYSTEM] Disconnecting..."))
			return

		case <-ticker.C:
			messages, err := getMessages(db, roomName, lastCheck, key)
			if err == nil && len(messages) > 0 {
				for _, msg := range messages {
					displayMessage(msg, username)
				}
				lastCheck = time.Now()
			}

		case input := <-inputChan:
			if input == "" {
				continue
			}

			if input == "/burn" {
				if err := burnRoom(db, roomName); err != nil {
					fmt.Printf("%s Failed to burn room: %v\n", redBold("[ERROR]"), err)
				} else {
					fmt.Println(redBold("[BURN] Room history wiped"))
				}
				continue
			}

			if err := sendMessage(db, roomName, username, input, key); err != nil {
				fmt.Printf("%s Failed to send message: %v\n", redBold("[ERROR]"), err)
			}
		}
	}
}
