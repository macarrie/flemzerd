+++
title = "Mediacenters"
description = "Trigger post download operations on your media center"
date = 2018-05-24T14:52:41Z
weight = 80
draft = false
bref = "Once a media is downloaded, operations may be done on your mediacenter automatically (refresh library per example)"
toc = true
+++

## Mediacenters overview
---

The final element of the flemzerd download chain is the media center. Once a media is downloadeda and moved to the library path, media center library needs to be refreshed. Flemzerd allows to define a mediacenter to perform this kind of refresh.

## Available Mediacenters
---

### Kodi
---
 Use a Kodi instance as media center.

#### How to use
---
* Enable `kodi` MediaCenter in configuration file
{{< highlight toml >}}
[mediacenters]
    [mediacenters.kodi]
        address = "address"
        port = 9090
{{< /highlight >}}
The port for RPC calls that Kodi uses to allow remote control in not presented in the interface. The default port is 9090 (but can be changed in Kodi configuration files).
