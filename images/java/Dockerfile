FROM metrue/fx-java-base

# Adding source, compile and package into a fat jar
ADD src /code/src
RUN ["mvn", "package"]

EXPOSE 3000
CMD ["/usr/lib/jvm/java-8-openjdk-amd64/bin/java", "-jar", "target/fx-app-java-0.1.0.jar"]
