# Video Service

A service for managing uploads and fetches for video content.

## Quick Start

Must have [Go compiler][1], [GNU Make][2], and [protoc][3]. Additional dependencies required for code generation will be installed via the make targets if they are not present in your $GOPATH. See the [makefile][4] for an exhaustive list of dependencies.

```sh
    git clone git@github.com:Uplite/video-service.git
    cd video-service
    make

    // Assuming you have a valid environment:
    ./bin/video-service-writer

    // OR
    ./bin/video-service-reader
```

#### Environment

Assure that all required environment variables are present. See the example `.env.example` file for required environment variables.

#### Code Generation

Generate the gRPC code via the `make generate` target.

#### Compilation

To compile, run the `make build` or `make` target. This will generate code and compile the programs in the [cmd](./cmd) directory.

#### Container

There are Dockerfiles for both the reader and the writer available inside of the [build/package/](./build/package) directory.

#### Testing

To test, run the `make test` or `COUNT=n make test` target. The default `make test` will run with the race detector enabled and bust the test cache each run with a defualt `-count=1` test flag, but receives an argument should you desire to pass one.


[1]: https://go.dev/
[2]: https://www.gnu.org/software/make/
[3]: https://grpc.io/docs/protoc-installation/
[4]: ./makefile
