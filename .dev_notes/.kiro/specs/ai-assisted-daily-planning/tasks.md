# Implementation Plan: AI-Assisted Daily Planning

## Overview

Implementation approach for a personal learning project focused on AI integration patterns. Tasks are designed to be incremental and build understanding of conversational AI systems while keeping complexity manageable.

## Tasks

- [ ] 1. Project Setup and Configuration
  - Initialize Go module with Genkit and Bubbletea dependencies
  - Set up basic project structure (cmd/, internal/, pkg/)
  - Create configuration system with JSON file support
  - Add basic CLI framework with Bubbletea
  - _Requirements: 9.1, 9.2, 9.3_

- [ ] 2. Configuration Management
  - [ ] 2.1 Implement interactive configure command
    - Create interactive configuration TUI with Bubbletea
    - Prompt for Jira URL, API token, vault path, calendar URL with forms
    - Store configuration in ~/.planner/config.json
    - Add validation and error handling in the TUI
    - _Requirements: 9.1, 9.2, 9.3_

  - [ ]* 2.2 Write unit tests for configuration
    - Test configuration validation
    - Test file creation and loading
    - _Requirements: 9.1, 9.2, 9.3_

- [ ] 3. File Management Foundation
  - [ ] 3.1 Implement basic Obsidian file operations
    - Create functions to read/write daily notes
    - Implement surgical editing for Goals, Meetings, Bonus Items sections
    - Add weekly note reading capability
    - _Requirements: 2.1, 2.2, 2.3_

  - [ ]* 3.2 Write property test for surgical editing
    - **Property 1: File Safety**
    - **Validates: Requirements 2.2, 2.3**

  - [ ] 3.3 Add Markdown parsing utilities
    - Parse existing goals from daily notes
    - Identify and preserve non-target sections
    - _Requirements: 2.2, 2.3_

- [ ] 4. External Data Integration
  - [ ] 4.1 Implement calendar data fetching
    - Create iCal URL fetcher using HTTP client
    - Parse ICS format to extract events (title, description, date only)
    - _Requirements: 3.2, 4.1, 4.2_

  - [ ]* 4.2 Write property test for calendar processing
    - **Property 3: Calendar Processing**
    - **Validates: Requirements 4.1, 4.2**

  - [ ] 4.3 Implement Jira API integration
    - Create Jira REST API client with token authentication
    - Fetch assigned tickets with required fields
    - _Requirements: 3.3, 5.1_

  - [ ] 4.4 Add PII sanitization
    - Implement regex-based sanitization using configured patterns
    - Apply to Jira data before LLM processing
    - _Requirements: 3.4, 5.2, 9.4, 9.5_

  - [ ]* 4.5 Write property test for PII sanitization
    - **Property 2: PII Sanitization**
    - **Validates: Requirements 3.4, 5.2, 9.4, 9.5**

- [ ] 5. Checkpoint - Basic Data Pipeline
  - Ensure configuration, file operations, and data fetching work
  - Test with real Jira and calendar data
  - Verify PII sanitization is working correctly

- [ ] 6. LLM Integration with Genkit
  - [ ] 6.1 Set up Genkit with Ollama
    - Initialize Genkit Go SDK
    - Configure connection to local Ollama instance
    - Test basic LLM connectivity
    - _Requirements: 8.1, 8.2_

  - [ ] 6.2 Create planning prompt templates
    - Design system prompt for daily planning
    - Create templates for different conversation types
    - Handle structured output for day plans
    - _Requirements: 1.1, 8.3, 8.4_

  - [ ] 6.3 Implement basic planning workflow
    - Gather context from all data sources
    - Send sanitized data to LLM with user prompt
    - Parse LLM response into DayPlan struct
    - _Requirements: 1.1, 3.5, 8.3_

- [ ] 7. Conversational Interface
  - [ ] 7.1 Implement interactive chat interface
    - Create Bubbletea model for conversational planning
    - Handle initial prompts and display responses with rich formatting
    - Show proposed changes in formatted preview panes
    - _Requirements: 1.1, 1.3, 1.4_

  - [ ]* 7.2 Write property test for conversational flow
    - **Property 4: Conversational Flow**
    - **Validates: Requirements 1.1, 1.4**

  - [ ] 7.3 Add iterative conversation and plan acceptance
    - Allow users to refine plans through continued conversation
    - Implement clear accept/reject interface with keybindings
    - Show real-time preview of what will be written to files
    - Handle user iterations with conversation history
    - _Requirements: 1.4, 1.5_

- [ ] 8. Advanced Planning Features
  - [ ] 8.1 Implement carryover task detection
    - Read previous day's daily note
    - Identify incomplete tasks in Goals and Bonus Items
    - Present carryover options to user
    - _Requirements: 6.1, 6.2, 6.3_

  - [ ] 8.2 Add contextual task integration
    - Handle prompts like "I was just asked to do xyz"
    - Extract new tasks and integrate into Goals section
    - Support status changes ("this jira is blocked")
    - _Requirements: 11.1, 11.2, 11.3_

  - [ ] 8.3 Implement capacity analysis
    - Analyze current goals against meeting density
    - Only suggest Jira tasks for 2+ hour blocks
    - Provide qualitative capacity assessments
    - _Requirements: 12.1, 12.2, 12.3, 12.5_

- [ ] 9. Meeting Management
  - [ ] 9.1 Add meeting classification
    - Determine which meetings need Obsidian links
    - Create links in [[YYYY/MM/Meeting Name|Meeting Name]] format
    - Avoid creating actual files (let Obsidian handle on-demand)
    - _Requirements: 4.3, 4.4_

  - [ ] 9.2 Implement meeting rendering
    - Render meetings without time information
    - Handle meeting descriptions appropriately
    - _Requirements: 4.2_

- [ ] 10. Unplanned Work Tracking
  - [ ] 10.1 Add capacity assessment
    - Track items added to Bonus Items section
    - Generate qualitative assessments ("a little", "some", etc.)
    - Write assessments to weekly journal note
    - _Requirements: 7.1, 7.2, 7.5_

  - [ ] 10.2 Implement rolling average calculation
    - Calculate capacity trends over time
    - Consider multi-day unplanned work patterns
    - Use for future capacity buffering
    - _Requirements: 7.3, 7.4_

- [ ] 11. Integration and Polish
  - [ ] 11.1 Add command routing and main interface
    - Create main Bubbletea application with command selection
    - Implement navigation between configure, chat, and overview modes
    - Add help system and keybinding documentation
    - _Requirements: 1.1_

  - [ ]* 11.2 Write integration tests
    - Test end-to-end planning workflow
    - Test with mock data for reliability
    - _Requirements: All_

  - [ ] 11.3 Add error handling and logging
    - Implement graceful error handling
    - Add clear logging for debugging
    - Handle network failures and API issues
    - _Requirements: All_

- [ ] 12. Final Checkpoint
  - Test complete workflow with real data
  - Verify all core properties work correctly
  - Document any limitations or future improvements

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Focus on learning AI integration patterns over enterprise complexity
- Property tests validate the 4 core correctness properties
- Personal project approach prioritizes simplicity and learning