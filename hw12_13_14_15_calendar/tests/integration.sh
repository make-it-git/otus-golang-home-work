#!/bin/bash

set -o pipefail
set -o errexit

uuid=`uuidgen`
curl --fail -d '{"id": "'$uuid'", "startTime": "2022-01-01T15:00:00.0Z", "endTime": "2022-01-01T16:00:00.0Z", "title": "test", "ownerId": 1, "description": "test"}' 127.0.0.1:8080/events/ > /dev/null
curl --fail -X PUT -d '{"startTime": "2022-01-10T15:00:00.0Z", "endTime": "2022-01-10T16:00:00.0Z", "title": "test 2", "ownerId": 1, "description": "test 2"}' 127.0.0.1:8080/events/$uuid > /dev/null
curl --fail 127.0.0.1:8080/events/day/2022-01-10 > /tmp/http.response

~/go/bin/grpcurl -plaintext -d '{"date": "2022-01-10T16:00:00.0Z"}' 127.0.0.1:8081 event.EventService/ListDay > /tmp/grpc.response

grep $uuid /tmp/http.response
grep $uuid /tmp/grpc.response

rm -f /tmp/http.response /tmp/grpc.response

echo "PASSED"