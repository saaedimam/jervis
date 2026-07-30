#!/bin/bash
load_config() {
  local config_file="/Users/ioriimasu/dev/jervis/config/notion_databases.json"
  if [ ! -f "$config_file" ]; then
    echo "Error: Configuration file $config_file not found."
    return 1
  fi
  export JERVIS_PAGE=$(jq -r '.parent_page' "$config_file")
  export ARCH_DB=$(jq -r '.architecture' "$config_file")
  export PKG_DB=$(jq -r '.packages' "$config_file")
  export FILE_DB=$(jq -r '.files' "$config_file")
  export SPEC_DB=$(jq -r '.specifications' "$config_file")
  export ADR_DB=$(jq -r '.adrs' "$config_file")
  export API_DB=$(jq -r '.apis' "$config_file")
  export TIMELINE_DB=$(jq -r '.timeline' "$config_file")
  export GATES_DB=$(jq -r '.gates' "$config_file")
  export HANDOFF_DB=$(jq -r '.handoffs' "$config_file")
  export COMMIT_DB=$(jq -r '.commits' "$config_file")
  export MEMORY_DB=$(jq -r '.memory' "$config_file")
  export MILESTONES_DB=$(jq -r '.milestones' "$config_file")
  export DASHBOARD_DB=$(jq -r '.dashboard' "$config_file")
  export DEPS_DB=$(jq -r '.dependencies' "$config_file")
  export GRAPH_METADATA_DB=$(jq -r '.graph_metadata' "$config_file")
}
