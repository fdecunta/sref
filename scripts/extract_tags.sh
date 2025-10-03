#!/usr/bin/bash

FILE="$1"

if [ ! -f "$FILE" ]; then
    echo "Error: "$1" is not a file"
    exit
fi

pandoc "$1" -t plain | grep -oE '@[[:alnum:]_]+[0-9]{4}[a-zA-Z]?'
