#!/bin/bash
# scripts/init_notion_milestones.sh

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

MILESTONES_DB="3ab1b27f-dcba-8152-bd83-e18e3848ed92"

create_milestone() {
  local phase=$1
  local status=$2
  local name=$3

  echo "Creating Milestone: $phase - $name"
  ntn api /v1/pages \
    "parent[database_id]=$MILESTONES_DB" \
    "properties[Phase][title][0][text][content]=$phase" \
    "properties[Status][select][name]=$status" \
    "properties[Name][rich_text][0][text][content]=$name"
}

create_milestone "Phase 1.1" "Complete" "Core Runtime Foundation"
create_milestone "Phase 1.2" "Complete" "Event Bus Subsystem"
create_milestone "Phase 1.3" "Complete" "Permission Engine"
create_milestone "Phase 1.4" "Complete" "Observer Subsystem"
create_milestone "Phase 1.4.1" "Pending" "Observer Foundation & Contracts"
create_milestone "Phase 1.5" "Complete" "Scheduler Component"
create_milestone "Phase 1.6" "Complete" "Session Management Engine"
create_milestone "Phase 2.1" "Complete" "Working Memory & Timeline Ledger"
create_milestone "Phase 2.2" "Complete" "Knowledge Store Driver"
create_milestone "Phase 3.1" "Complete" "Planner Service"
create_milestone "Phase 3.2" "Complete" "Projects Service"
create_milestone "Phase 3.3" "Complete" "Habits Service"
create_milestone "Phase 3.4" "Complete" "Meetings Service"
create_milestone "Phase 3.5" "Implementation" "Notion Integration Service"
create_milestone "Phase 3.6" "Implementation" "Calendar Integration Service"
