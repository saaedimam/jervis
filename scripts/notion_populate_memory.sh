#!/bin/bash
# Populate Engineering Memory with lessons learned

MEMORY_DB="38a76b5b-b20e-498e-b6e9-e643c2ae7d8b"

add_memory() {
  local type="$1"
  local title="$2"
  local content="$3"
  local pkg="$4"
  local adr="$5"
  local date="$6"
  
  local json=$(cat <<EOF
{
  "parent": {"database_id": "$MEMORY_DB"},
  "properties": {
    "Type": {"select": {"name": "$type"}},
    "Title": {"rich_text": [{"text": {"content": "$title"}}]},
    "Content": {"rich_text": [{"text": {"content": "$content"}}]},
    "Related Package": {"rich_text": [{"text": {"content": "$pkg"}}]},
    "Related ADR": {"rich_text": [{"text": {"content": "$adr"}}]},
    "Date": {"date": {"start": "$date"}}
  }
}
EOF
)
  
  curl -s -X POST "https://api.notion.com/v1/pages" \
    -H "Authorization: Bearer $NOTION_API_KEY" \
    -H "Notion-Version: 2025-09-03" \
    -H "Content-Type: application/json" \
    -d "$json" | jq -r '.id'
}

# Lessons Learned
add_memory "Lesson Learned" "Runtime Ownership First" "AI providers must never own the runtime. The system must work deterministically without any AI. This prevents vendor lock-in and ensures predictable behavior." "internal/runtime/*" "ADR-0002" "2026-07-28"
add_memory "Lesson Learned" "Pure Synchronous Phase 1" "Phase 1 primitives must use pure synchronous value semantics without channels or goroutines. This eliminates entire categories of concurrency bugs during foundation building." "internal/runtime/*" "" "2026-07-28"
add_memory "Lesson Learned" "100% Coverage Gate" "Require 100% test coverage on runtime packages. Use race detector in CI. This catches subtle concurrency issues early." "internal/runtime/*" "" "2026-07-28"

# Patterns
add_memory "Pattern" "Value Object Pattern" "Use immutable value objects with New() constructors, validation, and defensive copies. Zero values represent uninitialized state." "internal/runtime/permissions/*" "" "2026-07-28"
add_memory "Pattern" "AggregateError Pattern" "Use AggregateError for Continue-on-Error policies. Panic isolation per handler with recovery wrapping." "internal/runtime/eventbus/*" "" "2026-07-28"
add_memory "Pattern" "Wildcard Matching" "Use pure Go pattern matching (*, prefix.*, prefix*) instead of regex for event subscription matching. Deterministic and fast." "internal/runtime/eventbus/registry" "" "2026-07-28"

# Design Rules
add_memory "Design Rule" "Observer Never Calls AI" "The Observer subsystem must never call AI providers. It is a passive read-only event observation layer." "internal/runtime/observer/*" "" "2026-07-28"
add_memory "Design Rule" "Event Bus Never Calls AI" "The Event Bus must never call AI providers. It routes events between runtime components only." "internal/runtime/eventbus/*" "" "2026-07-28"
add_memory "Design Rule" "Memory Never Depends on AI" "The Memory Engine must never depend on AI providers. It manages state history independently." "internal/memory/*" "" "2026-07-28"
add_memory "Design Rule" "Services Never Depend on Provider" "Services must work with any AI provider through abstraction. No vendor-specific dependencies." "internal/services/*" "" "2026-07-28"

# Performance Notes
add_memory "Performance Note" "Synchronous Event Dispatch" "Synchronous event dispatch eliminates goroutine overhead for Phase 1. Future phases may add async variants." "internal/runtime/eventbus/dispatcher" "" "2026-07-28"
add_memory "Performance Note" "Defensive Copies" "All slice/array getters return defensive copies. Trade memory for safety in foundation packages." "internal/runtime/*" "" "2026-07-28"

echo "Engineering memory populated"
