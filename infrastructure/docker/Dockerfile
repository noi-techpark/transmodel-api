FROM golang:1.22.2-bookworm as base

# libc6 missing on host. Remove this once we are on newer hosts (read: kubernetes)
RUN apt-get update && apt-get install -y libc6
RUN apt-get clean
EXPOSE 8080

FROM base as dev
WORKDIR /code
# Air for auto rebuild on changes
RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]

FROM base as test
WORKDIR /code
CMD ["go", "test", "./..."]

FROM alpine:3.20 as build
# again, libc6 missing on host. Remove this once we are on newer hosts (read: kubernetes)
RUN apk add gcompat
WORKDIR app
EXPOSE 8080
COPY src/main main
COPY src/config/ config/
COPY src/openapi3.yaml openapi3.yaml
CMD [ "./main"]
