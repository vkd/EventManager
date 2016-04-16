#!/bin/bash

cd src/EventManager
go install
cd ../..

if [ -e EventManager.pid ]; then
    kill $(cat EventManager.pid)
fi

bin/EventManager
