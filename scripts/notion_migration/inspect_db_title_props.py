import os
import json
import urllib.request

api_key = os.environ.get("NOTION_API_KEY")
headers = {
    "Authorization": f"Bearer {api_key}",
    "Notion-Version": "2022-06-28"
}

with open("config/notion_databases.json", "r") as f:
    dbs = json.load(f)

for name, db_id in dbs.items():
    if not db_id or name in ["parent_page", "graph_metadata"]:
        continue
    req = urllib.request.Request(f"https://api.notion.com/v1/databases/{db_id}", headers=headers)
    try:
        with urllib.request.urlopen(req) as resp:
            data = json.loads(resp.read().decode())
            props = data.get("properties", {})
            title_prop = None
            for p_name, p_val in props.items():
                if p_val.get("type") == "title":
                    title_prop = p_name
                    break
            print(f"DB `{name}` ({db_id}): Title Property Name = '{title_prop}'")
    except Exception as e:
        print(f"DB `{name}` ({db_id}): Failed ({e})")
