#!/bin/bash

cd src/EventManager
go install
cd ../..

if [ -e pids/EventManager.pid ]; then
    kill $(cat pids/EventManager.pid)
fi

bin/EventManager
