FROM migrate/migrate:v4.15.2

COPY database /database

ENTRYPOINT ["sh", "-c", "migrate -path /database -database postgres://${USER}:${PASS}@${URL}:${PORT}/${DB}?sslmode=disable up"]