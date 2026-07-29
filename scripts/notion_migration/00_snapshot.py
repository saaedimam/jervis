import os
import json
import subprocess
from datetime import datetime
from common.client import NotionClient
from common.ids import DB_IDS
from common.logging import setup_logger

logger = setup_logger("snapshot")

def get_git_info():
    commit = subprocess.check_output(["git", "rev-parse", "HEAD"]).decode().strip()
    return commit

def run():
    logger.info("Starting Stage 0: Repository Snapshot")
    client = NotionClient()
    
    snapshot = {
        "timestamp": datetime.utcnow().isoformat() + "Z",
        "git_commit": get_git_info(),
        "database_ids": DB_IDS,
        "schemas": {}
    }
    
    for name, db_id in DB_IDS.items():
        if not db_id or name in ["parent_page", "graph_metadata"]:
            continue
        try:
            logger.info(f"Snapshotting schema for {name} ({db_id})...")
            db_data = client.get_database(db_id)
            snapshot["schemas"][name] = db_data.get("properties", {})
        except Exception as e:
            logger.error(f"Failed to snapshot {name}: {e}")
            
    os.makedirs("context/migration", exist_ok=True)
    with open("context/migration/snapshot.json", "w") as f:
        json.dump(snapshot, f, indent=2)
        
    logger.info("Stage 0 complete. Rollback metadata saved to context/migration/snapshot.json.")
    return True

if __name__ == "__main__":
    run()
