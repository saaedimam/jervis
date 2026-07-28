#!/bin/bash
# scripts/setup_github.sh
# Declarative, idempotent setup for the jervis GitHub repository.
# Usage: ./scripts/setup_github.sh [--plan]

set -e

SPEC_FILE=".github/governance.yaml"
SCHEMA_FILE=".github/governance.schema.json"
DRY_RUN=false

if [[ "$1" == "--plan" ]]; then
    DRY_RUN=true
    echo "🏗️  PLANNING MODE: No changes will be applied to GitHub."
fi

print_step() { echo -e "\033[1;34m[STEP]\033[0m $1"; }
print_ok() { echo -e "\033[1;32m[OK]\033[0m $1"; }
print_warn() { echo -e "\033[1;33m[WARN]\033[0m $1"; }
print_plan() { echo -e "\033[1;35m[PLAN]\033[0m $1"; }

# 0. Prerequisite Check & Validation
if ! command -v yq &> /dev/null; then echo "yq not found"; exit 2; fi
if ! command -v jq &> /dev/null; then echo "jq not found"; exit 2; fi

print_step "Validating Governance Specification..."
# Basic structural validation using yq
API_VERSION=$(yq '.apiVersion' "$SPEC_FILE")
if [[ "$API_VERSION" != "governance.jervis.io/v1" ]]; then
    echo "Error: Unsupported apiVersion $API_VERSION"; exit 1
fi

# 1. Load Spec
REPO_NAME=$(yq '.repository.name' "$SPEC_FILE")
REPO_OWNER=$(yq '.repository.owner' "$SPEC_FILE")
REPO="$REPO_OWNER/$REPO_NAME"
DESCRIPTION=$(yq '.repository.description' "$SPEC_FILE")
VISIBILITY=$(yq '.repository.visibility' "$SPEC_FILE")

# 2. Idempotent Repository Creation
print_step "Checking GitHub repository status for $REPO..."
if ! gh repo view "$REPO" &>/dev/null; then
    if [[ "$DRY_RUN" == "true" ]]; then
        print_plan "Create repository $REPO ($VISIBILITY)"
    else
        print_step "Creating repository $REPO..."
        gh repo create "$REPO" --"$VISIBILITY" --description "$DESCRIPTION"
    fi
else
    print_ok "Repository $REPO already exists."
fi

# 3. Idempotent Remote Configuration
if [[ "$DRY_RUN" == "false" ]]; then
    print_step "Checking local git remote..."
    if ! git remote get-url origin &>/dev/null; then
        print_step "Adding origin remote..."
        git remote add origin "https://github.com/$REPO.git"
    else
        print_ok "Remote 'origin' already configured."
    fi
fi

# 4. Sync Repository Features & Settings
print_step "Synchronizing Repository Settings..."
DELETE_BRANCH=$(yq '.repository.settings.delete_branch_on_merge' "$SPEC_FILE")
ALLOW_SQUASH=$(yq '.repository.settings.allow_squash_merge' "$SPEC_FILE")
ALLOW_MERGE=$(yq '.repository.settings.allow_merge_commit' "$SPEC_FILE")
ALLOW_REBASE=$(yq '.repository.settings.allow_rebase_merge' "$SPEC_FILE")

if [[ "$DRY_RUN" == "true" ]]; then
    print_plan "Update settings: delete_branch=$DELETE_BRANCH, squash=$ALLOW_SQUASH, merge=$ALLOW_MERGE, rebase=$ALLOW_REBASE"
else
    gh repo edit "$REPO" \
        --enable-discussions="$(yq '.repository.features.discussions' "$SPEC_FILE")" \
        --enable-wiki="$(yq '.repository.features.wiki' "$SPEC_FILE")" \
        --enable-projects="$(yq '.repository.features.projects' "$SPEC_FILE")" \
        --delete-branch-on-merge="$DELETE_BRANCH" \
        --allow-squash-merge="$ALLOW_SQUASH" \
        --allow-merge-commit="$ALLOW_MERGE" \
        --allow-rebase-merge="$ALLOW_REBASE"
fi

# 5. Labels
print_step "Synchronizing Labels..."
existing_labels=""
if [[ "$DRY_RUN" == "false" ]] || gh repo view "$REPO" &>/dev/null; then
    existing_labels=$(gh label list --repo "$REPO" --limit 100 2>/dev/null | awk '{print $1}' || echo "")
fi

yq -c '.labels[]' "$SPEC_FILE" | while read -r label; do
  name=$(echo "$label" | jq -r '.name')
  color=$(echo "$label" | jq -r '.color')
  if echo "$existing_labels" | grep -q "^$name$"; then
    if [[ "$DRY_RUN" == "true" ]]; then print_plan "Edit label: $name ($color)"; else
        gh label edit "$name" --color "$color" --repo "$REPO" --silent
    fi
  else
    if [[ "$DRY_RUN" == "true" ]]; then print_plan "Create label: $name ($color)"; else
        gh label create "$name" --color "$color" --repo "$REPO" --silent
    fi
  fi
done

# 6. Milestones
print_step "Synchronizing Milestones..."
existing_milestones=""
if [[ "$DRY_RUN" == "false" ]] || gh repo view "$REPO" &>/dev/null; then
    existing_milestones=$(gh api "repos/$REPO/milestones" --jq '.[].title' 2>/dev/null || echo "")
fi

yq -r '.milestones[]' "$SPEC_FILE" | while read -r m; do
  if echo "$existing_milestones" | grep -q "^$m$"; then
    print_ok "Milestone '$m' exists."
  else
    if [[ "$DRY_RUN" == "true" ]]; then print_plan "Create milestone: $m"; else
        gh api -X POST "repos/$REPO/milestones" -f title="$m" --silent
    fi
  fi
done

# 7. Branch Protections
apply_protection() {
  local branch=$1
  if [[ "$DRY_RUN" == "true" ]]; then
    print_plan "Apply protection to $branch"
  else
    print_step "Applying Branch Protection: $branch..."
    yq -o=json ".branches.$branch.protection" "$SPEC_FILE" | gh api -X PUT "repos/$REPO/branches/$branch/protection" \
      -H "Accept: application/vnd.github+json" \
      --input - --silent || print_warn "Failed to apply protection to $branch."
  fi
}

apply_protection "main"
apply_protection "develop"

if [[ "$DRY_RUN" == "true" ]]; then
    print_ok "Plan complete. Run without --plan to apply."
else
    print_ok "GitHub Ecosystem Setup Complete!"
fi
