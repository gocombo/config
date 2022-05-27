
cover_dir=.cover
cover_profile=${cover_dir}/profile.out
cover_html=${cover_dir}/coverage.html

.DEFAULT_GOAL := all

all: test

bin/golangci-lint: .golangci-version
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(shell cat .golangci-version)

lint: bin/golangci-lint
	bin/golangci-lint run

${cover_dir}:
	mkdir -p ${cover_dir}

test: lint ${cover_dir}
	go test -coverprofile=${cover_profile} ./...
	go tool cover -html=${cover_profile} -o ${cover_html}
