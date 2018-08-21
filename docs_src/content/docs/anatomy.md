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

The following scenario describes the movie download process. is This process spans from the moment a flemzerd user decides to watch a movie, to the moment the movie is playing on the screen.

* A user is at works and decides he wants to watch a specific movie tonight. He adds it on one of its **[Watchlists](/docs/watchlists)** (for example trakt)
* flemzerd checks at regular intervals watchlists content. He detects a new movie in one of its Watchlists and begins to process it.
* **[Providers](/docs/providers)** (for example TheMovieDB) fetches details for the movie. It includes full title, synopsis, release date, duration, ...
* Torrents are then searched using **[Indexers](/docs/indexers)**. Filters are applied to the torrent list to remove bad torrents (bad movie, wrong year
* flemzerd sends this torrent list to **[Downloaders](/docs/downloaders)** which will begin the actual download
* When the download is finished, flemzerd moves the movie to the library (path defined in the configuration). **[Mediacenters](/docs/mediacenters)** libraries are then refreshed.
* During the whole process, notifications are sent using **[Notifiers](/docs/notifiers)** to alert user of the differents events occuring (new movie detected, movie download succeeded, movie downloaded failed, ...)
* The download is now complete. The only thing left to do when getting home is sitting back and enjoying the movie ^^

The download of TV shows episodes is almost identical. The difference is the following:

* The users adds a TV show he watches.
* flemzerd watches at regular intervals for new episodes.
* Once a new episode aired, flemzerd begins the whole download process.

This is convenient when following a lot of different running TV shows at once. In fact, looking for new episodes and finding torrents everyday can be tiresome.


## More detailed description of modules types
---

* [Watchlists](/docs/watchlists) : **(required)** Watchlists define what flemzerd is tracking. These list can be defined in the configuration as static lists or dynamic lists defined with external services
* [Providers](/docs/providers) : **(required)** Providers are responsible for loading media details (title, release dates, synopsis, informations about actors...). For TV shows, Providers are very important because they retrieve informations about new episodes. They use external services such as TheTVDB or TheMovieDB.
* [Indexers](/docs/indexers) : **(required)** Indexers search for torrents for a specific media. Indexers prepare a list of what to download for Downloaders. Indexers are mandatory if you want to perform any media download.
* [Downloaders](/docs/downloaders) : **(required)** Downloaders perform the actual download of torrents. Same as Indexers, Downloaders are required only to download media. If no Downloader/Indexer is operational, flemzerd can still send notifications for new episodes for tracked TV shows.
* [Mediacenters](/docs/mediacenters) : Mediacenters module establish a link between flemzerd and mediacenters. They are used to refresh your mediacenter library when a download is finished. The mediacenter will then load informations about the new downloaded media and be presented with backgrounds, fanarts, related informations, etc... Defining mediacenters modules is optional.
* [Notifiers](/docs/notifiers) : Notifiers sends notifications for events occuring in flemzerd: 
    * New episode aired
    * New movie detected in watchlists
    * Download started
    * Download successful
    * Download failed

  Notifiers are also optional. If no notifier is defined, no notifications will be sent.


## Error handling
---

Since lots of external services and tools are used by flemzerd, errors can occur. 

* **Watchlists**: If no watchlists are defined, nothing will happen since no media have been added into flemzerd. If Watchlists are temporarly unavailable, new additions to wxtchlists will not be detected. New episodes will still be detected since corresponding show is already tracked by flemzerd.
* **Providers**: An error when trying to get information about movies will completely stop the download chain. No informationsabout new media from watchlists can be retrieved, so torrents cannot be retried and no media can be downloaded
* **Indexers**: You can define many Indexers to multiply chances to find good torrents. If flemzerd does not find any torrent using all indexers, torrent search will be retried at the next check.
* **Downloaders**: A lot of issues occur when downloading torrents. If the downloader is not available, easy, the download process stops. If a torrent download fails (no space left on device, invalid torrent), the next torrent in the list is downloaded. When the limit of failed torrents is reached (defined in the configuration), the download is marked as failed.
* **Notifiers**: Nothing dramatic. When flemzerd cannot contact a notifier, notifications will not be able to be sent for this notifier. Notifications will still be sent for the other notifiers defined in the configuration.
* **Mediacenters**: Nothing of top importance either. An unreachable media center won't have its library refreshed at the end of a successful download.


For more details on errors, see the [Troubleshooting](/docs/troubleshooting) page.
