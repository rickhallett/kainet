---
description: Run build, check output, and fix errors and lint warnings intelligently
allowed-tools: [Bash, Read, Edit, Write, Glob, Grep]
---

Run the build process and systematically fix any errors or warnings:

1. Execute `bun run build` and capture all output
2. Analyze build errors and TypeScript errors - fix them in order of severity
3. Run lint check and analyze all warnings
4. For each lint warning:
   - Attempt specific code fixes (preferred approach)
   - Track retry count per lint rule
   - After 5 failed attempts on the same rule, add a specific eslint ignore rule
5. After all fixes, run build again to verify success
6. Report summary of all fixes made

**Fix Strategy:**
- Prioritize actual code fixes over eslint rules
- Fix unused imports/variables by removing them
- Fix type errors with proper TypeScript types
- Fix formatting issues with proper code style
- Only add eslint rules as last resort after 5 attempts

**Rules Addition Strategy:**
- Prefer inline `// eslint-disable-next-line [rule]` for specific cases
- Add global rules to `eslint.config.mjs` only for persistent project-wide issues
- Document why the rule was added

Work methodically through all issues until build succeeds cleanly.
