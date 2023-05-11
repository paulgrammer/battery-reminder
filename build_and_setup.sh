#!/bin/bash

# Install Homebrew if it's not installed
if ! command -v brew &>/dev/null; then
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi

# Check if Go is installed
go_version=$(go version 2>/dev/null)
if [[ $? -ne 0 || $go_version != *"go1.20"* ]]; then
    echo "Go is not installed or not the right version, installing..."
    brew install go@1.20
    echo 'export PATH="/usr/local/opt/go@1.20/bin:$PATH"' >>~/.bash_profile
    source ~/.bash_profile
fi

# Check if Make is installed
if ! command -v make &>/dev/null; then
    echo "Make could not be found, installing..."
    brew install make
fi

# Build the program
make build

# Check if setup.sh is executable
if [ -x setup.sh ]; then
    ./setup.sh
else
    echo "setup.sh does not exist or is not executable"
fi
