# k8s-ldap-auth

## What
This is a webhook token authentication plugin implementation for ldap backend inspired from Daniel Weibel article "Implementing LDAP authentication for Kubernetes" at https://itnext.io/implementing-ldap-authentication-for-kubernetes-732178ec2155

k8s-ldap-auth provides two routes: `/auth` for the actual authentication from the CLI tool and `/token` for the token validation from the kube-api-server.

## Configuration
Access rights to clusters and resources will not be implemented with this authentication hook, kubernetes RBAC will do that for you.

### Cluster

### Client

The same user definition can be used on different clusters if they share this authentication hook.

```yml
---
apiVersion: v1
kind: Config
preferences: {}

users:
- name: my-user
  user:
    exec:
      command: "k8s-ldap-auth"

      apiVersion: "client.authentication.k8s.io/v1beta1"

      env:
      - name: "AUTH_API_VERSION"
        value: "client.authentication.k8s.io/v1beta1"

      args:
      - "authenticate"
      - "--endpoint=https://k8s-ldap/auth"

      installHint: |
        k8s-ldap-auth is required to authenticate to the current cluster.
        It can be installed from https://github.com/vbouchaud/k8s-ldap-auth.

      provideClusterInfo: false

clusters:
- cluster:
    server: https://kube.cluster.local:6443
  name: my-cluster

contexts:
- context:
    cluster: my-cluster
    user: my-user
  name: my-user@my-cluster

current-context: my-user@my-cluster
```

## Usage
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
   --search-attributes PROPERTY   Repeatable. User PROPERTY to fetch. Everything beside 'uid', 'dn', 'cn' (mandatory fields) will be stored in extra values in the UserInfo object. (default: "uid", "dn", "cn") [$LDAP_USER_SEARCHATTR]
   --search-scope SCOPE           The SCOPE of the search. Can take to values base object: 'base', single level: 'single' or whole subtree: 'sub'. (default: "sub") [$LDAP_USER_SEARCHSCOPE]
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
ORG=vbouchaud PLATFORM="linux/arm/v7,linux/amd64" make docker
```
`ORG` defaults to my private docker registry

`PLATFORM` defaults to `linux/arm/v7,linux/arm64/v8,linux/amd64`

## Deployment

## What's next
 - Group search for ldap not supporting memberof attribute
