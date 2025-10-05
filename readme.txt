# Step 1
    go mod init github.com/suk-chanthea/ezra
    go mod tiny
# Build & start in dev
    docker compose -f docker-compose.override.yml up --build
# terminal of postgres
    docker exec -it ezra-postgres-dev bash
    psql -U postgres -d ezradb