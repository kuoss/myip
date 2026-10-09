FROM golang:1.27 AS builder
ARG VERSION

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X 'main.Version=$VERSION'" -o /myip

FROM gcr.io/distroless/static-debian12:latest
COPY --from=builder /myip /myip

USER nobody
CMD ["/myip"]
