FROM gocv/opencv:4.13.0

ENV CGO_ENABLED=1

RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN go install github.com/githubnemo/CompileDaemon@latest
RUN CGO_ENABLED=0 go install github.com/swaggo/swag/cmd/swag@latest \
    && CGO_ENABLED=0 go install github.com/go-swagger/go-swagger/cmd/swagger@latest

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

CMD swag init -g cmd/api/api.go \
  && swagger validate docs/swagger.yaml \
  && /wait \
  && go run db/migrations/entry.go \
  && CompileDaemon \
       --build="go build -o /app/cmd/api/api /app/cmd/api/api.go" \
       --command="/app/cmd/api/api" \
       --color