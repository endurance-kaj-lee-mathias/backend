@default:
    just --list

deps:
    go mod tidy

run:
    go run cmd/*.go
