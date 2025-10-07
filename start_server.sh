#!/bin/bash

set -e

stop_server () {
    pid=$(ps axu | grep ./enders_game | head -n 1 | grep -oP '^\S+\s+\K\S+')
    lines=$(ps axu | grep ./enders_game | wc -l)
    if [ "$lines" != "1" ]; then
        echo "Killing active server"
        kill "$pid"
    fi
    rm -r ./bin
}

read -p "Start/Restart (1) or Terminate (2)?: " option

if [ "$option" == "2" ]; then
    stop_server
    exit
fi

if [ -d "./bin" ]; then
    stop_server
fi

mkdir bin

echo "Building Ender's Game Server"

cd apps/server/code
go build -o enders_game .
mv enders_game ../../../bin/
cd ../../../bin/

echo "Starting Ender's Game Server"
nohup ./enders_game &

echo "Server successfully started"
exit
