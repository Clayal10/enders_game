#!/bin/bash

stop () {
    pid=$(ps axu | grep ./$1 | head -n 1 | grep -oP '^\S+\s+\K\S+')
    lines=$(ps axu | grep ./$1 | wc -l)
    if [ "$lines" != "1" ]; then
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

    cd cmd/client/code
    go build -o colonel_graph .
    mv colonel_graph ../../../bin/
    cd ../../../bin/

    while true; do
        echo "Starting Client"
        ./colonel_graph > client.out
        echo "Client stopped, restarting in 30 seconds..."
        sleep 30    
    done

    exit
}
start_server(){
    if [ ! -d "./bin" ]; then
        mkdir bin
    fi

    echo "Building Ender's Game Server"

    cd cmd/server/code
    go build -o enders_game .
    mv enders_game ../../../bin/
    cd ../../../bin/

    while true; do
        echo "Starting Ender's Game Server"
        ./enders_game > server.out
        echo "Server stopped, restarting in 30 seconds..."
        sleep 30
    done
}

check=0

while [ $check == 0 ]; do
    read -p "Client (1) or Server (2)?: " option
    if [ $option == 1 ]; then
        stop_option colonel_graph
        start_client
        break
    fi
    if [ $option == 2 ]; then
        stop_option enders_game
        start_server
        break
    fi
done

