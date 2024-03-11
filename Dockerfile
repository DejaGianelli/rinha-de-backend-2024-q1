FROM golang:1.22

WORKDIR /app

COPY ./app/go.mod ./app/go.sum ./

RUN go mod download && go mod verify

COPY ./app/ .

RUN go build -o /rinha-backend-2024

EXPOSE 8080

CMD [ "/rinha-backend-2024" ]