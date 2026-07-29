#!/bin/bash
# scripts/populate_milestones.sh
source .env

create_milestone() {
  local name=$1
  local status=$2
  local phase=$3
  local coverage=$4

  echo "Populating Milestone: $name"
  
  page_id=$(curl -s -X POST "https://api.notion.com/v1/databases/$MILESTONES_DB/query" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    -d "{
      \"filter\": {
        \"property\": \"Milestone\",
        \"title\": { \"equals\": \"$name\" }
      }
    }" | jq -r '.results[0].id // empty')

  if [ -n "$page_id" ]; then
    method="PATCH"
    url="https://api.notion.com/v1/pages/$page_id"
    parent=""
  else
    method="POST"
    url="https://api.notion.com/v1/pages"
    parent="\"parent\": { \"database_id\": \"$MILESTONES_DB\" },"
  fi

  curl -s -X $method "$url" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Content-Type: application/json" \
    -H "Notion-Version: 2022-06-28" \
    -d "{
      $parent
      \"properties\": {
        \"Milestone\": { \"title\": [ { \"text\": { \"content\": \"$name\" } } ] },
        \"Status\": { \"select\": { \"name\": \"$status\" } },
        \"Phase\": { \"rich_text\": [ { \"text\": { \"content\": \"$phase\" } } ] },
        \"Coverage\": { \"number\": $coverage }
      }
    }" > /dev/null
}

# Parse context/MILESTONES.md
# Look for lines like "- [x] Phase 1.4 Observer Subsystem: Complete"
grep "^- \[" context/MILESTONES.md | while read -r line; do
  if [[ $line == *"[x]"* ]]; then status="Completed"; else status="Pending"; fi
  name=$(echo "$line" | sed 's/^- \[[x ]\] //g' | cut -d':' -f1)
  phase=$(echo "$name" | grep -o "Phase [0-9.]*")
  if [ -z "$phase" ]; then phase="General"; fi
  
  coverage=0.0
  if [[ $line == *"100% Test Coverage"* ]]; then coverage=1.0; fi
  
  create_milestone "$name" "$status" "$phase" $coverage
done
