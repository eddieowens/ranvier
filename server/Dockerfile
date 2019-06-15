FROM golang as builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM scratch

COPY --from=builder /build/main /app/

COPY --from=builder /build/config /app/config

WORKDIR /app

CMD ["./main"]