fmt:
	gofmt -w=true -s $(find . -type f -name '*.go' -not -path "./vendor/*")

clean:
	rm -rf _dist/*

run:
	go build -i -o _dist/trovehero cmd/game/game.go && ./_dist/trovehero

profile:
	go build -i -o _dist/trovehero cmd/profiling/profiling.go && ./_dist/trovehero -cpuprofile=_dist/trovehero.prof

analyze:
	go tool pprof _dist/trovehero cmd/profiling/profiling.go _dist/trovehero.prof -http