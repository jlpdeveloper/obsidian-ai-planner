# Obsidian AI Planner

A personal, conversational AI-assisted daily planning system that helps you manage your workday by integrating your Obsidian vault, Jira, and Calendar.

## Overview

The Obsidian AI Planner is designed to pull data from multiple sources—including your weekly goals in Obsidian, upcoming calendar events via iCal, and assigned Jira tickets—to help you build a structured daily plan. It uses a conversational interface to propose and apply updates directly to your Obsidian daily notes.

## Key Features

- **Conversational Planning:** Interact with an AI each morning to summarize your day and prioritize tasks.
- **Multi-Source Integration:**
    - **Obsidian:** Reads weekly goals and manages daily notes.
    - **Jira:** Pulls assigned tickets (with built-in PII sanitization).
    - **Calendar:** Read-only access to events via iCal.
- **Local-First & Safe:**
    - Uses local LLMs (via [Ollama](https://ollama.com/)).
    - Performs "surgical" edits on your Markdown files, preserving your existing content outside of the managed sections (**Goals**, **Meetings**, **Bonus Items**).
- **Markdown-First:** Your notes remain the single source of truth.

## Architecture

The system follows a predictable "Pull, Sanitize, Plan, Write" pipeline:

1.  **Pull:** Gathers context from Obsidian, Jira, and Calendar.
2.  **Sanitize:** Redacts sensitive information (PII) before sending data to the LLM.
3.  **Plan:** The LLM (acting as a planner) proposes modifications based on the context and user intent.
4.  **Write:** The system applies approved "surgical" updates to the daily note sections.

## Tech Stack

- **Language:** Go
- **AI Integration:** [Genkit](https://github.com/firebase/genkit) / Ollama
- **UI:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) (planned TUI)
- **Data Format:** Markdown

## Roadmap

- **v1:** Basic ingestion (Calendar/Jira), PII sanitization, and conversational daily note management.
- **v2:** Smarter meeting classification, task carryover logic, and capacity trend analysis.
- **v3:** Calendar mutation (automatic focus time blocking) and historical analysis.


## CI/CD

This project uses GitHub Actions for continuous integration:
- **Build & Test:** Automatically runs on every pull request to ensure code quality and prevent regressions.
- **Release Management:** Uses Release Drafter to automate release notes and versioning.

## Project Tracking
This project is also an experiment in using [Fizzy](https://app.fizzy.do/6130918/public/boards/EcoT3co8ffcS6zpMuFbfFYbX), 
all work will be tracked there.