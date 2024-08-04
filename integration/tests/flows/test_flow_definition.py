import pytest


def test_create_success(cleanup, client, flow_definition_body):
    response = client.post("/flows/definitions/", json=flow_definition_body)
    response_body = response.json()

    expected_body =  {
        'id': 2,
        'reference_id': 'b375fd71-1f1c-4796-9fb6-0a35ea729f4c',
        'created_at': '2024-08-04T15:01:48.779398264-03:00',
        'updated_at': '2024-08-04T15:01:48.779398324-03:00',
        'name': 'media',
        'description': '',
        'status': 'inactive', 'version': 2,
        'input_schema': {
            '$id': f'https://fs.dev.makecodes.dev/schemas/flow_definitions/{response_body["reference_id"]}/input.json',
            '$schema': 'https://json-schema.org/draft/2020-12/schema',
            'additionalProperties': False,
            'properties': {
                'age': {'description': 'Age in years which must be equal to or greater than zero.', 'minimum': 0, 'type': 'integer'},
                'firstName': {'description': "The person's first name.", 'type': 'string'},
                'lastName': {'description': "The person's last name.", 'type': 'string'},
            },
            'required': ['firstName', 'lastName'],
            'title': 'Flow Definition',
            'type': 'object',
        },
        'output_schema': {
            '$id': f'https://fs.dev.makecodes.dev/schemas/flow_definitions/{response_body["reference_id"]}/output.json',
            '$schema': 'https://json-schema.org/draft/2020-12/schema',
            'additionalProperties': False,
            'properties': {
                'fullName': {'description': 'The person full name.', 'type': 'string'},
            },
            'required': ['fullName'],
            'title': 'Flow Definition',
            'type': 'object',
        },
        'metadata': {},
    }

    assert response.status_code == 201
    assert response_body['name'] == flow_definition_body['name']
    assert response_body['input_schema'] == expected_body['input_schema']
    assert response_body['output_schema'] == expected_body['output_schema']

    response = client.get("/flows/definitions/")
    assert response.status_code == 200
    assert len(response.json()) == 1

    first = response.json()[0]

    response = client.get(f"/flows/definitions/{first['reference_id']}")
    assert response.status_code == 200
    assert response.json() == first

    response = client.delete(f"/flows/definitions/{first['reference_id']}")
    assert response.status_code == 204

    response = client.get("/flows/definitions/")
    assert response.status_code == 200
    assert len(response.json()) == 0
