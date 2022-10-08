#!/bin/bash

set -o pipefail
set -o errexit

uuid=`uuidgen`
now_date=`date '+%Y-%m-%d'`
now=`date '+%Y-%m-%dT%H:%M:%S.0Z'`
later=`date -d '+3 hours' '+%Y-%m-%dT%H:%M:%S.0Z'`
notify_at=`date -d '-1 hours' '+%Y-%m-%dT%H:%M:%S.0Z'`
curl --fail -d '{"id": "'$uuid'", "startTime": "'$now'", "endTime": "'$later'", "title": "test", "ownerId": 1, "description": "test"}' 127.0.0.1:8080/events/ > /dev/null
curl --fail -X PUT -d '{"startTime": "'$now'", "endTime": "'$later'", "title": "test 2", "ownerId": 1, "description": "description 2", "notificationTime":"'$notify_at'"}' 127.0.0.1:8080/events/$uuid > /dev/null
curl --fail 127.0.0.1:8080/events/day/$now_date > /tmp/http.response

~/go/bin/grpcurl -plaintext -d '{"date": "'$now'"}' 127.0.0.1:8081 event.EventService/ListDay > /tmp/grpc.response

grep $uuid /tmp/http.response
grep $uuid /tmp/grpc.response

rm -f /tmp/http.response /tmp/grpc.response

echo "PASSED"