#!/bin/bash
# Jervis Knowledge Synchronization Agent
# Monitors repository and syncs changes to Notion
# Runs every 5 minutes

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
NOTION_PAGE_ID="3ab1b27f-dcba-81d0-8b35-ed766e2e8420"
MILESTONES_DB="7cb28553-9f9b-4a6d-857b-5f9c724ddf22"
ADRS_DB="6b1d8415-c863-4943-bbe4-4381672f48f0"
API_FREEZE_DB="61b41663-2c66-45a1-87c9-5f04b3fd82b6"
DASHBOARD_DB="6b2ad157-1284-45a0-808b-cd352f73ab0a"
SESSIONS_DB="3d199938-9c6a-4eb5-a99c-98ccd03b4581"
MASTER_CONTEXT_PAGE="3ab1b27f-dcba-81e1-add1-eaa5a277ef10"

STATE_FILE="$REPO_PATH/.jervis_sync_state"
LOG_FILE="$REPO_PATH/.jervis_sync_log"

log() {
  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] $1" >> "$LOG_FILE"
}

# Get file hash for change detection
get_hash() {
  if [ -f "$1" ]; then
    md5 -q "$1" 2>/dev/null || md5sum "$1" | cut -d' ' -f1
  else
    echo "missing"
  fi
}

# Check if file changed
file_changed() {
  local file="$1"
  local current_hash=$(get_hash "$file")
  local stored_hash=""
  
  if [ -f "$STATE_FILE" ]; then
    stored_hash=$(grep "^$file:" "$STATE_FILE" 2>/dev/null | cut -d: -f2)
  fi
  
  if [ "$current_hash" != "$stored_hash" ]; then
    return 0  # Changed
  else
    return 1  # Not changed
  fi
}

# Update stored hash
update_hash() {
  local file="$1"
  local hash=$(get_hash "$file")
  
  if [ -f "$STATE_FILE" ]; then
    grep -v "^$file:" "$STATE_FILE" > "$STATE_FILE.tmp" 2>/dev/null || true
    mv "$STATE_FILE.tmp" "$STATE_FILE"
  fi
  
  echo "$file:$hash" >> "$STATE_FILE"
}

# Sync PROJECT_CONTEXT.md
sync_project_context() {
  local file="$REPO_PATH/context/PROJECT_CONTEXT.md"
  if file_changed "$file"; then
    log "PROJECT_CONTEXT.md changed - syncing"
    # In production: parse and update Notion page content
    update_hash "$file"
  fi
}

# Sync SESSION_CONTEXT.md
sync_session_context() {
  local file="$REPO_PATH/context/SESSION_CONTEXT.md"
  if file_changed "$file"; then
    log "SESSION_CONTEXT.md changed - syncing"
    
    # Extract session info and update Sessions database
    local session_id=$(grep "Session ID" "$file" | head -1 | sed 's/.*`\(.*\)`/\1/')
    local phase=$(grep "Current Phase" "$file" | head -1 | sed 's/.*Phase \([0-9.]*\).*/\1/')
    local date=$(date -u +%Y-%m-%d)
    
    # Check if session exists, create if not
    curl -s -X POST "https://api.notion.com/v1/pages" \
      -H "Authorization: Bearer $NOTION_API_KEY" \
      -H "Notion-Version: 2025-09-03" \
      -H "Content-Type: application/json" \
      -d "{
        \"parent\": {\"database_id\": \"$SESSIONS_DB\"},
        \"properties\": {
          \"Session ID\": {\"title\": [{\"text\": {\"content\": \"$session_id\"}}]},
          \"Date\": {\"date\": {\"start\": \"$date\"}},
          \"Phase\": {\"select\": {\"name\": \"Phase $phase\"}},
          \"Status\": {\"select\": {\"name\": \"In Progress\"}},
          \"Summary\": {\"rich_text\": [{\"text\": {\"content\": \"Auto-synced from SESSION_CONTEXT.md\"}}]}
        }
      }" > /dev/null 2>&1 || true
    
    # Update Dashboard
    curl -s -X PATCH "https://api.notion.com/v1/databases/$DASHBOARD_DB" \
      -H "Authorization: Bearer $NOTION_API_KEY" \
      -H "Notion-Version: 2025-09-03" \
      -H "Content-Type: application/json" \
      -d "{\"description\": [{\"text\": {\"content\": \"Last sync: $date\"}}]}" > /dev/null 2>&1 || true
    
    update_hash "$file"
  fi
}

# Sync MILESTONES.md
sync_milestones() {
  local file="$REPO_PATH/context/MILESTONES.md"
  if file_changed "$file"; then
    log "MILESTONES.md changed - syncing"
    # Parse markdown and update Milestones database
    # For now, just log the change
    update_hash "$file"
  fi
}

# Sync API_FREEZE.md
sync_api_freeze() {
  local file="$REPO_PATH/context/API_FREEZE.md"
  if file_changed "$file"; then
    log "API_FREEZE.md changed - syncing"
    update_hash "$file"
  fi
}

# Update MASTER_CONTEXT with latest info
update_master_context() {
  local current_phase=$(grep "Current Phase" "$REPO_PATH/context/SESSION_CONTEXT.md" | sed 's/.*Phase \([0-9.]*\).*/\1/' 2>/dev/null || echo "UNKNOWN")
  local session=$(grep "Session ID" "$REPO_PATH/context/SESSION_CONTEXT.md" | head -1 | sed 's/.*`\(.*\)`/\1/' 2>/dev/null || echo "UNKNOWN")
  local latest_commit=$(cd "$REPO_PATH" && git log --oneline -1 2>/dev/null | cut -d' ' -f1 || echo "UNKNOWN")
  local date=$(date -u +%Y-%m-%d)
  
  # Append update to MASTER_CONTEXT page
  curl -s -X PATCH "https://api.notion.com/v1/blocks/$MASTER_CONTEXT_PAGE/children" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "{
      \"children\": [
        {\"object\": \"block\", \"type\": \"paragraph\", \"paragraph\": {\"rich_text\": [{\"text\": {\"content\": \"🔄 Last Sync: $date | Phase: $current_phase | Session: $session | Commit: $latest_commit\"}}]}}
      ]
    }" > /dev/null 2>&1 || true
}

# Main sync cycle
main() {
  log "Starting sync cycle"
  
  # Ensure state file exists
  touch "$STATE_FILE"
  
  # Sync all monitored files
  sync_project_context
  sync_session_context
  sync_milestones
  sync_api_freeze
  
  # Update master context periodically (every 5th run)
  if [ $(( $(date +%s) / 300 % 5 )) -eq 0 ]; then
    update_master_context
  fi
  
  log "Sync cycle complete"
}

# Run main
main
