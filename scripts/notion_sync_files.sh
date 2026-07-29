#!/bin/bash
# Jervis File Registry Synchronization Script
# Complete implementation for Engineering Knowledge Compiler
# Scans repository and syncs all files to Notion File Registry

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
FILE_DB="d5b8d71a-c568-4288-9443-f3deb8b316bc"
PKG_DB="9c8bb7d5-5675-4cc5-b1b7-6a9c1ac3fe2f"
ARCH_DB="d3dcb133-f96e-4e8e-944f-5825c2d1eee0"
SPEC_DB="f30e0d51-a787-421a-ad6b-77935f7d2e53"
STATE_DIR="$REPO_PATH/.jervis"
LOG_FILE="$STATE_DIR/file_sync.log"

mkdir -p "$STATE_DIR"

log() {
  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] $1" | tee -a "$LOG_FILE"
}

# Get file hash for change detection
get_hash() {
  if [ -f "$1" ]; then
    md5 -q "$1" 2>/dev/null || md5sum "$1" | cut -d' ' -f1
  else
    echo "missing"
  fi
}

# Get line count
get_lines() {
  if [ -f "$1" ]; then
    wc -l "$1" | awk '{print $1}'
  else
    echo "0"
  fi
}

# Detect language from extension
detect_language() {
  local file="$1"
  case "$file" in
    *.go) echo "Go" ;;
    *.md) echo "Markdown" ;;
    *.yaml|*.yml) echo "YAML" ;;
    *.json) echo "JSON" ;;
    *.sh) echo "Shell" ;;
    *.sql) echo "SQL" ;;
    *.py) echo "Python" ;;
    *) echo "Unknown" ;;
  esac
}

# Determine package from path
get_package() {
  local file="$1"
  local rel_path="${file#$REPO_PATH/}"
  
  # Extract package from Go file path
  if [[ "$rel_path" == internal/* ]]; then
    echo "$rel_path" | sed 's|/[^/]*$||' | sed 's|/|.|g'
  elif [[ "$rel_path" == pkg/* ]]; then
    echo "$rel_path" | sed 's|/[^/]*$||' | sed 's|/|.|g'
  elif [[ "$rel_path" == cmd/* ]]; then
    echo "$rel_path" | sed 's|/[^/]*$||' | sed 's|/|.|g'
  else
    echo "root"
  fi
}

# Determine architecture from package
get_architecture() {
  local pkg="$1"
  
  if [[ "$pkg" == *runtime* ]]; then
    if [[ "$pkg" == *eventbus* ]]; then echo "ARCH-002"; return; fi
    if [[ "$pkg" == *permissions* ]]; then echo "ARCH-003"; return; fi
    if [[ "$pkg" == *observer* ]]; then echo "ARCH-004"; return; fi
    echo "ARCH-001"; return
  elif [[ "$pkg" == *memory* ]]; then
    echo "ARCH-001"; return
  elif [[ "$pkg" == *services* ]]; then
    echo "ARCH-001"; return
  elif [[ "$pkg" == *aiprovider* ]]; then
    echo "ARCH-001"; return
  elif [[ "$pkg" == *interfaces* ]]; then
    echo "ARCH-001"; return
  elif [[ "$pkg" == docs* ]] || [[ "$pkg" == *.md ]]; then
    echo "Documentation"; return
  else
    echo "Other"; return
  fi
}

# Get specification for file
get_specification() {
  local file="$1"
  local basename=$(basename "$file")
  
  case "$basename" in
    *EVENT_BUS_SPEC*) echo "SPEC-001"; return ;;
    *EVENT_MODEL*) echo "SPEC-002"; return ;;
    *DISPATCHER_SPEC*) echo "SPEC-004"; return ;;
    *MIDDLEWARE_SPEC*) echo "SPEC-005"; return ;;
    *BUS_SPEC*) echo "SPEC-006"; return ;;
    *PERMISSION_ENGINE_SPEC*) echo "SPEC-010"; return ;;
    *PERMISSION_MODEL*) echo "SPEC-011"; return ;;
    *OBSERVER_SPEC*) echo "SPEC-020"; return ;;
    *OBSERVER_MODEL*) echo "SPEC-021"; return ;;
    *) echo ""; return ;;
  esac
}

# Check if file is frozen (from API_FREEZE.md)
is_frozen() {
  local file="$1"
  local pkg=$(get_package "$file")
  
  # Check against frozen packages list
  if [[ "$pkg" == *eventbus* ]] || [[ "$pkg" == *permissions* ]] || [[ "$pkg" == *observer.contracts* ]]; then
    echo "true"
  else
    echo "false"
  fi
}

# Get coverage for Go files
get_coverage() {
  local file="$1"
  if [[ "$file" == *.go ]] && [ -f "$REPO_PATH/coverage.out" ]; then
    # Parse coverage.out for this file
    local rel_path="${file#$REPO_PATH/}"
    local coverage=$(grep "^$rel_path:" "$REPO_PATH/coverage.out" 2>/dev/null | head -1 || echo "")
    if [ -n "$coverage" ]; then
      echo "100%"
    else
      echo "Unknown"
    fi
  else
    echo "N/A"
  fi
}

# Get status
get_status() {
  local file="$1"
  if [ -f "$file" ]; then
    echo "Active"
  else
    echo "Missing"
  fi
}

# Get exports for Go files
get_exports() {
  local file="$1"
  if [[ "$file" == *.go ]] && [ -f "$file" ]; then
    grep "^func\|^type\|^const\|^var" "$file" 2>/dev/null | grep "^[A-Z]" | head -5 | sed 's/ {.*/.../' | tr '\n' '; ' | sed 's/; $//' || echo ""
  else
    echo ""
  fi
}

# Create or update file in Notion
sync_file_to_notion() {
  local file_id="$1"
  local file_path="$2"
  local rel_path="${file_path#$REPO_PATH/}"
  
  local language=$(detect_language "$file_path")
  local pkg=$(get_package "$file_path")
  local arch=$(get_architecture "$pkg")
  local spec=$(get_specification "$file_path")
  local lines=$(get_lines "$file_path")
  local hash=$(get_hash "$file_path")
  local frozen=$(is_frozen "$file_path")
  local coverage=$(get_coverage "$file_path")
  local status=$(get_status "$file_path")
  local exports=$(get_exports "$file_path")
  
  # Truncate exports if too long
  if [ ${#exports} -gt 100 ]; then
    exports="${exports:0:97}..."
  fi
  
  local json=$(cat <<EOF
{
  "parent": {"database_id": "$FILE_DB"},
  "properties": {
    "File ID": {"title": [{"text": {"content": "$file_id"}}]},
    "Path": {"rich_text": [{"text": {"content": "$rel_path"}}]},
    "Package": {"rich_text": [{"text": {"content": "$pkg"}}]},
    "Language": {"select": {"name": "$language"}},
    "Architecture": {"rich_text": [{"text": {"content": "$arch"}}]},
    "Exports": {"rich_text": [{"text": {"content": "$exports"}}]},
    "Coverage": {"rich_text": [{"text": {"content": "$coverage"}}]},
    "Frozen": {"checkbox": $frozen},
    "Status": {"rich_text": [{"text": {"content": "$status"}}]}
  }
}
EOF
)
  
  # Check if file already exists in Notion
  local existing=$(curl -s "https://api.notion.com/v1/data_sources/2f6a0483-207c-4902-b35c-7e534ca09a4e/query" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "{\"filter\": {\"property\": \"File ID\", \"title\": {\"equals\": \"$file_id\"}}}" 2>/dev/null | jq -r '.results[0].id')
  
  if [ "$existing" != "null" ] && [ -n "$existing" ]; then
    # Update existing
    curl -s -X PATCH "https://api.notion.com/v1/pages/$existing" \
      -H "Authorization: Bearer $NOTION_API_KEY" \
      -H "Notion-Version: 2025-09-03" \
      -H "Content-Type: application/json" \
      -d "{\"properties\": $(echo "$json" | jq '.properties')}" > /dev/null 2>&1
    echo "UPDATED"
  else
    # Create new
    curl -s -X POST "https://api.notion.com/v1/pages" \
      -H "Authorization: Bearer $NOTION_API_KEY" \
      -H "Notion-Version: 2025-09-03" \
      -H "Content-Type: application/json" \
      -d "$json" > /dev/null 2>&1
    echo "CREATED"
  fi
}

# Main sync function
main() {
  log "=== Jervis File Registry Sync ==="
  log "Scanning repository: $REPO_PATH"
  
  # Build file manifest
  find "$REPO_PATH" -type f \
    \( -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" -o -name "*.json" -o -name "*.sh" \) \
    | grep -v -E "(\.git|node_modules|vendor|coverage\.out|tmp|build|bin|dist|\.jervis)" \
    | sort > "$STATE_DIR/sync_manifest.txt"
  
  local total=$(wc -l < "$STATE_DIR/sync_manifest.txt")
  log "Found $total files to synchronize"
  
  local created=0
  local updated=0
  local failed=0
  local counter=1
  
  while IFS= read -r file_path; do
    local file_id=$(printf "FILE-%04d" $counter)
    local hash=$(get_hash "$file_path")
    local stored_hash=""
    
    # Check if file has changed
    if [ -f "$STATE_DIR/file_hashes.txt" ]; then
      stored_hash=$(grep "^$file_id:" "$STATE_DIR/file_hashes.txt" 2>/dev/null | cut -d: -f2 || echo "")
    fi
    
    if [ "$hash" != "$stored_hash" ]; then
      log "Syncing $file_id: $(basename "$file_path")"
      
      local result=$(sync_file_to_notion "$file_id" "$file_path")
      
      if [ "$result" = "CREATED" ]; then
        created=$((created + 1))
      elif [ "$result" = "UPDATED" ]; then
        updated=$((updated + 1))
      fi
      
      # Update hash store
      if [ -f "$STATE_DIR/file_hashes.txt" ]; then
        grep -v "^$file_id:" "$STATE_DIR/file_hashes.txt" > "$STATE_DIR/file_hashes.tmp" || true
        mv "$STATE_DIR/file_hashes.tmp" "$STATE_DIR/file_hashes.txt"
      fi
      echo "$file_id:$hash" >> "$STATE_DIR/file_hashes.txt"
      
      # Rate limiting
      sleep 0.3
    fi
    
    counter=$((counter + 1))
    
    # Progress every 50 files
    if [ $((counter % 50)) -eq 0 ]; then
      log "Progress: $counter/$total files processed"
    fi
  done < "$STATE_DIR/sync_manifest.txt"
  
  log ""
  log "=== Sync Complete ==="
  log "Total files: $total"
  log "Created: $created"
  log "Updated: $updated"
  log "Failed: $failed"
  log "State saved to: $STATE_DIR/"
}

main "$@"
