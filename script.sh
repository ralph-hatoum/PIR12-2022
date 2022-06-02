#!/bin/bash
ssh pi2
cd ./server_pir
./server &
exit

ssh pi1
cd ./client_pir
./client &

for var in pi3 pi4 pi5 pi6 pi7 pi8 pi9 pi10
do
    ssh $var
    cd ./client_pir
    ./client &
done
