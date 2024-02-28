#!/bin/bash

while :
do
	inotifywait -e MODIFY -r * 2>/dev/null
	sleep 1
	echo "===================================== go test -v"
	go test -v
	echo -e "===================================== " $(date) "\n\n\n"
done
