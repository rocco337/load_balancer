#!/bin/bash

for ((request=1;request<=10;request++))
do
    for ((x=1;x<=5;x++))
    do
        curl --cookie "SessionId=12345" localhost:8090 &
    done
done

wait