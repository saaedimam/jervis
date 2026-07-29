#!/bin/bash
# Add Dashboard entries to Notion

DASHBOARD_DB="6b2ad157-1284-45a0-808b-cd352f73ab0a"
DATE=$(date -u +%Y-%m-%d)

add_dashboard() {
  local name="$1"
  local value="$2"
  local category="$3"
  
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "{
      \"parent\": {\"database_id\": \"$DASHBOARD_DB\"},
      \"properties\": {
        \"Name\": {\"title\": [{\"text\": {\"content\": \"$name\"}}]},
        \"Value\": {\"rich_text\": [{\"text\": {\"content\": \"$value\"}}]},
        \"Category\": {\"select\": {\"name\": \"$category\"}},
        \"Last Updated\": {\"date\": {\"start\": \"$DATE\"}}
      }
    }" | jq -r '.id'
}

add_dashboard "Current Phase" "Phase 1.4.2" "Status"
add_dashboard "Overall Progress" "23/32 milestones (72%)" "Progress"
add_dashboard "Current Session" "2026-07-29-session-17" "Session"
add_dashboard "Latest Commit" "91ba7ad" "Session"
add_dashboard "Latest Commit Message" "chore: finalize governance verification and audit artifacts" "Session"
add_dashboard "Repository Health" "Healthy - All tests pass" "Health"
add_dashboard "Test Coverage" "100% on completed phases" "Health"
add_dashboard "Go Files Count" "84" "Status"
add_dashboard "Next Step" "Implement Phase 1.4.2 Observer Registry" "Status"
