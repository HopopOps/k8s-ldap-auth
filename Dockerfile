FROM golang:1.17.7-alpine AS build
# TODO: dynamically get this value
ENV GOVERSION=1.17.6

WORKDIR /usr/src
RUN apk add --no-cache git=2.34.1-r0 gcc=10.3.1_git20211027-r0 build-base=0.5-r2

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
