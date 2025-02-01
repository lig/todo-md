#!/bin/bash

PROJECT_DIR="$(dirname $0)/.."

go build -o $PROJECT_DIR/build/ $PROJECT_DIR/...
