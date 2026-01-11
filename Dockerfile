FROM golang:1.25-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o rss-proxy

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/rss-proxy .
EXPOSE 800
CMD ["./rss-proxy"]
