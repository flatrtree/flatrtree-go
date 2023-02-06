SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.PHONY: default clean proto


GO_PROTO_PATH=internal/flatrtree.pb.go
GO_VTPROTO_PATH=internal/flatrtree_vtproto.pb.go

default: $(GO_VTPROTO_PATH)

clean:
	rm $(GO_PROTO_PATH) $(GO_VTPROTO_PATH)

$(GO_VTPROTO_PATH): $(GO_PROTO_PATH)

$(GO_PROTO_PATH):
	protoc -I=spec/ \
		--go_out=. --plugin=protoc-gen-go=$(GOBIN)/protoc-gen-go \
		--go-vtproto_out=. --plugin protoc-gen-go-vtproto="$(GOBIN)/protoc-gen-go-vtproto" \
		--go-vtproto_opt=features=unmarshal+size \
		flatrtree.proto
