import json
import os

def load_db_ids():
    config_path = os.path.join(os.path.dirname(__file__), "../../../config/notion_databases.json")
    with open(config_path, "r") as f:
        return json.load(f)

DB_IDS = load_db_ids()

def get_id(key):
    return DB_IDS.get(key)
