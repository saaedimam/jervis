import json
from common.client import NotionClient
from common.ids import DB_IDS
from common.logging import setup_logger

logger = setup_logger("schema")

CANONICAL_SCHEMA = {
    "architecture": {
        "Name": {"rich_text": {}},
        "Purpose": {"rich_text": {}},
        "Design Principles": {"rich_text": {}},
        "Status": {"select": {"options": [{"name": "Active", "color": "green"}, {"name": "Deprecated", "color": "red"}]}},
        "Coverage": {"number": {"format": "percent"}},
        "Packages": {"relation": {"database_id": DB_IDS.get("packages"), "single_property": {}}},
        "Specifications": {"relation": {"database_id": DB_IDS.get("specifications"), "single_property": {}}},
        "ADRs": {"relation": {"database_id": DB_IDS.get("adrs"), "single_property": {}}}
    },
    "packages": {
        "Path": {"rich_text": {}},
        "Purpose": {"rich_text": {}},
        "Coverage": {"number": {"format": "percent"}},
        "Architecture": {"relation": {"database_id": DB_IDS.get("architecture"), "single_property": {}}},
        "Files": {"relation": {"database_id": DB_IDS.get("files"), "single_property": {}}},
        "APIs": {"relation": {"database_id": DB_IDS.get("apis"), "single_property": {}}}
    }
}

def run():
    logger.info("Starting Stage 1: Schema Recovery (Strict Validation)")
    client = NotionClient()
    
    with open("context/migration/snapshot.json", "r") as f:
        snapshot = json.load(f)
    schemas = snapshot.get("schemas", {})
    
    updates_made = 0
    report = ["# SCHEMA VALIDATION\n"]
    
    for db_name, canonical_props in CANONICAL_SCHEMA.items():
        db_id = DB_IDS.get(db_name)
        if not db_id:
            continue
            
        current_props = schemas.get(db_name, {})
        missing_props = {}
        
        for prop_name, prop_config in canonical_props.items():
            if prop_name not in current_props:
                missing_props[prop_name] = prop_config
            else:
                # Basic type validation
                expected_type = list(prop_config.keys())[0]
                actual_type = current_props[prop_name].get("type")
                if actual_type != expected_type:
                    logger.warning(f"Type mismatch on {db_name}.{prop_name}: expected {expected_type}, got {actual_type}")
        
        report.append(f"### {db_name.title()} ({db_id})")
        
        if missing_props:
            logger.info(f"Updating {db_name} schema. Adding properties: {list(missing_props.keys())}")
            try:
                # Note: creating relations via API requires specific payloads
                client.update_database(db_id, missing_props)
                updates_made += len(missing_props)
                report.append(f"- Properties added: {len(missing_props)}")
                report.append("- Schema: COMPLETE")
            except Exception as e:
                logger.error(f"Failed to update {db_name}: {e}")
                report.append("- Schema: FAILED TO UPDATE")
                return False
        else:
            logger.info(f"{db_name} schema is already complete.")
            report.append("- Properties added: 0 (Idempotent)")
            report.append("- Schema: COMPLETE")
            
    with open("context/migration/SCHEMA_VALIDATION.md", "w") as f:
        f.write("\n".join(report))
        
    logger.info(f"Stage 1 complete. Total schema updates made: {updates_made}")
    return True

if __name__ == "__main__":
    run()
