from utils import client


queue_data = {
    "worker_name": "worker_1",
}
queue_response = client().post("/tasks/queue", json=queue_data)
print(queue_response.json())
