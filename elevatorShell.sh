#!/bin/bash

gnome-terminal --geometry=49x14+0-1080 -e 'sh -c "go run main.go -port=10001 -id=0"' &&
gnome-terminal --geometry=49x14+1920-1080 -e 'sh -c "go run main.go -port=10002 -id=1"' &&
gnome-terminal --geometry=49x14+1920-0 -e 'sh -c "go run main.go -port=10003 -id=2"' 

#x-terminal-emulator -e go run main.go -port=10001 -geometry 73x31+100+300 &&
#gnome-terminal -e go run main.go -port=10001
#gnome-terminal --geometry=49x14--20+30
#go run main.go -port=10001 
#x-terminal-emulator -e go run main.go -port=10002 &&
#x-terminal-emulator -e go run main.go -port=10003

#open new terminal window with ssh login to machine_name
#Terminal -e ssh <remote_machine>
echo "All good!"
