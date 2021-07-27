# k8s-ldap-auth

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/vbouchaud/k8s-ldap-auth?style=for-the-badge)](https://github.com/vbouchaud/k8s-ldap-auth/releases/latest)
[![License](https://img.shields.io/github/license/vbouchaud/k8s-ldap-auth?style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)

## What

This is a webhook token authentication plugin implementation for ldap backend inspired by Daniel Weibel article "Implementing LDAP authentication for Kubernetes" at https://itnext.io/implementing-ldap-authentication-for-kubernetes-732178ec2155

k8s-ldap-auth is released as a binary containing both client and server.

The server part provides two routes:
 - `/auth` for the actual authentication from the CLI tool
 - `/token` for the token validation from the kube-apiserver.

The user created from the TokenReview will contain both uid and groups from the LDAP user so you can use both for role binding.

The same k8s-ldap-auth server can be used to authenticate with multiple kubernetes cluster since the ExecCredential it provides contains a signed token that will eventually be used by a kube-apiserver in a TokenReview that will be sent back.

I actually use this setup on quite a few clusters with a growing userbase.

## Configuration

Access rights to clusters and resources will not be implemented in this authentication hook, kubernetes RBAC will do that for you.

`KUBERNETES_EXEC_INFO` is currently disregarded but might be used in future versions.

### Cluster

In the following example, I use the api version `client.authentication.k8s.io/v1beta1`. Feel free to put another better suited for your need.

The following authentication token webhook config file will have to exist on every control-plane. In the following configuration it's located at `/etc/kubernetes/webhook-auth-config.yml`:
```yml
---
apiVersion: v1
kind: Config

clusters:
  - name: authentication-server
    cluster:
      server: https://k8s-ldap-auth/token

users:
  - name: kube-apiserver

contexts:
  - context:
      cluster: authentication-server
      user: kube-apiserver
    name: kube-apiserver@authentication-server

current-context: kube-apiserver@authentication-server
```

#### New cluster

If you're creating a new cluster with kubeadm, you can add the following to your init configuration file:
```yml
---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
apiServer:
  extraArgs:
    authentication-token-webhook-config-file: "/etc/ldap-auth-webhook/config.yml"
    authentication-token-webhook-version: client.authentication.k8s.io/v1beta1
  extraVolumes:
  - name: "webhook-config"
    hostPath: "/etc/kubernetes/webhook-auth-config.yml"
    mountPath: "/etc/ldap-auth-webhook/config.yml"
    readOnly: true
    pathType: File
```

#### Existing cluster

If the cluster was created with kubeadm, edit the kubeadm configuration stored in the namespace `kube-system` to add the configuration from above: `kubectl --namespace kube-system edit configmaps kubeadm-config`
Editing this configuration does not actually update your api-server. It will however be used if you need to add a new control-plane with `kubeadm join`.

On every control plane, edit the manifest found at `/etc/kubernetes/manifests/kube-apiserver.yaml`:
```yml
spec:
  containers:
  - name: kube-apiserver
    command:
    - kube-apiserver
    # ...
    - --authentication-token-webhook-config-file=/etc/ldap-auth-webhook/config.yml
    - --authentication-token-webhook-version=v1beta1

    # ...

    volumeMounts:
    - mountPath: /etc/ldap-auth-webhook/config.yml
      name: webhook-config
      readOnly: true

  # ...

  volumes:
  - hostPath:
      path: /etc/kubernetes/webhook-auth-config.yml
      type: File
    name: webhook-config
```

### Client

The same user definition can be used on different clusters if they share this authentication hook.

```yml
users:
  - name: my-user
    user:
      exec:
        # In the following, we assume a binary called `k8s-ldap-auth` is
        # available in the path. You can instead put a path to the binary.
        # Windows paths do work with kubectl so the following would also work:
        # `C:\users\foo\Documents\k8s-ldap-auth.exe`.
        command: k8s-ldap-auth

        # This field is used by kubectl to fill a template TokenReview in
        # `$KUBERNETES_EXEC_INFO` environment variable. Not currently used, it's
        # might be in the future.
        apiVersion: client.authentication.k8s.io/v1beta1

        env:
          # This environment variable is used within `k8s-ldap-auth` to create
          # an ExecCredential. Future version of this authenticator might not
          # need it but you'll have to provide it for now.
          - name: AUTH_API_VERSION
            value: client.authentication.k8s.io/v1beta1

          # You can fill a USER environment variable to your username if you
          # want to overwrite the USER from you system or to an empty one if you
          # want the authenticator to ask for one at runtime.
          - name: USER
            value: ""

        args:
          - authenticate

          # This is the endpoint to authenticate against. Basically, the server
          # started with `k8s-ldap-auth server` plus the `/auth` route, used for
          # authentication.
          - --endpoint=https://k8s-ldap-auth/auth

        installHint: |
          k8s-ldap-auth is required to authenticate to the current context.
          It can be installed from https://github.com/vbouchaud/k8s-ldap-auth.

        # This parameter, when true, tells `kubectl` to fill the TokenReview in
        # the `$KUBERNETES_EXEC_INFO` environment variable with extra config
        # from the definition of the specific cluster currently targeted.
        # This is not used today but might be in the future to allow for custom
        # rules on a per-cluster basis.
        provideClusterInfo: false
```

## Usage

```
NAME:
   k8s-ldap-auth - A client/server for kubernetes webhook authentication.

USAGE:
   k8s-ldap-auth [global options] command [command options] [arguments...]

VERSION:
   v2.0.0

AUTHOR:
   Vianney Bouchaud <vianney@bouchaud.org>

COMMANDS:
   server, s, serve       start the authentication server
   authenticate, a, auth  perform an authentication through a /auth endpoint
   reset, r               delete the cached ExecCredential to force authentication at next call
   help, h                Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose LEVEL  The verbosity LEVEL - (rs/zerolog#Level values). (default: 3) [$VERBOSE]
   --help, -h       show help (default: false)
   --version, -v    print the version (default: false)
```

#### Server
```
NAME:
   k8s-ldap-auth server - start the authentication server

USAGE:
   k8s-ldap-auth server [command options] [arguments...]

OPTIONS:
   --host HOST                    The HOST the server will listen on. [$HOST]
   --port PORT                    The PORT the server will listen to. (default: 3000) [$PORT]
   --ldap-host HOST               The ldap HOST (and scheme) the server will authenticate against. (default: "ldap://localhost") [$LDAP_ADDR]
   --bind-dn DN                   The service account DN to do the ldap search. [$LDAP_BINDDN]
   --bind-credentials PASSWORD    The service account PASSWORD to do the ldap search, can be located in '/etc/k8s-ldap-auth/ldap/password'. [$LDAP_BINDCREDENTIALS]
   --search-base DN               The DN where the ldap search will take place. [$LDAP_USER_SEARCHBASE]
   --search-filter FILTER         The FILTER to select users. (default: "(&(objectClass=inetOrgPerson)(uid=%s))") [$LDAP_USER_SEARCHFILTER]
   --member-of-property PROPERTY  The PROPERTY where group entitlements are located. (default: "ismemberof") [$LDAP_USER_MEMBEROFPROPERTY]
   --search-attributes PROPERTY   Repeatable. User PROPERTY to fetch. Everything beside 'uid', 'dn' (mandatory fields) will be stored in extra values in the UserInfo object. (default: "uid", "dn") [$LDAP_USER_SEARCHATTR]
   --search-scope SCOPE           The SCOPE of the search. Can take to values base object: 'base', single level: 'single' or whole subtree: 'sub'. (default: "sub") [$LDAP_USER_SEARCHSCOPE]
   --private-key-file PATH        The PATH to the private key file [$PRIVATE_KEY_FILE]
   --public-key-file PATH         The PATH to the public key file [$PUBLIC_KEY_FILE]
   --token-ttl TTL                The TTL for newly generated tokens, in seconds (default: 43200) [$TTL]
   --help, -h                     show help (default: false)
```

#### Client

Even though it's not specified anywhere, the `--password` option and the equivalent `$PASSWORD` environment variable as well as the configfile containing a password were added for convenience sake, e.g. when running in an automated fashion, etc. If not provided, it will be asked at runtime. The same can be said for `--user` options and `$USER` environment variables.

```
NAME:
   k8s-ldap-auth authenticate - perform an authentication through a /auth endpoint

USAGE:
   k8s-ldap-auth authenticate [command options] [arguments...]

OPTIONS:
   --endpoint URI       The full URI the client will authenticate against.
   --user USER          The USER the client will connect as. [$USER]
   --password PASSWORD  The PASSWORD the client will connect with, can be located in '$XDG_CONFIG_HOME/k8s-ldap-auth/password'. [$PASSWORD]
   --help, -h           show help (default: false)
```

## Build

A stripped binary can be built with:
```
make k8s-ldap-auth
```

A stripped and compressed binary can be build with:
```
make release
```

Docker release multiarch image can be built and pushed with:
```
PLATFORM="linux/arm/v7,linux/amd64" make docker
```

`PLATFORM` defaults to `linux/arm/v7,linux/arm64/v8,linux/amd64`

## Distribution

### Docker
A docker image of this projet is available for arm/v7, arm64/v8, amd64 at [vbouchaud/k8s-ldap-auth](https://hub.docker.com/r/vbouchaud/k8s-ldap-auth) on docker hub. A mirror has been setup on quay.io at [vbouchaud/k8s-ldap-auth](https://quay.io/vbouchaud/k8s-ldap-auth)

### Binary
Binaries for the following OS and architectures are available on the release page:
 - linux/arm64
 - linux/arm
 - linux/amd64
 - darwin/arm64
 - darwin/amd64
 - windows/amd64

### Linux
#### Archlinux
[![AUR version](https://img.shields.io/aur/version/k8s-ldap-auth-bin?label=k8s-ldap-auth-bin&style=for-the-badge)](https://aur.archlinux.org/packages/k8s-ldap-auth-bin/)

[![AUR version](https://img.shields.io/aur/version/k8s-ldap-auth?label=k8s-ldap-auth&style=for-the-badge)](https://aur.archlinux.org/packages/k8s-ldap-auth/)

[![AUR version](https://img.shields.io/aur/version/k8s-ldap-auth-git?label=k8s-ldap-auth-git&style=for-the-badge)](https://aur.archlinux.org/packages/k8s-ldap-auth-git/)

## What's next

 - Group search for ldap not supporting `memberof` attribute ;
 - Helm chart ;
