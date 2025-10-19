#!/bin/bash

# Post-parse-prd hook: Automatically convert tasks to GitHub issues
# This script runs after a PRD is parsed and tasks are generated

echo "🔄 Post-PRD hook: Converting tasks to GitHub issues..."

# Check if we're in a Git repository
if [ ! -d ".git" ]; then
    echo "⚠️  Not in a Git repository, skipping GitHub issue creation"
    exit 0
fi

# Check if GitHub CLI is available
if ! command -v gh &> /dev/null; then
    echo "⚠️  GitHub CLI (gh) not found, skipping GitHub issue creation"
    echo "💡 Install GitHub CLI: https://cli.github.com/"
    exit 0
fi

# Check if we're authenticated with GitHub
if ! gh auth status &> /dev/null; then
    echo "⚠️  Not authenticated with GitHub, skipping GitHub issue creation"
    echo "💡 Run 'gh auth login' to authenticate"
    exit 0
fi

# Run the conversion script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

if [ -f "$PROJECT_ROOT/scripts/convert-prd-to-issues.sh" ]; then
    echo "🚀 Running PRD to GitHub issues conversion..."
    "$PROJECT_ROOT/scripts/convert-prd-to-issues.sh"
else
    echo "⚠️  Conversion script not found: $PROJECT_ROOT/scripts/convert-prd-to-issues.sh"
fi
