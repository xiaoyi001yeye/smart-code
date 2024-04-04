FROM quay.io/quarkus/centos-quarkus-maven:latest as builder
WORKDIR /src
COPY pom.xml /src/
COPY src /src/src
RUN mvn -f /src/pom.xml clean package


FROM openjdk:17-jdk-alpine
WORKDIR /app
COPY --from=builder /src/target/app.jar /app/app.jar
CMD ["java", "-jar", "/app/app.jar"]