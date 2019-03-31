#!/bin/bash

while :
do
    gnome-terminal --geometry=49x14+100-1080 -e 'sh -c "ElevatorServer"'
    sleep 1 
    #echo "Starting the elevator"
    go run main.go -port=15657 -id=$1
done


