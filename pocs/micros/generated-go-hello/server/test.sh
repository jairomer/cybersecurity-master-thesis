#!/bin/bash

export TOKEN=$(curl -X POST -v "http://localhost:8000/login" -H "Content-Type: application/json" --data '{"user":"test1", "password":"test"}' | jq .token)
export TOKEN=${TOKEN//\"/}

# Assert successful access
curl -v -H "Authorization: Bearer $TOKEN" -X GET "http://localhost:8000/hello/world"

curl -v -H "Authorization: Bearer $TOKEN" -X GET "http://localhost:8000/hello/world?country=France"

export TOKEN=$(curl -X POST -v "http://localhost:8000/login" -H "Content-Type: application/json" --data '{"user":"test2", "password":"test"}' | jq .token)
export TOKEN=${TOKEN//\"/}

# Assert unsuccessful access
curl -v -H "Authorization: Bearer $TOKEN" -X GET "http://localhost:8000/hello/world"
curl -v -H "Authorization: Bearer $TOKEN" -X GET "http://localhost:8000/hello/world?country=France"
