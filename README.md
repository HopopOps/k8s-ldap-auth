# k8s-ldap-auth

## What
This is a webhook token authentication plugin implementation for ldap backend inspired from Daniel Weibel article "Implementing LDAP authentication for Kubernetes" at https://itnext.io/implementing-ldap-authentication-for-kubernetes-732178ec2155

k8s-ldap-auth provides two routes: `/auth` for the actual authentication from the CLI tool and `/token` for the token validation from the kube-api-server.

The user created from the TokenReview will contain both uid and groups from the LDAP user so you can use both for role binding.

## Configuration
Access rights to clusters and resources will not be implemented with this authentication hook, kubernetes RBAC will do that for you. `KUBERNETES_EXEC_INFO` is currently disregarded but might be used in future versions.

### Cluster

In the following example, I use the api version `client.authentication.k8s.io/v1beta1`. Feel free to put another better suited for your need if needed.

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

If the cluster was created with kubeadm, edit the kubeadm configuration file stored in namespace kube-system to add the configuration from above: `kubectl --namespace kube-system edit configmaps kubeadm-config`
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
        command: kube-ldap-auth

        apiVersion: client.authentication.k8s.io/v1beta1

        env:
          - name: AUTH_API_VERSION
            value: client.authentication.k8s.io/v1beta1

        args:
          - authenticate
          - --endpoint=https://k8s-ldap-auth/auth

        installHint: |
          k8s-ldap-auth is required to authenticate to the current context.
          It can be installed from https://github.com/vbouchaud/k8s-ldap-auth.

        provideClusterInfo: false

```

## Usage

```
NAME:
   k8s-ldap-auth - A client/server for kubernetes webhook authentication.

USAGE:
   k8s-ldap-auth [global options] command [command options] [arguments...]

VERSION:
   v0.1.0

AUTHOR:
   Vianney Bouchaud <vianney@bouchaud.org>

COMMANDS:
   server, s, serve       start the authentication server
   authenticate, a, auth  perform an authentication through a /auth endpoint
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
   --help, -h                     show help (default: false)
```

#### Client
```
NAME:
   k8s-ldap-auth authenticate - perform an authentication through a /auth endpoint

USAGE:
   k8s-ldap-auth authenticate [command options] [arguments...]

OPTIONS:
   --endpoint URI       The full URI the client will authenticate against.
   --user USER          The USER the client will connect as. [$USER]
   --password PASSWORD  The PASSWORD the client will connect with, can be located in '/home/vianney/.config/k8s-ldap-auth/password'. [$PASSWORD]
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

## What's next
 - Group search for ldap not supporting memberof attribute
 - Helm chart
 - Persisting certificates for jwt signing and validation so we can have multiple replicas
