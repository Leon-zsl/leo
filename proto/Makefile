TARG=.

GOFILES=\
	common.pb.go\
	account.pb.go\
	world.pb.go\
	module.pb.go

all: build

build: $(GOFILES)

%.pb.go: %.proto
	protoc --go_out=$(TARG) $<
