#!/bin/bash

for start in $(find ~/.vim -name start); do
    for dir in $start/*; do
        echo $dir
        cd $dir
        git pull
    done
done
