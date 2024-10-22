"""

"""
from utils import client


flow_instance_data = {
    "input_data": {
        "id": 1,
    },
}

flow_definition_ref_id = "f29c8333-2e2e-4690-90e1-aa86fdaa78ef"

response_flow_instance = client().post(f"/flows/instances/{flow_definition_ref_id}/start", json=flow_instance_data)
assert response_flow_instance.status_code == 201
