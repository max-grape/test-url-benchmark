ARG GO_IMG
ARG FROM_IMG

FROM $GO_IMG as build
ARG CWD
ARG GOOS
ARG GOARCH
WORKDIR $CWD
COPY . .
ENV GO111MODULE=off
RUN GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 go build -v -o app

FROM $FROM_IMG
ARG CWD
COPY --from=build $CWD/app .
COPY ./search.html ./search.html
RUN apk add curl
HEALTHCHECK --interval=10s --timeout=2s CMD curl 127.0.0.1:8080/health | grep -w "healthy"
CMD ["./app"]
