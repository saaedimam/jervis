#!/bin/bash
# Populate Jervis Milestones in Notion

MILESTONES_DB="7cb28553-9f9b-4a6d-857b-5f9c724ddf22"

add_milestone() {
  local name="$1"
  local phase="$2"
  local status="$3"
  local date="$4"
  local coverage="$5"
  
  local json
  if [ -n "$date" ]; then
    json=$(cat <<EOF
{
  "parent": {"database_id": "$MILESTONES_DB"},
  "properties": {
    "Name": {"title": [{"text": {"content": "$name"}}]},
    "Phase": {"select": {"name": "$phase"}},
    "Status": {"select": {"name": "$status"}},
    "Coverage": {"rich_text": [{"text": {"content": "$coverage"}}]},
    "Completion Date": {"date": {"start": "$date"}}
  }
}
EOF
)
  else
    json=$(cat <<EOF
{
  "parent": {"database_id": "$MILESTONES_DB"},
  "properties": {
    "Name": {"title": [{"text": {"content": "$name"}}]},
    "Phase": {"select": {"name": "$phase"}},
    "Status": {"select": {"name": "$status"}},
    "Coverage": {"rich_text": [{"text": {"content": "$coverage"}}]}
  }
}
EOF
)
  fi
  
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "$json" | jq -r '.id'
}

add_milestone "Architecture Frozen" "Phase 1.1" "Done" "2026-07-28" "N/A"
add_milestone "Engineering Specification Frozen" "Phase 1.1" "Done" "2026-07-28" "N/A"
add_milestone "Governance Specification Frozen" "Phase 1.1" "Done" "2026-07-28" "N/A"
add_milestone "Repository Infrastructure Bootstrap" "Phase 1.1" "Done" "2026-07-28" "N/A"
add_milestone "Phase 1.1 Core Runtime Foundation" "Phase 1.1" "Done" "2026-07-28" "100%"
add_milestone "Phase 1.2 Event Bus Spec Frozen" "Phase 1.2" "Done" "2026-07-28" "N/A"
add_milestone "Phase 1.2.8 Event Bus Facade" "Phase 1.2" "Done" "2026-07-28" "100%"
add_milestone "Phase 1.3 Permission Engine" "Phase 1.3" "Done" "2026-07-28" "100%"
add_milestone "Phase 1.4.0 Observer Architecture" "Phase 1.4" "Done" "2026-07-28" "N/A"
add_milestone "Phase 1.4.1 Observer Foundation" "Phase 1.4" "Done" "2026-07-29" "100%"
add_milestone "Phase 1.4.2 Observer Registry" "Phase 1.4" "In Progress" "" "TBD"
add_milestone "Phase 1.4.3 Observer Dispatcher" "Phase 1.4" "Pending" "" "TBD"
add_milestone "Phase 1.4.4 Observer Facade" "Phase 1.4" "Pending" "" "TBD"
add_milestone "Phase 1.5 Scheduler" "Phase 1.5" "Pending" "" "TBD"
add_milestone "Phase 1.6 Session Engine" "Phase 1.6" "Pending" "" "TBD"
add_milestone "Phase 2 Memory Engine" "Phase 2" "Pending" "" "TBD"
add_milestone "Phase 3 Domain Services" "Phase 3" "Pending" "" "TBD"
add_milestone "Phase 4 AI Providers" "Phase 4" "Pending" "" "TBD"
add_milestone "Phase 5 Interfaces" "Phase 5" "Pending" "" "TBD"
