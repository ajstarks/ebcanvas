#!/bin/bash
for i in $(cat cl)
do
	cd $i
	echo -n "$i "
	go build $* -ldflags="-s -w" . 2>/dev/null
	cd ..
done
echo
