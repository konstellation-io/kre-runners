# KRT Files Downloader 

Kubernetes executes this program as an init container into the workflow runners PODs.
The goal is to download the KRT file stored into a MongoDB GridFS and to extract it to
a shared directory.

## Required environment variables

The following variables are **required**:

* KRT_VERSION_ID: The admin-api stores the KRT file in GridFS using this value as identifier. It is necessary to download the KRT file.
* KRT_BASE_PATH: The location where the KRT files will be extracted.
* KRT_MONGO_URI: The MongoDB address.
* KRT_MONGO_DB_NAME: The MongoDB database name.
* KRT_MONGO_BUCKET: The bucket where the KRT files are.

## Run tests

If you want to run the tests locally execute:

```sh
go test ./...
```

## Linters

`golangci-lint` is a fast Go linters runner. It runs linters in parallel, uses caching, supports yaml config, has
integrations with all major IDE and has dozens of linters included.

As you can see in the `.golangci.yml` config file of this repo, we enable more linters than the default and 
have more strict settings.

To run `golangci-lint` execute: 
```
golangci-lint run
```
