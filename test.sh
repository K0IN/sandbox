#/bin/sh
go build main.go

sudo chown root:root main
sudo chmod u+s main

./main try "ls"

