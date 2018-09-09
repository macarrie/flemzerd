+++
title = "Development setup"
date = 2018-05-23T14:33:39Z
weight = 90
draft = false
description = "How to contribute to flemzerd"
bref = "This page describes how to setup a flemzerd development environment"
toc = true
+++


## Code organization
---

The flemzer code is organized as described below:

* `configuration/`: Configuration package. Responsible for configuration loading and parsing.
* `db/`: Interface to database. The database is a sqlite3 located in `/var/lib/flemzerd/flemzer.db`
* `docs/`: Compiled documentation files
* `docs\_src/`: Documentation source files
* `downloaders/`: Downloader package. Contains downloader collection handling and different downloader implementations
* `helpers/`: Various helpers
* `indexers/`: Indexer package. Contains indexer collection handling and different indexers implementations
* `install/`: Setup files for delivery and package creation. Contains install scripts and configuration files
* `logging/`: Logger package.
* `mediacenters/`: Mediacenters package. Contains mediacenter collection handling and different mediacenters implementations
* `mocks/`: Mock implementations of modules for tests.
* `notifiers/`: Notifiers package. Contains notifier collection handling and different notifiers implementations
* `objects/`: Objects declarations used all over the app
* `providers/`: Provider package. Contains providers collection handling and different providers implementations
* `scheduler/`: Scheduler packages. Brings all modules together to watch and download media.
* `server/`: Server setup for webui
* `testdata/`: Mock configuration files for unit tests
* `vidocq/`: Vidocq package. Handles interaction with external dependency Vidocq.
* `watchlists/`: Watchlist package. Contains watchlist collection handling and different watchlists implementations
* `go.mod`: Dependency declaration files
* `main.go`: Entry point of the app. Parse command line flags, and launches scheduler and webserver.
* `Makefile`


## Running the backend server
---
 The Makefile contains rules to run the server while listening to changes. Any change in the go code source files recompiles and restarts the server:
{{< highlight bash >}}
make start
{{< /highlight >}}


## Running the UI dev server
---

From the `server/ui/` folder, launch the webserver:
{{< highlight bash >}}
npm start
{{< /highlight >}}


This listens to changes in the `server/ui/` and recompiles the Web UI at every change.

## Update documentation
---

From the `docs\_src/` folder, launch the documentaton web server with the following command:
{{< highlight bash >}}
hugo server
{{< /highlight >}}
