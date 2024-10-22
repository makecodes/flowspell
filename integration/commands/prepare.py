"""
This script is used to populate the system with data for testing purposes.
It is used to create flow definitions, task definitions, and start a flow instance.
"""
from utils import client


pre_create = True
if pre_create:
    # 1. Garantindo que o sistema está limpo
    response_flow_definitions = client().get("/flows/definitions/")
    assert response_flow_definitions.json() == []

    response_flow_instances = client().get("/flows/instances/")
    assert response_flow_instances.json() == []

    response_task_definitions = client().get("/tasks/definitions/")
    assert response_task_definitions.json() == []

    flow_definition_data = {
        "name": "flow_collect_user_data",
        "status": "active",
        "input": {
            "properties": {
                "id": {
                    "type": "integer",
                    "description": "The person's identifier.",
                },
            },
            "required": ["id"],
        },
        "output": {
            "properties": {
                "id": {
                    "type": "integer",
                    "description": "The person's identifier.",
                },
                "firstName": {
                    "type": "string",
                    "description": "The person first name.",
                },
                "lastName": {
                    "type": "string",
                    "description": "The person last name.",
                },
                "age": {
                    "description": "Age in years which must be equal to or greater than zero.",
                    "type": "integer",
                    "minimum": 0
                },
            },
            "required": ["firstName", "lastName", "age", "id"],
        }
    }
    response_flow_definition = client().post("/flows/definitions/", json=flow_definition_data)
    assert response_flow_definition.status_code == 201
    flow_definition_ref_id = response_flow_definition.json()["reference_id"]

    # 2. Creating task definitions
    task_definition_data = {
        "name": "flow_1_task_1_add_name",
        "flow_definition_ref_id": flow_definition_ref_id,
        "flow_definition_id": response_flow_definition.json()["id"],
        "input": {
            "properties": {},
            "required": [],
        },
        "output": {
            "properties": {
                "firstName": {
                    "type": "string",
                    "description": "The person first name.",
                },
                "lastName": {
                    "type": "string",
                    "description": "The person last name.",
                },
                "age": {
                    "description": "Age in years which must be equal to or greater than zero.",
                    "type": "integer",
                    "minimum": 0
                },
            },
            "required": ["firstName", "lastName", "age"],
        }
    }
    response_task_definition = client().post("/tasks/definitions/", json=task_definition_data)
    task_def_ref_id = response_task_definition.json()["reference_id"]

    print("Task definition created: ", response_task_definition.json())
