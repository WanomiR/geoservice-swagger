#!/bin/bash
cd proxy && swag init -g ./cmd/main.go && cd ..
docker compose up --force-recreate --build
