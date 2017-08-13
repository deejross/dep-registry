# Go Dep Registry

## Introduction
Go Dep Registry hopes to be the future registry for Golang packages; both public and private packages. With the help of Sam Boyer (Vivid Cortex) and Yoav Landman (JFrog), we hope to shape the future of the Go packaging ecosystem using `dep` as the package management solution.

## Current Status
**Pre-Alpha**: This service is still being architected and developed and is **not** recommended for production use.

## Architecture
The architecture of the registry is broken down into two parts: metadata and binary storage. Both types of storage have their own backends that are pluggable and can be swapped out at any time.

### MetaStore
Metadata about packages and their versions are stored using MetaStore. This contains the import path, description of the package, availalbe versions, and the package's main landing page for providing more information about the package.

Supported backends:
* BoltDB

### BinStore
Binary releases, in tar, tgz, and zip formats, are stored here. A UUID4 is generated per version and is stored in the metadata for the package. Specific version binaries can be retrieved using this UUID4.

Supported backends:
* BoltDB

## Backends
The idea surrounding the dep registry is to allow metadata and binary data storage backends to be swappable. Currently, BoltDB is the only supported backend for both. It was important that the first and default backends for the registry not require any dependencies to allow for easier development and testing of the registry. More backends are planned for the future. If you would like to add support for a particular backend, please submit a PR. The interface for creating new backends is very simple.

## Configuration
Not yet implemented, please stand by...