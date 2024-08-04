import pytest


def test_create_success(cleanup, client, task_definition_body, flow_definition):
    response = client.get("/tasks/definitions/")
    response_body = response.json()
    assert len(response_body) == 0

    task_definition_body["flow_definition_ref_id"] = flow_definition["reference_id"]
    response = client.post("/tasks/definitions/", json=task_definition_body)
    response_body = response.json()

    assert response.status_code == 201


