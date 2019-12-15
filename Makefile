# Copyright (c) 2019, Sylabs Inc. All rights reserved.
# This software is licensed under a 3-clause BSD license. Please consult the
# LICENSE.md file distributed with the sources of this project regarding your
# rights to use or distribute this software.

all: syvalidate

syvalidate: 
	cd cmd/syvalidate; go build syvalidate.go

install: all
	go install ./...

test: install
	go test ./...

uninstall:
	@rm -f $(GOPATH)/bin/syvalidate

clean:
	@rm -f main
	@rm -f cmd/syvalidate/syvalidate 

distclean: clean
