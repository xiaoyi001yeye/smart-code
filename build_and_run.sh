#!/bin/bash

docker compose down smartcodeql
docker compose rm smartcodeql
docker compose build --no-cache smartcodeql
docker compose up -d