FROM golang:1.16.6-alpine AS build
ENV GOVERSION=1.16.6
WORKDIR /usr/src
RUN apk add --no-cache git=2.32.0-r0 gcc=10.3.1_git20210424-r2 build-base=0.5-r2
ARG PKG
ARG APPNAME
ARG COMMITHASH
ARG BUILDTIME
ARG VERSION
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o k8s-ldap-auth -ldflags "\
        -X ${PKG}/version.APPNAME=${APPNAME} \
        -X ${PKG}/version.VERSION=${VERSION} \
        -X ${PKG}/version.GOVERSION=${GOVERSION} \
        -X ${PKG}/version.BUILDTIME=${BUILDTIME} \
        -X ${PKG}/version.COMMITHASH=${COMMITHASH} \
        -s -w"

FROM alpine:3
EXPOSE 3000
WORKDIR /usr/src
COPY --from=build /usr/src/k8s-ldap-auth /usr/bin/
CMD ["/usr/bin/k8s-ldap-auth"]
