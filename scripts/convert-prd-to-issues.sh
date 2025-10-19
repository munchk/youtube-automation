#!/bin/bash

# Script to convert PRD tasks to GitHub issues
# Usage: ./scripts/convert-prd-to-issues.sh [prd-file]

set -e

# Default PRD file
PRD_FILE="${1:-prd.txt}"
PRD_PATH=".taskmaster/docs/$PRD_FILE"

echo "🔄 Converting PRD to GitHub Issues..."
echo "📄 PRD file: $PRD_PATH"

# Check if PRD file exists
if [ ! -f "$PRD_PATH" ]; then
    echo "❌ PRD file not found: $PRD_PATH"
    echo "Available PRD files:"
    ls -la .taskmaster/docs/*.txt 2>/dev/null || echo "No PRD files found"
    exit 1
fi

# Check if Taskmaster is initialized
if [ ! -d ".taskmaster" ]; then
    echo "🔧 Initializing Taskmaster..."
    task-master init --yes
fi

# Parse PRD and generate tasks
echo "📋 Parsing PRD and generating tasks..."
task-master parse-prd "$PRD_PATH" --force --research

# Get tasks count
TASK_COUNT=$(task-master list --json | jq '.tasks | length' 2>/dev/null || echo "0")
echo "📊 Generated $TASK_COUNT tasks"

if [ "$TASK_COUNT" -eq 0 ]; then
    echo "⚠️  No tasks found to convert"
    exit 0
fi

# Create GitHub issues for each task
echo "🎫 Creating GitHub issues..."

task-master list --json | jq -r '.tasks[] | select(.status == "pending") | @base64' | while read -r task_b64; do
    if [ -n "$task_b64" ]; then
        task=$(echo "$task_b64" | base64 -d)
        
        task_id=$(echo "$task" | jq -r '.id')
        task_title=$(echo "$task" | jq -r '.title')
        task_description=$(echo "$task" | jq -r '.description')
        task_priority=$(echo "$task" | jq -r '.priority')
        task_dependencies=$(echo "$task" | jq -r '.dependencies | join(", ")')
        
        # Create issue body
        issue_body="## Task $task_id: $task_title

**Priority:** $task_priority
**Dependencies:** ${task_dependencies:-None}

### Description
$task_description

### Acceptance Criteria
- [ ] Implementation completed
- [ ] Tests written and passing
- [ ] Documentation updated
- [ ] Code reviewed and approved

### Technical Notes
- This task was automatically generated from PRD
- Please ensure all dependencies are completed before starting
- Update this issue with progress and any blockers"

        # Set labels based on priority
        case "$task_priority" in
            "high")
                labels="enhancement,auto-generated,high-priority"
                ;;
            "medium")
                labels="enhancement,auto-generated,medium-priority"
                ;;
            "low")
                labels="enhancement,auto-generated,low-priority"
                ;;
            *)
                labels="enhancement,auto-generated"
                ;;
        esac

        # Create GitHub issue
        echo "Creating issue for task $task_id: $task_title"
        
        if gh issue create \
            --title "$task_title" \
            --body "$issue_body" \
            --label "$labels" \
            --repo "$(git remote get-url origin | sed 's/.*github.com[:/]\([^.]*\).*/\1/')" 2>/dev/null; then
            echo "✅ Created issue for task $task_id"
        else
            echo "❌ Failed to create issue for task $task_id"
        fi
    fi
done

echo "🎉 PRD to GitHub issues conversion completed!"
echo "📋 Check your repository's Issues tab to see the new issues"
