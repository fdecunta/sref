DOIS="$1"
FILE="$2"

cat "$1" | xargs -I % sref -file="$2" -input=%
