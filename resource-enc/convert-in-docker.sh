#!/bin/bash

cd `dirname $0`

echo "converting start!"

# Set path
mkdir ../resource
origin_dir=$(cd ../resource-origin && pwd)
result_dir=$(cd ../resource && pwd)
echo $origin_dir

# Search files
files=`find $origin_dir -type f`
for origin_file in $files; do
    rel_path=${origin_path#origin_dir}
    result_file=${origin_file/$origin_dir/$result_dir}
    origin_file_ext=${origin_file##*.}
    result_file_noext=${result_file%.*}

    case "$origin_file_ext" in
    "json"|"jpg"|"png" )
        echo "copying file($rel_path)..."
        cp origin_file result_file
        ;;
    "mp3"|"wav"|"org" )
        echo "converting to 'dca' file..."
        result_file=${result_file_noext}.dca
        ffmpeg -i $origin_file -f s16le -ar 48000 -ac 2 pipe:1 | dca > test.dca
        ;;
    esac
done