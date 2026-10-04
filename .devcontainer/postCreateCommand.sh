#!/bin/bash

echo "Hello, World!"

# Dependencies for building/running ebitengine.
sudo apt update -y
sudo apt install libx11-6 \
    libgl1-mesa-glx \
    libgles2-mesa-dev \
    -y