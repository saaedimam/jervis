import sys
import yaml

def main():
    if len(sys.argv) < 3:
        print("Usage: python parse_yaml_mapping.py <mapping.yaml> <package_path>")
        sys.exit(1)

    yaml_file = sys.argv[1]
    pkg_path = sys.argv[2]

    try:
        with open(yaml_file, 'r') as f:
            data = yaml.safe_load(f)
    except Exception as e:
        print(f"Error reading YAML: {e}", file=sys.stderr)
        sys.exit(1)

    packages = data.get('packages', {})
    
    # Try exact match, then substring match
    match = packages.get(pkg_path)
    if not match:
        for p, v in packages.items():
            if pkg_path.startswith(p):
                match = v
                break
                
    if not match:
        print("ARCH-001\n") # Default fallback
        sys.exit(0)

    arch = match.get('architecture', 'ARCH-001')
    spec = match.get('specification', '')
    
    print(f"{arch}\n{spec}")

if __name__ == '__main__':
    main()
