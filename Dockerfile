FROM golang:1.21

WORKDIR /app

COPY ./ ./

RUN go mod tidy

RUN ls -la

RUN go build -o /user-app

CMD ["/user-app"]