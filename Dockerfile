
FROM golang:1.17-bullseye as build

WORKDIR /go/src/edge-device-plugin
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install -ldflags="-s -w" cmd/edge-device-plugin/*.go

FROM debian:bullseye-slim
COPY --from=build /go/src/edge-device-plugin/edge-device-plugin.sh /bin/edge-device-plugin.sh
COPY --from=build /go/bin/main /bin/edge-device-plugin

CMD ["/bin/edge-device-plugin.sh"]