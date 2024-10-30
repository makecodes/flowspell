"""
This script is used to start a flow instance.
"""
import os
import json

from utils import client

if not os.path.exists("task_definition.json"):
    raise Exception("task_definition.json not found")

with open("task_definition.json", "r") as f:
    task_definition_data = json.load(f)

flow_instance_data = {
    "input_data": {
        "id": 1,
    },
}

flow_definition_ref_id = task_definition_data["flow_definition_ref_id"]

response_flow_instance = client().post(f"/flows/instances/{flow_definition_ref_id}/start", json=flow_instance_data)
assert response_flow_instance.status_code == 201
