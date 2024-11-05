FROM golang:1.23.1

WORKDIR /user

COPY . /user

RUN go build -o user_service /user

EXPOSE 8083

ENTRYPOINT ["./user_service"]
