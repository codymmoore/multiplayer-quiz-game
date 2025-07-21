# multiplayer-quiz-game
Multiplayer quiz game

![quizchief-client-server-diagram drawio](https://github.com/user-attachments/assets/9a24aa44-c7b8-40ee-9041-f2031d5de19f)

## Dependencies
- [go-chi/chi](https://github.com/go-chi/chi)
- [go-chi/jwtauth](https://github.com/go-chi/jwtauth)
- [sqlc](https://sqlc.dev)
- [pressly/goose](https://github.com/pressly/goose)

## Local Development
### Running Services Locally
- From the `server` directory, run the following:
```
docker compose up --build [-d] # -d runs services in the background
```

This will:
- Build and start all services defined in `docker-compose.yaml`
- Rebuild images if there are code or configuration changes

### Shutting Down Local Environment
- From the `server` directory, run the following:
```
docker compose down [--volumes] # --volumes will remove persistent data as well as stopping containers
```

### Troubleshooting
- View running containers:
```
docker compose ps
```
- View logs:
```
docker compose logs <service_name>
```
- Full reset:
```
docker compose down --volumes --rmi all
```
