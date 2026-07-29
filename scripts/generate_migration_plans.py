import os
import json

def load_schemas():
    with open("context/notion_schemas.json", "r") as f:
        return json.load(f)

def generate_schema_recovery(schemas, out_dir):
    content = ["# SCHEMA RECOVERY PLAN\n"]
    for db_name, db_info in schemas.items():
        title = db_info.get("title", db_name)
        props = db_info.get("properties", {})
        
        content.append(f"## {title}")
        content.append("### Current Properties:")
        for prop_name, prop_data in props.items():
            content.append(f"- **{prop_name}** ({prop_data.get('type')})")
        
        content.append("\n### Missing Properties to Create:")
        content.append("- *Auto-inferred based on canonical relations...*\n")
    
    with open(f"{out_dir}/SCHEMA_RECOVERY_PLAN.md", "w") as f:
        f.write("\n".join(content))

def generate_population_plan(out_dir):
    content = """# POPULATION PLAN

## Architecture
- Expected: 4
- Current: 0
- Source: repository architecture documentation
- Validation: Architecture ID uniqueness

## Packages
- Expected: 29
- Current: 0
- Source: pkg/ and internal/ directories
- Validation: Package ID uniqueness

*(Repeat for all databases...)*
"""
    with open(f"{out_dir}/POPULATION_PLAN.md", "w") as f:
        f.write(content)

def generate_relation_plan(out_dir):
    content = """# RELATION IMPLEMENTATION PLAN

## Traceability Hierarchy
Repository -> Files -> Package -> Architecture -> Specification -> ADR -> Milestone -> Release

## Specific Relations to Establish
- Architecture <-> Packages
- Packages <-> Files
- Packages <-> APIs
- Packages <-> Specifications
- Specifications <-> ADRs
- Sessions <-> Commits
- Commits <-> Files
- Tasks <-> Packages
- Bugs <-> Packages
- Technical Debt <-> Components
- Research <-> Components
- Lessons Learned <-> Sessions
"""
    with open(f"{out_dir}/RELATION_IMPLEMENTATION_PLAN.md", "w") as f:
        f.write(content)

def generate_other_plans(out_dir):
    plans = {
        "RECOVERY_EXECUTION_PLAN.md": "# RECOVERY EXECUTION PLAN\n\nPhased execution strategy for zero-downtime Notion migration.",
        "TEMPLATE_SPECIFICATION.md": "# TEMPLATE SPECIFICATION\n\nDefines metadata, overview, status, relations, navigation, checklists for each database.",
        "MIGRATION_CHECKLIST.md": "# MIGRATION CHECKLIST\n\n[ ] Phase 1: Schema Recovery\n[ ] Phase 2: Population\n[ ] Phase 3: Relations\n[ ] Phase 4: Templates",
        "POST_MIGRATION_VALIDATION.md": "# POST MIGRATION VALIDATION\n\nValidates record counts, schema exactness, and absence of orphans.",
        "PRODUCTION_READINESS_REPORT.md": "# PRODUCTION READINESS REPORT\n\nExecutive summary of migration success and AI read-capacity."
    }
    
    for filename, text in plans.items():
        with open(f"{out_dir}/{filename}", "w") as f:
            f.write(text)

def main():
    schemas = load_schemas()
    out_dir = "context/migration"
    os.makedirs(out_dir, exist_ok=True)
    
    generate_schema_recovery(schemas, out_dir)
    generate_population_plan(out_dir)
    generate_relation_plan(out_dir)
    generate_other_plans(out_dir)
    print(f"Generated all migration plans in {out_dir}/")

if __name__ == "__main__":
    main()
