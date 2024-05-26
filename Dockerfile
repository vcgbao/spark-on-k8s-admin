# build frontend
FROM node:20.12.2-alpine3.19 as frontend-builder
WORKDIR /app
COPY modules/web/package*.json ./
RUN npm install
COPY modules/web .
RUN npm run build


# build backend
FROM golang:1.22.2-alpine AS backend-build
WORKDIR /app
COPY modules/api/go.mod modules/api/go.sum ./
RUN go mod download
COPY modules/api .
RUN GOOS=linux go build -o spark-on-k8s-admin .


FROM alpine:3.19.1
WORKDIR /app
COPY --from=frontend-builder /app/build /var/www
COPY --from=backend-build /app/spark-on-k8s-admin .
RUN apk --no-cache add ca-certificates tzdata
ENTRYPOINT ["/app/spark-on-k8s-admin", "-config=/config/config.yaml"]
