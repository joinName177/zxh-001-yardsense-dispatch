FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/yardsense ./cmd/yardsense
FROM alpine:3.22
RUN adduser -D -u 10001 app
USER app
WORKDIR /app
COPY --from=build /out/yardsense /app/yardsense
EXPOSE 8080
ENTRYPOINT ["/app/yardsense"]
