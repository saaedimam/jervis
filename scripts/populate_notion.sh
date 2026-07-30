#!/bin/bash
# scripts/populate_notion.sh

# Load environment variables
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

if [ -z "$NOTION_TOKEN" ]; then
    echo "Error: NOTION_TOKEN is not set."
    exit 1
fi

# Sync static pages
echo "Syncing Dashboard..."
ntn pages edit "$DASHBOARD_ID" --content "$(cat context/DASHBOARD.md)"

echo "Syncing Agent Instructions..."
ntn pages edit "$AGENT_INSTRUCTIONS_ID" --content "$(cat context/AGENT_INSTRUCTIONS.md)"

echo "Syncing MASTER_CONTEXT..."
ntn pages edit "$MASTER_CONTEXT_ID" --content "$(cat context/MASTER_CONTEXT.md)"

# Populate Context Database
echo "Populating Context Database..."
create_context_entry() {
    local title=$1
    local file=$2
    echo "Creating entry for $title..."
    ntn api /v1/pages \
        "parent[database_id]=$CONTEXT_DB" \
        "properties[Name][title][0][text][content]=$title" \
        "properties[Type][select][name]=Documentation"
}

# Add some initial context entries (ids from .env)
# Note: This is a simplified version. A full sync would check for existing entries.
create_context_entry "PROJECT_CONTEXT" "context/PROJECT_CONTEXT.md"
create_context_entry "API_FREEZE" "context/API_FREEZE.md"
create_context_entry "MILESTONES" "context/MILESTONES.md"

echo "Done!"
