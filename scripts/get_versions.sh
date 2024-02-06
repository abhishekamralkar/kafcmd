#!/bin/bash

# Logic to determine the semantic version based on commit messages
# Example: Extract the first version-like string from commit messages
version=$(git log --merges --format=%B -n 1 | grep -o -P '(?<=Version: )(\d+\.\d+\.\d+)')

if [ -z "$version" ]; then
  # Default to a version if none is found
  version="1.0.0"
fi

echo "$version"

