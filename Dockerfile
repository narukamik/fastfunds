FROM golang:1.25.3

WORKDIR /app

COPY go.mod go.sum main.go ./
RUN go mod download && go mod verify

COPY docs/ docs/
COPY internal/ internal/

# Ensures that Swagger documentation is always updated
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init

RUN CGO_ENABLED=0 GOOS=linux go build -o fastfunds-api main.go

# Wait for Postregres to be up
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE 8080

CMD ["/wait-for-it.sh", "db:5432", "--", "/app/fastfunds-api"]

