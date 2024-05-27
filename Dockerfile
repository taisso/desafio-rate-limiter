FROM golang:1.22 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /app/main .
COPY --from=build /app/.env .
ENTRYPOINT [ "./main" ]