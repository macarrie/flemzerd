+++
title = "How it works"
description = "Different modules types existing in the flemzerd daemon and their role"
date = 2018-05-24T13:41:26Z
weight = 20
draft = false
bref = "flemzerd functions by orcherstrating different modules. This page describes these modules and how the interact together by giving an example of a media being processed"
toc = true
+++

## Walkthrough of a media download
---

The following scenario describes how a movie download process is done, from the moment a flemzerd user decides to watch a movie, to the moment the movie is playing on the screen.

* A user is at works and decides he wants to watch a specific movie tonight. He adds it on one of its **[Watchlists](/docs/watchlists)** (for example trakt)
* flemzerd regularly checks watchlist content. He detects a new movie has been added and begins to process it.
* Details for the movie are fetched using **[Providers](/docs/providers)** (for example TheMovieDB). It includes full title, synopsis, release date, duration, ...
* Torrents are then searched using **[Indexers](/docs/indexers)**.
* Once a torrent list for the movie has been retrieved, sorted by seeders and filtered to remove errors (torrents for a different movie, or torrents for the same movie of a different year), this list is sent to **[Downloaders](/docs/downloaders)** which will begin the actual download
* When the download is finished, the movie is moved to the library path defined in the configuration and **[Mediacenters](/docs/mediacenters)** library are refreshed (kodi).
* During the whole process, notifications are sent using **[Notifiers](/docs/notifiers)** to alert user of the differents events occuring (new movie detected, movie download succeeded, movie downloaded failed, ...)
* The download is now complete. The only thing left to do when getting home is sitting back and enjoying the movie ^^

The download of TV shows episodes is almost identical. The difference is the following:

* The users adds a TV show he regularly watches
* The TV Show is not downloaded directly (this can mean downloading a lot of episodes, so this is not done automatically). Instead, flemzerd watches regularly for new episodes.
* Once an episode is aired, flemzerd detects a new episode is out and begins the whole download process

This is convenient when following a lot of different running TV shows at once (a lot means more than one for me), since looking for new episodes and finding torrents everyday can be tiresome.


## More detailed description of modules types
---

* [Watchlists](/docs/watchlists) : **(required)** Watchlists define what flemzerd is tracking. These list can be defined locally as static lists or dynamic lists defined on external services
* [Providers](/docs/providers) : **(required)** Providers are responsible for loading media details (title, release dates, synopsis, informations about actors...). For TV shows, Providers are very important because they retrieve informations about new episodes. They use external services such as TheTVDB or TheMovieDB.
* [Indexers](/docs/indexers) : **(required)** Indexers search for torrents for a specific media. Indexers are used to prepare a list of what to download for Downloaders. Indexers are required if you want to perform any media download.
* [Downloaders](/docs/downloaders) : **(required)** Downloaders perform the actual download of torrents. Same as Indexers, Downloaders are required only to download media. If no downloader/indexer is operational, flemzerd can still be used just as a way to send notifications for new episodes for tracked TV shows.
* [Mediacenters](/docs/mediacenters) : Mediacenters module establish a link between flemzerd and mediacenters. They are used to refresh your mediacenter library when a download is finished. The mediacenter will then load informations about the new downloaded media and be presented with backgrounds, fanarts, related informations, etc... Defining mediacenters modules is optional.
* [Notifiers](/docs/notifiers) : Notifiers sends notifications (duh) for events occuring in flemzerd: New episode aired, New movie detected in watchlists, Download started, Download successful, Download failed.
  Notifiers are also optional. If no notifier is defined, no notifications will be sent.


## Error handling
---

Since lots of external services and tools are used by flemzerd, errors can occur. 

* **Watchlists**: If no watchlists are defined, nothing will happen since no media have been added into flemzerd. If Watchlists are temporarly unavailable, new additions to wxatchlists will not be detected but new episodes will still be detected since corresponding show is already tracked by flemzerd.
* **Providers**: An error when trying to get information about movies will completely stop the download chain. No informations about new media from watchlists can be retrieved, so torrents cannot be retried and no media can be downloaded
* **Indexers**: Multiple indexers can be defined to multiply chances to find good torrents. If no torrents are found using all indexers, torrent search will be retried at the next check: episodes or movies that just aired may not have torrents yet
* **Downloaders**: A lot of issues have to be dealt with when downloading torrents. If the downloader is not available, easy, download is stopped. If a torrent download fails (no space left on device, invalid torrent), the next torrent in the list is scheduled to download. In the case when the limit of failed torrents is reached (defined in the configuration), the download is marked as failed.
* **Notifiers**: Nothing dramatic. When a notifier cannot be joined, notifications will not be able to be sent for this notifier. Notifications will still be sent for the other notifiers defined in the configuration.
* **Mediacenters**: Nothing of top importance either. An unreachable media center won't have its library refreshed at the end of a successful download.


For more details on errors, see the [Troubleshooting](/docs/troubleshooting) page.
