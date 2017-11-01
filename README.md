# Flemzerd

Flemzerd is an automation tool (like a very lightweight Sonarr) for handling TV Shows.
It watches your tv shows for new episodes, downloads them in the client of your choices, and updates your media center library if needed.

# Current status

Flemzerd is still under heavy developpement. It is absolutely not ready for use yet.

# What is it ?

TODO

# Setup

The only way to get flemzerd yet is to build it for the sources.
For this, you will need to have go 1.9 installed and $HOME/go/bin in your PATH

```bash
# Clone repo
git clone github.com/macarrie/flemzerd ~/go/src/github.com/macarrie/flemzerd
cd ~/go/src/github.com/macarrie/flemzerd

# Get dep tool to install dependencies
go get -u github.com/golang/dep/cmd/dep

# Install dependencies
dep ensure

# Install flemzerd
go install
```

# Usage

```
Usage of flemzerd:
    -h, --help: Shows this help message
    -c, --config="": Configuration file path to use
    -d, --debug=false: Start in debug mode
```

# Configuration

A sample configuration file is present in the repo (flemzerd.yaml).
TODO: Explanation on what to put in the configuration file and available paths
