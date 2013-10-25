#!/bin/bash

cleanup() {
    rm -f nodes keys python_results go_results test_go_hashring
    exit 0
}
trap cleanup INT TERM EXIT


# Regression tests.
./generate_hashring_test_data.py nodes keys
go build -o test_go_hashring test_go_hashring.go

./test_python_hashring.py nodes keys > python_results
./test_go_hashring nodes keys > go_results

diff python_results go_results

# Unit tests.
cd hashring
go test

cleanup


