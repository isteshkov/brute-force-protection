FROM golang:1.16-buster as builder
RUN go get github.com/go-delve/delve/cmd/dlv
WORKDIR /app
COPY . .
RUN make deps
RUN go build -gcflags="all=-N -l" -o brute_force_protection .

FROM debian:buster-slim
RUN apt-get update && apt-get install -yq ca-certificates
COPY --from=builder /app/brute-force-protection /go/bin/dlv /
CMD /dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec /brute_force_protection --continue
