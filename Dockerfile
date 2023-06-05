FROM golang:1.20

ENV GOPATH=/
COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o test-task-ozon-migrate ./cmd/migrate/main.go
RUN go build -o test-task-ozon ./cmd/converter/main.go
CMD ["./test-task-ozon-migrate"]
CMD ["./test-task-ozon"]