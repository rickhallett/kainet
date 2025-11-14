# kaiNET Development Roadmap

## Current State (v1.2.0)

- AES-256-GCM end-to-end encryption
- room-based message isolation
- cross-platform CLI (macOS, Linux, Windows)
- web client with CRT terminal aesthetic
- commands: /burn, /exit
- one-shot installation scripts
- ~1 second message delivery

## Upcoming Release (v1.3.0) - Room Switching

**Status**: specification complete, ready for implementation

**Features**:
- `/switch <room-name>` command
- seamless room transitions without restart
- context-based goroutine management
- boot sequence replay on switch

**Estimated effort**: 2-3 hours
**Complexity**: 6-8/10 (MODERATE)

**See**: docs/specs/add-switch.spec.md

---

## Short Term (v1.4.0 - v1.5.0)

### v1.4.0 - Enhanced Room Management

**Priority**: HIGH
**Estimated effort**: 3-4 hours

**Features**:
1. **Room discovery**
   - `/rooms` command to list available rooms
   - show message count and last activity per room
   - filter by rooms user has visited

2. **Room history**
   - track recently visited rooms
   - `/switch -` to return to previous room
   - persist room history locally

3. **Room aliases**
   - `/alias <shortcut> <room-name>`
   - quick switch with `/sw @work`, `/sw @home`
   - list aliases with `/aliases`

4. **Room metadata**
   - display active users count
   - show room creation date
   - room description field

**Technical considerations**:
- add local config file (~/.kainet/config.json)
- new database table for room metadata
- cache room list for performance

---

### v1.5.0 - Message Features

**Priority**: MEDIUM
**Estimated effort**: 4-5 hours

**Features**:
1. **Message reactions**
   - `/react <message-id> <emoji>` or simple +1/-1
   - display reactions inline with messages
   - remove reaction with `/unreact`

2. **Message editing**
   - `/edit <message-id> <new-text>`
   - show "edited" indicator
   - keep edit history in database

3. **Message search**
   - `/search <query>` to find messages
   - search within current room
   - search across all rooms with `/search --all <query>`

4. **Message threading**
   - `/reply <message-id> <text>`
   - show thread indicator
   - `/thread <message-id>` to view full thread

5. **Message formatting**
   - support markdown: **bold**, *italic*, `code`
   - multiline messages with shift+enter (web client)
   - code blocks with syntax highlighting

**Technical considerations**:
- message schema changes (add parent_id, reactions, edited_at)
- markdown rendering library for web client
- message ID display in terminal (dim, small)

---

## Medium Term (v1.6.0 - v1.8.0)

### v1.6.0 - User Presence & Status

**Priority**: MEDIUM
**Estimated effort**: 5-6 hours

**Features**:
1. **Active user tracking**
   - heartbeat mechanism (update every 30s)
   - display active users in room header
   - show when users join/leave room

2. **User status**
   - `/status <message>` to set custom status
   - `/dnd` to silence notifications
   - `/away` to mark as inactive

3. **Typing indicators**
   - show "user is typing..."
   - debounced, disappears after 3s
   - optional, can be disabled

4. **Read receipts**
   - track last read message per user
   - show unread count on room list
   - `/mark-read` to mark all as read

**Technical considerations**:
- new table: user_presence (username, room, last_seen, status)
- polling interval for presence updates
- cleanup stale presence records

---

### v1.7.0 - File Sharing & Media

**Priority**: MEDIUM-LOW
**Estimated effort**: 8-10 hours

**Features**:
1. **File uploads**
   - `/upload <file-path>` to share files
   - encrypt files before upload
   - store in turso blob storage or R2

2. **Image preview**
   - display images inline (web client)
   - show image thumbnails (CLI as ASCII art?)
   - `/download <file-id>` to retrieve

3. **File expiration**
   - auto-delete files after 7 days
   - `/burn` also removes files
   - configurable retention policy

4. **File size limits**
   - max 10MB per file
   - show upload progress
   - compress images automatically

**Technical considerations**:
- file storage strategy (turso blobs vs external CDN)
- encryption for files (AES-256-GCM)
- bandwidth costs on free tier
- ASCII art library for CLI image preview

---

### v1.8.0 - Advanced Security Features

**Priority**: MEDIUM
**Estimated effort**: 6-8 hours

**Features**:
1. **Password-protected rooms**
   - require password to join room
   - PBKDF2 key derivation
   - store hashed room passwords

2. **Room permissions**
   - room admin role
   - admin can kick users
   - admin can lock room (read-only)

3. **Message expiration**
   - `/expire <time>` to set auto-delete timer
   - messages auto-delete after timer
   - "disappearing messages" mode

4. **Encrypted backups**
   - `/backup <password>` to export messages
   - encrypted JSON export
   - `/restore <file> <password>` to import

5. **Audit log**
   - track room joins, leaves, burns
   - admin can view with `/audit`
   - tamper-evident log

**Technical considerations**:
- add permissions table
- background job for message expiration
- backup file format (JSON, encrypted)

---

## Long Term (v2.0.0+)

### v2.0.0 - Mobile Apps

**Priority**: LOW
**Estimated effort**: 40-60 hours (major undertaking)

**Features**:
1. **Native iOS app**
   - Swift/SwiftUI implementation
   - push notifications
   - keychain integration for credentials

2. **Native Android app**
   - Kotlin/Compose implementation
   - firebase cloud messaging
   - encrypted shared preferences

3. **Cross-platform (Flutter/React Native)**
   - single codebase for iOS/Android
   - reuse encryption logic
   - share turso backend

**Technical considerations**:
- mobile-optimized UI
- battery efficiency (polling vs push)
- app store distribution
- mobile-specific features (camera, location sharing)

---

### v2.1.0 - Voice & Video

**Priority**: LOW
**Estimated effort**: 30-40 hours

**Features**:
1. **Voice messages**
   - record and send audio clips
   - play in-app with waveform visualization
   - transcription with whisper API

2. **Voice calls**
   - WebRTC peer-to-peer calls
   - end-to-end encrypted
   - call initiation with `/call <username>`

3. **Video calls**
   - WebRTC video support
   - screen sharing
   - picture-in-picture mode

**Technical considerations**:
- WebRTC signaling server
- STUN/TURN servers for NAT traversal
- bandwidth requirements
- not compatible with CLI (web/mobile only)

---

### v2.2.0 - Advanced Integrations

**Priority**: LOW
**Estimated effort**: 10-15 hours

**Features**:
1. **Webhooks**
   - send messages via HTTP POST
   - integrate with automation tools
   - GitHub/GitLab notifications

2. **Bot API**
   - create bots with simple API
   - custom commands
   - scheduled messages

3. **Bridge integrations**
   - bridge to IRC, Discord, Slack
   - two-way message sync
   - unified interface

4. **Email notifications**
   - email digest of missed messages
   - configurable frequency
   - reply via email to post message

**Technical considerations**:
- webhook authentication
- rate limiting
- bridge architecture (separate service)

---

## Developer Experience Improvements

### Immediate Priority

1. **Testing infrastructure**
   - implement agentic testing workflow
   - unit tests for encryption
   - integration tests for commands
   - CI/CD pipeline with GitHub Actions

2. **Development documentation**
   - architecture overview
   - database schema documentation
   - API documentation
   - contribution guide

3. **Local development**
   - docker-compose for local turso
   - hot reload for web client
   - CLI debug mode with verbose logging

4. **Release automation**
   - automated binary builds
   - changelog generation
   - version bump script
   - deploy to homebrew tap

### Future Enhancements

1. **Plugin system**
   - load plugins from ~/.kainet/plugins
   - custom commands
   - message transformers
   - theme system

2. **Configuration management**
   - ~/.kainet/config.toml
   - per-room settings
   - keybindings customization
   - color scheme customization

---

## Performance & Scalability

### Short Term

1. **Optimize polling**
   - exponential backoff when no activity
   - long polling instead of 1s interval
   - WebSocket support for web client

2. **Message pagination**
   - lazy load old messages
   - `/history --before <id>` to load more
   - infinite scroll in web client

3. **Database optimization**
   - add indexes for common queries
   - archive old messages
   - vacuum/optimize commands

### Long Term

1. **Self-hosted option**
   - docker image with embedded sqlite
   - kubernetes helm chart
   - migration guide from turso

2. **Federation**
   - connect multiple kaiNET instances
   - ActivityPub protocol support
   - cross-instance rooms

3. **Horizontal scaling**
   - multiple database replicas
   - load balancer for web client
   - redis for presence/typing indicators

---

## Platform Support

### Short Term

1. **Package managers**
   - homebrew formula
   - apt/deb package
   - chocolatey for Windows
   - snap package

2. **ARM support**
   - Raspberry Pi builds
   - ARM Linux support
   - test on various architectures

### Long Term

1. **Browser extension**
   - chrome/firefox extension
   - quick access from toolbar
   - notification support

2. **Desktop apps**
   - Electron wrapper for web client
   - native menubar app (macOS)
   - system tray integration (Windows/Linux)

---

## Documentation & Community

### Immediate

1. **User guide**
   - getting started tutorial
   - command reference
   - troubleshooting guide
   - FAQ

2. **Video tutorials**
   - installation walkthrough
   - basic usage demo
   - advanced features showcase

### Future

1. **Community features**
   - public demo instance
   - showcase page with screenshots
   - blog with use cases
   - community Discord/forum

2. **Security audit**
   - third-party security review
   - penetration testing
   - bug bounty program

---

## Release Schedule (Proposed)

| Version | Target Date | Focus |
|---------|-------------|-------|
| v1.3.0  | Q1 2025     | room switching |
| v1.4.0  | Q1 2025     | room management |
| v1.5.0  | Q2 2025     | message features |
| v1.6.0  | Q2 2025     | user presence |
| v1.7.0  | Q3 2025     | file sharing |
| v1.8.0  | Q3 2025     | advanced security |
| v2.0.0  | Q4 2025     | mobile apps |

---

## Key Metrics to Track

1. **Adoption**
   - GitHub stars
   - downloads per release
   - active installations (telemetry opt-in)

2. **Performance**
   - message delivery latency
   - database query times
   - binary size

3. **Reliability**
   - crash reports
   - error rates
   - uptime

4. **Security**
   - vulnerability reports
   - encryption audit status
   - security incident response time

---

## Decision Points

### Architecture Decisions Needed

1. **Polling vs WebSockets**
   - Current: 1s polling
   - Option A: keep polling, add exponential backoff
   - Option B: migrate to WebSockets (requires server)
   - Option C: hybrid (WebSocket for web, polling for CLI)
   - **Recommendation**: Option C for best of both worlds

2. **File Storage**
   - Option A: turso blob storage
   - Option B: cloudflare R2
   - Option C: IPFS distributed storage
   - **Recommendation**: Option B (R2) for cost/performance

3. **Mobile Strategy**
   - Option A: native apps (Swift + Kotlin)
   - Option B: React Native
   - Option C: Flutter
   - **Recommendation**: Option A for best performance, or C for faster development

4. **Authentication**
   - Current: shared credentials in binary
   - Option A: keep current (simplest)
   - Option B: user accounts with username/password
   - Option C: magic link authentication
   - **Recommendation**: Option A for now, B for v2.0

---

## Community Input Needed

Before committing to major features, gather feedback on:

1. Most wanted features (survey/poll)
2. Pain points with current version
3. Use cases and workflows
4. Platform priorities (mobile vs desktop vs web)
5. Security requirements (compliance needs)

---

## Notes

- maintain backward compatibility where possible
- prioritize security over convenience
- keep CLI and web clients feature-equivalent
- preserve lightweight philosophy (no bloat)
- all features should work on free tier of turso
- document breaking changes clearly
- provide migration guides for major versions

---

## Quick Wins (Low Effort, High Impact)

1. **Command aliases**: `/sw` for `/switch`, `/b` for `/burn` - 30 min
2. **Message timestamps**: show time for each message - 1 hour
3. **Color themes**: support light/dark terminal themes - 2 hours
4. **Export messages**: `/export` to save chat as text file - 2 hours
5. **Username colors**: assign consistent colors per user - 1 hour
6. **Notification sound**: beep on new message (optional) - 1 hour
7. **Unread indicator**: show unread count in room list - 2 hours
8. **Copy message**: `/copy <message-id>` to clipboard - 1 hour
9. **Message count**: show total messages in room header - 30 min
10. **Uptime**: show connection duration in status - 30 min

---

## Experimental Ideas (Research Needed)

1. **E2E encrypted voice notes** using WebRTC data channels
2. **Blockchain-based message verification** for audit trail
3. **Zero-knowledge proofs** for room access without password sharing
4. **Peer-to-peer mode** without database (direct connection)
5. **AI assistant** integration for message summaries
6. **Steganography** to hide messages in images
7. **Tor integration** for anonymous rooms
8. **Multi-party computation** for shared secrets
9. **Dead man's switch** to auto-burn after inactivity
10. **Quantum-resistant encryption** (post-quantum crypto)
