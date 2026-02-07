### Obsidian Integration Plan (Refined)

This document outlines the focused strategy for reading specific sections from Obsidian daily and weekly notes to provide data for an AI agent.

#### Assumptions & File Discovery

- **Existence**: Files are assumed to exist for the current period.
- **Naming Convention**:
    - **Daily Note**: `{{YEAR}}/YYYY-MM-DD.md` (e.g., `2026/2026-02-07.md`).
    - **Weekly Note**: `{{YEAR}}/Week WW.md` (e.g., `2026/Week 06.md`).
- **Path Resolution**: The vault path is retrieved from the configuration. The year subfolder and filename are derived from the current system date.

#### Data Extraction (Reading)

The goal is to extract specific sections to be used as context for an AI agent.

**1. Daily Note: Tasks Section**
- **Target**: Everything under the `## Tasks` header.
- **Sub-sections to include**:
    - `### Goals`
    - `### Meetings`
    - `### Bonus Items`
- **Logic**: 
    - Locate the `## Tasks` header.
    - Capture all lines until the next Level 2 header (`## `) or end of file.
    - Preserve Markdown formatting (task checkboxes `- [ ]`).

**2. Weekly Note: Weekly Goals Section**
- **Target**: Everything under the `## Weekly Goals` header.
- **Logic**:
    - Locate the `## Weekly Goals` header.
    - Capture all lines until the next Level 2 header (`## `) or end of file.
    - Preserve Markdown formatting.

#### Processing Strategy

- **Line-by-Line Parsing**: A simple state-based parser is sufficient.
    - `State: SeekHeader` -> Look for target `## Header`.
    - `State: Capture` -> Collect lines until another `## ` header is encountered.
- **Output Format**: The extracted text will be passed as a string/context to the AI agent.

#### Future Considerations (Writing)

While the immediate focus is reading, writing will follow a similar pattern:
- **Daily Tasks**: Appending or updating items within the `### Goals` or `### Meetings` sub-sections.
- **Idempotency**: Ensuring that the AI agent's output doesn't create duplicate entries when written back.
