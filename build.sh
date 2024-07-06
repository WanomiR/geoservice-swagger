#!/bin/bash
cd proxy
swag init
cd ..
docker-compose up --force-recreate --build
