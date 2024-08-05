import json

from utils import client


# 1. Garantindo que o sistema está limpo
response_flow_definitions = client().get("/flows/definitions/")
assert response_flow_definitions.json() == []

response_flow_instances = client().get("/flows/instances/")
assert response_flow_instances.json() == []

response_task_definitions = client().get("/tasks/definitions/")
assert response_task_definitions.json() == []

flow_definition_data = {
    "name": "flow_1",
    "input": {
        "properties": {
            "firstName": {
                "type": "string",
                "description": "The person's first name.",
            },
            "lastName": {
                "type": "string",
                "description": "The person's last name.",
            },
        },
        "required": ["firstName", "lastName"],
    },
    "output": {
        "properties": {
            "fullName": {
                "type": "string",
                "description": "The person full name.",
            }
        },
        "required": ["fullName"],
    }
}
response_flow_definition = client().post("/flows/definitions/", json=flow_definition_data)
assert response_flow_definition.status_code == 201
flow_definition_ref_id = response_flow_definition.json()["reference_id"]

# 2. Creating task definitions
task_definition_data = {
    "name": "flow_1_task_1_add_age",
    "flow_definition_ref_id": flow_definition_ref_id,
    "input": {
        "properties": {
            "age": {
                "description": "Age in years which must be equal to or greater than zero.",
                "type": "integer",
                "minimum": 0
            }
        },
        "required": ["age"]
    },
    "output": {
        "properties": {
            "fullName": {
                "type": "string",
                "description": "The person full name.",
            },
            "age": {
                "description": "Age in years which must be equal to or greater than zero.",
                "type": "integer",
                "minimum": 0
            },
        },
        "required": ["description"]
    }
}
response_task_definition = client().post("/tasks/definitions/", json=task_definition_data)
task_definition_ref_id = response_task_definition.json()["reference_id"]
print(json.dumps(response_task_definition.json(), indent=2))

# 2.
task_definition_data = {
    "name": "flow_1_task_1_add_age",
    "flow_definition_ref_id": flow_definition_ref_id,

    "input": {
        "properties": {
            "age": {
                "description": "Age in years which must be equal to or greater than zero.",
                "type": "integer",
                "minimum": 0
            }
        },
        "required": ["age"]
    },
    "output": {
        "properties": {
            "fullName": {
                "type": "string",
                "description": "The person full name.",
            },
            "age": {
                "description": "Age in years which must be equal to or greater than zero.",
                "type": "integer",
                "minimum": 0
            },
        },
        "required": ["description"]
    }
}

