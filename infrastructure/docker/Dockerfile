FROM golang:1.23-bookworm AS base

FROM base AS build-env
WORKDIR /app
COPY src/. .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# BUILD published image
FROM alpine:latest AS build
WORKDIR /app
COPY --from=build-env /app/main .
COPY src/config/ config/
COPY src/openapi3.yaml openapi3.yaml
# Temporary for initial PoC testing. File has to be generated periodically in production
COPY src/provider/testdata/skyalps.xml /shared/skyalps.xml
EXPOSE 8080
ENTRYPOINT [ "./main"]

FROM base AS dev
WORKDIR /code
# Air for auto rebuild on changes
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

FROM base AS test
WORKDIR /code
CMD ["go", "test", "./..."]