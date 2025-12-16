.PHONY: run build

run:
	go run ./cmd/Raidplanner

build:
	go build -o civapi ./cmd/Raidplanner