# create-prd

Create a Product Requirements Document (PRD) for a feature or component

## Usage

```
create-prd <feature-name> <description>
```

## Arguments

- `feature-name`: The name of the feature/component (will be used in filename)
- `description`: Brief description of what needs to be documented

## Process

1. Analyze the current implementation of the feature/component if it exists
2. Identify problems, requirements, and design specifications
3. Create a comprehensive PRD document
4. Save to `@specs/` directory with filename format: `<last number + 1><feature-name>.prd.md` e.g. `001-first-feature.prd.md`

You must scan the specs directory and find the next available number (e.g. 001,002 and 003 already exist, and so this new prd is "004-<feature-name>.prd.md).

## PRD Structure

The PRD should include:

1. **Executive Summary** - Brief overview of the feature
2. **Problem Statement** - Current issues and pain points
3. **Requirements** - Functional and technical specifications
   - User requirements
   - Technical requirements
   - Design requirements
4. **Implementation Notes** - Code examples and technical approach
5. **Responsive Design** - Mobile and desktop layouts
7. **Animation Specifications** - If applicable
8. **Success Metrics** - How to measure success (if applicable)
9. **Future Enhancements** - Potential improvements

## Examples

```
create-prd navigation-menu "Update navigation menu with mobile hamburger"
```
Creates: `@specs/navigation-menu.prd.md`

```
create-prd user-auth "Design user authentication flow"
```
Creates: `@specs/user-auth.prd.md`

## File Naming Convention

- Use kebab-case for filenames
- Always end with `.prd.md`
- Place in `@specs/` directory
- Examples:
  - `feed-design.prd.md`
  - `navigation-update.prd.md`
  - `auth-flow.prd.md`

## Notes

- Include visual mockups using ASCII diagrams where helpful
- Add code snippets for implementation guidance
- Consider mobile-first design approach
- Document all edge cases and error states
- Include version number and date in the document