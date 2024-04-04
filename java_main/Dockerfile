#build stage
FROM maven:3.8.5-openjdk-17 as builder
LABEL Name=smartcodeql-java Version=0.0.1
WORKDIR /src
COPY java_main/pom.xml .
COPY java_main/src .
COPY java_main/settings.xml .
RUN mvn -f ./pom.xml clean package -s ./settings.xml

#run stage


FROM openjdk:17-jdk-alpine
WORKDIR /app
COPY --from=builder /src/target/app.jar /app/app.jar
CMD ["java", "-jar", "/app/app.jar"]