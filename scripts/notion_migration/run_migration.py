#!/usr/bin/env python3
import sys
import os

def main():
    print("JERVIS Notion Migration Orchestrator")
    print("=" * 40)
    print("Run individual stages to proceed. According to the approved strategy,")
    print("stages must be run sequentially and require manual validation.\n")
    print("Usage:")
    print("  python3 -m scripts.notion_migration.00_snapshot")
    print("  python3 -m scripts.notion_migration.01_schema")
    print("  python3 -m scripts.notion_migration.02_templates")
    print("  ...")

if __name__ == "__main__":
    main()
