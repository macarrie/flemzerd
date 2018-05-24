# Flemzerd [![Build Status](https://travis-ci.org/macarrie/flemzerd.svg?branch=master)](https://travis-ci.org/macarrie/flemzerd/) [![Code Coverage](https://codecov.io/gh/macarrie/flemzerd/branch/master/graph/badge.svg)](https://codecov.io/gh/macarrie/flemzerd)

Flemzerd is an automation tool (like a very lightweight Sonarr) for handling TV Shows.
It watches your tv shows for new episodes, downloads them in the client of your choices, and updates your media center library if needed.

## Current status

Flemzerd is still under heavy developpement/testing.

## What is it ?

Flemzerd is a daemon intended to track and handle your media library. We are often doing the same thing by hand: 
* Look regularly for new episodes of tv show
* Once an episode is available, look for a way to download it
* Launch download in your download client and wait for the download to end
* Move the episode wherever you store your media and manually update your media center (refresh Kodi library for example)

Flemzerd is intended to automate these tasks.


## Dependencies

* go 1.9
* dep (depency manager for go)
* sqlite3
* sass

## Getting started

A Makefile is present at the project root to install/update flemzerd
* Install or update
```bash
make install
make update
```
* Setup and edit configuration file in one of the following locations:
 ** /etc/flemzerd/flemzerd.yml (created by install script)
 ** ~/.config/flemzerd/flemzerd.yml
* Start flemzerd daemon
```bash
systemctl start flemzerd
```
