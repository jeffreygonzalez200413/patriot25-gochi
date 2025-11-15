import requests

url = "http://127.0.0.1:8765/respond"

payload = {
    "userMessage": "i finished two tasks but i'm still tired",
    "state": {
        "mood": "sad",
        "personality": "judgmental",
        "completionRate": 0.8,
        "totalInteractions": 3,
    },
}

resp = requests.post(url, json=payload)
print(resp.status_code)
print(resp.json())
