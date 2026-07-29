# SCHEMA RECOVERY PLAN

## 🏛️ Architecture Registry
### Current Properties:
- **Purpose** (rich_text)
- **Repository Path** (rich_text)
- **Status** (select)
- **Specifications** (relation)
- **Coverage** (number)
- **Risks** (rich_text)
- **Files** (relation)
- **Related ADRs** (relation)
- **Packages** (relation)
- **Subsystem** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 📦 Package Registry
### Current Properties:
- **Exported APIs** (rich_text)
- **Architecture** (relation)
- **Purpose** (rich_text)
- **Dependencies** (rich_text)
- **Repository Path** (rich_text)
- **Coverage** (number)
- **Implementing Specification** (relation)
- **Files** (relation)
- **APIs** (relation)
- **Status** (select)
- **Package** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 📄 File Registry
### Current Properties:
- **Status** (select)
- **Architecture** (relation)
- **Specification** (relation)
- **Commits** (relation)
- **Path** (rich_text)
- **Exported APIs** (relation)
- **Coverage** (number)
- **Package** (relation)
- **Frozen** (checkbox)
- **Language** (select)
- **Line Count** (number)
- **File ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 📋 Specification Registry
### Current Properties:
- **Repository Path** (rich_text)
- **Architecture** (relation)
- **Related Packages** (relation)
- **Files** (relation)
- **Version** (rich_text)
- **Status** (select)
- **Frozen** (checkbox)
- **Specification** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 📋 ADRs
### Current Properties:
- **Affected Components** (rich_text)
- **Status** (select)
- **Consequences** (rich_text)
- **Repository Path** (rich_text)
- **Related Architecture** (relation)
- **Decision** (rich_text)
- **ADR** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## ⚡ API Registry
### Current Properties:
- **Frozen** (checkbox)
- **API Name** (rich_text)
- **Line Number** (number)
- **Status** (select)
- **File** (relation)
- **Package** (relation)
- **Type** (select)
- **API ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 📅 Engineering Timeline
### Current Properties:
- **Name** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## ✅ Quality Gates
### Current Properties:
- **Name** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 🤖 AI Handoff
### Current Properties:
- **Name** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 🔄 Commit Intelligence
### Current Properties:
- **Author** (rich_text)
- **Date** (date)
- **Message** (rich_text)
- **Files** (relation)
- **Hash** (rich_text)
- **Session** (relation)
- **Branch** (rich_text)
- **Commit ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 🧠 Engineering Memory
### Current Properties:
- **Name** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 🎯 Milestones
### Current Properties:
- **Status** (select)
- **Dependencies** (rich_text)
- **Phase** (rich_text)
- **Coverage** (number)
- **Milestone** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 📊 Dashboard
### Current Properties:
- **Name** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## 🔗 Dependency Graph
### Current Properties:
- **Name** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Tasks
### Current Properties:
- **Name** (rich_text)
- **Status** (select)
- **Task ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Releases
### Current Properties:
- **Version** (rich_text)
- **Date** (date)
- **Release ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Bugs
### Current Properties:
- **Description** (rich_text)
- **Status** (select)
- **Bug ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Risks
### Current Properties:
- **Description** (rich_text)
- **Impact** (select)
- **Risk ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Technical Debt
### Current Properties:
- **Priority** (select)
- **Description** (rich_text)
- **Debt ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Research
### Current Properties:
- **Summary** (rich_text)
- **Topic** (rich_text)
- **Research ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Prompt Library
### Current Properties:
- **Purpose** (rich_text)
- **Context** (rich_text)
- **Prompt ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Lessons Learned
### Current Properties:
- **Mistake** (rich_text)
- **Situation** (rich_text)
- **Resolution** (rich_text)
- **Lesson ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Patterns
### Current Properties:
- **Name** (rich_text)
- **Description** (rich_text)
- **Pattern ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Performance Benchmarks
### Current Properties:
- **Metric** (rich_text)
- **Value** (rich_text)
- **Benchmark ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Coding Standards
### Current Properties:
- **Description** (rich_text)
- **Rule** (rich_text)
- **Standard ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Glossary
### Current Properties:
- **Definition** (rich_text)
- **Term** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Assets
### Current Properties:
- **Name** (rich_text)
- **Type** (select)
- **Asset ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## External Dependencies
### Current Properties:
- **Version** (rich_text)
- **Name** (rich_text)
- **Dependency ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Roadmap
### Current Properties:
- **Name** (rich_text)
- **Quarter** (select)
- **Item ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*

## Engineering Decisions
### Current Properties:
- **Decision** (rich_text)
- **Problem** (rich_text)
- **Rationale** (rich_text)
- **Decision ID** (title)

### Missing Properties to Create:
- *Auto-inferred based on canonical relations...*
