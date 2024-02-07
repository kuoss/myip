FROM golang:1.22 AS build
ARG VERSION
WORKDIR /temp/
COPY . ./
RUN go mod download -x
RUN go build -ldflags="-X 'main.Version=$VERSION'" -o /app/myip

FROM gcr.io/distroless/base-debian12
COPY --from=build /app/myip /app/
CMD ["/app/myip"]
