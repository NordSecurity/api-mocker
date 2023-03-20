FROM golang:1.20.2-alpine3.17 as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o api-mocker .

FROM alpine:3.11.3 as api-mocker
COPY --from=builder /build/api-mocker .
COPY --from=builder /build/rules/test-rule.json /rules/test-rule.json

ENTRYPOINT ["./api-mocker","-host","0.0.0.0", "-rules", "rules/test-rule.json"]
