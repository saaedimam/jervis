#!/bin/bash
# scripts/notion_sync_jq.sh
# File synchronization script for Jervis Notion integration

# Load centralized configuration
REPO_PATH="/Users/ioriimasu/dev/jervis"
source "$REPO_PATH/scripts/load_config.sh"
load_config

# Use centralized credential and database ID management
# SPEC_DB, CONTEXT_DB, SESSION_DB are now imported from load_config.sh via notion_databases.json

sync_file() {
  local file=$1
  local db_id=$2
  local db_name=$3
  local prop_name=$4

  echo "Syncing $file to database $db_name ($db_id)..."

  title=$(basename "$file" .md)

  # Search for an existing page with the same Repository Path to ensure idempotency
  local filter=$(jq -n --arg file "$file" '{"property":"Repository Path","rich_text":{"equals":$file}}')
  local query_res=$(ntn datasources query "$db_id" --filter "$filter" --json < /dev/null 2>/dev/null || echo "error")

  local page_id=""
  if [ "$query_res" != "error" ]; then
    page_id=$(echo "$query_res" | jq -r '.results[0].id // empty')
  fi

  if [ -n "$page_id" ] && [ "$page_id" != "null" ]; then
    echo "Page already exists (ID: $page_id). Updating in-place..."

    # 1. Update title and properties
    jq -n \
      --arg title "$title" \
      --arg prop_name "$prop_name" \
      '{
        "properties": {
          ($prop_name): { "title": [ { "text": { "content": $title } } ] }
        }
      }' | ntn api /v1/pages/"$page_id" -X PATCH -d @- > /dev/null

    # 2. Retrieve children to update the existing code block
    local children_res=$(ntn api /v1/blocks/"$page_id"/children < /dev/null 2>/dev/null)
    local code_block_id=$(echo "$children_res" | jq -r '.results[] | select(.type == "code") | .id' | head -1)

    if [ -n "$code_block_id" ] && [ "$code_block_id" != "null" ]; then
      # Update the existing code block content using range-based chunking with 1900-char buffer
      jq -n \
        --rawfile content "$file" \
        '{
          "type": "code",
          "code": {
            "language": "markdown",
            "rich_text": [ range(0; $content|length; 1900) as $i | { "text": { "content": $content[$i:$i+1900] } } ]
          }
        }' | ntn api /v1/blocks/"$code_block_id" -X PATCH -d @- > /dev/null
    else
      # Append new code block if none existed
      jq -n \
        --rawfile content "$file" \
        '{
          "children": [
            {
              "object": "block",
              "type": "code",
              "code": {
                "language": "markdown",
                "rich_text": [ range(0; $content|length; 1900) as $i | { "text": { "content": $content[$i:$i+1900] } } ]
              }
            }
          ]
        }' | ntn api /v1/blocks/"$page_id"/children -X PATCH -d @- > /dev/null
    fi
    echo "Update complete."
  else
    echo "No existing page found. Creating new page..."
    # Use jq to build the JSON payload with 1900-character range chunking
    jq -n \
      --arg db_id "$db_id" \
      --arg title "$title" \
      --arg file_path "$file" \
      --arg prop_name "$prop_name" \
      --rawfile content "$file" \
      '{
        "parent": { "database_id": $db_id },
        "properties": {
          ($prop_name): { "title": [ { "text": { "content": $title } } ] },
          "Repository Path": { "rich_text": [ { "text": { "content": $file_path } } ] }
        },
        "children": [
          {
            "object": "block",
            "type": "code",
            "code": {
              "language": "markdown",
              "rich_text": [ range(0; $content|length; 1900) as $i | { "text": { "content": $content[$i:$i+1900] } } ]
            }
          }
        ]
      }' | ntn api /v1/pages -d @- > /dev/null
    echo "Creation complete."
  fi
}

# 1. Sync Specifications
find . -maxdepth 1 -name "*SPECIFICATION.md" | while read -r f; do
  sync_file "$f" "$SPEC_DB" "Specifications" "Specification"
done

# 2. Sync Context
find context -maxdepth 1 -name "*.md" | while read -r f; do
  sync_file "$f" "$CONTEXT_DB" "Context" "Name"
done

# 3. Sync Sessions
find context/sessions -maxdepth 1 -name "*.md" | while read -r f; do
  sync_file "$f" "$SESSION_DB" "Session Logs" "Session ID"
done