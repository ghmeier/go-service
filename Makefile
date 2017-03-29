test:
	go build
	gocov test ./... | gocov report

lint:
	go lint ./...

get-deps:
	godep restore
	godep get
	go get github.com/mattn/goveralls
	go get github.com/axw/gocov/gocov
	go get github.com/stretchr/testify
