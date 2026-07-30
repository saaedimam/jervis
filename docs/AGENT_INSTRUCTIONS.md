# Jervis: Agent Instructions

## 1. Role
You are a development agent for the Jervis Project. Your goal is to maintain the project's architectural integrity and follow the frozen specifications.

## 2. Source of Truth
- **Repository**: `/Users/ioriimasu/dev/jervis`
- **Notion OS**: `https://www.notion.so/3ab1b27fdcba81d08b35ed766e2e8420`

## 3. Key Databases
- **Architecture**: `https://www.notion.so/3ab1b27fdcba81c889bec6e885f81f5c`
- **Specifications**: `https://www.notion.so/3ab1b27fdcba81788316c73d93b940e7`
- **Context**: `https://www.notion.so/3ab1b27fdcba813c848ed922a96996d9`
- **Session Logs**: `https://www.notion.so/3ab1b27fdcba817d9488c0b11b753f7c`

## 4. Operational Protocols
- **Sync**: Always sync your session log to the Session Logs database after a significant change.
- **Freeze**: Before modifying any frozen API, check the API Freeze database.
- **Decisions**: Record any architectural decisions in the ADRs database.
