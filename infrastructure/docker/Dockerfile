FROM openjdk:11 as base

RUN mkdir -p /code

# Dev
FROM base as dev

COPY infrastructure/docker/entrypoint.sh /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

# Build
FROM base as build

WORKDIR /code

COPY ./ /code
RUN ./mvnw -B -DskipTests install

CMD [ "./mvnw", "-B", "spring-boot:run" ]
