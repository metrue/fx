#!/usr/bin/env bash

set -e

# start fx server
./build/fx serve > server_output 2>&1 &
sleep 20 # waiting fx server to pulling resource done

./build/fx up examples/functions/func.js
./build/fx up examples/functions/func.rb
./build/fx up examples/functions/func.py
./build/fx up examples/functions/func.go
./build/fx up examples/functions/func.php
./build/fx up examples/functions/func.jl
./build/fx up examples/functions/func.java
./build/fx up examples/functions/func.d
