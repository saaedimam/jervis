import os
import json
import urllib.request
import urllib.error

notion_key = os.environ.get("NOTION_API_KEY")
if not notion_key:
    print("NOTION_API_KEY missing")
    exit(1)

req = urllib.request.Request(
    "https://api.notion.com/v1/search",
    data=json.dumps({"filter": {"value": "database", "property": "object"}}).encode("utf-8"),
    headers={
        "Authorization": f"Bearer {notion_key}",
        "Notion-Version": "2022-06-28",
        "Content-Type": "application/json"
    }
)

try:
    with urllib.request.urlopen(req) as response:
        data = json.loads(response.read().decode())
        for result in data.get("results", []):
            title = result.get("title", [{}])
            name = title[0].get("plain_text", "Untitled") if title else "Untitled"
            print(f"{name}: {result['id']}")
except urllib.error.URLError as e:
    print(f"Error: {e}")
