#!/bin/bash

if pgrep translator &> /dev/null; then
        killall translator
fi

./translator &>> translator.log &
