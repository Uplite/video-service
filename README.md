# Video Service

A service for managing uploads and fetches for video content.

## Quick Start

Must have [Go compiler][1], [GNU Make][2], and [protoc][3]. Additional dependencies required for code generation will be installed via the make targets if they are not present in your $GOPATH. See the [makefile][4] for an exhaustive list of dependencies.

#### Environment

Assure that all required environment variables are present. See the example `.env.example` file for required environment variables.

#### Code Generation

Generate the gRPC code via the `make generate` target

#### Compilation

To compile, run the `make build` or `make` target. This will generate code and compile the programs in the [cmd](./cmd) directory.


[1]: https://go.dev/
[2]: https://www.gnu.org/software/make/
[3]: https://grpc.io/docs/protoc-installation/
[4]: ./makefile
