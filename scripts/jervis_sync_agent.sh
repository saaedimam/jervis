#!/bin/bash
# Jervis Sync Agent Script
# Agent-driven synchronization for Jervis

# Load centralized configuration
REPO_PATH="/Users/ioriimasu/dev/jervis"
source "$REPO_PATH/scripts/load_config.sh"
load_config

REPO_PATH="/Users/ioriimasu/dev/jervis"
STATE_FILE="$REPO_PATH/.jervis_agent_state"
LOG_FILE="$REPO_PATH/.jervis_agent.log"

log() {
  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] $1" >> "$LOG_FILE"
}

error() {
  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] ERROR: $1" >> "$LOG_FILE"
  echo "[ERROR] $1" >&2
  exit 1
}

# Sync session context to the Context database
sync_session_context() {
  local session_file="$1"
  if [ ! -f "$session_file" ]; then
    log "Session file not found: $session_file"
    return 1
  fi

  local session_name=$(basename "$session_file" .md)
  local content=$(cat "$session_file")

  # Prepare the data for Notion API
  data=$(jq -n \
    --arg name "$session_name" \
    --arg content "$content" \
    '{
      "parent": { "database_id": env.JERVIS_PAGE },
      "properties": {
        "Name": { "title": [ { "text": { "content": $name } } ] },
        "Repository Path": { "rich_text": [ { "text": { "content": env.session_file } } ] }
      },
      "children": [
        {
          "object": "block",
          "type": "code",
          "code": {
            "language": "markdown",
            "rich_text": [ { "text": { "content": $content } } ]
          }
        }
      ]
    }')

  # Use ntn API to create the page
  echo "$data" | ntn api /v1/pages -X POST -d @- > /dev/null 2>&1 || {
    error "Failed to create session page for $session_file"
  }
}

# Update dashboard with last sync time
update_dashboard() {
  local date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  curl -s -X PATCH "https://api.notion.com/v1/databases/$DASHBOARD_DB" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    -d "{\"description\": [{\"text\": {\"content\": \"Last sync: $date\"}}]}" > /dev/null 2>&1 || true
}

# Main function
main() {
  # Sync session context (example: sync the latest session)
  local session_file="$REPO_PATH/context/SESSION_CONTEXT.md"
  if [ -f "$session_file" ]; then
    sync_session_context "$session_file"
  fi

  # Update dashboard
  update_dashboard
}

main