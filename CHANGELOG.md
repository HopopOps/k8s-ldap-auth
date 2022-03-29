# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [4.0.0] - 2022-03-29
### Common
#### Modified
- Made a bunch of dep updates, all through dependabot so full list will be easy to get (`git log v3.2.1..HEAD --author=dependabot --pretty=oneline`).

### Server
#### Modified
- Changing `--search-attributes` to `--extra-attributes`, fixing extra attributes not being processed and added to the user object. This is the change triggering a major update since the interface changed.

## [3.2.1] - 2021-11-10
### Client
#### Modified
- The entry in the credential manager is now specific to the k8s-ldap-auth server address, allowing for different k8s-ldap-auth to be used against different ldap servers.

## [3.2.0] - 2021-10-29
### Client
#### Added
- Password is now stored into the OS credential manager upon successful interactive authentication.

## [3.1.0] - 2021-08-19
### Server
#### Added
- `/health` now serves the application status.

## [3.0.0] - 2021-08-17
### Server
#### Added
- User username properties mapping with ldap can now be set with specific parameters or environment variable.

#### Changed
- Parameter `--member-of-property` is now `--memberof-property` (style consistency change)
- TokenReview user won't contain a list of group cn only anymore but their full dn to prevent name collision

### Client
#### Changed
- Cache file and folder containing the ExecCredential are now only readable by the owner.

## [2.0.1] - 2021-07-27
### Common
#### Added
- Added PIE compilation for binary hardening.
- Added trimpath option for reproducible builds.

## [2.0.0] - 2021-07-26
### Server
#### Added
- Token longevity can now be configured (in seconds). Default to 43200 (12 hours).

#### Changed
- Token generated now only contains uid. Groups and DN are added to the TokenReview when kube-apiserver dial k8s-ldap-auth.

### Client
#### Added
- There is now a reset command to ease the removal of cached token and force re-authentication on next invocation.

## [1.0.0] - 2021-07-22
### Server
#### Added
- `/auth` route for ldap authentication, returning an ExecCredential
- `/token` route for apiserver TokenReview validation
- Loading key pair for jwt signing and validation from files
- Generating an arbitrary key pair for jwt signing and validation if none is given
- TokenReview contains user id and groups from LDAP

### Client
#### Added
- Password and username can be given from standard input, environment variables or files.
