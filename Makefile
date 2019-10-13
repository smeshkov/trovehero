fmt:
	gofmt -w=true -s $(find . -type f -name '*.go' -not -path "./vendor/*")

clean:
	rm -rf dist/*

run:
	go build -i -o dist/trovehero && ./dist/trovehero