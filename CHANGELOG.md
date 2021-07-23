# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Server
#### Added
- Token longevity can now be configured (in seconds). Default to 43200 (12 hours).

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
