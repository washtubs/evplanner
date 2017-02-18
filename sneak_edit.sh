#!/usr/bin/env zsh

address=$1

tmpfile=$(mktemp /tmp/ev-planner-XXXXX)

curl $address/lock -X PUT
curl -s $address/read > $tmpfile

echo "sneak edit" > $tmpfile

curl $address/write -X PUT --data-binary "@$tmpfile"
curl $address/unlock -X PUT
rm $tmpfile
