#!/bin/bash
cd $(dirname $(readlink -f $0))
go test ./... && go test ./... -bench=. -benchmem