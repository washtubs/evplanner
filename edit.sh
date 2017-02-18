#!/usr/bin/env zsh

address=$1

original=$(mktemp /tmp/ev-planner-XXXXX)
tmpfile=$(mktemp /tmp/ev-planner-XXXXX)

maybe_locked=$(curl -s -o /dev/null -w "%{http_code}" $address/lock -X PUT)
if [ $maybe_locked = 433 ]; then
    echo Already locked for editing. Aborting.
    exit 1
fi
if [ $maybe_locked = 403 ]; then
    echo Took too long to lock. Aborting.
    exit 1
fi
curl -s $address/read > $original
cat $original > $tmpfile
$EDITOR $tmpfile
response=$(curl -s -o /dev/null -w "%{http_code}" $address/write -X PUT --data-binary "@$tmpfile")
echo response $response
if [ $response = 403 ]; then
    # no longer locked
    curl $address/lock -X PUT
    tmpfile2=$(mktemp /tmp/ev-planner-XXXXX)
    curl -s $address/read > $tmpfile2
    if { diff $original $tmpfile2 &> /dev/null }; then
        echo "Timed out but the file didn't change so just push it again."
        curl -s $address/write -X PUT --data-binary "@$tmpfile"
    else
        echo "You took too long and the file differs. Try to merge the changes. Press enter to bring up a vimdiff"
        read ok
        vimdiff $tmpfile $tmpfile2
        curl -s $address/write -X PUT --data-binary "@$tmpfile"
    fi
    rm $tmpfile2
    curl -s $address/unlock -X PUT
else
    curl -s $address/unlock -X PUT
fi
rm $tmpfile
rm $original
