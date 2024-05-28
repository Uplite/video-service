BINDIR := ./bin
WRITER := video-service-writer
READER := video-service-reader

GOBIN    := $(shell go env GOPATH)/bin
GOSRC    := $(shell find . -type f -name '*.go' -print) go.mod go.sum
PROTOSRC := $(shell find . -type f -name '*.proto' -print)

GOGEN   := $(GOBIN)/protoc-gen-go
GOGRPC  := $(GOBIN)/protoc-gen-go-grpc
VTPROTO := $(GOBIN)/protoc-gen-go-vtproto

PROTODIR := ./api
PROTOGEN := $(PROTODIR)/pb
PROTODEF := $(patsubst $(PROTODIR)/%,%,$(PROTOSRC))

VTFLAGS := marshal+unmarshal+size
LDFLAGS := -w -s

COUNT ?= 1

# -----------------------------------------------------------------
#  build

.PHONY: all
all: build

.PHONY: build-writer
build-writer: $(BINDIR)/$(WRITER)

.PHONY: build-reader
build-reader: $(BINDIR)/$(READER)

.PHONY: build
build: $(BINDIR)/$(READER) $(BINDIR)/$(WRITER)

$(BINDIR)/$(READER): $(GOSRC)
	go build -trimpath -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(READER) ./cmd/$(READER)

$(BINDIR)/$(WRITER): $(GOSRC)
	go build -trimpath -ldflags '$(LDFLAGS)' -o $(BINDIR)/$(WRITER) ./cmd/$(WRITER)

# -----------------------------------------------------------------
#  test

.PHONY: test
test:
	go test -race -v -count=$(COUNT) ./...

# -----------------------------------------------------------------
#  generate

.PHONY: generate
generate: $(VTPROTO) $(GOGEN) $(GOGRPC) $(PROTODIR)/pb/.protogen

$(PROTOGEN)/.protogen: $(PROTOSRC)
	protoc --proto_path=$(PROTODIR)                                  \
		--go_out=.         --plugin protoc-gen-go=$(GOGEN)           \
		--go-grpc_out=.    --plugin protoc-gen-go-grpc=$(GOGRPC)     \
		--go-vtproto_out=. --plugin protoc-gen-go-vtproto=$(VTPROTO) \
		--go-vtproto_opt=features=$(VTFLAGS)                         \
		$(PROTODEF)
	@touch $(PROTOGEN)/.protogen

# -----------------------------------------------------------------
#  dependencies

$(GOGEN):
	( cd /; go install google.golang.org/protobuf/cmd/protoc-gen-go@latest)

$(GOGRPC):
	( cd /; go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest)

$(VTPROTO):
	( cd /; go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest)

# -----------------------------------------------------------------
#  misc

.PHONY: clean
clean:
	rm -rf $(PROTOGEN)
