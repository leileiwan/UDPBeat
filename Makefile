PROGRAME=sense-beat
VERSION=v1
PREFIX?=.
BASE_DESCRIPTION=The UDP Beat for k8s operaoter
PROGRAMEPATH=github.com/sense-beat
SERVERNAME=server
CLIENTNAME=client

.PHONY:build
build:clean build-server build-client

# build server
.PHONY:build-server
build-server:
	@echo $(BASE_DESCRIPTION)
	go build -o $(PROGRAME)$(PREFIX)$(SERVERNAME)$(PREFIX)$(VERSION) cmd/$(PROGRAME)$(PREFIX)$(VERSION)/$(SERVERNAME)/main.go

# build client
.PHONY:build-client
build-client:
	@echo $(BASE_DESCRIPTION)
	go build -o $(PROGRAME)$(PREFIX)$(CLIENTNAME)$(PREFIX)$(VERSION) cmd/$(PROGRAME)$(PREFIX)$(VERSION)/$(CLIENTNAME)/main.go


# clean
.PHONY:clean
clean:
ifeq ($(PROGRAME)$(PREFIX)$(SERVERNAME)$(PREFIX)$(VERSION),$(wildcard $(PROGRAME)$(PREFIX)$(SERVERNAME)$(PREFIX)$(VERSION)))
	rm $(PROGRAME)$(PREFIX)$(SERVERNAME)$(PREFIX)$(VERSION)
endif
ifeq ($(PROGRAME)$(PREFIX)$(CLIENTNAME)$(PREFIX)$(VERSION),$(wildcard $(PROGRAME)$(PREFIX)$(CLIENTNAME)$(PREFIX)$(VERSION)))
	rm $(PROGRAME)$(PREFIX)$(CLIENTNAME)$(PREFIX)$(VERSION)
endif


.PHONY:test
test:test_UDPBeat

.PHONY:test_UDPBeat
test_UDPBeat:
	
	@echo ............begin test the package UDPBeat............
	@go test -v github.com/sense-beat/pkg/UDPBeat
	@echo ............end test the package UDPBeat............
	@echo
	@echo