
FROM golang:1.17-bullseye as build

WORKDIR /go/src/edge-device-plugin
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install -ldflags="-s -w"

FROM debian:bullseye-slim
COPY --from=build /go/bin/edge-device-plugin /bin/edge-device-plugin

CMD ["/bin/edge-device-plugin"]