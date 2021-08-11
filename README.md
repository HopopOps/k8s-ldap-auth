# k8s-ldap-auth

[![GitHub release (latest by date)](https://img.shields.io/github/v/release/vbouchaud/k8s-ldap-auth?style=for-the-badge)](https://github.com/vbouchaud/k8s-ldap-auth/releases/latest)
[![License](https://img.shields.io/github/license/vbouchaud/k8s-ldap-auth?style=for-the-badge)](https://opensource.org/licenses/MPL-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/vbouchaud/k8s-ldap-auth?style=for-the-badge)](https://goreportcard.com/report/github.com/vbouchaud/k8s-ldap-auth)

## What

This is a webhook token authentication plugin implementation backed by LDAP.

k8s-ldap-auth is released as a binary containing both client and server.

The server part provides two routes:
 - `/auth` for the actual authentication from the CLI tool
 - `/token` for the token validation from the kube-apiserver.

The user created from the TokenReview will contain both uid and groups from the LDAP user so you can use both for role binding.

The same k8s-ldap-auth server can be used to authenticate with multiple kubernetes cluster since the ExecCredential it provides contains a signed token that will eventually be used by a kube-apiserver in a TokenReview that will be sent back.

I actually use this setup on quite a few clusters with a growing userbase.

Access rights to clusters and resources will not be implemented in this authentication hook, kubernetes RBAC will do that for you.

`KUBERNETES_EXEC_INFO` is currently disregarded but might be used in future versions.

## Usage

You can see the commands and their options with:
```
k8s-ldap-auth --help
# or
k8s-ldap-auth [command] --help
```

Pretty much all options can be set using environment variables and a few also read their values from files.

### Server

Create the password file for the bind-dn:
```
echo -n "bind_P@ssw0rd" > /etc/k8s-ldap-auth/ldap/password
```

The server can then be started with:
```
k8s-ldap-auth serve \
  --ldap-host="ldaps://ldap.company.local" \
  --bind-dn="uid=k8s-ldap-auth,ou=services,ou=company,ou=local" \
  --search-base="ou=people,ou=company,ou=local"
```

Note that if the server do not know of any key pair it will create one at launch but will not persist it.
If you want your jwt tokens to be valid accross server instances: after restarts or behind a load-balancer, you should provide a key pair.

Key pair can be created with openssl:
```
openssl genrsa -out key.pem 4096
openssl rsa -in key.pem -outform PEM -pubout -out public.pem
```

Then, the server can be started with:
```sh
k8s-ldap-auth serve \
  --ldap-host="ldaps://ldap.company.local" \
  --bind-dn="uid=k8s-ldap-auth,ou=services,ou=company,ou=local" \
  --search-base="ou=people,ou=company,ou=local" \
  --private-key-file="path/to/key.pem"
  --public-key-file="path/to/public.pem"
```

Now for the cluster configuration.

In the following example, I use the api version `client.authentication.k8s.io/v1beta1`. Feel free to put another better suited for your need.

The following authentication token webhook config file will have to exist on every control-plane. In the following configuration it's located at `/etc/kubernetes/webhook-auth-config.yml`:
```yml
---
apiVersion: v1
kind: Config

clusters:
  - name: authentication-server
    cluster:
      server: https://<server address>/token

users:
  - name: kube-apiserver

contexts:
  - context:
      cluster: authentication-server
      user: kube-apiserver
    name: kube-apiserver@authentication-server

current-context: kube-apiserver@authentication-server
```

##### New cluster

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

##### Existing cluster

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

#### Client

Even though it's not specified anywhere, the `--password` option and the equivalent `$PASSWORD` environment variable as well as the configfile containing a password were added for convenience sake, e.g. when running in an automated fashion, etc. If not provided, it will be asked at runtime. The same can be said for the `--user` options and `$USER` environment variables.

Authentication can be achieved with the following command you can execute to test your installation:
```
k8s-ldap-auth auth --endpoint="https://<server address>/auth"
```

You can now configure `kubectl` to use `k8s-ldap-auth` to authenticate to clusters by editing your kube config file and adding the following user:
```yml
users:
  - name: my-user
    user:
      exec:
        # In the following, we assume a binary called `k8s-ldap-auth` is
        # available in the path. You can instead put the full path to the binary.
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

This user can be used by setting the `--user` attribute for `kubectl`:
```
kubectl --user my-user get nodes
```

You can also create contexts with it:
```yaml
contexts:
  - name: context1
    context:
      cluster: cluster1
      user: my-user
  - name: context2
    context:
      cluster: cluster2
      user: my-user

current-context: context1
```

And then:
```
kubectl --context context2 get nodes
kubectl get nodes
```

### RBAC
Before you can actually get some result, you will have to upload some rolebindings to the cluster. As stated before, `k8s-ldap-auth` provides the apiserver with an ExecCredential containing both LDAP username and groups so both can be used in ClusterRoleBindings and RoleBindings.

Beware: group DNs, username and user id are all set to lowercase in the TokenReview.

#### Example

Given the following ldap users:

```
# User Alice
dn: uid=alice,ou=people,ou=company,ou=local
ismemberof: cn=somegroup,ou=groups,ou=company,ou=local

# User Bob
dn: uid=bob,ou=people,ou=company,ou=local
ismemberof: cn=somegroup,ou=groups,ou=company,ou=local

# User Carol
dn: uid=carol,ou=people,ou=company,ou=local
```

If I want to bind `cluster-admin` ClusterRole to the user `carol`, I can create a ClusterRoleBinding as following:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: custom-cluster-admininistrators
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: User
  name: carol
```

Let's say I want to bind the `view` ClusterRole so that all user in the group `somegroup` will have view access to a given namespace, I can create a RoleBinding such as:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: namespace-users
  namespace: somenamespace
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: view
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: cn=somegroup,ou=groups,ou=company,ou=local
```

Note: Kubernetes comes with some basic [predefined roles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles) for you to use.

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
Docker images of this projet are available for arm/v7, arm64/v8 and amd64 at [vbouchaud/k8s-ldap-auth](https://hub.docker.com/r/vbouchaud/k8s-ldap-auth) on docker hub and on quay.io at [vbouchaud/k8s-ldap-auth](https://quay.io/vbouchaud/k8s-ldap-auth).

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
[![AUR version](https://img.shields.io/aur/version/k8s-ldap-auth?label=k8s-ldap-auth&style=for-the-badge)](https://aur.archlinux.org/packages/k8s-ldap-auth/)

[![AUR version](https://img.shields.io/aur/version/k8s-ldap-auth-bin?label=k8s-ldap-auth-bin&style=for-the-badge)](https://aur.archlinux.org/packages/k8s-ldap-auth-bin/)

[![AUR last modified](https://img.shields.io/aur/last-modified/k8s-ldap-auth-git?label=k8s-ldap-auth-git&style=for-the-badge)](https://aur.archlinux.org/packages/k8s-ldap-auth-git/)

### Darwin
#### With `brew`

`k8s-ldap-auth.rb` is not in the official repository, you have to download [the formula](https://raw.githubusercontent.com/vbouchaud/k8s-ldap-auth/master/distribution/darwin/brew/k8s-ldap-auth.rb) and then specify its path when calling brew:
```
curl -O https://raw.githubusercontent.com/vbouchaud/k8s-ldap-auth/master/distribution/darwin/brew/k8s-ldap-auth.rb
brew install --formula ./k8s-ldap-auth.rb
```

## Inspiration

I originaly started this project after reading Daniel Weibel's article "Implementing LDAP authentication for Kubernetes" (https://learnk8s.io/kubernetes-custom-authentication or https://itnext.io/implementing-ldap-authentication-for-kubernetes-732178ec2155).

## What's next

 - Group search for ldap not supporting `memberof` attribute ;
 - Helm chart ;
