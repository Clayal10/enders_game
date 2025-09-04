#!/bin/bash

set -e

read -p "Kill active server? (y/n): " option

if [ "$option" == "y" || "$option" == "Y" ]; then
    echo "Killing active server"
    oldPID=$(cat pid.txt)
    kill oldPID
    exit
fi

echo "Building Enders Game Server"

cd apps/server/config
cp Config.json ../../../

cd ../code
go build -o enders_game .
mv enders_game ../../../
cd ../../../

pID=$(./enders_game &)

echo pID > pid.txt

echo "Server successfully started"