+++
title = "Downloaders"
description = "Perform the actual media download"
date = 2018-05-24T14:52:18Z
weight = 60
draft = false
bref = "Downloaders modules are responsible for performing downloads of media tracked by flemzerd."
toc = true
+++

## Downladers overview
---

Flemzerd uses Downloaders to handle torrent downloads. Downloaders in flemzerd are external software that can handle torrent download (Transmission, or BitTorrent for example)

## Available Downloaders
---

### Transmission
---
As of now, only Transmission download client can be used as a flemzerd Downloader. The Transmission daemon can be use locally to the flemzerd daemon or on an external server.

#### How to use
---
* Enable `transmission` Downloader in configuration file
{{< highlight toml >}}
[downloaders]
    [downloaders.transmission]
        address = "localhost"
        port = 9091
        user = "username"
        password = "password"
{{< /highlight >}}
* Setup authentication. If no authentication is needed for Transmission, leave the `user` and `password` keys empty or omit these keys.
