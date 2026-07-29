#!/bin/zsh
# scripts/populate_notion_all.sh
export PATH="/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
source .env

echo "🚀 Starting Full Notion Project OS Population..."

# Ensure all scripts are executable
chmod +x scripts/*.sh

# Phase 1: Core Registries
./scripts/populate_architecture.sh
./scripts/populate_packages.sh
./scripts/populate_milestones.sh
./scripts/populate_adrs.sh
./scripts/populate_prompts.sh

# Phase 2: Knowledge Graph (Source Code & APIs)
./scripts/get_notion_maps.sh
./scripts/populate_files.sh
./scripts/populate_apis.sh
./scripts/populate_commits.sh

# Phase 3: Documentation Sync
./scripts/notion_sync.sh

# Phase 4: Dashboard & Metrics
./scripts/populate_dashboard.sh
./scripts/update_dashboard_metrics.sh

echo "🎉 Notion Project OS Population Complete!"
