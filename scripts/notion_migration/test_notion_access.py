import os
import json
import urllib.request

api_key = os.environ.get("NOTION_API_KEY")
headers = {
    "Authorization": f"Bearer {api_key}",
    "Notion-Version": "2022-06-28",
    "Content-Type": "application/json"
}

req = urllib.request.Request(
    "https://api.notion.com/v1/search",
    data=json.dumps({"filter": {"value": "database", "property": "object"}}).encode("utf-8"),
    headers=headers,
    method="POST"
)

try:
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode())
        print(f"Found {len(data.get('results', []))} databases shared with integration key.")
        for item in data.get("results", []):
            title = item.get("title", [{}])[0].get("plain_text", "Untitled") if item.get("title") else "Untitled"
            print(f"- Title: {title} | ID: {item.get('id')}")
except Exception as e:
    print(f"Search failed: {e}")
