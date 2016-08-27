all: bin/greeter_client bin/greeter_multiplex_server

setup-gopath:
	mkdir -p .gopath
	if [ ! -L .gopath/src ]; then ln -s "$(CURDIR)/vendor" .gopath/src; fi
	if [ ! -L .gopath/src/github.com/gdm85/grpc-go-multiplex ]; then ln -s "$(CURDIR)" .gopath/src/github.com/gdm85/grpc-go-multiplex; fi

bin/add-torrent: setup-gopath
	mkdir -p bin
	cd examples && GO15VENDOREXPERIMENT=1 GOPATH="$(CURDIR)/.gopath" GOBIN="$(CURDIR)/bin/" go install add-torrent.go

bin/greeter_multiplex_server: setup-gopath
	mkdir -p bin
	cd greeter_multiplex_server && GO15VENDOREXPERIMENT=1 GOPATH="$(CURDIR)/.gopath" GOBIN="$(CURDIR)/bin/" go install greeter_multiplex_server.go

bin/greeter_client: setup-gopath
	mkdir -p bin
	cd greeter_client && GO15VENDOREXPERIMENT=1 GOPATH="$(CURDIR)/.gopath" GOBIN="$(CURDIR)/bin/" go install greeter_client.go

clean:
	rm -f bin/greeter_client bin/greeter_multiplex_server

.PHONY: all clean
