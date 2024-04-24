FROM golang:1.22.2-bookworm as base

# install if docker host doesn't has libc6 installed
RUN apt-get update && apt-get install -y libc6
RUN apt-get clean
EXPOSE 8080

FROM base as build
WORKDIR app
COPY src/main main
CMD [ "./main"]

FROM base as dev
WORKDIR /code
# Air for auto rebuild on changes
RUN go install github.com/cosmtrek/air@latest
RUN air init
CMD ["air"]
