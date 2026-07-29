import os
import json
import urllib.request
import urllib.error
import sys

def main():
    api_key = os.environ.get("NOTION_API_KEY")
    if not api_key:
        print("Missing NOTION_API_KEY")
        sys.exit(1)

    try:
        with open("config/notion_databases.json", "r") as f:
            dbs = json.load(f)
    except Exception as e:
        print(f"Error loading config: {e}")
        sys.exit(1)

    schemas = {}
    for key, db_id in dbs.items():
        if not db_id or key in ["parent_page", "graph_metadata"]:
            continue
            
        print(f"Fetching {key}...")
        req = urllib.request.Request(
            f"https://api.notion.com/v1/databases/{db_id}",
            headers={
                "Authorization": f"Bearer {api_key}",
                "Notion-Version": "2022-06-28"
            }
        )
        try:
            with urllib.request.urlopen(req) as response:
                data = json.loads(response.read().decode())
                schemas[key] = {
                    "id": data.get("id"),
                    "title": data.get("title", [{}])[0].get("plain_text", "") if data.get("title") else "",
                    "properties": data.get("properties", {})
                }
        except urllib.error.URLError as e:
            print(f"Failed to fetch {key}: {e}")

    os.makedirs("context", exist_ok=True)
    with open("context/notion_schemas.json", "w") as f:
        json.dump(schemas, f, indent=2)
    print("Saved schemas to context/notion_schemas.json")

if __name__ == "__main__":
    main()
