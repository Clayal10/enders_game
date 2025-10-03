#!/bin/bash

set -e

stop () {
    pid=$(ps axu | grep ./$1 | head -n 1 | grep -oP '^\S+\s+\K\S+')
    lines=$(ps axu | grep ./$1 | wc -l)
    if [ "$lines" == "2" ]; then
        echo "Killing $1"
        kill "$pid"
    fi
}

stop_option(){
    read -p "Start/Restart (1) or Terminate (2)?: " option

    if [ "$option" == "2" ]; then
        stop $1
        exit
    fi

    if [ -d "./bin" ]; then
        stop $1
    fi
}

start_client(){
    if [ ! -d "./bin" ]; then
        mkdir bin
    fi
    

    echo "Building Client"

    cd apps/client/code
    go build -o enders_game_client .
    mv enders_game_client ../../../bin/
    cd ../../../bin/

    echo "Starting Client"
    nohup ./enders_game_client &

    echo "Client successfully started"
    exit
}
start_server(){
    if [ ! -d "./bin" ]; then
        mkdir bin
    fi

    echo "Building Ender's Game Server"

    cd apps/server/code
    go build -o enders_game .
    mv enders_game ../../../bin/
    cd ../../../bin/

    echo "Starting Ender's Game Server"
    nohup ./enders_game &

    echo "Server successfully started"
    exit
}

check=0

while [ $check == 0 ]; do
    read -p "Client (1) or Server (2)?: " option
    if [ $option == 1 ]; then
        stop_option enders_game_client
        start_client
        break
    fi
    if [ $option == 2 ]; then
        stop_option enders_game
        start_server
        break
    fi
done






