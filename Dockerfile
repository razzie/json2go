FROM golang:1.21 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make

FROM alpine
WORKDIR /
COPY --from=builder /workspace/json2go .
ENTRYPOINT ["/json2go"]
