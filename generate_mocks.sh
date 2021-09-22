#!/bin/bash

set -o errexit

mockgen -source mocks/mocks.go -destination mocks/generated_mocks.go -package mocks
