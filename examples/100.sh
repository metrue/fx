#!/bin/bash

funcs=" "
for i in $(seq 128)
do
    funcs=${funcs}"client/functions/func.go "
done

./client/build/fx up ${funcs}
