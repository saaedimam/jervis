import os
import sys
import yaml
import json
import urllib.request
import urllib.error

def main():
    if len(sys.argv) < 3:
        sys.exit(1)
        
    yaml_file = sys.argv[1]
    db_id = sys.argv[2]
    api_key = os.environ.get("NOTION_API_KEY")
    
    with open(yaml_file, 'r') as f:
        data = yaml.safe_load(f)
        
    if not data:
        sys.exit(0)
        
    # The YAML top-level key is the list of items
    items = list(data.values())[0] if isinstance(data, dict) else data
    
    for item in items:
        # Mapping item dict to Notion properties dynamically
        properties = {}
        for k, v in item.items():
            # Basic mapping logic based on key names
            if k == 'id':
                properties[k.title()] = {"title": [{"text": {"content": str(v)}}]}
            elif k in ['status', 'impact', 'priority']:
                properties[k.title()] = {"select": {"name": str(v)}}
            elif k == 'date':
                properties[k.title()] = {"date": {"start": str(v)}}
            else:
                properties[k.title()] = {"rich_text": [{"text": {"content": str(v)}}]}
                
        req = urllib.request.Request(
            "https://api.notion.com/v1/pages",
            data=json.dumps({
                "parent": {"database_id": db_id},
                "properties": properties
            }).encode("utf-8"),
            headers={
                "Authorization": f"Bearer {api_key}",
                "Notion-Version": "2022-06-28",
                "Content-Type": "application/json"
            },
            method="POST"
        )
        try:
            with urllib.request.urlopen(req) as response:
                pass
        except urllib.error.URLError as e:
            print(f"Error creating item: {e}")

if __name__ == "__main__":
    main()
