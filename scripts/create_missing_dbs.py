import os
import json
import urllib.request
import urllib.error
import sys

def load_config():
    with open("config/notion_databases.json", "r") as f:
        return json.load(f)

def save_config(config):
    with open("config/notion_databases.json", "w") as f:
        json.dump(config, f, indent=2)

def create_database(api_key, parent_id, title, properties):
    req = urllib.request.Request(
        "https://api.notion.com/v1/databases",
        data=json.dumps({
            "parent": {"type": "page_id", "page_id": parent_id},
            "title": [{"type": "text", "text": {"content": title}}],
            "properties": properties
        }).encode("utf-8"),
        headers={
            "Authorization": f"Bearer {api_key}",
            "Notion-Version": "2022-06-28",
            "Content-Type": "application/json"
        },
        method="POST"
    )
    
    try:
        with urllib.request.urlopen(req) as response:
            data = json.loads(response.read().decode())
            return data["id"]
    except urllib.error.URLError as e:
        body = e.read().decode("utf-8") if hasattr(e, "read") else str(e)
        print(f"Error creating '{title}': {body}")
        return None

def main():
    api_key = os.environ.get("NOTION_API_KEY")
    if not api_key:
        print("Missing NOTION_API_KEY")
        sys.exit(1)
        
    config = load_config()
    parent_id = config.get("parent_page")
    if not parent_id:
        print("Missing parent_page in config")
        sys.exit(1)

    databases_to_create = {
        "tasks": ("Tasks", {"Task ID": {"title": {}}, "Name": {"rich_text": {}}, "Status": {"select": {"options": [{"name": "Todo", "color": "gray"}, {"name": "In Progress", "color": "blue"}, {"name": "Done", "color": "green"}]}}}),
        "releases": ("Releases", {"Release ID": {"title": {}}, "Version": {"rich_text": {}}, "Date": {"date": {}}}),
        "bugs": ("Bugs", {"Bug ID": {"title": {}}, "Description": {"rich_text": {}}, "Status": {"select": {"options": [{"name": "Open", "color": "red"}, {"name": "Resolved", "color": "green"}]}}}),
        "risks": ("Risks", {"Risk ID": {"title": {}}, "Description": {"rich_text": {}}, "Impact": {"select": {"options": [{"name": "High", "color": "red"}, {"name": "Medium", "color": "yellow"}, {"name": "Low", "color": "green"}]}}}),
        "tech_debt": ("Technical Debt", {"Debt ID": {"title": {}}, "Description": {"rich_text": {}}, "Priority": {"select": {"options": [{"name": "High", "color": "red"}, {"name": "Medium", "color": "yellow"}, {"name": "Low", "color": "green"}]}}}),
        "research": ("Research", {"Research ID": {"title": {}}, "Topic": {"rich_text": {}}, "Summary": {"rich_text": {}}}),
        "prompt_library": ("Prompt Library", {"Prompt ID": {"title": {}}, "Purpose": {"rich_text": {}}, "Context": {"rich_text": {}}}),
        "lessons_learned": ("Lessons Learned", {"Lesson ID": {"title": {}}, "Situation": {"rich_text": {}}, "Mistake": {"rich_text": {}}, "Resolution": {"rich_text": {}}}),
        "patterns": ("Patterns", {"Pattern ID": {"title": {}}, "Name": {"rich_text": {}}, "Description": {"rich_text": {}}}),
        "performance_benchmarks": ("Performance Benchmarks", {"Benchmark ID": {"title": {}}, "Metric": {"rich_text": {}}, "Value": {"rich_text": {}}}),
        "coding_standards": ("Coding Standards", {"Standard ID": {"title": {}}, "Rule": {"rich_text": {}}, "Description": {"rich_text": {}}}),
        "glossary": ("Glossary", {"Term": {"title": {}}, "Definition": {"rich_text": {}}}),
        "assets": ("Assets", {"Asset ID": {"title": {}}, "Name": {"rich_text": {}}, "Type": {"select": {"options": [{"name": "Image", "color": "blue"}, {"name": "Document", "color": "gray"}]}}}),
        "external_dependencies": ("External Dependencies", {"Dependency ID": {"title": {}}, "Name": {"rich_text": {}}, "Version": {"rich_text": {}}}),
        "roadmap": ("Roadmap", {"Item ID": {"title": {}}, "Name": {"rich_text": {}}, "Quarter": {"select": {"options": [{"name": "Q1", "color": "blue"}, {"name": "Q2", "color": "green"}, {"name": "Q3", "color": "yellow"}, {"name": "Q4", "color": "red"}]}}}),
        "engineering_decisions": ("Engineering Decisions", {"Decision ID": {"title": {}}, "Problem": {"rich_text": {}}, "Decision": {"rich_text": {}}, "Rationale": {"rich_text": {}}})
    }

    updated = False
    for key, (title, properties) in databases_to_create.items():
        if key not in config or not config[key]:
            print(f"Creating '{title}'...")
            db_id = create_database(api_key, parent_id, title, properties)
            if db_id:
                config[key] = db_id
                updated = True
                print(f"Created {title}: {db_id}")
            else:
                print(f"Failed to create {title}")
        else:
            print(f"Database '{title}' already configured ({config[key]})")
            
    if updated:
        save_config(config)
        print("Config updated.")

if __name__ == "__main__":
    main()
