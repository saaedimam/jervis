#!/bin/bash
# Generic script to push YAML data to Notion databases

set -e

REPO_PATH="/Users/ioriimasu/dev/jervis"
source "$REPO_PATH/scripts/load_config.sh"
load_config

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <yaml_file> <database_id_key>"
  exit 1
fi

YAML_FILE="$1"
DB_KEY="$2"

if [ ! -f "$YAML_FILE" ]; then
  echo "Error: $YAML_FILE not found"
  exit 1
fi

# Find DB ID from config based on DB_KEY
DB_ID=$(jq -r ".$DB_KEY" "$REPO_PATH/config/notion_databases.json")
if [ -z "$DB_ID" ] || [ "$DB_ID" == "null" ]; then
  echo "Database key $DB_KEY not found in config"
  exit 1
fi

python3 "$REPO_PATH/scripts/push_yaml_to_notion.py" "$YAML_FILE" "$DB_ID"
echo "Successfully pushed $YAML_FILE to Notion database $DB_KEY"
