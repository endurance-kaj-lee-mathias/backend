@default:
    just --list

deps:
    go mod tidy

run:
    go run cmd/*.go

up:
    goose up

down:
    goose down
