#!/bin/bash
set -eux
ex="$1"
out=/tmp/"$(basename "$ex")".png
go run ./bin --topo "$ex" | /usr/bin/dot /dev/stdin -Tpng -o "$out"
xdg-open "$out"
