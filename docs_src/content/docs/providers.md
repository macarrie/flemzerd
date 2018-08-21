+++
title = "Providers"
description = "Retrieve detailed informations about movies, TV shows and episodes"
date = 2018-05-24T14:52:00Z
weight = 40
draft = false
bref = "Providers are in flemzerd responsible for gettings detailed informations for tracked movies, TV shows and episodes."
toc = true
+++

## Providers overview
---

### Providers role
---

flemzerd uses Providers to get detailed informations about media. The most important information from Providers is the release date. It allows flemzerd to launch detect new episodes and movies.

flemzerd also uses Providers modules to get a lot of details for the UI: plot details, pictures, episodes and season list...

### Different providers types
---

Providers are external services such as [TheMovieDB](https://www.themoviedb.org/) or [TheTVDB](https://www.thetvdb.com/) that collect all these precious informations. Since some sites are sometimes specialized, flemzerd defines two types of Providers:
* TV Providers: retrieve informations for TV shows and episodes
* Movie Providers: retrieve informations for movies.

Providers can be TV Providers and Movie Providers at the same time. For example, The TVDB Provider is only a TV Provider because [TheTVDB](https://www.thetvdb.com/) only provides information for TV shows.
The TMDB Provider is a TV and Movie Provider because [TheMovieDB](https://www.themoviedb.org/) contains informations about both media types.
**
## Available Providers
---

The Providers used by the flemzerd daemon are defined in the configuration files. In this configuration file, multiple you can define multiple Providers with the following constraints:
* If TV Shows tracking is enabled, you must define at least one TV Provider.
* If movie tracking is enabled, you must define at least one Movie Provider.

You can also define multiple Providers of the same type. flemzerd uses the additional Providers as backups.

### TheMovieDB
---
 **Type**: TVProvider, MovieProvider

Uses [TheMovieDB](https://www.themoviedb.org/) to get movies and TV shows informations.

#### How to use
---
* Enable `tmdb` Provider in configuration file
{{< highlight toml >}}
[providers]
    tmdb = []
{{< /highlight >}}
* Define TMDB API key
    TMDB needs an key to perform API requests. This key is passed to flemzerd by defining the `FLZ_TMDB_API_KEY` environment variable. When defined during flemzerd compilation, this variable is compiled into the binary.
    In [packages](https://github.com/macarrie/flemzerd/releases) found on GitHub, this key is precompiled into the binary.

### TheTVDB
---
 **Type**: TVProvider

Uses [TheTVDB](https://www.thetvdb.com/) to get TV shows informations.

#### How to use
---
* Enable `tvdb` Provider in configuration file
{{< highlight toml >}}
[providers]
    tvdb = []
{{< /highlight >}}
* Define TMDB API key
    TVDB needs an key to perform API requests. This key is passed to flemzerd by defining the `FLZ_TVDB_API_KEY` environment variable. When defined during flemzerd compilation, this variable is compiled into the binary.
    In [packages](https://github.com/macarrie/flemzerd/releases) found on GitHub, this key is precompiled into the binary.
