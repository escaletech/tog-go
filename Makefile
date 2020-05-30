TEST ?= ./...

GOCMD=$(if $(shell which richgo),richgo,go)

test:
	$(GOCMD) test -v $(TEST) -race -covermode=atomic

test-ci:
	$(GOCMD) test -v $(TEST) -race -covermode=atomic -coverpkg=./... -coverprofile=coverage.out

test-watch:
	reflex -s --decoration=none -r \.go$$ -- make test TEST=$(TEST)

release:
	@bash -c "$$(curl -s https://raw.githubusercontent.com/escaletech/releaser/master/tag-and-push.sh)"
