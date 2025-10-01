Notification Hub 🔔


A lightweight real-time notification service built with Go (Gin) and Redis, packaged with Docker and deployed on AWS EC2 with Caddy as reverse proxy for HTTPS.

It provides:


- A WebSocket endpoint for client subscriptions.

- REST API endpoints to send and acknowledge notifications.

- A simple HTML client UI (client.html) to test connections.

- Continuous delivery with GitLab CI/CD → builds Docker images and pushes them to GitLab Container Registry.

- Secure deployment at https://notification.benexpense.lat.


---

🚀 Features

- WebSocket: real-time push notifications (/api/ws)

- Send notification: POST /api/notification/send

- Acknowledge notification: POST /api/notification/acknowledge

- Health check: GET /healthz

- UI Testing page: client.html is served at /


---

🛠️ Tech Stack

- Go + Gin — backend web framework

- Redis — message broker + persistence

- Docker & Docker Compose — containerized services

- GitLab CI/CD — builds & pushes images to registry

- Caddy — reverse proxy + TLS via Let’s Encrypt

- AWS EC2 + Route 53 — hosting and DNS


---

📂 Project Structure

	.
	├── main.go             # Gin server, routes, WebSocket
	├── controller/         # Request handlers
	├── client.html         # Simple browser testing UI
	├── Dockerfile          # Builds Go binary and includes UI
	├── docker-compose.yml  # Redis + Go + Caddy stack
	├── Caddyfile           # HTTPS + reverse proxy config
	└── .gitlab-ci.yml      # CI/CD pipeline config


---

⚡ API Endpoints

🔹 Health

	GET /healthz

🔹 WebSocket

	GET /api/ws?user_id=123&client_name=test

🔹 Send notification

	POST /api/notification/send
	Content-Type: application/json
	
	{
	  "user_id": "123",
	  "payload": { "text": "Hello world!" }
	}

Response:


	{"success":true,"message_id":"<unique-id>"}

🔹 Acknowledge notification

	POST /api/notification/acknowledge
	Content-Type: application/json
	
	{
	  "user_id": "123",
	  "message_ids": ["<message-id>"]
	}

Response:


	{"success":true}
