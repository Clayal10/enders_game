#!/bin/bash

set -e

read -p "Kill active server? (y/n): " option

if [[ "$option" == "y" || "$option" == "Y" ]]; then
    echo "Killing active server"
    pid=$(ps | grep enders_game | grep -o '[0-9]\+ ' | head -n 1)
    kill "$pid"
    rm -r ./bin
    exit
fi

mkdir bin

echo "Building Ender's Game Server"

cd apps/server/config
cp Config.json ../../../bin/

cd ../code
go build -o enders_game .
mv enders_game ../../../bin/
cd ../../../bin/

echo "Starting Ender's Game Server"
./enders_game &

echo "Server successfully started"
exit