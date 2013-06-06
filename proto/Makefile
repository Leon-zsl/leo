TARG=.

GOFILES=\
	account.pb.go\
	world.pb.go

all: build

build: $(GOFILES)

%.pb.go: %.proto
	../tool/protoc --go_out=$(TARG) $<
