FROM golang:latest AS BuildStage
COPY ./ .
RUN go mod download
RUN go build -o /example-app main.go

FROM debian:bookworm-slim
RUN apt-get update -y
RUN apt-get install ca-certificates -y
WORKDIR /
COPY --from=BuildStage /example-app .
EXPOSE 8080
CMD ["/example-app"]