#!/usr/bin/bash

my_query=$(echo "$1" | jq -sRr @uri)
curl "https://api.crossref.org/works?query=$my_query" | jq '.message.items[0]'
