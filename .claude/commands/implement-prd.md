# implement-prd

Implement a Product Requirements Document (PRD) using Test-Driven Development approach

## Usage

```
implement-prd <prd-filename> [thinking-mode]
```

## Arguments

- `prd-filename`: Name of the PRD file in @specs/ (without .prd.md extension)
- `thinking-mode` (optional): 
  - `think` - Standard implementation with basic reasoning
  - `think-harder` - Enhanced analysis with deeper consideration
  - `ultrathink` - Use mcp__sequential-thinking__sequentialthinking for complex reasoning

## Process

### 1. Setup Phase
- Read the PRD from `@specs/<prd-filename>.prd.md`
- Create implementation report at `@specs/<prd-filename>-report.md`
- Initialize git status check
- Create TodoWrite list with all tasks

### 2. Task Breakdown
Extract tasks from PRD and break down into:
- Primary tasks (from main requirements)
- Subtasks (atomic, testable units)
- Dependencies between tasks

### 3. TDD Implementation Loop

For each task:

#### a. Pre-Implementation
```bash
git status  # Verify clean working tree
```

#### b. Test Creation (if applicable)
- Write failing test for the feature
- Run test to confirm failure
- Commit test with message: `test: add test for <feature>`

#### c. Implementation
- Implement minimal code to pass test
- Verify implementation works
- Run any existing tests

#### d. Commit Atomic Change
```bash
git add <specific-files>
git commit -m "<type>: <description>"
```

Commit message types:
- `feat:` New feature
- `fix:` Bug fix
- `style:` Formatting, missing semicolons, etc.
- `refactor:` Code restructuring
- `test:` Adding tests
- `docs:` Documentation only
- `chore:` Maintenance

#### e. Verify & Document
- Run `git status` to confirm clean staging
- Update report with task completion
- Mark task as completed in TodoWrite

### 4. Task Order

1. **Setup & Configuration**
   - Dependencies installation
   - Configuration updates
   - Type definitions

2. **Core Functionality**
   - Data structures
   - Business logic
   - API integrations

3. **UI Components**
   - Layout changes
   - Visual components
   - Animations

4. **Testing & Validation**
   - Unit tests
   - Integration tests
   - Manual verification

5. **Documentation**
   - Code comments
   - README updates
   - update ARCHITECTURE.md (invoke .claude/commands/architect.md ultrathink)
   - Report finalization

### 5. Report Structure

The implementation report should include:

```markdown
# Implementation Report: <Feature Name>
## Date: <Current Date>
## PRD: <prd-filename>.prd.md

## Tasks Completed
- [x] Task 1: Description
  - Commit: <hash> <message>
  - Files: <files-changed>
- [x] Task 2: Description
  - Commit: <hash> <message>
  - Files: <files-changed>

## Testing Summary
- Tests written: <count>
- Tests passing: <count>
- Coverage: <percentage>

## Challenges & Solutions
- Challenge 1: Description
  - Solution: How it was resolved

## Performance Metrics
- Before: <metrics>
- After: <metrics>

## Next Steps
- Future enhancements
- Technical debt identified
```

## Thinking Modes

### Standard (`think`)
- Analyze requirements
- Plan implementation
- Execute tasks sequentially

### Enhanced (`think-harder`)
- Deep analysis of edge cases
- Consider multiple implementation approaches
- Optimize for performance and maintainability

### Ultra (`ultrathink`)
When available, use sequential thinking MCP:
```javascript
mcp__sequential-thinking__sequentialthinking({
  thought: "Analyzing PRD requirements...",
  nextThoughtNeeded: true,
  thoughtNumber: 1,
  totalThoughts: 10
})
```

## Example Workflow

```bash
# Start implementation
implement-prd feed-design ultrathink

# Agent actions:
1. Read @specs/feed-design.prd.md
2. Create @specs/feed-design-report.md
3. Break down into tasks:
   - Task 1: Update countdown timer
     - Subtask 1.1: Add JavaScript calculation
     - Subtask 1.2: Update HTML structure
     - Subtask 1.3: Add CSS animations
   - Task 2: Fix member counter layout
     - Subtask 2.1: Update grid system
     - Subtask 2.2: Adjust responsive breakpoints
4. For each subtask:
   - Implement change
   - Test functionality
   - Commit with message
   - Verify staging clean
5. Update report with results
```

## Git Commit Guidelines

### Message Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Examples
```bash
feat(feed): add dynamic countdown timer to September 15

Replaces static "NOW" text with calculated days remaining.
Updates daily at midnight local time.

feat(feed): update member counter to match countdown width

Removes full-width behavior and aligns with grid system.

style(feed): equalize statistics block widths

Updates grid layout for consistent visual hierarchy.
```

## Verification Checklist

Before marking task complete:
- [ ] Code follows PRD specifications
- [ ] Tests pass (if applicable)
- [ ] No linting errors
- [ ] Git staging area is clean
- [ ] Commit message follows conventions
- [ ] Report updated with task details

## Error Handling

If implementation fails:
1. Document error in report
2. Attempt alternative approach
3. If blocked, note in report and continue with next task
4. Create TODO for resolution

## Notes

- Always verify git status before and after commits
- Keep commits atomic and focused
- Document any deviations from PRD in report
- Use thinking mode appropriate to complexity
- Test each change before committing
- Update TodoWrite list throughout process