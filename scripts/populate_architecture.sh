#!/bin/bash
# scripts/populate_architecture.sh
source .env

create_arch() {
  local name=$1
  local status=$2
  local purpose=$3
  local path=$4
  local coverage=$5

  echo "Populating Architecture: $name"
  
  # Check if page exists
  page_id=$(curl -s -X POST "https://api.notion.com/v1/databases/$ARCHITECTURE_DB/query" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Notion-Version: 2022-06-28" \
    -H "Content-Type: application/json" \
    -d "{
      \"filter\": {
        \"property\": \"Subsystem\",
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
    parent="\"parent\": { \"database_id\": \"$ARCHITECTURE_DB\" },"
  fi

  curl -s -X $method "$url" \
    -H "Authorization: Bearer $NOTION_TOKEN" \
    -H "Content-Type: application/json" \
    -H "Notion-Version: 2022-06-28" \
    -d "{
      $parent
      \"properties\": {
        \"Subsystem\": { \"title\": [ { \"text\": { \"content\": \"$name\" } } ] },
        \"Status\": { \"select\": { \"name\": \"$status\" } },
        \"Purpose\": { \"rich_text\": [ { \"text\": { \"content\": \"$purpose\" } } ] },
        \"Repository Path\": { \"rich_text\": [ { \"text\": { \"content\": \"$path\" } } ] },
        \"Coverage\": { \"number\": $coverage }
      }
    }" > /dev/null
}

# Runtime Layer
create_arch "Build Info" "Frozen" "Metadata about the binary build (version, commit, date)." "internal/runtime/buildinfo" 1.0
create_arch "Config" "Frozen" "System configuration and environment management." "internal/runtime/config" 1.0
create_arch "Contracts" "Frozen" "Core runtime interfaces and abstract definitions." "internal/runtime/contracts" 1.0
create_arch "Errors" "Frozen" "Canonical error types and error handling logic." "internal/runtime/errors" 1.0
create_arch "Event Bus" "Frozen" "In-process synchronous event routing engine." "internal/runtime/eventbus" 1.0
create_arch "Lifecycle" "Frozen" "System startup, shutdown, and component health monitoring." "internal/runtime/lifecycle" 1.0
create_arch "Observer" "Frozen" "Passive read-only event observation subsystem." "internal/runtime/observer" 1.0
create_arch "Permissions" "Frozen" "Capability-based access control (CBAC) engine." "internal/runtime/permissions" 1.0
create_arch "Scheduler" "Frozen" "Deterministic task scheduling and cron management." "internal/runtime/scheduler" 1.0
create_arch "Session" "Frozen" "User session state and context isolation." "internal/runtime/session" 1.0
create_arch "Types" "Frozen" "Shared domain types and primitives." "internal/runtime/types" 1.0
create_arch "Version" "Frozen" "Semantic versioning and compatibility checking." "internal/runtime/version" 1.0

# Memory Layer
create_arch "Working Memory" "Frozen" "In-memory sliding window FIFO context storage." "internal/memory/working" 1.0
create_arch "Timeline Ledger" "Frozen" "Immutable append-only chronological event ledger." "internal/memory/timeline" 1.0
create_arch "Memory Store" "Implementation" "Persistence driver layer (SQLite)." "internal/memory/store" 0.0
