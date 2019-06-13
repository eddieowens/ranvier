FROM golang as builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine

COPY --from=builder /build/main /app/

COPY --from=builder /build/config /app/config

RUN mkdir -p ~/.ssh

RUN touch ~/.ssh/known_hosts

WORKDIR /app

CMD ["./main"]