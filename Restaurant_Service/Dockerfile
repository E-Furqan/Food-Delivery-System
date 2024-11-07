FROM golang:1.23.1 

WORKDIR /app

# Install PostgreSQL client and netcat (for checking the database status)
# Install PostgreSQL client and netcat (choose netcat-openbsd or netcat-traditional)
RUN apt-get update && apt-get install -y postgresql-client netcat-openbsd

COPY go.mod go.sum ./

RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

# Copy the migration script
COPY run_migrations.sh /app/run_migrations.sh

# Make the script executable
RUN chmod +x /app/run_migrations.sh
RUN chmod -R 755 /app/Migration
RUN chmod +x run_migrations.sh

RUN go build -o restaurantservice ./main.go

ENTRYPOINT ["/app/run_migrations.sh"]

CMD ["./restaurantservice"]