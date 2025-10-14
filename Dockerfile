FROM node:24 AS build-vue

WORKDIR /app/frontend

COPY frontend/package*.json ./

RUN npm install

COPY frontend/. .

RUN npm run build

FROM golang:latest AS build-go

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint

FROM gcr.io/distroless/static-debian11

WORKDIR /

COPY --from=build-go /entrypoint /entrypoint
COPY --from=build-vue /app/frontend/dist /frontend/dist

EXPOSE 8080

ENTRYPOINT ["/entrypoint"]
