# Step 1
    go mod init github.com/suk-chanthea/ezra
    go mod tiny
# Build & start in dev
### build & up
    docker compose -f docker-compose.override.yml up --build -d
### only up
    docker compose -f docker-compose.override.yml up -d
# terminal of postgres
    docker exec -it ezra-postgres-dev bash
    psql -U postgres -d ezradb