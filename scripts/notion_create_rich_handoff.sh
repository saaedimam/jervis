#!/bin/bash
# Rich AI Handoff with complete context
source .env
HANDOFF_DB="$AI_HANDOFF_DB"

BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "main")
COMMIT=$(git log --oneline -1 2>/dev/null | cut -d' ' -f1 || echo "UNKNOWN")
COMMIT_MSG=$(git log --format=%B -n 1 2>/dev/null | head -1 || echo "N/A")
DATE=$(date -u +%Y-%m-%d)
TIME=$(date -u +%H:%M:%SZ)

# Get modified packages from last commit
MODIFIED_PKGS=$(cd /Users/ioriimasu/dev/jervis && git diff-tree --no-commit-id --name-only -r HEAD 2>/dev/null | grep -E '^internal/.*\.go$' | sed 's|/[^/]*$||' | sort -u | tr '\n' ', ' | sed 's/, $//' || echo "None")

# Get affected specs
AFFECTED_SPECS="SPEC-001..022"

# Get architecture impact
ARCH_IMPACT="ARCH-001..004"

# Calculate risk level
RISK_LEVEL="Low"
if echo "$COMMIT_MSG" | grep -qi "breaking\|BREAKING"; then
  RISK_LEVEL="High"
fi

json=$(cat <<EOF
{
  "parent": {"database_id": "$HANDOFF_DB"},
  "properties": {
    "Handoff ID": {"title": [{"text": {"content": "HANDOFF-001"}}]},
    "Session": {"rich_text": [{"text": {"content": "SESSION-017"}}]},
    "Branch": {"rich_text": [{"text": {"content": "$BRANCH"}}]},
    "Commit": {"rich_text": [{"text": {"content": "$COMMIT"}}]},
    "Current Goal": {"rich_text": [{"text": {"content": "Implement P0 canonical engineering knowledge graph with complete traceability"}}]},
    "Completed": {"rich_text": [{"text": {"content": "✅ Architecture Registry (ARCH-001..004) with 4 components; ✅ Package Registry (PKG-001..029) with 29 packages; ✅ File Registry (FILE-0001..0023) with 23 tracked files; ✅ Specification Registry (SPEC-001..022) with 12 specifications; ✅ Engineering Memory with 11 lessons/patterns/rules; ✅ AI Handoff system; ✅ Commit Intelligence; ✅ Dependency Graph with belongs_to, implements, approved_by, affects relations"}}]},
    "Blocked": {"rich_text": [{"text": {"content": "None - proceeding to implementation"}}]},
    "Next Task": {"rich_text": [{"text": {"content": "Implement Phase 1.4.2 Observer Registry (PKG-026..028) following SPEC-020..022"}}]},
    "Modified Files": {"rich_text": [{"text": {"content": "$MODIFIED_PKGS"}}]},
    "Risks": {"rich_text": [{"text": {"content": "$RISK_LEVEL - Foundation complete, ready for implementation"}}]}
  }
}
EOF
)

curl -s -X POST "https://api.notion.com/v1/pages" \
  -H "Authorization: Bearer $NOTION_API_KEY" \
  -H "Notion-Version: 2025-09-03" \
  -H "Content-Type: application/json" \
  -d "$json" | jq -r '.id'

echo "Rich AI Handoff created"
echo ""
echo "HANDOFF-001 includes:"
echo "- Current Branch: $BRANCH"
echo "- HEAD Commit: $COMMIT"
echo "- Modified Packages: $MODIFIED_PKGS"
echo "- Affected Specs: $AFFECTED_SPECS"
echo "- Architecture Impact: $ARCH_IMPACT"
echo "- Risk Level: $RISK_LEVEL"
echo "- Required Reading: PROJECT_CONTEXT.md, SESSION_CONTEXT.md, Architecture Registry, Package Registry"
