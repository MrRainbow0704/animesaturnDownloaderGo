#!/usr/bin/env sh

go mod tidy
go build -o ./bin/download-anime ./cmd