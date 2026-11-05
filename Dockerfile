FROM golang:1.23-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o wosaas .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /build/wosaas .
COPY --from=builder /build/conf ./conf

ENV TZ=Asia/Shanghai
EXPOSE 8600 8081

ENTRYPOINT ["./wosaas"]
CMD ["-conf", "./conf/app.toml"]
