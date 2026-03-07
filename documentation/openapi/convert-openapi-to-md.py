#!/usr/bin/env python3

import ast
import json
import re
import shutil
from pathlib import Path


ROOT_DIR = Path(__file__).resolve().parents[2]
OPENAPI_DIR = Path(__file__).resolve().parent
DOCS_API_DIR = ROOT_DIR / "docs" / "md" / "api"
PUBLIC_API_DOCS_PATH = ROOT_DIR / "public" / "api-docs.json"
PACKAGE_JSON_PATH = ROOT_DIR / "package.json"

HTTP_METHOD_ORDER = {
    "get": 0,
    "post": 1,
    "put": 2,
    "patch": 3,
    "delete": 4,
    "options": 5,
    "head": 6,
}

SOURCE_ORDER = {
    "path": 0,
    "query": 1,
    "body": 2,
    "header": 3,
    "cookie": 4,
}

SYSTEM_CONFIGS = [
    {
        "name": "core",
        "spec_path": OPENAPI_DIR / "openapi-v3-foundry-rest-api-relay-core.yaml",
        "title": "Core API Reference",
        "description": "Core FoundryVTT REST API Relay endpoints.",
        "skip_tags": {"Docs", "default"},
        "tag_order": [
            "Roll",
            "Entity",
            "Search",
            "Encounter",
            "Structure",
            "Macros",
            "Sheet",
            "Utilities",
            "Session",
            "FileSystem",
        ],
    },
    {
        "name": "dnd5e",
        "spec_path": OPENAPI_DIR / "openapi-v3-foundry-rest-api-relay-dnd5e.yaml",
        "title": "DnD5e API Reference",
        "description": "D&D 5e specific FoundryVTT REST API Relay endpoints.",
        "skip_tags": set(),
        "tag_order": ["DnD5e"],
    },
]

TAG_SLUG_OVERRIDES = {
    "DnD5e": "dnd5e",
    "FileSystem": "fileSystem",
    "Macros": "macro",
    "Utilities": "utility",
}

TAG_TITLE_OVERRIDES = {
    "DnD5e": "DnD5e",
    "FileSystem": "File System",
    "Macros": "Macro",
    "Utilities": "Utility",
}

TAG_SUMMARY_OVERRIDES = {
    "DnD5e": "D&D 5e specific endpoints.",
    "Encounter": "Combat encounter management endpoints.",
    "Entity": "Entity creation, update, deletion, and item transfer endpoints.",
    "FileSystem": "File upload, download, and filesystem browsing endpoints.",
    "Macros": "Macro lookup and execution endpoints.",
    "Roll": "Dice rolling and roll history endpoints.",
    "Search": "Search endpoints for Foundry content.",
    "Session": "Headless session lifecycle endpoints.",
    "Sheet": "Actor sheet rendering endpoints.",
    "Structure": "World structure, folder, and compendium endpoints.",
    "Utilities": "Miscellaneous utility endpoints.",
}


class YamlParser:
    def __init__(self, text: str):
        self.lines = text.splitlines()
        self.index = 0

    def parse(self):
        self._skip_blanks()
        if self.index >= len(self.lines):
            return None
        return self._parse_block(self._indent_of(self.lines[self.index]))

    def _parse_block(self, indent: int):
        self._skip_blanks()
        if self.index >= len(self.lines):
            return None

        current_indent = self._indent_of(self.lines[self.index])
        if current_indent < indent:
            return None
        if self._stripped(self.lines[self.index]).startswith("-"):
            return self._parse_sequence(current_indent)
        return self._parse_mapping(current_indent)

    def _parse_mapping(self, indent: int):
        mapping = {}
        while self.index < len(self.lines):
            self._skip_blanks()
            if self.index >= len(self.lines):
                break

            raw_line = self.lines[self.index]
            current_indent = self._indent_of(raw_line)
            stripped = self._stripped(raw_line)

            if current_indent < indent:
                break
            if current_indent != indent or stripped.startswith("-"):
                break

            key, value = self._split_key_value(stripped)
            self.index += 1

            if self._is_block_scalar(value):
                mapping[key] = self._parse_block_scalar(indent, value[0])
                continue

            if value == "":
                self._skip_blanks()
                if self.index < len(self.lines) and self._indent_of(self.lines[self.index]) > indent:
                    mapping[key] = self._parse_block(self._indent_of(self.lines[self.index]))
                else:
                    mapping[key] = None
                continue

            mapping[key] = self._parse_scalar(value)

        return mapping

    def _parse_sequence(self, indent: int):
        items = []
        while self.index < len(self.lines):
            self._skip_blanks()
            if self.index >= len(self.lines):
                break

            raw_line = self.lines[self.index]
            current_indent = self._indent_of(raw_line)
            stripped = self._stripped(raw_line)

            if current_indent < indent or current_indent != indent or not stripped.startswith("-"):
                break

            item_content = stripped[1:].lstrip()
            self.index += 1

            if item_content == "":
                self._skip_blanks()
                if self.index < len(self.lines) and self._indent_of(self.lines[self.index]) > indent:
                    items.append(self._parse_block(self._indent_of(self.lines[self.index])))
                else:
                    items.append(None)
                continue

            if self._looks_like_key_value(item_content):
                key, value = self._split_key_value(item_content)
                item = {}

                if self._is_block_scalar(value):
                    item[key] = self._parse_block_scalar(indent, value[0])
                elif value == "":
                    self._skip_blanks()
                    if self.index < len(self.lines) and self._indent_of(self.lines[self.index]) > indent:
                        item[key] = self._parse_block(self._indent_of(self.lines[self.index]))
                    else:
                        item[key] = None
                else:
                    item[key] = self._parse_scalar(value)

                self._skip_blanks()
                if self.index < len(self.lines) and self._indent_of(self.lines[self.index]) > indent:
                    nested = self._parse_block(self._indent_of(self.lines[self.index]))
                    if isinstance(nested, dict):
                        item.update(nested)
                items.append(item)
                continue

            items.append(self._parse_scalar(item_content))

        return items

    def _parse_block_scalar(self, parent_indent: int, style: str):
        content_lines = []
        base_indent = None

        while self.index < len(self.lines):
            raw_line = self.lines[self.index]
            if raw_line.strip() == "":
                if base_indent is not None:
                    content_lines.append("")
                self.index += 1
                continue

            current_indent = self._indent_of(raw_line)
            if current_indent <= parent_indent:
                break

            if base_indent is None:
                base_indent = current_indent

            content_lines.append(raw_line[base_indent:])
            self.index += 1

        if style == ">":
            return self._fold_block_scalar(content_lines)
        return "\n".join(content_lines)

    @staticmethod
    def _fold_block_scalar(lines):
        folded = []
        pending_blank = False
        for line in lines:
            if line == "":
                pending_blank = True
                continue
            if folded and not pending_blank:
                folded.append(" ")
            elif pending_blank and folded:
                folded.append("\n")
            folded.append(line)
            pending_blank = False
        return "".join(folded).strip()

    @staticmethod
    def _split_key_value(line: str):
        parts = line.split(":", 1)
        if len(parts) != 2:
            raise ValueError(f"Invalid YAML mapping line: {line}")
        key = parts[0].strip().strip("'\"")
        value = parts[1].strip()
        return key, value

    @staticmethod
    def _looks_like_key_value(text: str):
        return bool(
            re.match(r"^(?:[^:\s][^:]*|['\"][^'\"]+['\"]):(?:\s|$)", text)
        )

    @staticmethod
    def _is_block_scalar(value: str):
        return value in {"|", "|-", ">", ">-"}

    @staticmethod
    def _parse_scalar(value: str):
        lowered = value.lower()
        if lowered in {"null", "~"}:
            return None
        if lowered == "true":
            return True
        if lowered == "false":
            return False
        if value in {"[]", "[ ]"}:
            return []
        if value in {"{}", "{ }"}:
            return {}
        if re.fullmatch(r"-?\d+", value):
            try:
                return int(value)
            except ValueError:
                return value
        if re.fullmatch(r"-?\d+\.\d+", value):
            try:
                return float(value)
            except ValueError:
                return value
        if value.startswith("'") and value.endswith("'"):
            return value[1:-1].replace("''", "'")
        if value.startswith('"') and value.endswith('"'):
            try:
                return ast.literal_eval(value)
            except (SyntaxError, ValueError):
                return value[1:-1]
        if value.startswith(("[", "{")) and value.endswith(("]", "}")):
            try:
                return ast.literal_eval(value)
            except (SyntaxError, ValueError):
                return value
        return value

    @staticmethod
    def _indent_of(line: str):
        return len(line) - len(line.lstrip(" "))

    @staticmethod
    def _stripped(line: str):
        return line.strip()

    def _skip_blanks(self):
        while self.index < len(self.lines) and self.lines[self.index].strip() == "":
            self.index += 1


def extract_top_level_block_map(text: str):
    lines = text.splitlines(keepends=True)
    sections = []
    pattern = re.compile(r"^([A-Za-z0-9_]+):\s*$")

    for index, line in enumerate(lines):
        if line.startswith(" "):
            continue
        match = pattern.match(line.rstrip("\n"))
        if match:
            sections.append((match.group(1), index))

    block_map = {}
    for position, (name, start_index) in enumerate(sections):
        end_index = sections[position + 1][1] if position + 1 < len(sections) else len(lines)
        block_map[name] = "".join(lines[start_index:end_index])
    return block_map


def extract_child_mapping_blocks(block_text: str, indent: int):
    lines = block_text.splitlines(keepends=True)
    blocks = {}
    current_name = None
    current_lines = []

    for line in lines[1:]:
        if line.strip() == "":
            if current_name is not None:
                current_lines.append(line)
            continue

        current_indent = len(line) - len(line.lstrip(" "))
        stripped = line.strip()
        parts = stripped.split(":", 1)

        if current_indent == indent and len(parts) == 2 and parts[1].strip() == "":
            if current_name is not None:
                blocks[current_name] = "".join(current_lines)
            current_name = parts[0].strip().strip("'\"")
            current_lines = [line]
            continue

        if current_name is not None:
            current_lines.append(line)

    if current_name is not None:
        blocks[current_name] = "".join(current_lines)

    return blocks


def parse_paths_block(block_text: str):
    paths = {}
    for path_name, path_block in extract_child_mapping_blocks(block_text, 2).items():
        parsed = YamlParser("paths:\n" + path_block).parse()
        if not isinstance(parsed, dict):
            continue
        path_item = (parsed.get("paths") or {}).get(path_name)
        if isinstance(path_item, dict):
            paths[path_name] = path_item
    return paths


def load_openapi_spec(path: Path):
    text = path.read_text(encoding="utf-8")
    block_map = extract_top_level_block_map(text)
    spec = {}
    for section_name in ("info", "components", "tags", "paths"):
        block = block_map.get(section_name)
        if not block:
            continue
        if section_name == "paths":
            spec["paths"] = parse_paths_block(block)
            continue
        parsed = YamlParser(block).parse()
        if isinstance(parsed, dict):
            spec.update(parsed)
    return spec


def read_package_version():
    package_json = json.loads(PACKAGE_JSON_PATH.read_text(encoding="utf-8"))
    return package_json.get("version", "1.0.0")


def slugify_tag(tag_name: str):
    if tag_name in TAG_SLUG_OVERRIDES:
        return TAG_SLUG_OVERRIDES[tag_name]
    return re.sub(r"[^a-zA-Z0-9]+", "-", tag_name).strip("-").lower()


def display_tag_name(tag_name: str):
    return TAG_TITLE_OVERRIDES.get(tag_name, tag_name)


def normalize_description(text):
    if not text:
        return ""
    normalized = str(text).strip()
    if normalized.startswith("## "):
        normalized = normalized[3:].strip()
    normalized = re.sub(r"\s+", " ", normalized)
    return normalized.strip()


def escape_markdown_cell(value):
    text = "" if value is None else str(value)
    return text.replace("|", r"\|").replace("\n", " ").strip()


def resolve_ref(spec, ref):
    if not isinstance(ref, str) or not ref.startswith("#/"):
        return None

    node = spec
    for part in ref[2:].split("/"):
        if isinstance(node, dict):
            node = node.get(part)
        elif isinstance(node, list):
            node = node[int(part)]
        else:
            return None
    return node


def dereference_schema(spec, schema):
    if not isinstance(schema, dict):
        return schema
    if "$ref" in schema:
        resolved = resolve_ref(spec, schema["$ref"])
        if resolved is None:
            return schema
        return dereference_schema(spec, resolved)
    if "allOf" in schema:
        merged = {"type": "object", "properties": {}, "required": []}
        description_parts = []
        for part in schema["allOf"]:
            resolved_part = dereference_schema(spec, part)
            if not isinstance(resolved_part, dict):
                continue
            if resolved_part.get("description"):
                description_parts.append(str(resolved_part["description"]))
            if resolved_part.get("type") and resolved_part["type"] != "object":
                merged["type"] = resolved_part["type"]
            merged["properties"].update(resolved_part.get("properties", {}))
            merged["required"].extend(resolved_part.get("required", []))
        if schema.get("description"):
            description_parts.append(str(schema["description"]))
        if description_parts:
            merged["description"] = " ".join(description_parts)
        return merged
    return schema


def schema_type_name(spec, schema):
    resolved = dereference_schema(spec, schema)
    if not isinstance(resolved, dict):
        return "object"
    if resolved.get("type"):
        return str(resolved["type"])
    if "properties" in resolved:
        return "object"
    if "items" in resolved:
        return "array"
    return "object"


def schema_to_body_parameters(spec, schema, required=False):
    resolved = dereference_schema(spec, schema)
    if not isinstance(resolved, dict):
        return []

    schema_required = set(resolved.get("required", []))
    properties = resolved.get("properties")
    if isinstance(properties, dict) and properties:
        parameters = []
        for name, property_schema in properties.items():
            property_resolved = dereference_schema(spec, property_schema)
            description = ""
            if isinstance(property_resolved, dict):
                description = normalize_description(property_resolved.get("description", ""))
            parameters.append(
                {
                    "name": name,
                    "type": schema_type_name(spec, property_schema),
                    "required": name in schema_required or required,
                    "sources": {"body"},
                    "description": description,
                }
            )
        return parameters

    description = normalize_description(resolved.get("description", ""))
    if (
        resolved.get("type") == "object"
        and resolved.get("additionalProperties") is True
        and not description
    ):
        return []
    return [
        {
            "name": "body",
            "type": schema_type_name(spec, resolved),
            "required": required,
            "sources": {"body"},
            "description": description,
        }
    ]


def extract_operation_parameters(spec, operation):
    combined = {}
    auth_header_present = False

    def upsert(param):
        nonlocal auth_header_present
        key = param["name"]
        if key.lower() == "x-api-key":
            auth_header_present = True
            return
        existing = combined.get(key)
        if existing is None:
            combined[key] = param
            return
        if param["sources"] == {"header"} and existing["sources"] != {"header"}:
            return
        if existing["sources"] == {"header"} and param["sources"] != {"header"}:
            existing["sources"] = set(param["sources"])
        else:
            existing["sources"].update(param["sources"])
        existing["required"] = existing["required"] or param["required"]
        if len(param["description"]) > len(existing["description"]):
            existing["description"] = param["description"]
        if existing["type"] == "object" and param["type"] != "object":
            existing["type"] = param["type"]

    for parameter in operation.get("parameters", []) or []:
        resolved = dereference_schema(spec, parameter)
        if not isinstance(resolved, dict):
            continue
        location = str(resolved.get("in", "query"))
        schema = resolved.get("schema", {})
        upsert(
            {
                "name": str(resolved.get("name", "")),
                "type": schema_type_name(spec, schema),
                "required": bool(resolved.get("required", location == "path")),
                "sources": {location},
                "description": normalize_description(
                    resolved.get("description")
                    or (schema.get("description") if isinstance(schema, dict) else "")
                    or ""
                ),
            }
        )

    request_body = operation.get("requestBody")
    if isinstance(request_body, dict):
        required = bool(request_body.get("required", False))
        content = request_body.get("content", {}) or {}
        media = None
        for media_type in ("application/json", "multipart/form-data", "application/x-www-form-urlencoded"):
            if media_type in content:
                media = content[media_type]
                break
        if media is None and content:
            first_key = next(iter(content))
            media = content[first_key]
        if isinstance(media, dict) and "schema" in media:
            for body_parameter in schema_to_body_parameters(spec, media["schema"], required=required):
                upsert(body_parameter)

    parameters = list(combined.values())
    parameters.sort(
        key=lambda item: (
            0 if item["required"] else 1,
            min(SOURCE_ORDER.get(source, 99) for source in item["sources"]),
            item["name"].lower(),
        )
    )
    return parameters, auth_header_present


def select_response(operation):
    responses = operation.get("responses", {}) or {}
    if not isinstance(responses, dict) or not responses:
        return None

    for status_code in sorted(responses.keys()):
        if str(status_code).startswith("2"):
            return responses[status_code]
    first_key = next(iter(responses))
    return responses[first_key]


def response_summary(spec, operation):
    response = select_response(operation)
    if not isinstance(response, dict):
        return {"type": "object", "description": "Successful response"}

    content = response.get("content", {}) or {}
    schema = None
    if "application/json" in content:
        schema = content["application/json"].get("schema")
    elif content:
        first_key = next(iter(content))
        schema = content[first_key].get("schema")

    return {
        "type": schema_type_name(spec, schema) if schema else "object",
        "description": normalize_description(response.get("description", "Successful response")) or "Successful response",
    }


def display_operation_path(path: str, system_name: str):
    prefix = f"/{system_name}"
    if system_name != "core" and path.startswith(prefix):
        trimmed = path[len(prefix):]
        return trimmed if trimmed.startswith("/") else f"/{trimmed}"
    return path


def sources_to_text(sources):
    ordered = sorted(sources, key=lambda source: SOURCE_ORDER.get(source, 99))
    return ", ".join(ordered)


def extract_operations(spec, system_name: str):
    operations = []
    for path, path_item in (spec.get("paths", {}) or {}).items():
        if not isinstance(path_item, dict):
            continue
        for method, operation in path_item.items():
            if method.lower() not in HTTP_METHOD_ORDER or not isinstance(operation, dict):
                continue
            tags = operation.get("tags", []) or []
            primary_tag = tags[0] if tags else "default"
            parameters, auth_header_present = extract_operation_parameters(spec, operation)
            returns = response_summary(spec, operation)
            security_rules = operation.get("security") or []
            is_noauth = any("noauthAuth" in security for security in security_rules)
            operations.append(
                {
                    "system": system_name,
                    "tag": primary_tag,
                    "method": method.upper(),
                    "method_order": HTTP_METHOD_ORDER[method.lower()],
                    "path": path,
                    "display_path": display_operation_path(path, system_name),
                    "description": normalize_description(operation.get("description") or operation.get("summary") or ""),
                    "parameters": parameters,
                    "returns": returns,
                    "requires_auth": path not in {"/api/status", "/api/docs"} and (auth_header_present or not is_noauth),
                }
            )
    operations.sort(key=lambda item: (item["tag"].lower(), item["display_path"].lower(), item["method_order"]))
    return operations


def render_resource_markdown(resource_title: str, operations):
    lines = [
        "---",
        f"title: {resource_title}",
        f"sidebar_label: {resource_title}",
        "---",
        "",
        f"# {resource_title}",
        "",
    ]

    for index, operation in enumerate(operations):
        lines.append(f"## {operation['method']} {operation['display_path']}")
        lines.append("")

        if operation["description"]:
            lines.append(operation["description"])
            lines.append("")

        if operation["parameters"]:
            lines.append("### Parameters")
            lines.append("")
            lines.append("| Name | Type | Required | Source | Description |")
            lines.append("|------|------|----------|--------|-------------|")
            for parameter in operation["parameters"]:
                required_mark = "Yes" if parameter["required"] else "No"
                lines.append(
                    "| {name} | {type_name} | {required} | {source} | {description} |".format(
                        name=escape_markdown_cell(parameter["name"]),
                        type_name=escape_markdown_cell(parameter["type"]),
                        required=required_mark,
                        source=escape_markdown_cell(sources_to_text(parameter["sources"])),
                        description=escape_markdown_cell(parameter["description"]),
                    )
                )
            lines.append("")

        lines.append("### Returns")
        lines.append("")
        lines.append(
            f"**{escape_markdown_cell(operation['returns']['type'])}** - {escape_markdown_cell(operation['returns']['description'])}"
        )
        lines.append("")

        if index < len(operations) - 1:
            lines.append("---")
            lines.append("")

    return "\n".join(lines).rstrip() + "\n"


def render_system_index(system_config, resources):
    lines = [
        "---",
        f"title: {system_config['title']}",
        "sidebar_label: Overview",
        "---",
        "",
        f"# {system_config['title']}",
        "",
        system_config["description"],
        "",
        "## Available Endpoint Groups",
        "",
    ]

    for resource in resources:
        lines.append(
            f"- [{resource['title']}](/api/{system_config['name']}/{resource['slug']}/) - {resource['summary']}"
        )

    lines.append("")
    return "\n".join(lines)


def render_root_index():
    return """---
id: api-reference
title: API Reference
sidebar_label: API Reference
---

# FoundryVTT REST API Reference

This documentation is generated from the OpenAPI source files in `documentation/openapi`.

## Available Systems

- [Core](/api/core/) - Core FoundryVTT REST API Relay endpoints.
- [DnD5e](/api/dnd5e/) - D&D 5e specific endpoints.

## Authentication

Most data endpoints require authentication with an API key. Provide your API key in the `x-api-key` header with each request.
"""


def category_json(label: str, position: int):
    return json.dumps({"label": label, "position": position}, indent=2) + "\n"


def resource_summary(operations):
    if operations:
        override = TAG_SUMMARY_OVERRIDES.get(operations[0]["tag"])
        if override:
            return override
    descriptions = [operation["description"] for operation in operations if operation["description"]]
    if descriptions:
        return descriptions[0]
    return f"{len(operations)} documented endpoint(s)."


def build_api_docs_json(version: str, operations):
    endpoints = []
    for operation in sorted(operations, key=lambda item: (item["path"].lower(), item["method_order"])):
        required_parameters = []
        optional_parameters = []
        for parameter in operation["parameters"]:
            entry = {
                "name": parameter["name"],
                "type": parameter["type"],
                "description": parameter["description"],
                "location": sources_to_text(parameter["sources"]),
            }
            if parameter["required"]:
                required_parameters.append(entry)
            else:
                optional_parameters.append(entry)

        endpoint = {
            "method": operation["method"],
            "path": operation["path"],
            "description": operation["description"] or operation["returns"]["description"],
            "requiredParameters": required_parameters,
            "optionalParameters": optional_parameters,
        }
        if not operation["requires_auth"]:
            endpoint["authentication"] = False
        endpoints.append(endpoint)

    return {
        "version": version,
        "baseUrl": "https://your-relay-server.com",
        "authentication": {
            "required": True,
            "headerName": "x-api-key",
            "description": "API key must be included in the x-api-key header for all endpoints except /api/status and /api/docs",
        },
        "endpoints": endpoints,
    }


def reset_docs_directory():
    DOCS_API_DIR.mkdir(parents=True, exist_ok=True)
    for entry in DOCS_API_DIR.iterdir():
        if entry.name in {"index.md", "_category_.json"}:
            continue
        if entry.is_dir():
            shutil.rmtree(entry)
        else:
            entry.unlink()


def write_text(path: Path, content: str):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")


def main():
    version = read_package_version()
    reset_docs_directory()

    write_text(DOCS_API_DIR / "_category_.json", category_json("API Reference", 7))
    write_text(DOCS_API_DIR / "index.md", render_root_index())

    all_operations = []

    for position, system_config in enumerate(SYSTEM_CONFIGS, start=1):
        spec = load_openapi_spec(system_config["spec_path"])
        operations = extract_operations(spec, system_config["name"])
        all_operations.extend(operations)

        system_dir = DOCS_API_DIR / system_config["name"]
        write_text(system_dir / "_category_.json", category_json(system_config["title"].replace(" API Reference", ""), position))

        visible_resources = []
        grouped_operations = {}
        for operation in operations:
            grouped_operations.setdefault(operation["tag"], []).append(operation)

        ordered_tags = list(system_config["tag_order"]) + [
            tag for tag in grouped_operations.keys() if tag not in system_config["tag_order"]
        ]

        for tag in ordered_tags:
            if tag in system_config["skip_tags"] or tag not in grouped_operations:
                continue

            resource_ops = grouped_operations[tag]
            resource_title = display_tag_name(tag)
            resource_slug = slugify_tag(tag)
            write_text(system_dir / f"{resource_slug}.md", render_resource_markdown(resource_title, resource_ops))
            visible_resources.append(
                {
                    "title": resource_title,
                    "slug": resource_slug,
                    "summary": resource_summary(resource_ops),
                }
            )

        write_text(system_dir / "index.md", render_system_index(system_config, visible_resources))

    PUBLIC_API_DOCS_PATH.parent.mkdir(parents=True, exist_ok=True)
    PUBLIC_API_DOCS_PATH.write_text(
        json.dumps(build_api_docs_json(version, all_operations), indent=2),
        encoding="utf-8",
    )

    print(f"Generated OpenAPI markdown in {DOCS_API_DIR}")
    print(f"Generated API JSON in {PUBLIC_API_DOCS_PATH}")


if __name__ == "__main__":
    main()
