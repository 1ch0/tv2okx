FROM golang:1.20-alpine as builder

ENV GOPROXY https://goproxy.cn

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the go source
COPY cmd/server/main.go cmd/server/main.go
COPY pkg/ pkg/
COPY go.mod go.mod
COPY go.sum go.sum

# Build
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux \
    go build -a \
    -o server cmd/server/main.go


FROM alpine:3.19

WORKDIR /workspace

COPY --from=builder /workspace/server /workspace/

CMD ["/workspace/server"]