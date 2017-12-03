# Flemzerd [![Build Status](https://travis-ci.org/macarrie/flemzerd.svg?branch=master)](https://travis-ci.org/macarrie/flemzerd/)

Flemzerd is an automation tool (like a very lightweight Sonarr) for handling TV Shows.
It watches your tv shows for new episodes, downloads them in the client of your choices, and updates your media center library if needed.

## Current status

Flemzerd is still under heavy developpement. It is absolutely not ready for use yet.

## What is it ?

Flemzerd is a daemon intended to track and handle your tv show library. We are often doing the same thing by hand: 
* Look regularly for new episodes
* Once an episode is available, look for a way to download it
* Launch download in your download client and wait for the download to end
* Move the episode wherever you store your media and manually update your media center (refresh Kodi library for example)

Flemzerd is intended to automate these tasks.

### Anatomy

Flemzerd works by scheduling tasks between 4 different modules
* **Providers**: Retrieves informations about shows and episodes (Show informations, seasons and episodes details and air dates for episodes). Providers use external services, such as thetvdb, trakt or imdb for example.
* **Indexers**: Indexers looks for and provide torrent links for an episode. Once a new episode is found, flemzerd asks indexers to find a corresponding torrent.
* **Downloaders**: Downloaders are responsible to download torrents provided by the indexer
* **Notifiers**: Notifiers are responsible for alerting the user when an action is done by flemzerd. Notifications can occur when a new episode is found, has been downloaded, or when an error occurs in the episode processing chain.

Each module must be configured in the configuration file for the daemon to work fully.

## Setup

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

# Build flemzerd binary
go build

# Install and run
sudo ./install.sh
systemctl start flemzerd

# Enable service at startup if needed
systemctl enable flemzerd
```

## Usage

```
Usage of flemzerd:
    -h, --help: Shows this help message
    -c, --config="": Configuration file path to use
    -d, --debug=false: Start in debug mode
```

## Configuration

A sample configuration file is present in the repo (flemzerd.yaml).
This file is written in a yaml format and is hopefully self-explanatory enough. It contains sections about each module type (providers, notifiers, indexers, downloaders) that have to filled to enable corresponding modules.
