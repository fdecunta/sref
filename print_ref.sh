FILE=$1
DOI=$2

jq --arg doi "$DOI" '.[$doi]' "$FILE"
