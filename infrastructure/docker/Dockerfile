FROM golang:1.22.2-bookworm as base

# install if docker host doesn't has libc6 installed
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
WORKDIR app
COPY src/main main
COPY src/config/ config/
COPY src/openapi3.yaml openapi3.yaml
CMD [ "./main"]
