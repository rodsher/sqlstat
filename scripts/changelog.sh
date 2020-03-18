#!/bin/bash

# Generate CHANGELOG.md file using configuration file
git-chglog --config build/.chglog/config.yml -o CHANGELOG.md v1.0.0..
