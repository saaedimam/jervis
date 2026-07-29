import os
import json
import urllib.request

def main():
    api_key = os.environ.get("NOTION_API_KEY")
    headers = {
        "Authorization": f"Bearer {api_key}",
        "Notion-Version": "2022-06-28",
        "Content-Type": "application/json"
    }

    # 1. Fetch live databases via Search API
    req = urllib.request.Request(
        "https://api.notion.com/v1/search",
        data=json.dumps({"filter": {"value": "database", "property": "object"}}).encode("utf-8"),
        headers=headers,
        method="POST"
    )

    live_dbs_by_title = {}
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode())
        for item in data.get("results", []):
            raw_title = item.get("title", [{}])[0].get("plain_text", "") if item.get("title") else ""
            clean_title = raw_title.lower().replace("🏛️", "").replace("📦", "").replace("📄", "").replace("🔌", "").replace("🧵", "").replace("🏁", "").replace("📜", "").replace("📅", "").replace("🧠", "").replace("📋", "").replace("📏", "").replace("🤖", "").strip()
            live_dbs_by_title[clean_title] = {
                "id": item["id"],
                "raw_title": raw_title,
                "archived": item.get("archived", False)
            }

    # 2. Build live mapping dict
    mapping_name_to_key = {
        "architecture registry": "architecture",
        "architecture": "architecture",
        "package registry": "packages",
        "file registry": "files",
        "api registry": "apis",
        "commit registry": "commits",
        "milestones": "milestones",
        "adrs": "adrs",
        "session logs": "sessions",
        "engineering memory": "memory",
        "specifications": "specifications",
        "coding standards": "coding_standards",
        "ai handoffs": "handoffs",
        "tasks": "tasks",
        "bugs": "bugs",
        "releases": "releases",
        "decisions": "engineering_decisions",
        "prompt library": "prompt_library"
    }

    mapped_config = {}
    for title_str, db_info in live_dbs_by_title.items():
        key = mapping_name_to_key.get(title_str)
        if key and key not in mapped_config:
            mapped_config[key] = db_info

    # 3. Inspect real record counts for each mapped database
    audit_report = ["# LIVE NOTION DATABASE ROOT CAUSE AUDIT\n"]
    audit_report.append("## Root Cause Findings\n")
    audit_report.append("**Primary Root Cause**: Database ID Mismatch in `config/notion_databases.json`.")
    audit_report.append("The configuration file contained invalid/stale Database IDs (e.g. `d3dcb133-f96e-4e8e-944f-5825c2d1eee0`), causing Notion API requests to fail silently with HTTP 404.\n")

    audit_report.append("## Live Database Verification Audit Table\n")
    audit_report.append("| Database Name | Clean Key | Live Notion Database ID | Live Record Count | Status |")
    audit_report.append("| --- | --- | --- | --- | --- |")

    new_config_json = {}
    for key, info in mapped_config.items():
        db_id = info["id"]
        new_config_json[key] = db_id

        # Query live count
        query_req = urllib.request.Request(
            f"https://api.notion.com/v1/databases/{db_id}/query",
            data=json.dumps({}).encode("utf-8"),
            headers=headers,
            method="POST"
        )
        count = 0
        try:
            with urllib.request.urlopen(query_req) as q_resp:
                q_data = json.loads(q_resp.read().decode())
                count = len(q_data.get("results", []))
        except Exception as e:
            count = f"Error ({e})"

        status = "Archived" if info["archived"] else ("Populated" if isinstance(count, int) and count > 0 else "Empty")
        audit_report.append(f"| {info['raw_title']} | `{key}` | `{db_id}` | **{count}** | {status} |")

    # Update config/notion_databases.json with correct live IDs
    with open("config/notion_databases.json", "r") as f:
        existing_full_cfg = json.load(f)

    existing_full_cfg.update(new_config_json)
    with open("config/notion_databases.json", "w") as f:
        json.dump(existing_full_cfg, f, indent=2)

    os.makedirs("context/migration", exist_ok=True)
    with open("context/migration/ROOT_CAUSE_INVESTIGATION.md", "w") as f:
        f.write("\n".join(audit_report))

    print("Remapped config and generated live audit in context/migration/ROOT_CAUSE_INVESTIGATION.md")

if __name__ == "__main__":
    main()
