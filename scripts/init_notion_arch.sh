#!/bin/bash
# scripts/init_notion_db.sh

ARCH_DB="3ab1b27f-dcba-81c8-89be-c6e885f81f5c"

create_arch() {
  local name=$1
  local status=$2
  local purpose=$3
  local path=$4

  echo "Creating Architecture entry: $name"
  ntn api /v1/pages \
    "parent[database_id]=$ARCH_DB" \
    "properties[Subsystem][title][0][text][content]=$name" \
    "properties[Status][select][name]=$status" \
    "properties[Purpose][rich_text][0][text][content]=$purpose" \
    "properties[Repository Path][rich_text][0][text][content]=$path"
}

create_arch "Build Info" "Frozen" "Metadata about the binary build (version, commit, date)." "internal/runtime/buildinfo"
create_arch "Config" "Frozen" "System configuration and environment management." "internal/runtime/config"
create_arch "Contracts" "Frozen" "Core runtime interfaces and abstract definitions." "internal/runtime/contracts"
create_arch "Errors" "Frozen" "Canonical error types and error handling logic." "internal/runtime/errors"
create_arch "Event Bus" "Frozen" "In-process synchronous event routing engine." "internal/runtime/eventbus"
create_arch "Lifecycle" "Frozen" "System startup, shutdown, and component health monitoring." "internal/runtime/lifecycle"
create_arch "Observer" "Implementation" "Passive read-only event observation subsystem." "internal/runtime/observer"
create_arch "Permissions" "Frozen" "Capability-based access control (CBAC) engine." "internal/runtime/permissions"
create_arch "Scheduler" "Draft" "Deterministic task scheduling and cron management." "internal/runtime/scheduler"
create_arch "Session" "Draft" "User session state and context isolation." "internal/runtime/session"
create_arch "Types" "Frozen" "Shared domain types and primitives." "internal/runtime/types"
create_arch "Version" "Frozen" "Semantic versioning and compatibility checking." "internal/runtime/version"
