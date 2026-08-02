#!/bin/bash
# Event-driven synchronization script for Jervis
# This is an improved version that watches Git events

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
LOCK_FILE="$REPO_PATH/.jervis_sync.lock"
LOG_FILE="$REPO_PATH/.jervis_sync.log"

# Source centralized configuration
source "$REPO_PATH/scripts/load_config.sh"
load_config

log() {
  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] $1" >> "$LOG_FILE"
}

# Prevent concurrent runs
if [ -f "$LOCK_FILE" ]; then
  pid=$(cat "$LOCK_FILE" 2>/dev/null)
  if kill -0 "$pid" 2>/dev/null; then
    log "Sync already running (PID: $pid)"
    exit 0
  else
    rm -f "$LOCK_FILE"
  fi
fi
echo $$ > "$LOCK_FILE"
trap 'rm -f "$LOCK_FILE"' EXIT

# Parse Git diff for changed files
detect_changes() {
  cd "$REPO_PATH"

  # Get latest commit info
  local latest_commit=$(git log --oneline -1 | cut -d' ' -f1)
  local commit_msg=$(git log --format=%B -n 1 HEAD | head -1)
  local author=$(git log --format=%an -n 1 HEAD)
  local date=$(git log --format=%ai -n 1 HEAD | cut -d' ' -f1)
  local branch=$(git branch --show-current)

  log "Detected commit: $latest_commit"
  log "Branch: $branch"

  # Sync to Commit Intelligence (Commit Registry) using the correct, schema-validated fields
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    -d "{
      \"parent\": {\"database_id\": \"$COMMIT_DB\"},
      \"properties\": {
        \"Commit ID\": {\"title\": [{\"text\": {\"content\": \"COMMIT-$latest_commit\"}}]},
        \"Message\": {\"rich_text\": [{\"text\": {\"content\": \"$commit_msg\"}}]},
        \"Author\": {\"rich_text\": [{\"text\": {\"content\": \"$author\"}}]},
        \"Date\": {\"date\": {\"start\": \"$date\"}}
      }
    }" >> "$LOG_FILE" 2>&1 || true
}

# Sync architecture status updates
sync_architecture() {
  # Read current phase from SESSION_CONTEXT
  local current_phase=$(grep "Current Phase" "$REPO_PATH/context/SESSION_CONTEXT.md" 2>/dev/null | sed 's/.*Phase \([0-9.]*\).*/\1/' || echo "UNKNOWN")

  log "Current phase: $current_phase"

  # Update architecture status based on completion
  # This would query Milestones and update Architecture Registry
}

# Sync package coverage
sync_packages() {
  # Parse coverage.out if exists
  if [ -f "$REPO_PATH/coverage.out" ]; then
    local coverage=$(grep "^mode:" "$REPO_PATH/coverage.out" | head -1 || echo "N/A")
    log "Coverage data available"
  fi
}

# Main sync cycle
main() {
  log "=== Event-driven sync starting ==="

  # Detect changes from Git
  detect_changes

  # Sync dependent data
  sync_architecture
  sync_packages

  log "=== Sync complete ==="
}

main