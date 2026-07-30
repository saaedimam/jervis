#!/bin/bash
# Event-driven synchronization script for Jervis
# This is an improved version that watches Git events

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
LOCK_FILE="$REPO_PATH/.jervis_sync.lock"
LOG_FILE="$REPO_PATH/.jervis_sync.log"

JERVIS_PAGE="3ab1b27f-dcba-81d0-8b35-ed766e2e8420"
ARCH_DB="d3dcb133-f96e-4e8e-944f-5825c2d1eee0"
PKG_DB="9c8bb7d5-5675-4cc5-b1b7-6a9c1ac3fe2f"
SPEC_DB="f30e0d51-a787-421a-ad6b-77935f7d2e53"
HANDOFF_DB="c1e36ebb-a3fc-4aea-a3d2-ac8214e1e40a"
MEMORY_DB="38a76b5b-b20e-498e-b6e9-e643c2ae7d8b"
COMMIT_DB="69c5145a-b84c-43e5-83b2-05d746a80e26"
DEPS_DB="1de04b92-6fe3-4756-b85d-c9370f838a3b"

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
echo $$ > "$LOCK_FILE
trap 'rm -f "$LOCK_FILE"' EXIT

# Parse Git diff for changed files
detect_changes() {
  cd "$REPO_PATH"
  
  # Get latest commit info
  local latest_commit=$(git log --oneline -1 | cut -d' ' -f1)
  local commit_msg=$(git log --format=%B -n 1 HEAD | head -1)
  local author=$(git log --format=%an -n 1 HEAD)
  local date=$(git log --format=%ai -n 1 HEAD | cut -d' ' -f1)
  local files_changed=$(git diff-tree --no-commit-id --name-only -r HEAD | wc -l | tr -d ' ')
  
  # Detect which packages were touched
  local packages=$(git diff-tree --no-commit-id --name-only -r HEAD | grep -E '\.go$' | sed 's|/[^/]*$||' | sort -u | tr '\n' ',' | sed 's/,$//')
  
  # Detect architecture impact
  local arch_impact=""
  if echo "$packages" | grep -q "eventbus"; then
    arch_impact="ARCH-002 Event Bus"
  elif echo "$packages" | grep -q "permissions"; then
    arch_impact="ARCH-003 Permission Engine"
  elif echo "$packages" | grep -q "observer"; then
    arch_impact="ARCH-004 Observer"
  fi
  
  # Check if breaking change
  local breaking="false"
  if echo "$commit_msg" | grep -qi "breaking\|BREAKING"; then
    breaking="true"
  fi
  
  log "Detected commit: $latest_commit"
  log "Files changed: $files_changed"
  log "Packages: $packages"
  log "Architecture: $arch_impact"
  
  # Sync to Commit Intelligence
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "{
      \"parent\": {\"database_id\": \"$COMMIT_DB\"},
      \"properties\": {
        \"Commit ID\": {\"title\": [{\"text\": {\"content\": \"COMMIT-$latest_commit\"}}]},
        \"Hash\": {\"rich_text\": [{\"text\": {\"content\": \"$latest_commit\"}}]},
        \"Message\": {\"rich_text\": [{\"text\": {\"content\": \"$commit_msg\"}}]},
        \"Author\": {\"rich_text\": [{\"text\": {\"content\": \"$author\"}}]},
        \"Date\": {\"date\": {\"start\": \"$date\"}},
        \"Files Changed\": {\"number\": $files_changed},
        \"Packages\": {\"rich_text\": [{\"text\": {\"content\": \"$packages\"}}]},
        \"Architecture Impact\": {\"rich_text\": [{\"text\": {\"content\": \"$arch_impact\"}}]},
        \"Breaking Change\": {\"checkbox\": $breaking}
      }
    }" > /dev/null 2>&1 || true
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
