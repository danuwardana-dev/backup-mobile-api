FROM golang:1.24.1

ENV TZ=Asia/Jakarta`
RUN mkdir /app
ADD ../mirror/backend-mobile-api /app

WORKDIR /app

COPY  go.mod  .

RUN go mod tidy

RUN go build -o  backend-mobile-api ./app/

CMD [ "./backend-mobile-api", "rest" ]