# Run with full path to postman collection file as a single parameter, in Postman 2.1 format
# See "Postman to Wiki" at https://medium.com/@sreedharc/postman-to-wiki-e7d31c76db57 for more information
# Generates Markdown for Confluence Wiki from a given Collection file exported from Postman per Collections 2.1 schema

import json
import sys
import os
import re
from operator import itemgetter
from itertools import groupby

input_file = sys.argv[1]
working_dir = os.path.dirname(input_file)
serial_num = 1
output_dir = "generated"


def clean(s):
    # Python 3 replacement for string_escape decoding
    s = s.encode('latin1').decode('unicode_escape')
    s = s.replace('{{', '$').replace('}}', '').replace('\r\n', '')
    if s.startswith('"') and s.endswith('"'):
        return s[1:-1]
    return s


def clean_list(a_list):
    return [clean(i) for i in str(a_list).split(',')]


def escape_markup(value):
    return re.sub(r"([{\}\[\]])", r"\\\1", value)


def encode_full_url(url):
    return (
        re.sub(r"\s+", " ", url).replace("{{", "$").replace("}}", "")
            .replace("%", "%25").replace("'", "%27")
            .replace("{", "%7B").replace("}", "%7D")
            .replace("\"", "%22").replace(" ", "%20")
            .replace("[", "%5B").replace("]", "%5D")
    )


def format_raw_data(data):
    try:
        return 'json', json.dumps(json.loads(data), indent=2)
    except:
        return 'xml', data


# Generates markdown for the section describing Collection Variables
def gen_variables_section(list_of_variables):
    op = "## Variables Used in this Collection\n\n"
    op += "| Name | Description | Example |\n"
    op += "| ---- | ----------- | ------- |\n"
    for variable in sorted(list_of_variables, key=itemgetter('key')):
        op += f"| {variable['key']} | | {variable['value']} |\n"
    return op


# Generates markdown for request section
def gen_request_section(item):
    item = json.loads(json.dumps(item))
    # Heading, e.g. GET Users
    request_section = f"## **{item['request']['method']}** {item['name']}\n\n"

    # Description
    if "description" in item["request"]:
        request_section += f"{item['request']['description']}\n\n"

    request_section += "### Request\n\n"
    # URL format
    request_section += "#### Request URL\n\n"
    request_section += f"```\n{encode_full_url(item['request']['url']['raw'])}\n```\n\n"

    # Headers
    request_section += "#### Request Headers\n\n"
    request_section += "| Key | Value | Description |\n"
    request_section += "| --- | ----- | ----------- |\n"
    for header in item["request"]["header"]:
        description = header["description"] if "description" in header else " "
        request_section += f"| {header['key']} | {escape_markup(header['value'])} | {description} |\n"
    request_section += "\n"

    # Request Parameters
    add_param_section = False
    param_section = "#### Request Parameters\n\n"
    param_section += "| Parameter Type | Key | Value | Description |\n"
    param_section += "| -------------- | --- | ----- | ----------- |\n"

    url = item["request"]["url"]
    # Path Parameters
    if "path" in url:
        path_parameters = url["path"]
        for param in path_parameters:
            if "{{" in param:
                add_param_section = True
                param_section += f"| Path Parameter | {param.replace('{{', '$').replace('}}', '')} | | |\n"
    # Query String Parameters
    if "query" in url:
        qs_parameters = url["query"]
        for param in qs_parameters:
            add_param_section = True
            description = param["description"] if "description" in param else " "
            param_section += f"| Query String Parameter | {param['key']} | {escape_markup(param['value'])} | {description} |\n"

    print(param_section)

    if add_param_section:
        print("ADD PARAMETERS SECTION")
        request_section += param_section + "\n"

    if str(item["request"]["method"]) != "GET":
        request_section += "#### Request Payload\n\n"
        language = "json"
        # print(item["request"]["body"])
        if (not("body" in item["request"]) or item["request"]["body"] is None):
            # print("No Payload")
            payload = "{}"
        else:
            # print("CLEANED UP")
            if "raw" in item["request"]["body"]:
                try:
                    p = json.dumps(item["request"]["body"]["raw"])
                    x = clean(json.dumps(json.loads(p)))
                    # print(x)
                    if not x:
                        payload = "{}"
                    else:
                        language, data = format_raw_data(x)
                        payload = data
                except:
                    payload = item["request"]["body"]["raw"]
            else:
                payload = "{}"
        # print("PAYLOAD: \n" + payload)
        # Request payload as code snippet
        request_section += f"```{language}\n{payload}\n```\n\n"
    else:
        print("No payload for GET")

    return request_section


# Generates markdown for response section
def gen_response_section(item):
    response = json.loads(json.dumps(item["response"]))
    response_section = "### Response\n\n"
    # print(item["name"] + " - " + str(len(response)))
    if len(response) > 0:
        code = "200"
        if "code" in response[0]:
            code = str(response[0]["code"])
        else:
            print(response[0])

        status = "OK"
        if "status" in response[0]:
            status = str(response[0]["status"])
        else:
            print(response[0])

        response_section += f"#### Status: {code} {status}\n\n"

        language = "json"
        if "body" in response[0]:
            body = response[0]["body"]
            if body is None:
                body = "{}"
            # print("BODY: \n" + body)
            language, body = format_raw_data(body)
        else:
            body = " "
            print("Warning - response body is missing")

        # Response body as code snippet
        response_section += f"```{language}\n{body}\n```\n\n"
    return response_section


# Generates a section for the given resource
def gen_resource_section(name, group):
    global serial_num
    global working_dir
    global output_dir
    resource_sections = ""
    if not os.path.exists(os.path.join(working_dir, output_dir)):
        os.mkdir(os.path.join(working_dir, output_dir))
    for item in group:
        print("\t Processing " + str(item["request"]["method"]))
        resource_section = gen_request_section(item) + gen_response_section(item) + "\n"
        # Write to a separate output file
        # Sanitize the name by replacing invalid filename characters
        safe_name = name.replace("/", "").replace("\\", "-").replace(":", "-").replace("*", "-") \
                       .replace("?", "-").replace("\"", "-").replace("<", "-").replace(">", "-") \
                       .replace("|", "-")
        file_name = safe_name + "-" + str(item["request"]["method"]) + ".md"
        output_path = os.path.join(working_dir, output_dir, file_name)
        # Specify UTF-8 encoding when writing the file
        with open(output_path, 'w', encoding='utf-8') as output_file:
            output_file.write(resource_section)
        serial_num = serial_num + 1
        resource_sections += resource_section
    return resource_sections


# Read in the input as JSON
with open(input_file, encoding="utf-8") as json_data:  # Specify UTF-8 encoding
    input_data = json.load(json_data)
    # Extract all variables and print them out as a table
    v = gen_variables_section(input_data["variable"])
    
    # Now extract all requests ("items" in the collection) sorted by resource name ("name" in the collection)...
    items = sorted(input_data["item"], key=lambda x: (x['name'], x['request']['method']))
    
    # ... and generate documentation for supported methods by resource, in alphabetical order (GET, PATCH, POST, etc.)
    resource_full = "# API Documentation\n\n"
    resource_full += v + "\n"
    
    for name, group in groupby(items, key=lambda y: (y['name'])):
        print("Exporting: " + name)
        resource_full += gen_resource_section(name, group)
    
    # with open(os.path.join(working_dir, output_dir, "full.md"), 'w') as output_file:
    #     output_file.write(resource_full)
