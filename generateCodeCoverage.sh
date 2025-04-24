#!/bin/sh

REMOVEOUT=false

go test ./... -coverprofile=coverage.out >/dev/null

while getopts "ir" flag; do
 case "${flag}" in
   i) go-ignore-cov --file coverage.out ;; # Remove all ignored blocks from coverage output
   r) REMOVEOUT=true ;; # Remove coverage file when done
 esac
done


while read p || [ -n "$p" ]
do
sed -i '' "/${p//\//\\/}/d" ./coverage.out
done < ./exclude-from-code-coverage.txt

go tool cover -html=coverage.out

if $REMOVEOUT ; then rm coverage.out; fi