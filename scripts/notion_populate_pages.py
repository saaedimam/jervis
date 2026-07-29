import os
import sys
import json
import urllib.request
import urllib.error

def chunk_text(text, max_len=2000):
    """Chunk text to avoid Notion's 2000 character limit per block."""
    chunks = []
    while len(text) > max_len:
        split_idx = text.rfind('\n', 0, max_len)
        if split_idx == -1:
            split_idx = max_len
        chunks.append(text[:split_idx])
        text = text[split_idx:].lstrip('\n')
    if text:
        chunks.append(text)
    return chunks

def md_to_blocks(md_text):
    blocks = []
    lines = md_text.split('\n')
    current_text = []
    
    for line in lines:
        if line.startswith('# '):
            if current_text:
                for chunk in chunk_text('\n'.join(current_text)):
                    blocks.append({"object": "block", "type": "paragraph", "paragraph": {"rich_text": [{"type": "text", "text": {"content": chunk}}]}})
                current_text = []
            blocks.append({"object": "block", "type": "heading_1", "heading_1": {"rich_text": [{"type": "text", "text": {"content": line[2:]}}]}})
        elif line.startswith('## '):
            if current_text:
                for chunk in chunk_text('\n'.join(current_text)):
                    blocks.append({"object": "block", "type": "paragraph", "paragraph": {"rich_text": [{"type": "text", "text": {"content": chunk}}]}})
                current_text = []
            blocks.append({"object": "block", "type": "heading_2", "heading_2": {"rich_text": [{"type": "text", "text": {"content": line[3:]}}]}})
        elif line.startswith('### '):
            if current_text:
                for chunk in chunk_text('\n'.join(current_text)):
                    blocks.append({"object": "block", "type": "paragraph", "paragraph": {"rich_text": [{"type": "text", "text": {"content": chunk}}]}})
                current_text = []
            blocks.append({"object": "block", "type": "heading_3", "heading_3": {"rich_text": [{"type": "text", "text": {"content": line[4:]}}]}})
        else:
            current_text.append(line)
            
    if current_text:
        for chunk in chunk_text('\n'.join(current_text)):
            blocks.append({"object": "block", "type": "paragraph", "paragraph": {"rich_text": [{"type": "text", "text": {"content": chunk}}]}})
            
    return blocks

def append_blocks(api_key, page_id, blocks):
    req = urllib.request.Request(
        f"https://api.notion.com/v1/blocks/{page_id}/children",
        data=json.dumps({"children": blocks}).encode("utf-8"),
        headers={
            "Authorization": f"Bearer {api_key}",
            "Notion-Version": "2022-06-28",
            "Content-Type": "application/json"
        },
        method="PATCH"
    )
    try:
        with urllib.request.urlopen(req) as response:
            return True
    except urllib.error.URLError as e:
        print(f"Error appending blocks: {e}")
        return False

def main():
    if len(sys.argv) < 3:
        print("Usage: python notion_populate_pages.py <page_id> <markdown_file>")
        sys.exit(1)
        
    api_key = os.environ.get("NOTION_API_KEY")
    if not api_key:
        print("Missing NOTION_API_KEY")
        sys.exit(1)
        
    page_id = sys.argv[1]
    md_file = sys.argv[2]
    
    if not os.path.exists(md_file):
        print(f"File not found: {md_file}")
        sys.exit(1)
        
    with open(md_file, "r") as f:
        md_text = f.read()
        
    blocks = md_to_blocks(md_text)
    
    # Notion allows max 100 blocks per request
    for i in range(0, len(blocks), 100):
        batch = blocks[i:i+100]
        if not append_blocks(api_key, page_id, batch):
            print("Failed to append batch.")
            sys.exit(1)
            
    print(f"Successfully populated page {page_id}")

if __name__ == "__main__":
    main()
