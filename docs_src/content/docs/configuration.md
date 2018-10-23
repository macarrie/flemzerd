+++
title = "Configuration"
date = 2018-05-24T13:35:12Z
weight = 30
draft = false
description = "Available configuration options to make flemzerd fit your use case"
bref = "This page summarizes the different configuration options avaible to make the flemzerd daemon fit your use case"
toc = true
+++


## Configuration file loading
---

The flemzerd daemon loads its configuration file when starting. The configuration file is written in [TOML](https://github.com/toml-lang/toml) and loaded from the following places:

* `~/.config/flemzerd/flemzerd.toml`
* `/etc/flemzerd/flemzerd.toml`

The first file found is used as the configuration file. When starting, the flemzerd daemon displays in the logs the configuration file used.

### Using another configuration file
---

Configuration files can be loaded from other places than the default directories listed above. To specify another configuration file, use the `-d` flag when starting flemzerd:
{{< highlight bash >}}
./flemzerd -c /path/to/config.toml
{{< /highlight >}}

## Configuration file organization
---

The configuration file is organized in different sections, each section corresponding to a specific subject.

### System settings
---
{{< highlight toml >}}
[system]
    check_interval = 15
    torrent_download_attempts_limit = 5
    track_shows = true
    track_movies = true
    preferred_media_quality = "720p"
    excluded_release_types = "cam,screener,telesync,telecine"
{{< /highlight >}}

##### Options explanation
---
* `check_interval` (default: `15`)<br />
flemzerd checks regularly for new movies  and shows from watchlists and new episodes from tracked TV shows. This option sets the check interval (in minutes).
* `torrent_download_attempts_limit` (default: `5`)<br />
Failure to download torrent can happen. flemzerd waits by default for 5 torrent download failures before marking the download as failed. This options sets the number of download failures needed to consider the media download as failed
* `track_shows` (default: `true`)<br />
Track and download TV shows episodes found in Watchlists. Set to `false` to disable TV show handling.
* `track_movies` (default: `true`)<br />
Track and download movies found in Watchlists. Set to `false` to disable movies handling.
* `preferred_media_quality` (default: `720p`)<br />
  When getting torrents for media, flemzerd can sort them by quality, putting the one you chose in this option first. Possible values are: 480p, 576p, 720p, 900p, 1080p, 1440p, 2160p, 5k, 8k, 16k
* `excluded_release_types` (default: `cam, screener, telesync, telecine`)<br />
  When the following release types are detected in a torrent name, it will be excluded from download list. Possible values are: cam, screener, telesync, telecine, dvdrip, hdtv, webdl, blurayrip

### Interface settings
---
{{< highlight toml >}}
[interface]
    enabled = true
    port = 8400
{{< /highlight >}}

##### Options explanation
---
* `enabled` (default: `true`)<br />
   Enable web interface
* `port` (default: `8400`)<br />
  Access port to web interface

### Providers declaration
---
{{< highlight toml >}}
[providers]
    [providers.provider1]
        key1 = "value"
        key2 = 1
{{< /highlight >}}

For flemzerd to work, all used modules, including providers must be declared. This section is a list of the different Providers to use.
Each Provider defines a subsection of `[providers]`, and each subsection of `[providers]` must be an exisiting Provider with its configuration keys.

See more details about available Providers in the associated documentation page: [Providers](/docs/providers)

### Notifiers declaration
---
{{< highlight toml >}}
[notifiers]
    notifier1 = []
    notifier2 = []

    [notifiers.notifier3]
        key1 = "value"
        key2 = 1
{{< /highlight >}}

Notifiers used by the flemzerd daemon are defined in the `[notifiers]` section. Using notifiers is optional, so this section can be left empty of even not declared at all in the configuration file.
Each Notifier defines a subsection of `[notifiers]`, and each subsection of `[notifiers]` must be an exisiting Notifier with its configuration keys.

See more details about available Notifiers in the associated documentation page: [Notifiers](/docs/notifiers)

### Indexers declaration
---
{{< highlight toml >}}
[indexers]
    [[indexers.indexer1]]
        key1 = "value"
        key2 = 1

    [[indexers.indexer2]]
        key1 = "value"
        key2 = 1
{{< /highlight >}}

Indexers used in the flemzerd daemon are defined in the `[indexers]` section.
Each Indexer defines a subsection of `[indexers]`, and each subsection of `[indexers]` must be an exisiting Indexer with its configuration keys.

See more details about available Indexers in the associated documentation page: [Indexers](/docs/indexers)

### Indexers declaration
---
{{< highlight toml >}}
[downloaders]
    [downloaders.downloader1]
        key1 = "value"
        key2 = 1
{{< /highlight >}}

Downloaders used in the flemzerd daemon are defined in the `[downloaders]` section.
Each Downloader defines a subsection of `[downloaders]`, and each subsection of `[downloaders]` must be an exisiting Downloader with its configuration keys.
While multiple Downloaders can be defined and loaded into flemzerd, only the first Downloader available will be used for downloading torrents.

See more details about available Downloaders in the associated documentation page: [Downloaders](/docs/downloaders)

### Watchlists declaration
---
{{< highlight toml >}}
[watchlists]
    watchlist1 = []
    watchlist2 = [
        "TV_SHOW_1",
        "TV_SHOW_2",
        "TV_SHOW_3"
    ]
{{< /highlight >}}

Watchlists used in the flemzerd daemon are defined in the `[watchlists]` section.
Each watchlist defines a subsection of `[watchlists]`, and each subsection of `[watchlists]` must be an exisiting Watchlist with its configuration keys.

See more details about available Watchlists in the associated documentation page: [Watchlists](/docs/watchlists)

### Mediacenters declaration
---
{{< highlight toml >}}
[mediacenters]
    [mediacenters.mediacenter1] 
        key1 = "value"
        key2 = 1
{{< /highlight >}}

Mediacenters used in the flemzerd daemon are defined in the `[mediacenters]` section.
Each mediacenter defines a subsection of `[mediacenters]`, and each subsection of `[mediacenters]` must be an exisiting mediacenter with its configuration keys.

See more details about available Mediacenters in the associated documentation page: [Mediacenters](/docs/mediacenters)


### Notifications options
---
{{< highlight toml >}}
[notifications]
    enabled = true
    notify_new_episode = true
    notify_new_movie = true
    notify_download_start = true
    notify_download_complete = true
    notify_failure = true
{{< /highlight >}}

##### Options explanation
---
* `enabled` (default: `true`)<br />
  Global state of notifications. If set to `false`, no notifications will be sent, even if the others parameters in the `[notifications]` section are set to `true`
* `notify_new_episode` (default: `true`)<br />
  If set to `true`, send notification when a new episode from tracked shows has aired recently.
* `notify_new_movie` (default: `true`)<br />
  If set to `true`, send notification when a new movie has been found in Watchlists
* `notify_download_start` (default: `true`)<br />
  If set to `true`, send notification when a media download starts.
* `notify_download_complete` (default: `true`)<br />
  If set to `true`, send notification when a media download has ended successfully
* `notify_failure` (default: `true`)<br />
  If set to `true`, send notification when a media download failed


### Library options
---
{{< highlight toml >}}
[library]
    show_path = "/var/lib/flemzer/library/shows"
    movie_path = "/var/lib/flemzer/library/movies"
    custom_tmp_dir = "/var/lib/flemzerd/tmp"
{{< /highlight >}}

When downloading media, you can choose where flemzerd will put downloaded content. This can be useful if you have a predefined place where you put your media (mediacenter library for example).

The flemzerd daemon also uses a temporary folder when downloading items. This prevents your library to be filled with incomplete downloads or useless files created by an unsuccessful torrent download.

##### Options explanation
---
* `show_path` (default `var/lib/flemzer/library/shows`)<br />
  Path where downloaded TV show episodes will be placed
* `movie_path` (default `var/lib/flemzer/library/movies`)<br />
  Path where downloaded movies will be placed
* `custom_tmp_dir` (default `var/lib/flemzerd/tmp`)<br />
  Path where in progress downloads will be placed
