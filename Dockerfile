FROM golang:1.24 AS build
ARG VERSION

WORKDIR /go/src/myip
COPY . .

RUN go mod download -x

RUN CGO_ENABLED=0 go build -ldflags="-X 'main.Version=$VERSION'" -o /go/bin/myip


FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/myip /
CMD ["/myip"]
