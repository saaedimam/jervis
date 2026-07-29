import json
import os
import urllib.request
import urllib.error
import logging

class NotionClient:
    def __init__(self):
        self.api_key = os.environ.get("NOTION_API_KEY")
        if not self.api_key:
            raise ValueError("NOTION_API_KEY not set")
            
        self.base_url = "https://api.notion.com/v1"
        self.headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Notion-Version": "2022-06-28",
            "Content-Type": "application/json"
        }
        
    def _request(self, method, endpoint, data=None):
        url = f"{self.base_url}{endpoint}"
        body = json.dumps(data).encode("utf-8") if data else None
        req = urllib.request.Request(url, data=body, headers=self.headers, method=method)
        try:
            with urllib.request.urlopen(req) as response:
                return json.loads(response.read().decode())
        except urllib.error.URLError as e:
            err_body = e.read().decode("utf-8") if hasattr(e, "read") else str(e)
            logging.error(f"Notion API Error: {err_body}")
            raise RuntimeError(f"Notion API Request Failed: {err_body}")

    def get_database(self, db_id):
        return self._request("GET", f"/databases/{db_id}")
        
    def update_database(self, db_id, properties):
        return self._request("PATCH", f"/databases/{db_id}", {"properties": properties})
        
    def query_database(self, db_id, filter=None):
        payload = {}
        if filter:
            payload["filter"] = filter
        return self._request("POST", f"/databases/{db_id}/query", payload)
        
    def create_page(self, parent_db, properties, children=None):
        payload = {
            "parent": {"database_id": parent_db},
            "properties": properties
        }
        if children:
            payload["children"] = children
        return self._request("POST", "/pages", payload)
