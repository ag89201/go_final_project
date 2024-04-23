FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY app ./app


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./go_final_project

FROM ubuntu:latest
WORKDIR /app
COPY --from=builder /app/go_final_project app/go_final_project 
COPY web ./web
VOLUME /app/db
CMD ["app/go_final_project"]





