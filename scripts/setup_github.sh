#!/bin/bash
# scripts/setup_github.sh
# Declarative, idempotent setup for the jervis GitHub repository.

set -e

SPEC_FILE=".github/governance.yaml"

print_step() { echo -e "\033[1;34m[STEP]\033[0m $1"; }
print_ok() { echo -e "\033[1;32m[OK]\033[0m $1"; }
print_warn() { echo -e "\033[1;33m[WARN]\033[0m $1"; }

# 0. Prerequisite Check
if ! command -v yq &> /dev/null; then echo "yq not found"; exit 2; fi

# 1. Load Spec
REPO_NAME=$(yq '.repository.name' "$SPEC_FILE")
REPO_OWNER=$(yq '.repository.owner' "$SPEC_FILE")
REPO="$REPO_OWNER/$REPO_NAME"
DESCRIPTION=$(yq '.repository.description' "$SPEC_FILE")
VISIBILITY=$(yq '.repository.visibility' "$SPEC_FILE")

# 2. Idempotent Repository Creation
print_step "Checking GitHub repository status for $REPO..."
if ! gh repo view "$REPO" &>/dev/null; then
    print_step "Creating repository $REPO..."
    gh repo create "$REPO" --"$VISIBILITY" --description "$DESCRIPTION"
else
    print_ok "Repository $REPO already exists."
fi

# 3. Idempotent Remote Configuration
print_step "Checking local git remote..."
if ! git remote get-url origin &>/dev/null; then
    print_step "Adding origin remote..."
    git remote add origin "https://github.com/$REPO.git"
else
    print_ok "Remote 'origin' already configured."
fi

# 4. Sync Repository Features & Settings
print_step "Synchronizing Repository Settings..."
DELETE_BRANCH=$(yq '.repository.settings.delete_branch_on_merge' "$SPEC_FILE")
ALLOW_SQUASH=$(yq '.repository.settings.allow_squash_merge' "$SPEC_FILE")
ALLOW_MERGE=$(yq '.repository.settings.allow_merge_commit' "$SPEC_FILE")
ALLOW_REBASE=$(yq '.repository.settings.allow_rebase_merge' "$SPEC_FILE")

gh repo edit "$REPO" \
    --enable-discussions="$(yq '.repository.features.discussions' "$SPEC_FILE")" \
    --enable-wiki="$(yq '.repository.features.wiki' "$SPEC_FILE")" \
    --enable-projects="$(yq '.repository.features.projects' "$SPEC_FILE")" \
    --delete-branch-on-merge="$DELETE_BRANCH" \
    --allow-squash-merge="$ALLOW_SQUASH" \
    --allow-merge-commit="$ALLOW_MERGE" \
    --allow-rebase-merge="$ALLOW_REBASE"

# 5. Labels
print_step "Synchronizing Labels..."
existing_labels=$(gh label list --repo "$REPO" --limit 100 | awk '{print $1}')
yq -c '.labels[]' "$SPEC_FILE" | while read -r label; do
  name=$(echo "$label" | jq -r '.name')
  color=$(echo "$label" | jq -r '.color')
  if echo "$existing_labels" | grep -q "^$name$"; then
    gh label edit "$name" --color "$color" --repo "$REPO" --silent
  else
    gh label create "$name" --color "$color" --repo "$REPO" --silent
  fi
done

# 6. Milestones
print_step "Synchronizing Milestones..."
existing_milestones=$(gh api "repos/$REPO/milestones" --jq '.[].title')
yq -r '.milestones[]' "$SPEC_FILE" | while read -r m; do
  if echo "$existing_milestones" | grep -q "^$m$"; then
    print_ok "Milestone '$m' exists."
  else
    gh api -X POST "repos/$REPO/milestones" -f title="$m" --silent
  fi
done

# 7. Branch Protections
apply_protection() {
  local branch=$1
  print_step "Applying Branch Protection: $branch..."
  # Convert YAML protection to JSON and apply
  yq -o=json ".branches.$branch.protection" "$SPEC_FILE" | gh api -X PUT "repos/$REPO/branches/$branch/protection" \
    -H "Accept: application/vnd.github+json" \
    --input - --silent || print_warn "Failed to apply protection to $branch (check permissions or branch existence)."
}

apply_protection "main"
apply_protection "develop"

print_ok "GitHub Ecosystem Setup Complete!"
