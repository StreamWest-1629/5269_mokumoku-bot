#!/bin/bash

cd `dirname $0`

echo "converting start!"

# Set path
mkdir ../resource
origin_dir_gen=$(cd ../resource-origin && pwd)
result_dir_gen=$(cd ../resource && pwd)
echo $origin_dir

# Search files
dirs=`find $origin_dir_gen -type d`
files=`find $origin_dir_gen -type f`

for origin_dir in $dirs; do
    result_dir=${origin_dir/$origin_dir_gen/$result_dir_gen}
    mkdir "$result_dir"
done

for origin_file in $files; do
    rel_path=${origin_path#origin_dir_gen}
    result_file=${origin_file/$origin_dir_gen/$result_dir_gen}
    origin_file_ext=${origin_file##*.}
    result_file_noext=${result_file%.*}

    case "$origin_file_ext" in
    "json"|"jpg"|"png" )
        echo "copying file($rel_path)..."
        cp "$origin_file" "$result_file"
        ;;
    "mp3"|"wav"|"org"|"flac" )
        echo "converting to 'dca' file..."
        result_file_dca="${result_file_noext}.dca"
        ffmpeg -i "${origin_file}" -f s16le -ar 48000 -ac 2 pipe:1 | dca > "${result_file_dca}"
        echo "converted: $result_file_dca ($result_file, $result_file_noext)"
        ;;
    esac
done