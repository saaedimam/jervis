import os
import json
import time
import urllib.request
import urllib.error
import yaml
import subprocess

def get_headers():
    api_key = os.environ.get("NOTION_API_KEY")
    return {
        "Authorization": f"Bearer {api_key}",
        "Notion-Version": "2022-06-28",
        "Content-Type": "application/json"
    }

def get_db_ids():
    with open("config/notion_databases.json", "r") as f:
        return json.load(f)

def query_existing_titles(db_id, title_prop_name, headers):
    titles = set()
    has_more = True
    next_cursor = None
    while has_more:
        payload = {}
        if next_cursor:
            payload["start_cursor"] = next_cursor
        req = urllib.request.Request(
            f"https://api.notion.com/v1/databases/{db_id}/query",
            data=json.dumps(payload).encode("utf-8"),
            headers=headers,
            method="POST"
        )
        try:
            with urllib.request.urlopen(req) as resp:
                data = json.loads(resp.read().decode())
                for item in data.get("results", []):
                    props = item.get("properties", {})
                    t_val = props.get(title_prop_name, {})
                    if t_val.get("type") == "title":
                        t_list = t_val.get("title", [])
                        if t_list:
                            titles.add(t_list[0].get("plain_text", ""))
                has_more = data.get("has_more", False)
                next_cursor = data.get("next_cursor")
        except Exception as e:
            print(f"Error querying {db_id}: {e}")
            break
    return titles

def create_notion_row(db_id, title_prop_name, title_val, extra_props, headers):
    payload = {
        "parent": {"database_id": db_id},
        "properties": {
            title_prop_name: {
                "title": [{"text": {"content": title_val}}]
            }
        }
    }
    if extra_props:
        payload["properties"].update(extra_props)
    
    req = urllib.request.Request(
        "https://api.notion.com/v1/pages",
        data=json.dumps(payload).encode("utf-8"),
        headers=headers,
        method="POST"
    )
    try:
        with urllib.request.urlopen(req) as resp:
            time.sleep(0.25)
            return True
    except Exception as e:
        err_msg = e.read().decode("utf-8") if hasattr(e, "read") else str(e)
        print(f"Failed '{title_val}' in {db_id}: {err_msg}")
        return False

def populate_yaml_store(db_key, title_prop, yaml_file, root_yaml_key, extra_builder, headers, db_ids):
    db_id = db_ids.get(db_key)
    if not db_id:
        return 0, 0

    existing = query_existing_titles(db_id, title_prop, headers)
    if not os.path.exists(yaml_file):
        return len(existing), 0

    with open(yaml_file, "r") as f:
        data = yaml.safe_load(f) or {}

    items = data.get(root_yaml_key, [])
    created = 0

    for item in items:
        title_val = str(item.get("id") or item.get("name") or item.get("version") or item.get("description", "Untitled"))
        if title_val in existing:
            continue

        extra_props = extra_builder(item) if extra_builder else {}
        if create_notion_row(db_id, title_prop, title_val, extra_props, headers):
            created += 1
            existing.add(title_val)

    return len(existing), created

def main():
    headers = get_headers()
    db_ids = get_db_ids()

    print("Starting Live Supporting Databases Ingestion with Verified Schemas...")
    summary = {}

    # 1. Tasks
    cnt, created = populate_yaml_store("tasks", "Task", "data/tasks.yaml", "tasks", None, headers, db_ids)
    summary["tasks"] = (cnt, created)

    # 2. Bugs
    cnt, created = populate_yaml_store("bugs", "Bug ID", "data/bugs.yaml", "bugs", None, headers, db_ids)
    summary["bugs"] = (cnt, created)

    # 3. Releases
    cnt, created = populate_yaml_store("releases", "Version", "data/releases.yaml", "releases", None, headers, db_ids)
    summary["releases"] = (cnt, created)

    # 4. Tech Debt
    cnt, created = populate_yaml_store("tech_debt", "Debt ID", "data/tech_debt.yaml", "tech_debt", None, headers, db_ids)
    summary["tech_debt"] = (cnt, created)

    # 5. Risks
    cnt, created = populate_yaml_store("risks", "Risk ID", "data/risks.yaml", "risks", None, headers, db_ids)
    summary["risks"] = (cnt, created)

    print("\nLive Population Ingestion Results:")
    for db, (total, new_added) in summary.items():
        print(f"- `{db}`: {total} total live records ({new_added} newly created)")

if __name__ == "__main__":
    main()
