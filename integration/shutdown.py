from utils import client

# 1. Percorrendo os registros para excluir
response_task_definitions = client().get("/tasks/definitions/")
for task_definition in response_task_definitions.json():
    client().delete(f"/tasks/definitions/{task_definition['reference_id']}")

response_flow_definitions = client().get("/flows/definitions/")
for flow_definition in response_flow_definitions.json():
    client().delete(f"/flows/definitions/{flow_definition['reference_id']}")

# response_flow_instances = client().get("/flows/instances/")
# for flow_instance in response_flow_instances.json():
#     client().delete(f"/flows/instances/{flow_instance['reference_id']}")
