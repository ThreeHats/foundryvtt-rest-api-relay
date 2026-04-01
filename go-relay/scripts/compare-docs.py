#!/usr/bin/env python3
"""Compare Go-generated API docs against the Node-generated originals.

Usage:
    python3 scripts/compare-docs.py [--orig DIR] [--new DIR]

Defaults:
    --orig  ../../../public              (main repo)
    --new   ../../public                 (worktree public/)

Run from go-relay/:
    python3 scripts/compare-docs.py
"""

import argparse
import json
import os
import sys


def load_json(path):
    with open(path) as f:
        return json.load(f)


def compare_api_docs(orig_path, new_path):
    orig = load_json(orig_path)
    new = load_json(new_path)

    orig_eps = {f"{e['method']} {e['path']}": e for e in orig["endpoints"]}
    new_eps = {f"{e['method']} {e['path']}": e for e in new["endpoints"]}

    print("=" * 60)
    print("api-docs.json")
    print("=" * 60)
    print(f"  Original: {len(orig_eps)} endpoints")
    print(f"  New:      {len(new_eps)} endpoints")

    missing = sorted(set(orig_eps) - set(new_eps))
    added = sorted(set(new_eps) - set(orig_eps))

    ok = True
    if missing:
        ok = False
        print(f"\n  MISSING ({len(missing)}):")
        for m in missing:
            print(f"    \u274c {m}")
    else:
        print("  \u2705 All original endpoints present")

    if added:
        print(f"\n  Added ({len(added)}):")
        for a in added:
            print(f"    \u2795 {a}")

    # Param regressions on shared endpoints
    shared = sorted(set(orig_eps) & set(new_eps))
    param_issues = []
    for key in shared:
        o = orig_eps[key]
        n = new_eps[key]
        o_all = set(
            p["name"]
            for p in o.get("requiredParameters", []) + o.get("optionalParameters", [])
        )
        n_all = set(
            p["name"]
            for p in n.get("requiredParameters", []) + n.get("optionalParameters", [])
        )
        missing_params = o_all - n_all
        if missing_params:
            param_issues.append((key, missing_params))

    if param_issues:
        ok = False
        print(f"\n  Param regressions ({len(param_issues)}):")
        for key, mp in param_issues:
            print(f"    \u274c {key}: missing {mp}")
    else:
        print("  \u2705 No param regressions")

    return ok


def compare_openapi(orig_path, new_path):
    orig = load_json(orig_path)
    new = load_json(new_path)

    ot = set(t["name"] for t in orig.get("tags", []))
    nt = set(t["name"] for t in new.get("tags", []))
    mt = sorted(ot - nt)

    op = set(orig.get("paths", {}).keys())
    np = set(new.get("paths", {}).keys())
    mp = sorted(op - np)

    print("\n" + "=" * 60)
    print("openapi.json")
    print("=" * 60)

    ok = True
    if mt:
        ok = False
        print(f"  \u274c Missing tags: {mt}")
    else:
        print(f"  \u2705 All {len(ot)} tags present (new has {len(nt)})")

    if mp:
        ok = False
        print(f"  \u274c Missing paths: {mp}")
    else:
        print(f"  \u2705 All {len(op)} paths present (new has {len(np)})")

    return ok


def compare_asyncapi(orig_path, new_path):
    orig = load_json(orig_path)
    new = load_json(new_path)

    oc = set(orig.get("channels", {}).keys())
    nc = set(new.get("channels", {}).keys())
    mc = sorted(oc - nc)

    print("\n" + "=" * 60)
    print("asyncapi.json")
    print("=" * 60)
    print(f"  Channels: orig={len(oc)}, new={len(nc)}")

    ok = True
    if mc:
        # request/json is a known Node docgen bug (should be request/file-system)
        real_missing = [c for c in mc if c != "request/json"]
        if real_missing:
            ok = False
            print(f"  \u274c Missing channels: {real_missing}")
        if "request/json" in mc:
            print(
                "  \u26a0\ufe0f  'request/json' missing (known Node docgen bug — correct name is 'request/file-system')"
            )
    else:
        print(f"  \u2705 All channels present")

    return ok


def check_markdown(md_dir):
    print("\n" + "=" * 60)
    print("Markdown docs")
    print("=" * 60)

    if not os.path.isdir(md_dir):
        print(f"  \u274c Directory not found: {md_dir}")
        return False

    md_files = sorted(f for f in os.listdir(md_dir) if f.endswith(".md"))
    total_endpoints = 0
    total_examples = 0
    for f in md_files:
        with open(os.path.join(md_dir, f)) as fh:
            content = fh.read()
        endpoints = content.count("\n## ")
        examples = content.count("### Code Examples")
        total_endpoints += endpoints
        total_examples += examples

    print(f"  Files:     {len(md_files)}")
    print(f"  Endpoints: {total_endpoints} documented")
    print(f"  Examples:  {total_examples} code example sections")

    expected_files = [
        "auth.md", "canvas.md", "chat.md", "clients.md", "dnd5e.md",
        "effects.md", "encounter.md", "entity.md", "fileSystem.md",
        "macro.md", "roll.md", "scene.md", "scoped-keys.md", "search.md",
        "session.md", "sheet.md", "sheetImage.md", "structure.md",
        "utility.md", "websocket.md",
    ]
    missing_files = [f for f in expected_files if f not in md_files]
    if missing_files:
        print(f"  \u274c Missing files: {missing_files}")
        return False
    else:
        print(f"  \u2705 All expected markdown files present")
    return True


def main():
    parser = argparse.ArgumentParser(description="Compare Go-generated docs against Node originals")
    parser.add_argument("--orig", default=None, help="Path to original public/ dir (default: auto-detect)")
    parser.add_argument("--new", default=None, help="Path to new public/ dir (default: auto-detect)")
    args = parser.parse_args()

    # Auto-detect paths relative to this script
    script_dir = os.path.dirname(os.path.abspath(__file__))
    go_relay_dir = os.path.dirname(script_dir)  # go-relay/
    worktree_dir = os.path.dirname(go_relay_dir)  # worktree root

    if args.orig:
        orig_dir = args.orig
    else:
        # Try main repo
        main_repo = os.path.join(worktree_dir, "..", "..", "..", "public")
        if os.path.isdir(main_repo):
            orig_dir = main_repo
        else:
            orig_dir = os.path.join(worktree_dir, "public")

    if args.new:
        new_dir = args.new
    else:
        new_dir = os.path.join(worktree_dir, "public")

    orig_dir = os.path.abspath(orig_dir)
    new_dir = os.path.abspath(new_dir)
    md_dir = os.path.abspath(os.path.join(worktree_dir, "docs", "md", "api"))

    print(f"Original: {orig_dir}")
    print(f"New:      {new_dir}")
    print(f"Markdown: {md_dir}")
    print()

    all_ok = True
    all_ok &= compare_api_docs(
        os.path.join(orig_dir, "api-docs.json"),
        os.path.join(new_dir, "api-docs.json"),
    )
    all_ok &= compare_openapi(
        os.path.join(orig_dir, "openapi.json"),
        os.path.join(new_dir, "openapi.json"),
    )
    all_ok &= compare_asyncapi(
        os.path.join(orig_dir, "asyncapi.json"),
        os.path.join(new_dir, "asyncapi.json"),
    )
    all_ok &= check_markdown(md_dir)

    print("\n" + "=" * 60)
    if all_ok:
        print("\u2705 ALL CHECKS PASSED")
    else:
        print("\u274c SOME CHECKS FAILED")
    print("=" * 60)

    sys.exit(0 if all_ok else 1)


if __name__ == "__main__":
    main()
