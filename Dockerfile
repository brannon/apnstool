FROM golang:1.14 AS builder
WORKDIR /go/src/github.com/brannon/apnstool/
ADD . /go/src/github.com/brannon/apnstool/
RUN CGO_ENABLED=0 GOOS=linux go build -o apnstool .

FROM alpine:3.12
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/brannon/apnstool/apnstool .
CMD ["./apnstool"]
