# Requirements Document

## Introduction

The AI-Assisted Daily Planning System is a personal productivity tool that combines weekly goals, calendar events, and Jira tickets through conversational AI interaction to create and maintain structured daily notes in Obsidian. This is a single-user learning project focused on understanding AI integration patterns.

**Personal Project Scope:**
- Single user (no multi-tenancy)
- Local-first operation
- Learning-focused implementation
- Simple, maintainable code

**Full Product Considerations:**
- Would need user management and authentication
- Multi-tenant architecture
- Robust error handling and monitoring
- Scalable infrastructure

## Glossary

- **System**: The AI-Assisted Daily Planning System
- **Daily_Note**: A structured Markdown file in Obsidian containing Goals, Meetings, and Bonus Items sections
- **Weekly_Note**: An existing Markdown file containing weekly goals, created by a separate journal plugin
- **Orchestrator**: The core component that coordinates planning sessions and manages workflow
- **Calendar_Adapter**: Component responsible for fetching and parsing iCal calendar data
- **Jira_Adapter**: Component responsible for retrieving and sanitizing Jira ticket data
- **Obsidian_Adapter**: Component responsible for reading and writing Markdown files in the Obsidian vault
- **Sanitization_Layer**: Component that removes PII and sensitive data before LLM processing
- **Planning_Session**: A single execution cycle that gathers context, processes with LLM, and updates the daily note
- **Carryover_Tasks**: Incomplete tasks from the previous business day that need user decision
- **Focus_Time**: Calendar blocks for concentrated work (future capability)

## Requirements

### Requirement 1: Conversational Planning Interface

**User Story:** As a user, I want to interact with the system using natural language in an iterative conversation, so that I can refine my daily plan before it's written to my notes.

#### Acceptance Criteria

1. WHEN a user starts a chat with contextual information like "I was just asked to do xyz", THE System SHALL incorporate that context into the Goals section planning
2. WHEN the system processes a planning request, THE System SHALL present proposed changes in text format before writing any files
3. WHEN a user provides follow-up conversation, THE System SHALL allow refinement of priorities and plan modifications
4. WHEN a user is satisfied with the proposed plan, THE System SHALL require explicit acceptance before writing to files
5. WHEN the user accepts the plan, THE System SHALL apply changes to the Obsidian daily note and provide a summary

### Requirement 2: Daily Note Management

**User Story:** As a user, I want the system to create and update my daily Obsidian notes automatically, so that my planning is captured in my existing workflow.

#### Acceptance Criteria

1. WHEN a daily note does not exist for today, THE System SHALL create it using the predefined template
2. WHEN modifying a daily note, THE System SHALL only update the Goals, Meetings, and Bonus Items sections
3. WHEN updating sections, THE System SHALL preserve all other content in the note byte-for-byte
4. THE System SHALL perform surgical edits that replace only the specified sections
5. THE System SHALL ensure all note modifications are safe and idempotent

### Requirement 3: Multi-Source Data Integration

**User Story:** As a user, I want the system to automatically gather information from my weekly goals, calendar, and Jira tickets, so that my daily plan reflects all my commitments and priorities.

#### Acceptance Criteria

1. WHEN initiating a planning session, THE System SHALL read weekly goals from the weekly note
2. WHEN gathering calendar data, THE Calendar_Adapter SHALL fetch events from a private iCal URL
3. WHEN retrieving Jira data, THE Jira_Adapter SHALL fetch tickets assigned to the user
4. WHEN processing external data, THE Sanitization_Layer SHALL remove all PII before LLM processing
5. THE System SHALL combine data from all sources into a unified planning context

### Requirement 4: Calendar Event Processing

**User Story:** As a user, I want my calendar events to be reflected in my daily note as meetings without time details, so that I have a clean overview of my scheduled interactions.

#### Acceptance Criteria

1. WHEN processing calendar events, THE Calendar_Adapter SHALL extract event title and description
2. WHEN rendering meetings in the daily note, THE System SHALL exclude event times
3. WHEN a meeting warrants long-lived context, THE System SHALL create a link in the format [[YYYY/MM/Meeting Name|Meeting Name]]
4. WHEN creating meeting links, THE System SHALL only create the link and let Obsidian create the file on demand
5. THE System SHALL classify meetings to determine which ones receive links

### Requirement 5: Jira Ticket Integration

**User Story:** As a user, I want my assigned Jira tickets to be considered in my daily planning, so that my work commitments are properly reflected in my goals.

#### Acceptance Criteria

1. WHEN retrieving Jira tickets, THE Jira_Adapter SHALL fetch summary, description, status, priority, and due date fields
2. WHEN processing Jira data, THE Sanitization_Layer SHALL redact sensitive values like emails, account numbers, and usernames
3. WHEN sanitizing data, THE System SHALL ensure redaction is not reversible
4. THE System SHALL authenticate with Jira using API tokens
5. WHEN incorporating tickets into planning, THE System SHALL consider ticket priority and due dates

### Requirement 6: Task Carryover Management

**User Story:** As a user, I want incomplete tasks from previous days to be surfaced and managed, so that I can decide whether to continue, defer, or drop them.

#### Acceptance Criteria

1. WHEN incomplete tasks exist in the Goals or Bonus Items sections from the previous business day, THE System SHALL detect and surface them
2. WHEN carryover tasks are found, THE System SHALL ask the user whether to carry them forward, drop them, or explicitly defer them
3. WHEN the user chooses to carry forward tasks, THE System SHALL add them to today's Goals section
4. WHEN the user chooses to defer tasks, THE System SHALL handle the deferral appropriately
5. THE System SHALL require user confirmation before removing or deferring tasks

### Requirement 7: Unplanned Work Observation

**User Story:** As a user, I want the system to track unplanned work that arises during my day and use historical patterns to buffer against over-capacity, so that I can understand my capacity patterns and plan more realistically.

#### Acceptance Criteria

1. WHEN items are added to the Bonus Items section, THE System SHALL observe them as unplanned work signals
2. WHEN generating capacity assessments, THE System SHALL create qualitative descriptions like "a little", "some", "mostly unplanned", or "everything was unplanned"
3. WHEN calculating capacity trends, THE System SHALL use a rolling average approach to account for multi-day unplanned work patterns
4. WHEN unplanned work occurs, THE System SHALL consider that subsequent days are likely to have higher than average unplanned work
5. WHEN writing capacity assessments, THE System SHALL record them in the weekly journal note
6. THE System SHALL focus on narrative trend awareness and capacity buffering rather than precise metrics
7. WHEN the daily note is initially created, THE System SHALL leave the Bonus Items section empty

### Requirement 8: Local AI Integration

**User Story:** As a system administrator, I want to use local LLMs for development and operation, so that I maintain control over my data and can learn AI system design patterns.

#### Acceptance Criteria

1. THE System SHALL use Ollama for local LLM execution
2. WHEN selecting models, THE System SHALL ensure they support structured output and tool usage
3. WHEN processing with the LLM, THE System SHALL provide sanitized context and current Markdown sections
4. WHEN receiving LLM output, THE System SHALL expect structured plans and proposed section replacements
5. THE System SHALL treat the LLM as a planner that proposes changes rather than directly mutating state

### Requirement 9: Security and Data Protection

**User Story:** As a user, I want my sensitive data to be protected and API credentials to be stored securely, so that my personal and work information remains safe.

#### Acceptance Criteria

1. WHEN storing Jira API tokens, THE System SHALL use secure storage methods like OS keychains or credential managers
2. WHEN first-time setup is required, THE System SHALL provide a configure command similar to AWS CLI setup
3. WHEN configuration is needed, THE System SHALL store settings in a configuration file at ~/.planner/config.json
2. WHEN processing external data, THE Sanitization_Layer SHALL remove PII before any LLM interaction
3. WHEN sanitizing data, THE System SHALL use deterministic, regex-based approaches
4. WHEN configuring PII removal, THE System SHALL read an array of regex patterns from the configuration file
5. WHEN applying PII sanitization, THE System SHALL use the configured regex patterns to identify and redact sensitive data
6. THE System SHALL ensure the LLM never sees raw Jira data
7. THE System SHALL maintain a clear boundary between raw external data and LLM-safe representations

### Requirement 11: Contextual Task Integration and Status Management

**User Story:** As a user, I want to provide contextual information about new tasks, status changes, and priority shifts during chat, so that the system can intelligently adapt my daily plan.

#### Acceptance Criteria

1. WHEN a user provides contextual information like "I was just asked to do xyz", THE System SHALL extract the task and incorporate it into the Goals section
2. WHEN a user reports status changes like "this jira is blocked, it just hasn't been marked as such yet", THE System SHALL adjust the plan to reflect the current reality
3. WHEN a user indicates priority shifts like "we are deferring this jira for later", THE System SHALL remove or reschedule affected items
4. WHEN a user asks about capacity like "how much capacity is left", THE System SHALL analyze current goals against available time and suggest additional tasks
5. WHEN processing contextual inputs, THE System SHALL consider existing schedule and priorities for optimal adjustments

**User Story:** As a developer, I want clear separation between components and external systems, so that the system is maintainable and can evolve to cloud deployment.

### Requirement 12: Capacity Analysis and Gap Filling

**User Story:** As a user, I want to understand my remaining capacity and get suggestions for additional work, so that I can optimize my daily productivity without over-scheduling.

#### Acceptance Criteria

1. WHEN a user asks about remaining capacity, THE System SHALL analyze current goals against meeting density to identify meaningful work blocks
2. WHEN evaluating capacity, THE System SHALL only consider gaps of 2+ hours as suitable for Jira tasks
3. WHEN gaps smaller than 2 hours exist, THE System SHALL not suggest Jira work and may suggest quick administrative tasks instead
4. WHEN the user has manually added items to goals, THE System SHALL read the current daily note state and incorporate existing changes
5. WHEN suggesting additional work, THE System SHALL prioritize not over-filling the day and maintaining realistic expectations

**User Story:** As a developer, I want clear separation between components and external systems, so that the system is maintainable and can evolve to cloud deployment.

### Requirement 10: Architectural Modularity

**User Story:** As a developer, I want clear separation between components and external systems, so that the system is maintainable and can evolve to cloud deployment.

#### Acceptance Criteria

1. WHEN designing components, THE System SHALL isolate each external system behind a thin adapter
2. WHEN processing data, THE Orchestrator SHALL coordinate workflow without performing parsing or heuristics
3. WHEN handling calendar data, THE Calendar_Adapter SHALL treat ingestion as a pure function
4. WHEN managing Obsidian files, THE Obsidian_Adapter SHALL serve as the single source of truth for state
5. THE System SHALL follow stateless execution patterns compatible with serverless architectures