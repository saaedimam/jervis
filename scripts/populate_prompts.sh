#!/bin/bash
# scripts/populate_prompts.sh
source .env

create_prompt() {
  local name=$1
  local category=$2
  local content=$3

  echo "Populating Prompt: $name"
  
  page_id=$(curl -s -X POST "https://api.notion.com/v1/databases/$PROMPTS_DB/query" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    -d "{
      \"filter\": {
        \"property\": \"Prompt ID\",
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
    parent="\"parent\": { \"database_id\": \"$PROMPTS_DB\" },"
  fi

  curl -s -X $method "$url" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Content-Type: application/json" \
    -H "Notion-Version: 2022-06-28" \
    -d "{
      $parent
      \"properties\": {
        \"Prompt ID\": { \"title\": [ { \"text\": { \"content\": \"$name\" } } ] }
      },
      \"children\": [
        {
          \"object\": \"block\",
          \"type\": \"code\",
          \"code\": {
            \"language\": \"markdown\",
            \"rich_text\": [ { \"text\": { \"content\": $(jq -Rs . <<< "$content") } } ]
          }
        }
      ]
    }" > /dev/null
}

# Sync prompts/context_sync.md
content=$(cat prompts/context_sync.md)
create_prompt "Context Sync Workflow" "Workflow" "$content"
