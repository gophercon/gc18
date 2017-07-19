GOCMD=go
GOGET=$(GOCMD) get
GOBUILD=$(GOCMD) build
GOBUILDPROD=$(GOBUILD) -ldflags "-linkmode external -extldflags -static" 
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
DOCKER=docker
DOCKERCOMPOSE=docker-compose
SODA=buffalo db
BUFFALO=buffalo

build:
	$(BUFFALO) build -o bin/gc18

buildprod:
	$(GOBUILDPROD) -v -o gc18

clean:
	$(GOCLEAN) -n -i -x
	rm -f $(GOPATH)/bin/gcon
	rm -rf gc18

test:
	$(GOTEST) -v ./grifts -race
	$(GOTEST) -v ./models -race
	$(GOTEST) -v ./actions -race

db-up:
	@echo "Make sure you've run 'make db-setup' before this"
	$(DOCKER) run --name=gc18_db -d -p 5432:5432 -e POSTGRES_DB=gc18_development postgres

db-setup:
	$(DOCKER) run --name=gc18_db -d -p 5432:5432 -e POSTGRES_DB=gc18_development postgres
	sleep 6
	$(BUFFALO) db create -a
	$(BUFFALO) db migrate up
	$(DOCKER) ps | grep gc18_db

db-down:
	$(DOCKER) stop gc18_db
	$(DOCKER) rm gc18_db

