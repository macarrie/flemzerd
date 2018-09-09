+++
title = "Getting Started"
date = 2018-05-23T14:33:39Z
weight = 10
draft = false
description = "How to get a flemzerd instance running from scratch"
bref = "This page describes how to get a flemzerd instance running, starting from scratch, including pre-requisites and build instructions"
toc = true
+++


## Short presentation
---

Before describing how to get it running, let's describe in a few lines what flemzerd is about:

> flemzerd is a tool that automates your media download process. It also can be used as a notifier for tracked shows.

The goal of this tool is to ease away the unpleasing following sequence of tasks:

* Watch regularly for each the TV shows you are watching if a new episode  has aired
* Find a torrent for the episode
* Add torrent to your download client and look after the download
* When the download is finished, moved the downloaded files to where you store your media files
* Refresh your media center to see the downloaded files

## Install
---

* Download the latest release from GitHub: [Releases](https://github.com/macarrie/flemzerd/releases)
* Launch the install via the Makefile:
{{< highlight bash >}}
make install
{{< /highlight >}}

## Update
---

* Download the latest release from GitHub: [Releases](https://github.com/macarrie/flemzerd/releases)
* Launch the install via the Makefile:
{{< highlight bash >}}
make update
{{< /highlight >}}

## Using flemzerd
---

#### Setup environment
---

Flemzerd acts as an orchestrator for different external services (download client, torrent indexer, watchlist). For flemzerd to function properly, all these services must be started and available from the server running flemzerd, and then configured into flemzerd.

The services/configuration needed are the following:

* A Watchlist (manual or via trakt.tv)
* A Provider (TMDB or TVDB with API key configured)
* A download client (transmission)
* _(Optional)_ Notifiers (Telegram, Pushbullet)
* _(Optional)_ A Media center (kodi)

More information about general concepts of flemzerd, configuration options and modules can be found of the rest of this documentation.

#### How to start
---

* As a service (via systemctl)
{{< highlight bash >}}
systemctl start flemzerd
{{< /highlight >}}
    * To enable at startup
    {{< highlight bash >}}
    systemctl enable flemzerd
    {{< /highlight >}}
* As a standalone binary
{{< highlight bash >}}
/usr/bin/flemzerd -d
{{< /highlight >}}

## Build from source
---

If released packages are not for you (wrong architecture, want to follow latest dev), flemzerd can be built from source.

The following dependencies are needed:

* Go 1.11
* Nodejs
* npm
* Rust (for external dependency: [Vidocq](https://github.com/macarrie/vidocq))

To build, simply use the Makefile:
{{< highlight bash >}}
make build
{{< /highlight >}}

This will generate a package under the `package/` folder, which is the same as the ones that can be downloaded as releases on Github. Install and update work as described previously.
