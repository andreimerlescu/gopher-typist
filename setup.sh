#!/usr/bin/env bash

go get gioui.org
go get github.com/faiface/beep
go mod tidy
make build
make run
