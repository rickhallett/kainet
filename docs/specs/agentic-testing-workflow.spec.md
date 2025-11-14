# Agentic Testing Workflow Specification

## Overview

This specification outlines a comprehensive agentic testing strategy for kaiNET using Claude subagents, MCP browser automation tools (Chrome DevTools and Playwright), and the Turso CLI. The workflow enables automated multi-user interaction testing with full visibility into application state, database operations, and user workflows.

## Available Testing Tools

### 1. Claude Subagents (Task Tool)
**Capabilities:**
- Spawn independent agent processes
- Execute autonomous multi-step tasks
- Run CLI binaries with full control
- Coordinate parallel operations
- Access to all development tools

**Use cases:**
- Launch multiple kaiNET CLI instances
- Simulate different users concurrently
- Execute test scenarios autonomously
- Monitor and verify application behavior

### 2. Chrome DevTools MCP
**Available commands:**
- `mcp__chrome-devtools__navigate_page` - Navigate to URLs
- `mcp__chrome-devtools__take_snapshot` - Capture accessibility tree
- `mcp__chrome-devtools__click` - Click elements by UID
- `mcp__chrome-devtools__fill` - Fill form inputs
- `mcp__chrome-devtools__evaluate_script` - Execute JavaScript
- `mcp__chrome-devtools__list_console_messages` - View console logs
- `mcp__chrome-devtools__list_network_requests` - Monitor network traffic
- `mcp__chrome-devtools__wait_for` - Wait for text to appear

**Use cases:**
- Test web client interface
- Verify message display
- Monitor encryption operations
- Capture UI state at each step
- Validate user interactions

### 3. Playwright MCP
**Available commands:**
- `mcp__playwright__browser_navigate` - Navigate to pages
- `mcp__playwright__browser_snapshot` - Capture page state
- `mcp__playwright__browser_click` - Interact with elements
- `mcp__playwright__browser_type` - Type into inputs
- `mcp__playwright__browser_console_messages` - View console
- `mcp__playwright__browser_network_requests` - Monitor network
- `mcp__playwright__browser_tabs` - Manage multiple tabs

**Use cases:**
- Multi-tab testing (multiple users in browser)
- Automated form submissions
- Screenshot comparison
- Network request validation

### 4. Turso CLI
**Available commands:**
```bash
turso db shell <database> # Interactive SQL
turso db show <database>  # Database info
turso db tokens create    # Create tokens
```

**Use cases:**
- Verify message encryption in database
- Check room isolation
- Monitor message creation timestamps
- Validate data integrity
- Inspect database state between operations

## Testing Architecture

### High-Level Workflow

```
┌─────────────────────────────────────────────────────────────┐
│                    Orchestrator Agent                        │
│  (Main Claude instance - coordinates all testing)           │
└─────────────────┬───────────────────────────────────────────┘
                  │
      ┌───────────┼───────────┬──────────────┐
      │           │           │              │
      ▼           ▼           ▼              ▼
┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
│ Agent 1  │ │ Agent 2  │ │ Browser  │ │ Turso    │
│ (CLI)    │ │ (CLI)    │ │ Testing  │ │ Verifier │
│ User A   │ │ User B   │ │ Agent    │ │ Agent    │
└──────────┘ └──────────┘ └──────────┘ └──────────┘
      │           │           │              │
      └───────────┴───────────┴──────────────┘
                  │
                  ▼
          ┌────────────────┐
          │  Turso DB      │
          │  (Messages)    │
          └────────────────┘
```

### Agent Roles

**1. Orchestrator Agent (Main)**
- Launches all subagents
- Defines test scenarios
- Coordinates timing
- Collects results
- Generates reports

**2. CLI Testing Agents**
- Each represents one user
- Runs kaiNET binary
- Sends/receives messages
- Reports on received messages
- Monitors terminal output

**3. Browser Testing Agent**
- Controls web client
- Verifies UI state
- Captures screenshots
- Monitors console/network
- Tests web-specific features

**4. Database Verification Agent**
- Uses Turso CLI
- Queries database state
- Verifies encryption
- Checks room isolation
- Validates timestamps

## Test Scenarios

### Scenario 1: Two-User Message Exchange

**Objective:** Verify basic message sending and receiving between two CLI users.

**Setup:**
```yaml
Users:
  - alice (CLI instance 1)
  - bob (CLI instance 2)
Room: test-room-001
Duration: 2 minutes
```

**Steps:**

1. **Launch Phase** (Orchestrator)
   ```
   - Spawn Agent 1: Launch alice's CLI
   - Spawn Agent 2: Launch bob's CLI
   - Spawn Agent 3: Monitor database
   - Wait for both users to connect
   ```

2. **Interaction Phase** (Parallel)
   ```
   Agent 1 (alice):
   - Wait 2 seconds
   - Send: "Hello Bob, this is Alice"
   - Monitor for bob's response
   - Report when message received

   Agent 2 (bob):
   - Monitor for alice's message
   - When received, send: "Hi Alice, I got your message"
   - Report success

   Agent 3 (database):
   - Query messages table every second
   - Verify both messages appear
   - Check encryption is applied
   - Verify correct room association
   ```

3. **Verification Phase** (Orchestrator)
   ```
   - Collect reports from all agents
   - Verify message delivery times < 2 seconds
   - Confirm both users saw both messages
   - Check database shows encrypted payloads
   - Generate test report
   ```

**Expected Results:**
- Alice sees bob's message within 2 seconds
- Bob sees alice's message within 2 seconds
- Database contains 2 encrypted messages
- Both messages tagged with test-room-001
- Terminal output shows proper formatting

### Scenario 2: Web Client Multi-Room Testing

**Objective:** Test room switching and isolation in web client.

**Setup:**
```yaml
Browser: Chrome DevTools MCP
Users:
  - charlie (Tab 1)
  - diana (Tab 2)
Rooms:
  - room-alpha
  - room-beta
```

**Steps:**

1. **Setup Phase** (Browser Agent)
   ```javascript
   // Navigate to web client
   mcp__chrome-devtools__navigate_page("http://localhost:5173")

   // Open two tabs
   mcp__playwright__browser_tabs({action: "new"})

   // Tab 1: Login as charlie in room-alpha
   mcp__chrome-devtools__fill({
     fields: [
       {name: "username", value: "charlie"},
       {name: "room", value: "room-alpha"}
     ]
   })

   // Tab 2: Login as diana in room-beta
   // (switch tab and repeat)
   ```

2. **Message Sending** (Browser Agent)
   ```javascript
   // Tab 1 (charlie in room-alpha):
   mcp__chrome-devtools__type({
     element: "message input",
     text: "Message in alpha room"
   })

   // Tab 2 (diana in room-beta):
   mcp__chrome-devtools__type({
     element: "message input",
     text: "Message in beta room"
   })

   // Capture snapshots of both tabs
   mcp__chrome-devtools__take_snapshot()
   ```

3. **Isolation Verification** (Database Agent)
   ```sql
   -- Verify messages are in separate rooms
   SELECT room_name, COUNT(*) FROM messages
   GROUP BY room_name;

   -- Verify charlie's message only in alpha
   SELECT * FROM messages
   WHERE room_name = 'room-alpha'
   AND message LIKE '%charlie%';
   ```

4. **Cross-Room Check** (Browser Agent)
   ```javascript
   // Verify diana doesn't see charlie's message
   mcp__chrome-devtools__wait_for({
     text: "Message in alpha room",
     timeout: 3000
   })
   // Should timeout - message shouldn't appear
   ```

### Scenario 3: /burn Command Verification

**Objective:** Test room history deletion with multiple users.

**Setup:**
```yaml
Users:
  - eve (CLI)
  - frank (CLI)
  - grace (Web)
Room: burn-test-room
```

**Steps:**

1. **Populate Room** (All Agents Parallel)
   ```
   Agent 1 (eve): Send 5 messages
   Agent 2 (frank): Send 5 messages
   Agent 3 (grace): Send 5 messages via web
   Total: 15 messages in room
   ```

2. **Pre-Burn Verification** (Database Agent)
   ```sql
   SELECT COUNT(*) FROM messages
   WHERE room_name = 'burn-test-room';
   -- Should return 15
   ```

3. **Execute Burn** (Agent 1)
   ```
   eve: Type "/burn"
   eve: Wait for confirmation message
   eve: Report success
   ```

4. **Post-Burn Verification** (All Agents)
   ```
   Agent 1 (eve):
   - Verify screen shows "PURGED" message
   - Confirm history cleared from terminal

   Agent 2 (frank):
   - Verify messages disappear from view
   - Check no historical messages load

   Agent 3 (grace):
   - Refresh web page
   - Verify empty message history

   Database Agent:
   SELECT COUNT(*) FROM messages
   WHERE room_name = 'burn-test-room';
   -- Should return 0
   ```

### Scenario 4: /exit Command with Confirmation

**Objective:** Test graceful shutdown with confirmation flow.

**Setup:**
```yaml
User: henry (CLI)
Room: exit-test-room
```

**Steps:**

1. **Setup** (Agent)
   ```
   - Launch henry's CLI in exit-test-room
   - Send one test message
   - Verify connected
   ```

2. **Exit Flow - Decline** (Agent)
   ```
   - Send command: "/exit"
   - Wait for prompt: "confirm terminal shutdown? (yes/no):"
   - Respond: "no"
   - Verify: Still connected
   - Send another message to confirm
   ```

3. **Exit Flow - Accept** (Agent)
   ```
   - Send command: "/exit"
   - Wait for prompt
   - Respond: "yes"
   - Monitor for shutdown sequence:
     * "disconnecting secure channel..."
     * "clearing encryption keys..."
     * "closing connection..."
     * "SECURE TERMINAL CLOSED"
   - Verify process exits cleanly (exit code 0)
   ```

### Scenario 5: Encryption Validation

**Objective:** Verify end-to-end encryption and database security.

**Setup:**
```yaml
Users:
  - iris (CLI)
  - jack (Web)
Room: encryption-test
```

**Steps:**

1. **Send Clear Message** (Agent 1)
   ```
   iris: Send "THIS IS A SECRET MESSAGE"
   ```

2. **Database Inspection** (Database Agent)
   ```sql
   -- Query the encrypted message
   SELECT message FROM messages
   WHERE room_name = 'encryption-test'
   ORDER BY created_at DESC LIMIT 1;

   -- Verify the output is encrypted (base64, not plaintext)
   -- Should NOT contain: "THIS IS A SECRET MESSAGE"
   -- Should contain: encrypted blob (e.g., "aG4f8...")
   ```

3. **Decryption Verification** (Agent 2)
   ```
   jack (Web): Monitor for message arrival
   jack: Verify sees "THIS IS A SECRET MESSAGE" (decrypted)
   jack: Confirm NOT seeing encrypted blob
   ```

4. **Encryption Parameters** (Database Agent)
   ```
   - Verify all messages have:
     * Base64-encoded format
     * Consistent length patterns (AES-256-GCM)
     * No plaintext leakage
   ```

## Implementation Guide

### Phase 1: Setup Testing Environment

**Prerequisites:**
```bash
# Ensure kaiNET is built
./build-all.sh

# Start web dev server
cd web && bun run dev

# Verify Turso CLI access
turso db shell bt-phone-home
```

**MCP Server Configuration:**
```json
{
  "mcpServers": {
    "chrome-devtools": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-chrome-devtools"]
    },
    "playwright": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-playwright"]
    }
  }
}
```

### Phase 2: Create Test Orchestration Scripts

**Orchestrator Template:**
```markdown
Test: Two-User Message Exchange
================================

Objective: Verify message delivery between two CLI users

Steps:
1. Use Task tool to spawn Agent 1 (alice)
   - Launch CLI: ./bin/kainet-darwin-arm64 alice test-room
   - Wait for connection confirmation
   - Send message: "Hello from alice"
   - Monitor output for bob's response

2. Use Task tool to spawn Agent 2 (bob)
   - Launch CLI: ./bin/kainet-darwin-arm64 bob test-room
   - Wait for alice's message to appear
   - Send response: "Hello from bob"

3. Use Task tool to spawn Agent 3 (verifier)
   - Execute: turso db shell bt-phone-home
   - Run query: SELECT COUNT(*) FROM messages WHERE room_name='test-room'
   - Verify count = 2
   - Check encrypted field format

4. Collect results:
   - Agent 1 report: Did bob's message arrive?
   - Agent 2 report: Did alice's message arrive?
   - Agent 3 report: Database state valid?
   - Timing: All < 2 seconds?

Expected: ALL checks pass
```

### Phase 3: Agent Task Specifications

**CLI Agent Task Template:**
```
Agent Type: general-purpose
Task: Simulate kaiNET CLI user

Instructions:
1. Launch the binary: ./bin/kainet-darwin-arm64 {username} {room}
2. Monitor output until you see "SECURE CHANNEL ACTIVE"
3. Wait {delay} seconds
4. Send message: "{message_text}"
5. Monitor for incoming messages for {timeout} seconds
6. Record:
   - Connection time
   - Messages sent
   - Messages received
   - Any errors
7. Return detailed report

Expected output format:
{
  "connected": true/false,
  "connection_time_ms": 1234,
  "messages_sent": ["msg1", "msg2"],
  "messages_received": [
    {"from": "user", "text": "...", "timestamp": "..."}
  ],
  "errors": []
}
```

**Browser Agent Task Template:**
```
Agent Type: general-purpose
Task: Test web client functionality

Instructions:
1. Use mcp__chrome-devtools__navigate_page to http://localhost:5173
2. Take initial snapshot with mcp__chrome-devtools__take_snapshot
3. Fill login form:
   - username: {username}
   - room: {room}
4. Wait for connection using mcp__chrome-devtools__wait_for("SECURE CHANNEL")
5. Type message using mcp__chrome-devtools__type
6. Monitor console with mcp__chrome-devtools__list_console_messages
7. Capture final state with snapshot
8. Return report with:
   - Login success
   - Message sent confirmation
   - Console errors (if any)
   - Network requests summary
```

**Database Verification Agent Task Template:**
```
Agent Type: general-purpose
Task: Verify database state

Instructions:
1. Connect to Turso using: turso db shell bt-phone-home
2. Execute query: {sql_query}
3. Parse results
4. Verify:
   - Row count matches expected
   - Encryption format correct (base64)
   - Room isolation maintained
   - Timestamps reasonable
5. Return structured report

Expected output:
{
  "query": "SELECT ...",
  "row_count": 10,
  "encryption_valid": true,
  "room_isolation_valid": true,
  "timestamp_order_valid": true,
  "sample_rows": [...]
}
```

### Phase 4: Execution Workflow

**Step 1: Launch Orchestrator**
```
Main Claude instance initiates test suite
- Load test scenarios from specs
- Prepare agent configurations
- Set up result collection
```

**Step 2: Spawn Agents in Parallel**
```
Use Task tool with multiple parallel calls:

Task 1: CLI Agent (alice)
- Launch binary, send messages, monitor

Task 2: CLI Agent (bob)
- Launch binary, respond to messages

Task 3: Browser Agent (carol)
- Open web client, interact via UI

Task 4: Database Agent
- Monitor database state continuously
```

**Step 3: Coordinate Timing**
```
Orchestrator ensures proper sequencing:
- Wait for all agents to connect
- Signal message sending in order
- Coordinate verification steps
- Collect results when complete
```

**Step 4: Result Aggregation**
```
Orchestrator receives reports from all agents:
- Parse each agent's output
- Cross-reference timing data
- Verify consistency across agents
- Identify any failures
```

**Step 5: Generate Test Report**
```
Create comprehensive report:
- Test scenario name
- Pass/fail status
- Timing metrics
- Agent reports
- Database state snapshots
- Any errors or warnings
- Recommendations
```

## Detailed Test Cases

### Test Case 1: Basic Message Flow

```yaml
ID: TC-001
Name: Basic Two-User Message Exchange
Agents Required: 2 CLI
Duration: 30 seconds
```

**Orchestrator Instructions:**
```
1. Launch Agent 1 (alice):
   Prompt: "Launch kaiNET CLI as user 'alice' in room 'tc001-room'.
           Wait for connection. After 5 seconds, send message
           'Test message from alice'. Monitor for bob's response
           for 10 seconds. Report all received messages."

2. Launch Agent 2 (bob):
   Prompt: "Launch kaiNET CLI as user 'bob' in room 'tc001-room'.
           Wait for alice's message. When received, reply with
           'Received - bob'. Report timestamp of alice's message."

3. Wait for both agents to complete (max 30 seconds)

4. Verify:
   - Agent 1 received bob's message
   - Agent 2 received alice's message
   - Time between messages < 2 seconds
```

### Test Case 2: Room Isolation

```yaml
ID: TC-002
Name: Verify Room Message Isolation
Agents Required: 3 CLI, 1 Database
Duration: 1 minute
```

**Orchestrator Instructions:**
```
1. Launch Agent 1 (alice in room-A):
   - Send: "Message in room A"

2. Launch Agent 2 (bob in room-A):
   - Monitor: Should see alice's message

3. Launch Agent 3 (carol in room-B):
   - Monitor: Should NOT see alice's message
   - Send: "Message in room B"

4. Launch Database Agent:
   - Query room-A: Expect 1 message
   - Query room-B: Expect 1 message
   - Verify no cross-room leakage

5. Verify:
   - Bob sees alice's message
   - Carol does NOT see alice's message
   - Carol sees her own message only
   - Database confirms isolation
```

### Test Case 3: /burn Command Multi-User

```yaml
ID: TC-003
Name: Burn Command with Active Users
Agents Required: 3 CLI, 1 Database
Duration: 2 minutes
```

**Orchestrator Instructions:**
```
1. Setup Phase:
   - Launch alice, bob, carol in 'burn-test-room'
   - Each sends 3 messages (total 9 messages)
   - Verify all users see all 9 messages

2. Burn Phase:
   - Alice executes: /burn
   - Wait for "PURGED" confirmation

3. Verification Phase:
   - Bob: Check if messages disappeared
   - Carol: Check if messages disappeared
   - Database: Query message count (expect 0)

4. Recovery Test:
   - Alice sends: "After burn message"
   - Verify bob and carol receive it
   - Verify no old messages reappear
```

### Test Case 4: Concurrent Message Storm

```yaml
ID: TC-004
Name: High-Volume Concurrent Messages
Agents Required: 5 CLI, 1 Database
Duration: 3 minutes
```

**Orchestrator Instructions:**
```
1. Launch 5 users in same room:
   - user1, user2, user3, user4, user5

2. Message Storm (parallel):
   - Each user sends 10 messages rapidly
   - Messages numbered: "user1-msg-1", "user1-msg-2", etc.
   - Total: 50 messages

3. Verification:
   - Each user monitors for all 50 messages
   - Database confirms 50 messages stored
   - Check for any lost messages
   - Verify correct ordering

4. Metrics to collect:
   - Average message delivery time
   - Maximum delivery time
   - Any duplicates or losses
   - Database query performance
```

### Test Case 5: Web Client Network Failure

```yaml
ID: TC-005
Name: Web Client Network Resilience
Agents Required: 1 Browser, 1 CLI, 1 Database
Duration: 2 minutes
```

**Browser Agent Instructions:**
```
1. Navigate to web client
2. Login as 'web-user' in 'resilience-room'
3. Send initial message: "Before network issue"

4. Simulate network disruption:
   - Use DevTools Network throttling
   - Or disable network briefly

5. During disruption:
   - CLI user sends: "Message during outage"
   - Record: Does web user see it?

6. Restore network:
   - Verify message appears after reconnection
   - Check console for error handling
   - Verify no message loss

7. Database verification:
   - Confirm both messages in database
   - Check timestamps make sense
```

### Test Case 6: Encryption Integrity

```yaml
ID: TC-006
Name: End-to-End Encryption Validation
Agents Required: 2 CLI, 1 Database
Duration: 1 minute
```

**Orchestrator Instructions:**
```
1. Launch alice and bob in 'crypto-test-room'

2. Alice sends: "PLAINTEXT_SECRET_MESSAGE_12345"

3. Database Agent:
   - Query the message from database
   - Extract encrypted payload
   - Verify:
     * NOT plaintext visible
     * Base64 format
     * Reasonable length (>50 chars)
   - Record: actual encrypted value

4. Bob's perspective:
   - Verify sees: "PLAINTEXT_SECRET_MESSAGE_12345"
   - Confirm decryption successful

5. Security checks:
   - Ensure no plaintext in database
   - Verify encryption key not logged
   - Check room_name not encrypted (expected)
```

## Visibility and Monitoring

### Real-Time Monitoring Dashboard

**Information to track:**
```
┌─────────────────────────────────────────────────────────┐
│ Test Execution Dashboard                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│ Active Agents:                                          │
│  [•] Agent 1 (alice) - CLI - Connected - 3 msgs sent   │
│  [•] Agent 2 (bob)   - CLI - Connected - 2 msgs sent   │
│  [•] Agent 3 (carol) - Web - Connected - 1 msg sent    │
│  [•] Database Agent  - Monitoring - 6 rows detected    │
│                                                         │
│ Current Test: TC-002 (Room Isolation)                  │
│ Status: In Progress - Step 3/5                         │
│ Elapsed: 00:00:45                                       │
│                                                         │
│ Message Flow (last 10):                                │
│  14:23:01 alice → room-A: "Test message"               │
│  14:23:02 bob   ← room-A: received ✓                   │
│  14:23:03 carol → room-B: "Other message"              │
│  14:23:03 carol ← room-B: received ✓                   │
│                                                         │
│ Database State:                                         │
│  room-A: 1 message                                      │
│  room-B: 1 message                                      │
│  Total: 2 messages                                      │
│  Encryption: Valid ✓                                    │
│                                                         │
│ Checks:                                                 │
│  [✓] Room isolation working                            │
│  [✓] Message delivery < 2s                             │
│  [✓] Database consistency                              │
│  [ ] Final verification pending                        │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Agent Logging

Each agent should log:
```json
{
  "agent_id": "agent-1-alice",
  "timestamp": "2025-11-14T15:30:00Z",
  "event": "message_sent",
  "details": {
    "message": "Hello bob",
    "room": "test-room",
    "latency_ms": 145
  }
}
```

### Database State Snapshots

Periodic snapshots:
```sql
-- Snapshot query template
SELECT
  id,
  room_name,
  username,
  SUBSTR(message, 1, 20) as message_preview,
  created_at,
  LENGTH(message) as encrypted_length
FROM messages
WHERE created_at > datetime('now', '-5 minutes')
ORDER BY created_at DESC;
```

## Success Criteria

### Individual Test Success

A test passes when:
- All agents complete without errors
- Message delivery within acceptable time (< 2s)
- Database state matches expectations
- No encryption leaks detected
- No race conditions observed
- Clean agent shutdown

### Test Suite Success

Full suite passes when:
- 100% of test cases pass
- No critical issues identified
- Performance metrics within targets
- Security validation complete
- Cross-platform compatibility confirmed

## Troubleshooting Guide

### Common Issues

**Issue: Agent hangs during CLI launch**
```
Symptoms: Agent doesn't report connection
Diagnosis:
  - Check binary path is correct
  - Verify database credentials valid
  - Check network connectivity
Resolution:
  - Use timeout in agent task (30s)
  - Log stderr from CLI process
  - Verify Turso database is accessible
```

**Issue: Messages not appearing in database**
```
Symptoms: Sent messages don't show in queries
Diagnosis:
  - Check room name spelling
  - Verify encryption keys
  - Check database connection
Resolution:
  - Add database logging in app
  - Verify turso CLI can connect
  - Check auth token validity
```

**Issue: Browser agent can't find elements**
```
Symptoms: DevTools commands fail with "element not found"
Diagnosis:
  - Check page load complete
  - Verify element UID from snapshot
  - Check for dynamic rendering delays
Resolution:
  - Add wait_for commands
  - Use longer timeouts
  - Take snapshot before interaction
```

**Issue: Race condition between agents**
```
Symptoms: Inconsistent test results
Diagnosis:
  - Check timing coordination
  - Verify message ordering
  - Look for parallel conflicts
Resolution:
  - Add explicit synchronization points
  - Use message acknowledgment
  - Coordinate through orchestrator
```

## Performance Benchmarks

### Target Metrics

**Message Delivery:**
- Single message: < 1 second
- Under load (10 users): < 2 seconds
- Peak throughput: > 100 messages/minute

**Room Operations:**
- Room switch: < 2 seconds
- Burn command: < 500ms for any size

**Database:**
- Message insert: < 100ms
- Query 100 messages: < 200ms
- Room isolation query: < 50ms

**Web Client:**
- Initial load: < 2 seconds
- Message render: < 100ms
- Network reconnect: < 3 seconds

## Future Enhancements

### Automated Test Suite

**Goal:** Run full test suite on every release

**Implementation:**
- CI/CD integration
- Automated agent orchestration
- Result collection and reporting
- Regression detection

### Visual Regression Testing

**Goal:** Detect UI changes in web client

**Implementation:**
- Screenshot comparison with Playwright
- Baseline image storage
- Diff visualization
- Automatic failure detection

### Load Testing

**Goal:** Determine system capacity

**Implementation:**
- Spawn 50+ concurrent agents
- Measure message throughput
- Monitor database performance
- Identify bottlenecks

### Chaos Engineering

**Goal:** Test resilience to failures

**Implementation:**
- Random network disconnects
- Database connection drops
- Simulated latency
- Message corruption attempts

## Conclusion

This agentic testing workflow provides comprehensive validation of kaiNET functionality through:

1. **Multi-agent coordination** - Simulating real user interactions
2. **Full visibility** - Monitoring every layer (UI, logic, database)
3. **Automated verification** - Reducing manual testing burden
4. **Security validation** - Confirming encryption integrity
5. **Performance metrics** - Ensuring acceptable user experience

The combination of Claude subagents, browser automation MCPs, and Turso CLI creates a powerful testing framework that can validate complex multi-user scenarios with high confidence.

**Estimated setup time:** 4-6 hours
**Test execution time:** 5-10 minutes per full suite
**Maintenance effort:** LOW (tests are self-documenting)
**Coverage:** 90%+ of user workflows
