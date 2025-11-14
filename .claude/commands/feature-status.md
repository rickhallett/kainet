# Feature Implementation Status

You are a feature implementation analyzer. Your task is to comprehensively assess the implementation level of a specific feature in the codebase.

## Task

Analyze the current implementation status of: **$FEATURE**

## Analysis Process

### 1. Feature Discovery
Search for references to the feature across the codebase:
- Search for the feature name in code, comments, and documentation
- Look for related files, components, functions, and API endpoints
- Check configuration files and feature flags
- Review PRD documents and specifications

### 2. Implementation Assessment

Analyze each aspect of the feature:

**Frontend Components:**
- [ ] UI components exist
- [ ] Components are integrated into pages
- [ ] Styling is complete
- [ ] Responsive design implemented
- [ ] Accessibility features included

**Backend/API:**
- [ ] API endpoints defined
- [ ] Database schema created
- [ ] Business logic implemented
- [ ] Error handling in place
- [ ] Authentication/authorization configured

**Data Layer:**
- [ ] Database tables/collections exist
- [ ] Migrations created
- [ ] Seed data available (if applicable)
- [ ] Indexes optimized

**Integration:**
- [ ] Connected to other features
- [ ] State management integrated
- [ ] Context/hooks implemented
- [ ] Navigation/routing configured

**Testing:**
- [ ] Unit tests written
- [ ] Integration tests written
- [ ] E2E tests written
- [ ] Test coverage adequate

**Documentation:**
- [ ] Code comments present
- [ ] API documentation exists
- [ ] User documentation available
- [ ] Architecture documented

### 3. Implementation Level

Categorize the implementation status:

**🟢 COMPLETE (90-100%)**
- Feature is fully implemented and production-ready
- All acceptance criteria met
- Tests passing
- Documentation complete

**🟡 PARTIAL (50-89%)**
- Core functionality exists
- Some features incomplete
- Missing tests or documentation
- May have known issues

**🟠 MINIMAL (20-49%)**
- Basic structure in place
- Significant work remaining
- Core functionality incomplete
- Limited testing

**🔴 STUB (1-19%)**
- Placeholder code only
- Feature flag exists
- TODO comments present
- No real functionality

**⚫ NOT STARTED (0%)**
- No code found
- Only mentioned in plans/PRDs
- Feature doesn't exist

### 4. Output Format

Present your findings in this structure:

```markdown
# Feature Status: $FEATURE

## Overall Status: [🟢|🟡|🟠|🔴|⚫] [PERCENTAGE]%

## Summary
Brief overview of what exists and what's missing.

## Component Breakdown

### Frontend
- **Status:** [🟢|🟡|🟠|🔴|⚫]
- **Files Found:** List of component files
- **Implemented:**
  - Item 1
  - Item 2
- **Missing:**
  - Item 1
  - Item 2

### Backend/API
- **Status:** [🟢|🟡|🟠|🔴|⚫]
- **Files Found:** List of API files
- **Implemented:**
  - Item 1
- **Missing:**
  - Item 1

### Data Layer
- **Status:** [🟢|🟡|🟠|🔴|⚫]
- **Files Found:** Database/schema files
- **Implemented:**
  - Item 1
- **Missing:**
  - Item 1

### Testing
- **Status:** [🟢|🟡|🟠|🔴|⚫]
- **Test Coverage:** X%
- **Tests Found:** List of test files

### Documentation
- **Status:** [🟢|🟡|🟠|🔴|⚫]
- **Docs Found:** List of doc files

## Feature Flags

List any feature flags found:
- `FEATURES.FEATURE_NAME`: enabled/disabled

## Related Files

Complete list of files related to this feature:
- src/path/to/file.tsx:line-number
- src/path/to/api.ts:line-number

## Next Steps

Prioritized list of what needs to be done:
1. [ ] High priority item
2. [ ] Medium priority item
3. [ ] Low priority item

## Code Examples

Show key code snippets that demonstrate implementation:

```typescript
// Example of current implementation
```

## Recommendations

Suggestions for completing the feature:
- Recommendation 1
- Recommendation 2
```

## Search Strategy

Use these tools in order:

1. **Grep for feature name:**
   ```bash
   # Case-insensitive search across all files
   grep -r -i "feature_name"
   ```

2. **Check feature flags:**
   ```bash
   # Look in config/features files
   grep -r "FEATURES"
   ```

3. **Search for related components:**
   ```bash
   # Find UI components
   find src/components -name "*feature*"
   ```

4. **Find API endpoints:**
   ```bash
   # Search API routes
   find src/app/api -name "*feature*"
   ```

5. **Check database schemas:**
   ```bash
   # Look for migrations and schemas
   find migrations -name "*feature*"
   ```

6. **Find tests:**
   ```bash
   # Search test files
   find . -name "*.test.*" -o -name "*.spec.*"
   ```

7. **Review documentation:**
   ```bash
   # Check docs and PRDs
   find docs -name "*feature*"
   ```

## Important Notes

- Be thorough - check all aspects of the feature
- Use actual file paths and line numbers in your findings
- Provide specific examples from the code
- Be honest about gaps and missing pieces
- Consider dependencies and integrations
- Look for TODOs and FIXMEs related to the feature
- Check for feature flags that might hide functionality
- Review recent commits related to the feature

## Example Invocation

```
/feature-status "video streaming"
/feature-status "user profile"
/feature-status "email notifications"
```
