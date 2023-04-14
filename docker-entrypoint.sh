#!/bin/bash

dlv debug ./portfolio-context --output=./portfolio-context/__debug_bin --listen=0.0.0.0:9090 --accept-multiclient --continue --headless --allow-non-terminal-interactive & 
dlv debug ./broker-context --output=./broker-context/__debug_bin --listen=0.0.0.0:9091 --accept-multiclient --continue --headless --allow-non-terminal-interactive & 
wait -n
  
# Exit with status of process that exited first
exit $?
