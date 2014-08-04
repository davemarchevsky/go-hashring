#!/bin/bash

cleanup() {
    rm -f nodes keys python_results go_results test_go_hashring
    exit 0
}
trap cleanup INT TERM EXIT


./generate_hashring_test_data.py nodes keys
./test_python_hashring.py nodes keys > python_results
go test

diff python_results go_results

cleanup

