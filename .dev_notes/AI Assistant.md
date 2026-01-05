# AI-Assisted Daily Planning System – Specification

## 1. Purpose

The purpose of this project is to build a personal, AI-assisted daily planning system that helps plan workdays by combining:

- Weekly goals
- Calendar events
- Jira tickets

The system assists through **conversational interaction**, modifies a structured **daily note in Obsidian (Markdown)**, and prioritizes simplicity, local-first development, and learning how to build AI-powered systems that can later map cleanly to cloud/serverless architectures.

---

## 2. High-Level Goals

- Interact with an AI each morning using natural language (e.g., “What’s my day look like?”).
- Automatically create or update a daily Obsidian note.
- Allow the AI to modify specific sections of the daily note.
- Use local LLMs (via Ollama) for development.
- Learn practical AI system design, not just prompt usage.

---

## 3. Daily Note Structure

The daily note follows a predefined template with **three sections** that this system is allowed to read and modify:

1. **Goals**
2. **Meetings**
3. **Bonus Items**

All other sections of the note (if present) are out of scope and must not be modified.

The template itself is fixed and will be provided separately.

---

## 4. Core Interaction Model

### Morning Conversation

The primary interaction is conversational and initiated by the user, for example:

- “What’s my day look like?”
- “What do I need to get done today?”

When invoked, the system will:

1. Read weekly goals from the weekly note.
2. Read today’s calendar events.
3. Read Jira tickets assigned to the user.
4. Read or create today’s daily note.
5. Propose and apply updates to the daily note.
6. Summarize the plan conversationally to the user.

Follow-up conversation can refine priorities or modify the plan.

---

## 5. Input Sources

### 5.1 Weekly Note

- The weekly note already exists.
- It is created by a separate journal plugin.
- The system may read it for planning context but does not modify it directly.

---

### 5.2 Calendar

- Read-only access via a private iCal (ICS) URL.
- The system may read:
  - Event title
  - Description (if available)
- Event times are ignored for note rendering.
- Calendar mutation (blocking focus time) is a **long-term goal**, not part of v1 or v2.

---

### 5.3 Jira

- Read Jira tickets assigned to the user.
- Required fields:
  - Summary
  - Description
  - Status
  - Priority
  - Due date (if present)
- Authentication via API token is acceptable initially.

#### PII Handling

- Jira data must be sanitized before being sent to the LLM.
- Sensitive values (emails, account numbers, usernames, etc.) are redacted.
- Redaction is **not reversible**.
- Original data remains accessible in Jira when needed.

---

### 5.4 Obsidian Vault

The system must be able to:

- Locate an Obsidian vault on disk.
- Read Markdown files.
- Create a daily note if it does not exist.
- Modify existing Markdown files safely and idempotently.
- Avoid destructive rewrites.

---

## 6. Daily Note Section Behavior

### 6.1 Goals

- Contains planned tasks for the day.
- Tasks may originate from:
  - Jira tickets
  - Weekly goals
- The AI may:
  - Add tasks
  - Reorder tasks
  - Remove or defer tasks (with user confirmation)
- Incomplete tasks from the previous business day are detected and surfaced.

#### Carryover Flow

When incomplete tasks exist from the previous day, the system asks whether to:

- Carry them forward
- Drop them
- Explicitly defer them

---

### 6.2 Meetings

- Meetings are derived from calendar events.
- Meetings **do not include times**.
- Each meeting is listed once.

#### Meeting Links

- Only meetings that warrant long-lived context should receive links.
- Routine meetings (e.g., standups, releases) typically do not get links.
- When a link is created, it must follow this format:

  [[YYYY/MM/Meeting Name|Meeting Name]]

- Only the link is created; Obsidian will create the file on demand.
- This avoids creating dead files for canceled meetings.

Meeting classification may be:
- Rule-based
- LLM-assisted
- Or a hybrid of both

---

### 6.3 Bonus Items

- Represents unplanned work that arises during the day.
- Initially empty when the daily note is created.
- Items added here are used as a signal for capacity forecasting.

---

## 7. Unplanned Work Observation

- The system observes items added under **Bonus Items**.
- It generates a qualitative assessment of unplanned work, such as:
  - “a little”
  - “some”
  - “mostly unplanned”
  - “everything was unplanned”
- This assessment is written to the weekly journal note.
- No numeric metrics are required.
- The intent is narrative trend awareness, not precise forecasting.

---

## 8. Long-Term Goal: Calendar Mutation

A future capability (not v1 or v2):

- Automatically block focus time on the calendar.
- Focus blocks are based on:
  - Planned goals
  - Meeting density
  - Historical unplanned work patterns

The system architecture should be designed with this future capability in mind.

---

## 9. Technical Constraints

### 9.1 AI / LLM

- Local-first development using Ollama.
- Models must support structured output and tool usage.

---

### 9.2 Language Choices

Preferred languages, in order of practicality:

1. Python
2. Go
3. JavaScript / TypeScript
4. C#

Language choice should favor:
- Simple, readable libraries
- Ease of AI integration
- Maintainability over cleverness

---

### 9.3 Security

- Jira API tokens must be stored securely.
- Environment variables or OS keychains are acceptable for local development.

---

### 9.4 Cloud Migration (Learning Objective)

Although this project may remain local, the design should map cleanly to AWS serverless concepts:

- Stateless execution
- Clear separation of concerns
- File-based state (Markdown) compatible with S3 or Git-backed storage

---

## 10. Version Scope

### v1

- Calendar ingestion (read-only)
- Jira ingestion
- PII sanitization
- Daily note creation and modification
- Conversational morning planning
- Local-only execution

### v2

- Smarter meeting classification
- Better carryover logic
- Qualitative capacity trend usage

### v3 (Future)

- Calendar mutation (focus time blocking)
- Multi-day planning awareness
- Deeper historical analysis

---

## 11. Non-Goals

- Full task management replacement
- Autonomous decision-making without user confirmation
- Complex metrics or dashboards

The system is intended to remain lightweight, transparent, and Markdown-first.
