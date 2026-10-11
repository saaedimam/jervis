#!/usr/bin/env python3
import os
import json
import subprocess

# Load configuration
with open('config/notion_databases.json', 'r') as f:
    config = json.load(f)

TASKS_DB_ID = config['tasks']
PROJECTS_DB_ID = config['packages']

def get_jervis_tasks():
    """Get all tasks from the Jervis planner."""
    result = subprocess.run(['./bin/jervis', 'planner', '-list'], capture_output=True, text=True)
    lines = result.stdout.strip().split('\n')
    tasks = []
    for line in lines[1:]: # Skip header
        parts = line.split()
        tasks.append({
            'id': parts[0],
            'status': parts[1],
            'title': ' '.join(parts[2:]),
        })
    return tasks

def sync_tasks(tasks):
    """Sync tasks to the Notion Tasks database."""
    for task in tasks:
        print(f"Syncing task: {task['title']}")

        # Prepare the data for Notion API
        data = {
            "parent": {"database_id": TASKS_DB_ID},
            "properties": {
                "Task": {
                    "title": [
                        {
                            "text": {
                                "content": task['title']
                            }
                        }
                    ]
                },
                "Status": {
                    "select": {
                        "name": task['status']
                    }
                }
            }
        }

        # Convert data to JSON string for ntn command
        json_data = json.dumps(data)

        # Use ntn API directly to create page
        result = subprocess.run([
            'ntn', 'api', '/v1/pages',
            '-X', 'POST',
            '-d', json_data
        ], capture_output=True, text=True, env={**os.environ, 'NOTION_API_KEY': os.environ.get('NOTION_TOKEN', '')})

        if result.returncode != 0:
            print(f"Error syncing task {task['title']}: {result.stderr}")
        else:
            print(f"Successfully synced task: {task['title']}")

if __name__ == '__main__':
    print('Syncing core data to Notion...')
    tasks = get_jervis_tasks()
    sync_tasks(tasks)
    print('Core data sync complete.')