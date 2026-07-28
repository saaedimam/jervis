#!/bin/bash
set -e

REPO="ioriimasu/jervis"

echo "🚀 Starting GitHub Ecosystem Setup for $REPO"

# 1. Enable Core Features
echo "🔹 Enabling Discussions, Wiki, and Projects..."
gh repo edit $REPO --enable-discussions --enable-wiki --enable-projects

# 2. Configure Discussions Categories
echo "🔹 Configuring Discussion Categories..."
# Note: Categories can't be fully managed via CLI yet, but enabling discussions is a start.

# 3. Security Configuration
echo "🔹 Enabling Security Features..."
gh api -X PATCH repos/$REPO/import/settings -f private_vulnerability_reporting=enabled || true
# Dependabot and Secret Scanning usually require UI or specific API flags if not default
gh api -X PATCH repos/$REPO -f security_and_analysis='{"secret_scanning":{"status":"enabled"},"secret_scanning_push_protection":{"status":"enabled"}}' || true

# 4. Labels Taxonomy
echo "🔹 Creating Labels..."
labels=(
  "kind/bug:#d73a4a"
  "kind/feature:#0075ca"
  "kind/refactor:#e99695"
  "kind/docs:#008672"
  "kind/security:#cfd3d7"
  "kind/performance:#a2eeef"
  "kind/test:#fef2c0"
  "kind/architecture:#7057ff"
  "priority/high:#d93f0b"
  "priority/medium:#fbca04"
  "priority/low:#0e8a16"
  "status/blocked:#000000"
  "status/in-review:#cc317c"
  "status/needs-info:#5319e7"
  "phase/runtime:#bfd4f2"
  "phase/memory:#bfdadc"
  "phase/services:#c5def5"
  "phase/providers:#d4c5f9"
  "phase/interfaces:#f9d0c4"
)

for label in "${labels[@]}"; do
  IFS=":" read -r name color <<< "$label"
  gh label create "$name" --color "$color" --force || true
done

# 5. Version-based Milestones
echo "🔹 Creating Milestones..."
milestones=(
  "v0.1.0 Runtime"
  "v0.2.0 Memory"
  "v0.3.0 Services"
  "v0.4.0 Providers"
  "v0.5.0 Interfaces"
  "v1.0.0 Stable"
)

for m in "${milestones[@]}"; do
  gh api -X POST repos/$REPO/milestones -f title="$m" || true
done

# 6. Branch Protection (main)
echo "🔹 Setting Branch Protection for main..."
gh api -X PUT repos/$REPO/branches/main/protection \
  -H "Accept: application/vnd.github+json" \
  -F "required_status_checks[strict]=true" \
  -F "required_status_checks[contexts][]=test" \
  -F "required_status_checks[contexts][]=lint" \
  -F "enforce_admins=true" \
  -F "required_pull_request_reviews[required_approving_review_count]=1" \
  -F "required_pull_request_reviews[dismiss_stale_reviews]=true" \
  -F "required_pull_request_reviews[require_code_owner_reviews]=true" \
  -F "restrictions=null" \
  -F "required_linear_history=true" \
  -F "allow_force_pushes=false" \
  -F "allow_deletions=false" || true

# 7. Branch Protection (develop)
echo "🔹 Setting Branch Protection for develop..."
gh api -X PUT repos/$REPO/branches/develop/protection \
  -H "Accept: application/vnd.github+json" \
  -F "required_status_checks[strict]=true" \
  -F "required_status_checks[contexts][]=test" \
  -F "enforce_admins=false" \
  -F "required_pull_request_reviews[required_approving_review_count]=1" \
  -F "required_linear_history=true" \
  -F "allow_force_pushes=true" \
  -F "allow_deletions=false" || true

echo "✅ GitHub Ecosystem Setup Complete!"
