FROM golang:1.19.0-alpine AS build
# TODO: dynamically get this value
ENV GOVERSION=1.19.0

WORKDIR /usr/src
RUN apk add --no-cache \
    gcc=11.2.1_git20220219-r2 \
    build-base=0.5-r3
    # git=2.36.1-r0 \

ARG PKG
ARG APPNAME
ARG COMMITHASH
ARG BUILDTIME
ARG VERSION

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build \
      -a \
      -o k8s-ldap-auth \
      -ldflags "\
        -X ${PKG}/version.APPNAME=${APPNAME} \
        -X ${PKG}/version.VERSION=${VERSION} \
        -X ${PKG}/version.GOVERSION=${GOVERSION} \
        -X ${PKG}/version.BUILDTIME=${BUILDTIME} \
        -X ${PKG}/version.COMMITHASH=${COMMITHASH}" \
      main.go

FROM gcr.io/distroless/static:nonroot
EXPOSE 3000
WORKDIR /
COPY --from=build /usr/src/k8s-ldap-auth .
USER 65532:65532

ENTRYPOINT [ "/k8s-ldap-auth" ]
