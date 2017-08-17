# Go Dep Registry

## Introduction
Go Dep Registry hopes to be the future registry for Golang packages; both public and private packages. With the help of Sam Boyer (Vivid Cortex) and Yoav Landman (JFrog), we hope to shape the future of the Go packaging ecosystem using `dep` as the package management solution.

## Current Status
**Pre-Alpha**: This service is still being architected and developed and is **not** recommended for production use.

## Architecture
The architecture of the registry is broken down into three parts: authentication, metadata, and binary storage. All types have their own backends that are pluggable and can be swapped out at any time.

## Backends
Below is a list of supported backends. More backends are planned for the future. If you would like to add support for a particular backend, please submit a PR. The interface for creating new backends is very simple.

Backend configuration is contained in a single string, prefixed with a backend identifier. Anything after this identifier is up to the backend to parse. The list of supported backends below shows the identifier and proper usage.

### Auth
Authentication source for validating identity.

Supported backends:
* User/Password backed by BoltDB
    * `userpass://<filename>`

### MetaStore
Metadata about packages and their versions are stored using MetaStore. This contains the import path, description of the package, availalbe versions, and the package's main landing page for providing more information about the package.

Supported backends:
* BoltDB
    * `boltdb://<filename>`

### BinStore
Binary releases, in tar, tgz, and zip formats, are stored here. A UUID4 is generated per version and is stored in the metadata for the package. Specific version binaries can be retrieved using this UUID4.

Supported backends:
* BoltDB
    * `boltdb://<filename>`

## Configuration
The registry can be configured using either a JSON config file or environment variables:
* `auth_path` / `AUTH_PATH`: The auth backend connection string
* `binstore_path` / `BINSTORE_PATH`: The BinStore connection string
* `metastore_path` / `METASTORE_PATH`: The MetaStore connection string
* `signing_key` / `SIGNING_KEY`: The key used to sign auth tokens
* `token_ttl` / `TOKEN_TTL`: Time-to-live for tokens duration (i.e. 2h for 2 hours)
* `port` / `PORT`: The port the HTTP server will listen on

Configuration has sane defaults and will print a warning to `stdout` identifying any settings that need to be adjusted. Running without any configuration generates a new signing key at every start, invalidating any previously generated tokens. It will also default to using BoltDB for all backends.

To use a JSON config file, pass the filename as the first argument to the executable. A bare-minimum JSON config file might look like this:
```json
{
    "signing_key": "some-super-secret-key",
    "token_ttl": "24h"
}
```

## Contributions
Please help out by opening issues and submitting PR's. This could be the future of Go package management, so your input matters!