FROM golang:1.23.4

WORKDIR /app

COPY ../go.mod  ./
RUN go mod download && go mod verify

COPY .. .
RUN go build -v -o worker worker.go

ENTRYPOINT ["./worker","-p=8000"]

# Exposing server port
EXPOSE 8000