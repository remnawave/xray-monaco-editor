import json
import re
import sys

from typing import Iterator, TypedDict


class RawType(TypedDict):
    title: str
    description: str
    raw_properties: list[dict]


class JsonschemaType(TypedDict):
    title: str
    description: str
    # enable markdown in monaco
    # https://github.com/microsoft/monaco-editor/issues/1816
    markdownDescription: str
    properties: dict[str, dict]
    additionalProperties: bool


KNOWN_BAD_RESOLVES = (
    "FakeDnsObject",
    "metricsObject",
    "TransportObject",
    "noiseObject",
    "DnsServerObject",
    "xhttpSettings",
    "XHTTPObject",
    "PingConfigObject",
    "XHTTP: Beyond REALITY",
    "CostObject",
    "SockoptObject",
    "quicParamsObject",
)
USED_OBJECTS = set()

# A heading only starts a new definition if it is a top-level (`##`) section or
# if its title looks like a type name. Upstream splits some objects into
# descriptive `###` subsections (e.g. StreamSettingsObject into "Способы
# передачи" / "传输方式" etc.), and those must keep filling the enclosing
# object instead of stealing its properties into a definition nobody refs.
TYPE_NAME_RE = re.compile(r"^[A-Za-z][A-Za-z0-9_.\-]*$")

# Properties the docs renamed but that we want to keep documented under the old
# name. Keyed by object title, so unrelated objects with a same-named property
# (e.g. shadowsocks' `method`) are untouched.
PROPERTY_RENAMES = {
    # xray-core accepts both `method` and `network` (`method` wins when both are
    # set), but every existing config uses `network`, so only offer that one.
    "StreamSettingsObject": {"method": "network"},
}


def clean_prefix(line: str) -> bool:
    """
    Checks that the line:
      - starts with '>'
      - if a '`' appears **before** the first ':', there are only spaces between '>' and that '`';
      - otherwise returns True.
    """
    if not line.startswith(">"):
        return False

    body = line[1:]
    idx_tick = body.find("`")
    idx_colon = body.find(":")
    if idx_colon == -1:
        idx_colon = body.find("：")

    if 0 <= idx_tick < idx_colon:
        return all(ch == " " for ch in body[:idx_tick])

    return True


def starts_definition(hashes: str, title: str) -> bool:
    """
    Whether a markdown heading introduces a new definition rather than a
    subsection of the definition we are currently inside.

    `##` always does — that is where objects, and the root definition itself,
    live. Deeper headings only do when they name a type (`RuleObject`, `Peers`,
    `header-custom`); anything prose-shaped belongs to its parent object.
    """
    return len(hashes) == 2 or bool(TYPE_NAME_RE.match(title))


def finalize(current_obj: RawType) -> JsonschemaType:
    description = current_obj["description"]

    return {
        "title": current_obj["title"],
        "description": description,
        "markdownDescription": description,
        "properties": {x["name"]: x for x in current_obj["raw_properties"]},
        # turn off additionalProperties so that monaco will warn on
        # unknown properties. xray does allow for unknown
        # properties but most likely, setting them is a mistake. we
        # only do this if we have any props ourselves, otherwise
        # there is no point.
        "additionalProperties": not current_obj["raw_properties"],
    }


def parse(stdin: Iterator[str]) -> Iterator[JsonschemaType]:
    current_obj: RawType | None = None

    for line in stdin:
        heading = re.match(r"^(#{2,})\s+(.*?)\s*$", line)
        if heading and current_obj and not starts_definition(*heading.groups()):
            # A descriptive subsection of the object we are already in. Fold it
            # into the description so the properties below it keep landing on
            # the enclosing object.
            heading = None

        if heading:
            if current_obj:
                yield finalize(current_obj)

            current_obj = {
                "title": heading.group(2),
                "description": "",
                "raw_properties": [],
            }
        elif line.startswith("> ") and (":" in line or "：" in line) and current_obj:
            if ":" in line:
                name, ty = line[2:].split(":", 1)
            else:
                name, ty = line[2:].split("：", 1)

            if name == "Tony":
                continue

            if not clean_prefix(line):
                continue

            name = name.strip(" `")
            name = PROPERTY_RENAMES.get(current_obj["title"], {}).get(name, name)

            try:
                type_info = parse_type(ty)
            except Exception:
                # Not an actual property definition but a descriptive bullet
                # that happens to use the same "> `name`: text" syntax (e.g.
                # finalmask's "> `dns`: подделка под ..." enum docs). Treat it
                # as prose and fold it into the current description instead of
                # crashing the whole build.
                if current_obj["raw_properties"]:
                    current_obj["raw_properties"][-1]["description"] += line
                    current_obj["raw_properties"][-1]["markdownDescription"] += line
                else:
                    current_obj["description"] += line
                continue

            current_obj["raw_properties"].append(
                {
                    "name": name,
                    "description": "",
                    "markdownDescription": "",
                    **type_info,
                }
            )
        elif current_obj:
            if current_obj["raw_properties"]:
                current_obj["raw_properties"][-1]["description"] += line
                current_obj["raw_properties"][-1]["markdownDescription"] += line
            else:
                current_obj["description"] += line

    # Whatever section the input ended on still has to be emitted, otherwise the
    # last definition in the stream is silently dropped.
    if current_obj:
        yield finalize(current_obj)


def parse_type(input: str) -> dict:
    input = (
        input.replace('<Badge text="WIP" type="warning"/>', "")
        .replace('<Badge text="BETA" type="warning"/>', "")
        .replace("<br>", "")
        .strip()
    )

    if not input:
        return {}

    if input.startswith("\\[") and input.endswith("\\]"):
        return {"type": "array", "items": parse_type(input[2:-2])}

    # Add handling for incomplete escaped arrays
    if input.startswith("\\[") and input.endswith("\\"):
        # Extract the type between \[ and \
        inner_type = input[2:-1]
        return {"type": "array", "items": parse_type(inner_type)}

    if input.startswith("[") and input.endswith("]"):
        return {"type": "array", "items": parse_type(input[1:-1])}

    if (input.startswith("[") and input.endswith(")")) or input.endswith("Object"):
        name = input.split("]")[0].strip("[]")
        if name in KNOWN_BAD_RESOLVES:
            # If there is a dangling reference, monaco editor will turn off
            # all inline validation markers, as the root object has a warning.
            # So we catch all dangling references here and replace them with
            # object.
            return {"type": "object"}
        else:
            USED_OBJECTS.add(name)
            return {"$ref": f"#/definitions/{name}"}

    if input in ("true", "false", "true | false", "bool"):
        return {"type": "boolean"}

    if " | " in input:
        return {"anyOf": [parse_type(x) for x in input.split(" | ")]}

    if input in ("address", "address_port", "CIDR"):
        return {"type": "string"}

    if input in ("string", "number"):
        return {"type": input}

    if input == "int":
        return {"type": "integer"}

    if input.startswith("map"):
        return {"type": "object"}

    if input.startswith('"') and input.endswith('"'):
        return {"const": input[1:-1]}

    if input.startswith("a list of"):
        return {}

    if input == "string array" or input == "array" or input == "list":
        return {"type": "array", "items": {"type": "string"}}

    if input.startswith("string, any of"):
        return {"type": "string"}

    if input == "object":
        return {}

    if input == "float number":
        return {"type": "number"}

    if input == "{}":
        return {"type": "object"}

    if input == "struct":
        return {"type": "object"}

    # Handle inline object types like {"port": string, "interval": number}
    if input.startswith("{") and input.endswith("}"):
        return {"type": "object"}

    # Handle empty or whitespace-only input
    if not input.strip():
        return {}

    # Handle "null" type
    if input == "null":
        return {"type": "null"}

    # Handle dash-separated identifier values (e.g. "header-custom", "mkcp-original")
    if re.match(r'^[a-zA-Z0-9][-a-zA-Z0-9]*$', input):
        return {"const": input}

    raise Exception(f"Unknown type: '{input}'")


def main():
    # root definition
    root_definition = sys.argv[1]

    definitions = {}
    for definition in parse(sys.stdin):
        key = definition["title"]
        if key in definitions:
            # Handle multiple instances of
            # InboundConfigurationObject/OutboundConfigurationObject
            if "anyOf" not in definitions[key]:
                definitions[key] = {"anyOf": [definitions[key]]}
            definitions[key]["anyOf"].append(definition)
        else:
            definitions[key] = definition

    for name in USED_OBJECTS:
        assert name in definitions, f"Cannot resolve {name}, add to KNOWN_BAD_RESOLVES?"

    #     schema = {
    #         "$schema": "http://json-schema.org/draft-07/schema#",
    #         "$ref": "#/definitions/Основные модули конфигурации",
    #         "definitions": definitions
    #     }

    schema = {
        "$schema": "http://json-schema.org/draft-07/schema#",
        "$ref": f"#/definitions/{root_definition}",
        "definitions": definitions,
    }

    print(json.dumps(schema, indent=2, ensure_ascii=False))


if __name__ == "__main__":
    main()
