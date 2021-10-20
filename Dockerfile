# Initial stage: download modules
FROM golang:1.17 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.17 as builder

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 bot_joy

RUN mkdir -p /bot_joy
ADD . /bot_joy
WORKDIR /bot_joy

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go install

# Final stage: Run the binary
FROM scratch

# don't forget /etc/passwd from previous stage
COPY --from=builder /etc/passwd /etc/passwd
USER bot_joy

# and finally the binary
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/bot_joy /bot_joy

CMD ["/bot_joy"]