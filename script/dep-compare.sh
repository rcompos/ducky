#!/bin/bash
reqsfile="$HOME/tmp/chandler-go-reqs.txt"
compfile="$HOME/tmp/chandler-components.txt"

while IFS= read -r line
do
    # display $line or do somthing with $line
	printf '%s\n' "$line"
	#grep -o "$line" $reqfile 
done <"$compfile"
