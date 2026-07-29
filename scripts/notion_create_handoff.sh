#!/bin/bash
# Create AI Handoff entry for current session

HANDOFF_DB="c1e36ebb-a3fc-4aea-a3d2-ac8214e1e40a"

# Get current git info
BRANCH=$(cd /Users/ioriimasu/dev/jervis && git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "main")
COMMIT=$(cd /Users/ioriimasu/dev/jervis && git log --oneline -1 2>/dev/null | cut -d' ' -f1 || echo "UNKNOWN")
DATE=$(date -u +%Y-%m-%d)

json=$(cat <<EOF
{
  "parent": {"database_id": "$HANDOFF_DB"},
  "properties": {
    "Handoff ID": {"title": [{"text": {"content": "HANDOFF-001"}}]},
    "Session": {"rich_text": [{"text": {"content": "2026-07-29-session-17"}}]},
    "Branch": {"rich_text": [{"text": {"content": "$BRANCH"}}]},
    "Commit": {"rich_text": [{"text": {"content": "$COMMIT"}}]},
    "Current Goal": {"rich_text": [{"text": {"content": "Implement P0 canonical engineering knowledge graph with Architecture Registry, Package Registry, Specification Registry, and Engineering Memory"}}]},
    "Completed": {"rich_text": [{"text": {"content": "Created Architecture Registry (ARCH-001 through ARCH-004), Package Registry (PKG-001 through PKG-029), Specification Registry (SPEC-001 through SPEC-022), Engineering Memory, AI Handoff, Commit Intelligence databases"}}]},
    "Blocked": {"rich_text": [{"text": {"content": "None"}}]},
    "Next Task": {"rich_text": [{"text": {"content": "Continue with Phase 1.4.2 Observer Registry implementation"}}]},
    "Modified Files": {"rich_text": [{"text": {"content": "Notion databases created, sync scripts in /scripts/"}}]},
    "Risks": {"rich_text": [{"text": {"content": "Low - Foundation work complete, ready for implementation"}}]},
    "Required Reading": {"rich_text": [{"text": {"content": "PROJECT_CONTEXT.md, SESSION_CONTEXT.md, MILESTONES.md, Architecture Registry in Notion"}}]}
  }
}
EOF
)

curl -s -X POST "https://api.notion.com/v1/pages" \
  -H "Authorization: Bearer $NOTION_API_KEY" \
  -H "Notion-Version: 2025-09-03" \
  -H "Content-Type: application/json" \
  -d "$json" | jq -r '.id'

echo "AI Handoff created"
