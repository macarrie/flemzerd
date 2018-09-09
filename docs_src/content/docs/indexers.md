+++
title = "Indexers"
description = "Find download source for TV show episodes or movies"
date = 2018-05-24T14:52:09Z
weight = 50
draft = false
bref = "Indexer modules are responsible for searching torrents to download for a specific TV show episode or movie."
toc = true
+++

## Indexers overview
---

flemzerd uses Indexers to create a list of torrents corresponding to a media. When creating this list, flemzerd also applies filters and sorting to have a more accurate list of torrents.
To improve torrent list, the quality can be filtered according to preference set in the configuration file. (see `preferred_media_quality` parameter in the [configuration](/docs/configuration))

### Different providers types
---

Just like Providers, 2 types of indexers can be defined.
* TV Indexers: retrieve torrents for tvshows
* Movie Indexers: retrieve torrents for movies

An Indexer can be a TV Indexer, Movie Indexer or both at the same time.

## Available Indexers
---

The Indexers used by the flemzerd daemon are defined in the configuration file. In this configuration file, multiple you can define multiple Indexers with the following constraints:
* If TV Shows tracking is enabled, you must define at least one TV Indexer.
* If movie tracking is enabled, you must define at least one Movie Indexer.

### Torznab
---
**Type**: TVIndexer, MovieIndexer

Currently only one type of Indexer exists: Torznab.

Torznab is a Newznab-like API designed to query torrents. It is used by Jackett and Cardigann as an interface to torrents sites.

#### How to use
---
* Enable a `torznab` Indexer in configuration file with the URL to the Torznab indexer (usually Jackett or Cardigann), a name and the API key for the Indexer.
{{< highlight toml >}}
[indexers]
    [[indexers.torznab]]
        name = "Indexer name"
        url = "http://first-indexer:8080/torznab"
        apikey = "API_KEY"
{{< /highlight >}}
