FROM golang:1.16-buster as builder
WORKDIR /app
COPY . /app
RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN task deps
RUN task build
CMD /app/brute_force_protection test migrate && /app/brute_force_protection