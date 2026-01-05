# High-Level Architecture (CLI + Local AI)

## Architectural Philosophy

- **Local-first, text-first**: Markdown is the system of record.
- **LLM as planner, not owner**: The LLM reasons and proposes changes; your code applies them.
- **Pull, sanitize, plan, write**: Every run follows a predictable pipeline.
- **Composable boundaries**: Each external system is isolated behind a thin adapter.

---

## Core Components (Mental Model)

```
┌───────────┐
│   CLI     │  ← conversational interface
└────┬──────┘
     │
     ▼
┌────────────────────┐
│ Orchestrator       │  ← coordinates one "planning session"
└────┬─────┬─────┬───┘
     │     │     │
     ▼     ▼     ▼
Calendar  Jira  Obsidian
Adapter   MCP   Adapter
     │     │     │
     └─────┴─────┘
           │
           ▼
     Sanitized Context
           │
           ▼
       Ollama LLM
           │
           ▼
     Plan / Patch Output
           │
           ▼
   Markdown Writer

```

---

## 1. CLI (Chat Interface)

**Responsibility**

- User-facing interaction
- Sends user prompts into the system
- Displays summaries and follow-up questions

**Key Characteristics**

- Feels conversational, but is **stateless per invocation**
- Example commands:
    - `planner chat`
    - `planner today`
    - `planner repl`

**Why CLI first**

- Lowest friction
- Easy to debug
- Maps cleanly to:
    - Obsidian command later
    - HTTP API later
    - AWS Lambda later

The CLI does _not_ talk directly to Jira, calendar, or files — it only talks to the orchestrator.

---

## 2. Orchestrator (The Session Brain)

**Responsibility**

- Runs a single “planning session”
- Owns the workflow, not the logic details

**Typical flow**

1. Receive user intent from CLI
2. Gather context:
    - Weekly note
    - Daily note (or create it)
    - Calendar events
    - Jira issues
3. Sanitize external data
4. Call the LLM
5. Apply LLM-approved mutations to Markdown
6. Return a human-readable summary

**Important**

- The orchestrator does **no parsing or heuristics**
- It delegates everything to adapters

This makes it testable and cloud-friendly.

---

## 3. Calendar Adapter (iCal via curl)

**Responsibility**

- Fetch calendar data from a private iCal URL
- Parse ICS into internal events

**Key Design Choice**

- Treat calendar ingestion as a _pure function_:
    - Input: iCal URL
    - Output: list of calendar events

**Why this is nice**

- `curl` compatibility = trivial to debug
- Easy to replace with Google/Microsoft APIs later
- No authentication complexity in v1

**Output Shape (conceptual)**

- Event title
- Description
- Date (day-level granularity only)

No mutation, no state.

---

## 4. Jira Adapter via Custom MCP

This is the most interesting piece.

### Why MCP Makes Sense Here

- Jira access is _tool-like_, not conversational
    
- You want:
    - Deterministic data retrieval
    - Strong sanitization guarantees
    - Zero prompt leakage of raw data

### Responsibilities

- Call Jira (via CLI or REST)
- Retrieve assigned tickets
- Sanitize ticket data
- Return **LLM-safe** representations

### Boundary Rule (Very Important)

> **The LLM never sees raw Jira data.**

Instead:

```
Raw Jira → MCP Tool → Sanitized Ticket Model → LLM
```
This keeps:

- PII concerns isolated
- Prompts simpler
- Behavior deterministic

### Why a Custom MCP Is a Good Call

- You can:
    - Control exactly what fields exist
    - Evolve sanitization independently
    - Swap Jira CLI vs REST without touching prompts

Even if you don’t fully adopt MCP immediately, designing _as if_ this is an MCP tool is the right move.

---

## 5. Obsidian Adapter (Markdown I/O)

**Responsibility**

- Locate vault
- Read notes
- Apply safe, scoped edits

### Critical Design Constraint

The adapter must support **surgical edits**:

- Read full Markdown
- Identify the three sections:
    - Goals
    - Meetings
    - Bonus Items
- Replace _only those sections_
- Preserve all other content byte-for-byte

This adapter is your **single source of truth** for state.

No database. No cache. No memory store.

---

## 6. Sanitization Layer (Explicit Step)

Even though sanitization lives logically near Jira, conceptually it is its own step:

```
External Data
     ↓
Sanitization
     ↓
LLM Context

```

**Characteristics**

- Regex-based
- Deterministic
- One-way
- Applied before _any_ LLM interaction
    

This is a strong architectural win:

- You can unit-test it easily
- You can audit it
- You can expand it later

---

## 7. Ollama (LLM Runtime)

**Role**

- Planner and classifier
- Markdown editor _by instruction_, not by file access

**Important Constraints**

- LLM:
    - Does not read files
    - Does not call APIs
    - Does not mutate state directly

Instead:

- It receives:
    - Sanitized context
    - Current Markdown sections
    - Explicit instructions
- It returns:
    - Structured plan
    - Proposed section replacements
    - Questions when uncertain

This keeps the system:

- Safe
- Predictable
- Debbugable
    

---

## 8. Output Contract (LLM → System)

Rather than “write the whole note,” the LLM should output something like:

- Updated Goals section
- Updated Meetings section
- Notes about Bonus Items (if any)
- A conversational summary

Your code decides whether and how to apply it.

This is key for future:

- Undo
- Diff
- Partial acceptance
- UI previews

---

## 9. Background Services

### Ollama

- Runs continuously as a local service
- Treated as an external dependency
- CLI checks for availability before execution

### Everything Else

- Ephemeral
- No long-running planner process required

This maps cleanly to serverless later.

---

## 10. How This Maps to AWS (Later)

|Local Concept|AWS Equivalent|
|---|---|
|CLI invocation|Lambda invocation|
|Orchestrator|Lambda handler|
|Markdown files|S3 objects / Git repo|
|Jira MCP|Lambda tool / Step Function|
|Ollama|Bedrock / ECS / SageMaker|

You are accidentally designing something _very deployable_.