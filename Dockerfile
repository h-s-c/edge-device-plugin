FROM golang:1.25 AS build

WORKDIR /go/src/edge-device-plugin
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/edge-device-plugin cmd/edge-device-plugin/*.go

FROM gcr.io/distroless/static

COPY --from=build /go/bin/edge-device-plugin /edge-device-plugin

ENTRYPOINT ["/edge-device-plugin"]