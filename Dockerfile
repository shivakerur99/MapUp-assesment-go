FROM golang:1.12.0-alpine3.9
WORKDIR /test
COPY . /test
RUN go build /test
EXPOSE 8080
ENTRYPOINT [ "./dockersam" ]