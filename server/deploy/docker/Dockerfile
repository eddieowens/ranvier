FROM golang as builder

RUN mkdir /build

ADD . /build/

WORKDIR /build

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

RUN mkdir /.ssh

RUN ssh-keyscan -t rsa github.com >> /.ssh/known_hosts

FROM scratch

COPY --from=builder /build/main /app/

COPY --from=builder /build/config /app/config

COPY --from=builder /.ssh/known_hosts /.ssh/known_hosts

ENV SSH_KNOWN_HOSTS=/.ssh/known_hosts

EXPOSE 8080

WORKDIR /app

CMD ["./main"]