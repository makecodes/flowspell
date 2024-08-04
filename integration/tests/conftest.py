import os

import httpx
import pytest
from faker import Faker


@pytest.fixture
def base_url():
    return os.getenv("FLOWSPELL_TEST_URL", "http://localhost:8266")


@pytest.fixture
def client(base_url):
    return httpx.Client(base_url=base_url)


@pytest.fixture
def faker():
    return Faker()

@pytest.fixture
def cleanup(client):
    yield
    response = client.get('/tasks/definitions/')
    response_body = response.json()
    for task_definition in response_body:
        client.delete(f"/tasks/definitions/{task_definition['reference_id']}")

    response = client.get("/flows/definitions/")
    response_body = response.json()
    for flow_definition in response_body:
        client.delete(f"/flows/definitions/{flow_definition['reference_id']}")


@pytest.fixture
def flow_definition_body(faker):
    return {
        "name": faker.word(),
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
                "age": {
                    "description": "Age in years which must be equal to or greater than zero.",
                    "type": "integer",
                    "minimum": 0,
                }
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

@pytest.fixture
def task_definition_body(faker):
    return {
        "name": faker.word(),
        "flow_definition_ref_id": None,
        "input": {
            "properties": {
                "firstName": {
                    "type": "string",
                    "description": "The person's first name."
                },
                "lastName": {
                    "type": "string",
                    "description": "The person's last name."
                },
                "age": {
                    "description": "Age in years which must be equal to or greater than zero.",
                    "type": "integer",
                    "minimum": 0
                }
            },
            "required": ["firstName", "lastName"]
        },
        "output": {
            "properties": {
                "fullName": {
                    "type": "string",
                    "description": "The person full name."
                }
            },
            "required": ["fullName"]
        }
    }


@pytest.fixture
def flow_definition(client, flow_definition_body):
    response = client.post("/flows/definitions/", json=flow_definition_body)
    return response.json()
