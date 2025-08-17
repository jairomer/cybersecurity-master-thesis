#!/bin/bash

export TOKEN=$(curl -X POST -v "http://localhost:8000/login" -H "Content-Type: application/json" --data '{"user":"test", "password":"test"}' | jq .token)
export TOKEN=${TOKEN//\"/}

curl -v -H "Authorization: Bearer $TOKEN" -X GET "http://localhost:8000/hello/world"

curl -v -H "Authorization: Bearer $TOKEN" -X GET "http://localhost:8000/hello/world?country=France"

