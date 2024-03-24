#!/bin/bash

AIR_VERSION="v1.51.0"
GOLANGCI_LINT_VERSION="v1.56.2"

setup() {
    go-mod-download
    go-mod-verify
    go-mod-tidy
    install-air
    install-golangci-lint
}

go-mod-download() {
    echo "Downloading Go module dependencies..."
    go mod download
}

go-mod-verify() {
    echo "Verifying Go module dependencies..."
    go mod verify
}

go-mod-tidy() {
    echo "Tidying Go module dependencies..."
    go mod tidy
}

install-air() {
    echo "Installing air for hot reloading..."
    curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s $AIR_VERSION
}

install-golangci-lint() {
    echo "Installing golangci-lint for linting..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $GOLANGCI_LINT_VERSION
}

setup