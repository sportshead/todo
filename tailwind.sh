#!/bin/bash
# tailwindcss installer and wrapper
if [ ! -f ./tailwindcss ]; then
    if [[ "$(uname)" == "Darwin" ]]; then
        if [[ "$(arch)" == "arm64" ]]; then
            PLATFORM="macos-arm64"
        else
            PLATFORM="macos-x64"
        fi
    elif [[ "$(uname)" == "Linux" ]]; then
        if [[ "$(arch)" == "armv7l" ]]; then
            PLATFORM="linux-armv7"
        elif [[ "$(arch)" == "aarch64" ]]; then
            PLATFORM="linux-arm64"
        else
            PLATFORM="linux-x64"
        fi
    else
        echo "Unsupported operating system"
        exit 1
    fi

    URL="https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-$PLATFORM"
    echo "Downloading $URL to ./tailwindcss"
    curl -sL -o ./tailwindcss $URL
    chmod +x ./tailwindcss
fi

./tailwindcss "$@"