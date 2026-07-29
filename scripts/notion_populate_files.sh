#!/bin/bash
# Jervis File Registry Population Script v2
# Uses centralized config, Go AST parsing, declarative mappings, and incremental hashing.
# Outputs intermediate graph.json.

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
STATE_DIR="$REPO_PATH/.jervis"
LOG_FILE="$STATE_DIR/file_populate.log"

mkdir -p "$STATE_DIR"

log() {
  echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] $1" | tee -a "$LOG_FILE"
}

source "$REPO_PATH/scripts/load_config.sh"
load_config

if [ -z "$FILE_DB" ] || [ "$FILE_DB" == "null" ]; then
  log "Error: File Registry not found in config."
  exit 1
fi

# Build Go extractor if needed
if [ ! -f "$REPO_PATH/bin/go_extractor" ]; then
  log "Building Go extractor..."
  mkdir -p "$REPO_PATH/bin"
  go build -o "$REPO_PATH/bin/go_extractor" "$REPO_PATH/scripts/extract_go_metadata.go"
fi

# Initialize files manifest
find "$REPO_PATH" -type f \
  \( -name "*.go" -o -name "*.md" -o -name "*.yaml" -o -name "*.yml" -o -name "*.json" -o -name "*.sh" \) \
  | grep -v -E "(\.git|node_modules|vendor|coverage\.out|tmp|build|bin|dist|\.jervis)" \
  | sort > "$STATE_DIR/sync_manifest.txt"

total=$(wc -l < "$STATE_DIR/sync_manifest.txt" | tr -d ' ')
log "Found $total files to process"

created=0
updated=0
unchanged=0
counter=1

# Initialize graph.json array
GRAPH_OUT="$REPO_PATH/graph.json"
echo "[" > "$GRAPH_OUT"
first_item=true

# Read last commits file
COMMIT_LOOKUP="$STATE_DIR/commit_lookup.json"
if [ ! -f "$COMMIT_LOOKUP" ]; then echo "{}" > "$COMMIT_LOOKUP"; fi

while IFS= read -r file_path; do
  rel_path="${file_path#$REPO_PATH/}"
  file_id=$(printf "FILE-%04d" $counter)
  
  # Determine language
  case "$rel_path" in
    *.go) lang="Go" ;;
    *.md) lang="Markdown" ;;
    *.yaml|*.yml) lang="YAML" ;;
    *.json) lang="JSON" ;;
    *.sh) lang="Shell" ;;
    *) lang="Unknown" ;;
  esac

  # Hash content
  if command -v shasum &> /dev/null; then
    hash=$(shasum -a 256 "$file_path" | cut -d' ' -f1)
  else
    hash=$(sha256sum "$file_path" | cut -d' ' -f1)
  fi

  # Check incremental hash
  stored_hash=""
  if [ -f "$STATE_DIR/file_hashes.txt" ]; then
    stored_hash=$(grep "^$file_id:" "$STATE_DIR/file_hashes.txt" 2>/dev/null | cut -d: -f2 || echo "")
  fi

  # Determine Package Path
  pkg_path="root"
  if [[ "$rel_path" == internal/* ]] || [[ "$rel_path" == pkg/* ]] || [[ "$rel_path" == cmd/* ]]; then
    pkg_path=$(dirname "$rel_path")
  fi

  # AST parsing for Go
  exports=""
  imports=""
  api_count=0
  if [ "$lang" == "Go" ]; then
    metadata=$("$REPO_PATH/bin/go_extractor" "$file_path")
    exports_arr=$(echo "$metadata" | jq -r '.exports | join(";")')
    imports_arr=$(echo "$metadata" | jq -r '.imports | join(";")')
    api_count=$(echo "$metadata" | jq '.exports | length')
    
    if [ ${#exports_arr} -gt 1500 ]; then exports="${exports_arr:0:1497}..."; else exports="$exports_arr"; fi
    if [ ${#imports_arr} -gt 1500 ]; then imports="${imports_arr:0:1497}..."; else imports="$imports_arr"; fi
  fi

  # Mapping Architecture & Spec
  map_res=$(python3 "$REPO_PATH/scripts/parse_yaml_mapping.py" "$REPO_PATH/config/spec_mapping.yaml" "$pkg_path")
  arch_id=$(echo "$map_res" | head -n 1)
  spec_id=$(echo "$map_res" | tail -n 1)

  # Owner and commit
  last_author=$(git log -1 --format='%an' "$file_path")
  last_commit=$(git log -1 --format='%H' "$file_path")

  # Generate Notion JSON
  # Resolve relations if possible... For now, since Notion relations require Notion UUIDs,
  # we leave them empty or resolve them in Notion Sync. Wait, the instructions say 
  # "notion_populate_files.sh Responsibilities: Create/update File Registry records".
  # To link we need the Notion ID of the Package, Architecture, Spec.
  
  # Fetch Relation IDs from Notion
  # Package
  notion_pkg_id=$(curl -s "https://api.notion.com/v1/databases/$PKG_DB/query" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "{\"filter\": {\"property\": \"Package ID\", \"title\": {\"equals\": \"$pkg_path\"}}}" | jq -r '.results[0].id // empty')
  
  # Architecture
  notion_arch_id=$(curl -s "https://api.notion.com/v1/databases/$ARCH_DB/query" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "{\"filter\": {\"property\": \"Architecture ID\", \"title\": {\"equals\": \"$arch_id\"}}}" | jq -r '.results[0].id // empty')

  # Spec
  notion_spec_id=""
  if [ -n "$spec_id" ]; then
    notion_spec_id=$(curl -s "https://api.notion.com/v1/databases/$SPEC_DB/query" \
      -H "Authorization: Bearer $NOTION_API_KEY" \
      -H "Notion-Version: 2025-09-03" \
      -H "Content-Type: application/json" \
      -d "{\"filter\": {\"property\": \"Spec ID\", \"title\": {\"equals\": \"$spec_id\"}}}" | jq -r '.results[0].id // empty')
  fi

  # Notion payload
  props="{
    \"File ID\": {\"title\": [{\"text\": {\"content\": \"$file_id\"}}]},
    \"Path\": {\"rich_text\": [{\"text\": {\"content\": \"$rel_path\"}}]},
    \"Language\": {\"select\": {\"name\": \"$lang\"}},
    \"Exports\": {\"rich_text\": [{\"text\": {\"content\": \"$exports\"}}]},
    \"Imports\": {\"rich_text\": [{\"text\": {\"content\": \"$imports\"}}]},
    \"Owner\": {\"rich_text\": [{\"text\": {\"content\": \"$last_author\"}}]},
    \"API Count\": {\"number\": $api_count}
  }"

  if [ -n "$notion_pkg_id" ]; then
    props=$(echo "$props" | jq --arg id "$notion_pkg_id" '. + {"Package": {"relation": [{"id": $id}]}}')
  fi
  if [ -n "$notion_arch_id" ]; then
    props=$(echo "$props" | jq --arg id "$notion_arch_id" '. + {"Architecture": {"relation": [{"id": $id}]}}')
  fi
  if [ -n "$notion_spec_id" ]; then
    props=$(echo "$props" | jq --arg id "$notion_spec_id" '. + {"Specification": {"relation": [{"id": $id}]}}')
  fi

  json="{
    \"parent\": {\"database_id\": \"$FILE_DB\"},
    \"properties\": $props
  }"

  if [ "$hash" != "$stored_hash" ]; then
    log "Syncing $file_id: $rel_path"

    existing=$(curl -s "https://api.notion.com/v1/databases/$FILE_DB/query" \
      -H "Authorization: Bearer $NOTION_API_KEY" \
      -H "Notion-Version: 2025-09-03" \
      -H "Content-Type: application/json" \
      -d "{\"filter\": {\"property\": \"File ID\", \"title\": {\"equals\": \"$file_id\"}}}" | jq -r '.results[0].id // empty')

    if [ -n "$existing" ]; then
      curl -s -X PATCH "https://api.notion.com/v1/pages/$existing" \
        -H "Authorization: Bearer $NOTION_API_KEY" \
        -H "Notion-Version: 2025-09-03" \
        -H "Content-Type: application/json" \
        -d "{\"properties\": $props}" > /dev/null
      updated=$((updated + 1))
    else
      curl -s -X POST "https://api.notion.com/v1/pages" \
        -H "Authorization: Bearer $NOTION_API_KEY" \
        -H "Notion-Version: 2025-09-03" \
        -H "Content-Type: application/json" \
        -d "$json" > /dev/null
      created=$((created + 1))
    fi

    # Update hash
    if [ -f "$STATE_DIR/file_hashes.txt" ]; then
      grep -v "^$file_id:" "$STATE_DIR/file_hashes.txt" > "$STATE_DIR/file_hashes.tmp" || true
      mv "$STATE_DIR/file_hashes.tmp" "$STATE_DIR/file_hashes.txt"
    fi
    echo "$file_id:$hash" >> "$STATE_DIR/file_hashes.txt"
  else
    unchanged=$((unchanged + 1))
  fi

  # Append to graph.json
  graph_entry=$(jq -n \
    --arg id "$file_id" \
    --arg path "$rel_path" \
    --arg pkg "$pkg_path" \
    --arg arch "$arch_id" \
    --arg spec "$spec_id" \
    --arg lang "$lang" \
    --arg exports "$exports" \
    --arg imports "$imports" \
    '{id: $id, path: $path, package: $pkg, architecture: $arch, specification: $spec, language: $lang, exports: $exports, imports: $imports}')
    
  if [ "$first_item" = true ]; then
    echo "$graph_entry" >> "$GRAPH_OUT"
    first_item=false
  else
    echo ",$graph_entry" >> "$GRAPH_OUT"
  fi

  counter=$((counter + 1))
done < "$STATE_DIR/sync_manifest.txt"

echo "]" >> "$GRAPH_OUT"

log "=== File Registry Population Complete ==="
log "Total: $total"
log "Created: $created"
log "Updated: $updated"
log "Unchanged: $unchanged"
