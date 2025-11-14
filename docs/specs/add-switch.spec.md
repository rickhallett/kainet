# Implementation Specification: /switch Command

## Overview

Add a `/switch <room-name>` command that allows users to change chat rooms without restarting the terminal. The command should cleanly exit the current room, clear the screen, replay the boot sequence, and enter the new room seamlessly.

## Requirements

### Functional Requirements

1. **Command Syntax**
   - `/switch <room-name>` - Switch to specified room
   - Room name validation (non-empty, trimmed)
   - Case-sensitive room names

2. **User Experience**
   - Clear terminal screen
   - Gracefully stop current room's message polling
   - Replay boot sequence animation with new room name
   - Display last 20 messages from new room
   - Maintain same username across rooms

3. **State Management**
   - Preserve database connection
   - Preserve encryption key
   - Reset message tracking (lastID)
   - Clean goroutine lifecycle

4. **Error Handling**
   - Invalid room name (empty or whitespace)
   - Database connection issues
   - Graceful degradation on errors

### Non-Functional Requirements

1. **Performance**
   - Room switch should complete in < 2 seconds
   - No memory leaks from abandoned goroutines
   - No race conditions between old/new polling loops

2. **Security**
   - Maintain encryption throughout switch
   - No data leakage between rooms
   - Room isolation remains intact

## Technical Design

### Architecture Changes

#### 1. Refactor Main Loop Structure

**Current state (main.go:113-166):**
- Single monolithic loop in main()
- Direct goroutine spawning
- No abstraction for room lifecycle

**Proposed structure:**
```go
type ChatRoomResult struct {
    Action  string // "exit" or "switch"
    NewRoom string // only used when Action == "switch"
}

func main() {
    // ... existing setup code ...

    currentRoom := roomName

    for {
        result := runChatRoom(ctx, db, encKey, username, currentRoom)

        switch result.Action {
        case "exit":
            return
        case "switch":
            currentRoom = result.NewRoom
            clearScreen()
            showBootSequence(username, currentRoom)
            continue
        }
    }
}

func runChatRoom(parentCtx context.Context, db *sql.DB, encKey []byte,
                 username, roomName string) ChatRoomResult {
    // Extracted main loop logic
}
```

#### 2. Goroutine Lifecycle Management

**Key considerations:**
- Use context cancellation for clean shutdown
- Create child context for each room
- Cancel context before switching rooms
- Wait for goroutine to acknowledge cancellation

**Implementation approach:**
```go
func runChatRoom(parentCtx context.Context, ...) ChatRoomResult {
    // Create room-specific context
    ctx, cancel := context.WithCancel(parentCtx)
    defer cancel()

    // Goroutine with proper shutdown
    done := make(chan struct{})
    go func() {
        defer close(done)
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                return // Clean exit
            case <-ticker.C:
                // Poll for messages
            }
        }
    }()

    // Main input loop
    // ...

    // When switching:
    cancel()
    <-done // Wait for goroutine to finish

    return ChatRoomResult{Action: "switch", NewRoom: newRoomName}
}
```

#### 3. Screen Management

**Functions to implement:**

```go
// Clear screen using ANSI escape codes
func clearScreen() {
    fmt.Print("\033[H\033[2J")
}

// Extract boot sequence to reusable function
func showBootSequence(username, roomName string) {
    // Lines 250-268 from current main()
    printBanner()
    bootSequence()

    // Display room header
    green.Println("╔════════════════════════════════════════╗")
    green.Printf("║  SECURE CHANNEL ACTIVE                 ║\n")
    green.Printf("║  OPERATOR: %-28s║\n", strings.ToUpper(username))
    green.Printf("║  ROOM: %-32s║\n", strings.ToUpper(roomName))
    green.Printf("║  ENCRYPTION: AES-256-GCM               ║\n")
    green.Println("╚════════════════════════════════════════╗")
    fmt.Println()
}
```

#### 4. Command Handler Extension

**Modify handleCommand function (main.go:390-432):**

```go
func handleCommand(db *sql.DB, cmd string, lastID *int64, username,
                   roomName string) (ChatRoomResult, error) {
    // Clear input line
    fmt.Print("\033[A\033[2K\r")

    parts := strings.Fields(cmd)
    command := parts[0]

    switch command {
    case "/burn":
        // ... existing burn logic ...
        return ChatRoomResult{}, nil

    case "/exit":
        // ... existing exit logic ...
        // os.Exit(0) or return special result

    case "/switch":
        if len(parts) < 2 {
            return ChatRoomResult{}, fmt.Errorf("usage: /switch <room-name>")
        }

        newRoom := strings.TrimSpace(strings.Join(parts[1:], " "))
        if newRoom == "" {
            return ChatRoomResult{}, fmt.Errorf("room name cannot be empty")
        }

        if newRoom == roomName {
            dim.Println("already in room: " + roomName)
            return ChatRoomResult{}, nil
        }

        yellow.Printf("switching to room: %s...\n", newRoom)
        time.Sleep(300 * time.Millisecond)

        return ChatRoomResult{
            Action:  "switch",
            NewRoom: newRoom,
        }, nil

    default:
        return ChatRoomResult{}, fmt.Errorf("unknown command: %s", cmd)
    }
}
```

## Implementation Steps

### Phase 1: Refactoring (30-40 min)

1. **Extract boot sequence function**
   - Create `showBootSequence(username, roomName string)`
   - Move lines 250-268 into function
   - Test that startup still works

2. **Create ChatRoomResult type**
   - Define struct in main.go
   - Add before main() function

3. **Extract main loop to runChatRoom()**
   - Create function signature
   - Move lines 113-166 into function
   - Return ChatRoomResult instead of breaking loop
   - Update main() to call runChatRoom() in loop

4. **Add clearScreen() utility**
   - Simple ANSI escape code implementation

### Phase 2: Goroutine Management (20-30 min)

1. **Add context parameter**
   - Pass context to runChatRoom()
   - Create child context inside function
   - Defer cancel()

2. **Update polling goroutine**
   - Add done channel
   - Close done channel on exit
   - Use defer close(done)

3. **Wait for goroutine on switch**
   - Cancel context
   - Wait for <-done
   - Ensures clean shutdown

### Phase 3: Command Implementation (15-20 min)

1. **Update handleCommand signature**
   - Return (ChatRoomResult, error)
   - Update all call sites

2. **Add /switch case**
   - Parse room name from arguments
   - Validate input
   - Return switch result

3. **Update command help text**
   - Add /switch to help message (line 102)

### Phase 4: Integration (10-15 min)

1. **Wire up main loop**
   - Add for loop in main()
   - Handle switch result
   - Call clearScreen() and showBootSequence()

2. **Update existing commands**
   - Ensure /burn returns empty result
   - Ensure /exit properly terminates

### Phase 5: Testing (30-45 min)

1. **Basic functionality**
   - Switch to new room
   - Verify messages display correctly
   - Switch back to original room

2. **Edge cases**
   - Empty room name
   - Same room name
   - Special characters in room name
   - Rapid switching

3. **Goroutine cleanup**
   - Use `go test -race` to check for race conditions
   - Monitor goroutine count
   - Verify no leaks

4. **Integration tests**
   - Switch while receiving messages
   - Switch immediately after sending message
   - Multiple switches in sequence

## Code Changes Summary

### Files to modify

**main.go:**
- Add ChatRoomResult type (before main)
- Extract showBootSequence() function
- Extract runChatRoom() function
- Modify main() to use loop structure
- Update handleCommand() signature and return type
- Add /switch case in handleCommand()
- Update command help text (line 102)
- Add clearScreen() utility function

**Estimated line changes:**
- +80 lines (new functions and types)
- ~40 lines modified (refactoring existing code)
- Total: ~120 lines changed/added

## Testing Checklist

### Manual Testing

- [ ] Basic switch: `/switch test-room`
- [ ] Switch to same room (should show message)
- [ ] Empty room name: `/switch` (should error)
- [ ] Whitespace room name: `/switch    ` (should error)
- [ ] Room name with spaces: `/switch my room`
- [ ] Switch while messages are incoming
- [ ] Switch after sending a message
- [ ] Multiple rapid switches
- [ ] Switch then /burn
- [ ] Switch then /exit
- [ ] Verify screen clears properly
- [ ] Verify boot animation plays
- [ ] Verify new room messages load

### Race Condition Testing

```bash
go test -race ./...
```

### Memory Leak Testing

```go
// Test goroutine cleanup
func TestNoGoroutineLeaks(t *testing.T) {
    before := runtime.NumGoroutine()

    // Perform multiple room switches

    time.Sleep(2 * time.Second)
    after := runtime.NumGoroutine()

    if after > before+2 { // Allow some margin
        t.Errorf("Goroutine leak detected: before=%d, after=%d", before, after)
    }
}
```

## Edge Cases and Error Handling

### Input Validation

1. **Empty room name**
   - Error: "room name cannot be empty"
   - Do not switch

2. **Same room name**
   - Message: "already in room: <name>"
   - Do not switch

3. **Room name with special characters**
   - Allow all characters (database handles it)
   - Trim whitespace

### Database Errors

1. **Connection lost during switch**
   - Show error message
   - Attempt to reconnect
   - Fall back to exit if unable

2. **Query errors in new room**
   - Show error message
   - Allow retry or switch back

### Concurrency Issues

1. **Message arrives during switch**
   - Context cancellation stops polling
   - No messages lost (they remain in database)
   - New room will show them on load

2. **User sends message during switch**
   - Scanner should be blocked during switch
   - Or handle gracefully with error message

## Performance Considerations

### Benchmarks

Expected performance targets:
- Room switch total time: < 2 seconds
  - Context cancellation: < 100ms
  - Screen clear: < 10ms
  - Boot animation: ~1 second
  - History load: < 500ms

### Optimization Opportunities

1. **Parallel operations**
   - Start history query while animation plays
   - Preload new room data

2. **Animation speed**
   - Option to skip animation on switch
   - Faster animation for subsequent switches

## Future Enhancements

### Potential Additions

1. **Room history**
   - Track recently visited rooms
   - `/switch -` to go back to previous room

2. **Room list**
   - `/rooms` command to list available rooms
   - Shows room names with message counts

3. **Aliases**
   - `/sw` as shorthand for `/switch`

4. **Favorites**
   - Pin frequently used rooms
   - Quick switch with `/switch @1`, `/switch @2`

## Migration Path

### Version Compatibility

1. **v1.2.0 → v1.3.0**
   - Existing commands remain unchanged
   - New /switch command is additive
   - No breaking changes to API

2. **Database schema**
   - No changes required
   - Existing room isolation works as-is

3. **Configuration**
   - No new config required
   - Room name validation same as startup

## Web Client Considerations

### Parallel Implementation

The web client (web/src/main.js) will need similar changes:

1. **State management**
   - Add `currentRoom` to component state
   - Track room history in state

2. **Switch handler**
   ```javascript
   function handleSwitch(newRoom) {
       setCurrentRoom(newRoom);
       clearMessages();
       playBootAnimation(newRoom);
       loadRoomHistory(newRoom);
   }
   ```

3. **Polling updates**
   - Cancel existing interval
   - Start new interval for new room
   - Update room name in UI

**Web client complexity: LOWER** (4-5/10)
- No goroutine management
- React handles state transitions naturally
- Animation replay is simpler

## References

### Related Code Locations

- Main loop: main.go:113-166
- Command handler: main.go:390-432
- Boot sequence: main.go:250-268
- Message polling: main.go:117-133
- Help text: main.go:102

### Dependencies

- No new dependencies required
- Uses existing context package
- Standard library only

## Approval Checklist

Before implementing:

- [ ] Review architecture design
- [ ] Confirm goroutine cleanup approach
- [ ] Verify no breaking changes
- [ ] Estimate testing time
- [ ] Plan rollback strategy
- [ ] Document user-facing changes

## Estimated Effort

**Total implementation time: 2-3 hours**
- Refactoring: 30-40 minutes
- Goroutine management: 20-30 minutes
- Command implementation: 15-20 minutes
- Integration: 10-15 minutes
- Testing: 30-45 minutes

**Complexity rating: 6-8/10 (MODERATE)**

## Success Criteria

Implementation is complete when:

1. User can switch rooms with `/switch <name>`
2. Screen clears and boot animation replays
3. New room messages load correctly
4. No goroutine leaks detected
5. No race conditions in tests
6. All edge cases handled gracefully
7. Command help updated
8. Documentation updated
