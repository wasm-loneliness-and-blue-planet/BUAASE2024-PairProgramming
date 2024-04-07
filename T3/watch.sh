#!/bin/bash

while true
do
  echo "Running pnpm build..."
  pnpm run build

  if [ $? -ne 0 ]; then
    echo "Build failed! Aborting."
    exit 1
  fi
  
  echo "Build successful. Running tests..."
  pnpm run test
  
  if [ $? -ne 0 ]; then
    echo "Tests failed!"
  else
    echo "Tests passed."
  fi
  
  echo "Waiting 5 seconds before next run..."
  sleep 5
done
