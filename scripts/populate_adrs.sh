#!/bin/bash
# scripts/populate_adrs.sh
source .env

create_adr() {
  local id=$1
  local title=$2
  local status=$3
  local decision=$4

  echo "Populating ADR: $id - $title"
  
  page_id=$(curl -s -X POST "https://api.notion.com/v1/databases/$ADRS_DB/query" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    -d "{
      \"filter\": {
        \"property\": \"ADR\",
        \"title\": { \"equals\": \"$id: $title\" }
      }
    }" | jq -r '.results[0].id // empty')

  if [ -n "$page_id" ]; then
    method="PATCH"
    url="https://api.notion.com/v1/pages/$page_id"
    parent=""
  else
    method="POST"
    url="https://api.notion.com/v1/pages"
    parent="\"parent\": { \"database_id\": \"$ADRS_DB\" },"
  fi

  curl -s -X $method "$url" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Content-Type: application/json" \
    -H "Notion-Version: 2022-06-28" \
    -d "{
      $parent
      \"properties\": {
        \"ADR\": { \"title\": [ { \"text\": { \"content\": \"$id: $title\" } } ] },
        \"Status\": { \"select\": { \"name\": \"$status\" } },
        \"Decision\": { \"rich_text\": [ { \"text\": { \"content\": \"$decision\" } } ] }
      }
    }" > /dev/null
}

# ADR-0001
create_adr "ADR-0001" "Initial Architecture Baseline" "Superseded" "Initial project setup using TypeScript/Node.js."
# ADR-0002
create_adr "ADR-0002" "Architecture Reconciliation" "Accepted" "Transitioned to Go 1.22+, 5-tier hierarchy, and AI decoupling."
