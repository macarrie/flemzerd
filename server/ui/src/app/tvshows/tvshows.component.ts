import { Component, OnInit } from '@angular/core';

import { TvShow } from '../tvshow';
import { Episode } from '../episode';

import { TvshowsService } from '../tvshows.service';
import { EpisodeService } from '../episode.service';
import { UtilsService } from '../utils.service';
import { ConfigService } from '../config.service';

@Component({
    selector: 'app-tvshows',
    templateUrl: './tvshows.component.html',
    styleUrls: ['./tvshows.component.scss']
})
export class TvshowsComponent implements OnInit {
    trackedShows :TvShow[];
    removedShows :TvShow[];
    downloadingEpisodes :Episode[];
    refreshing :boolean;
    config :any;

    constructor(
        private configService :ConfigService,
        private tvshowsService :TvshowsService,
        private episodeService :EpisodeService,
        private utils :UtilsService
    ) {}

    ngOnInit() {
        this.getConfig();
        this.getTrackedTvShows();
        this.getRemovedTvShows();
        this.getDownloadingEpisodes();
    }

    getConfig() {
        this.configService.getConfig().subscribe(config => {
            this.config = config;
        })
    }

    refreshTvShows() :void {
        this.refreshing = true;
        this.tvshowsService.refreshTvShows().subscribe(data => {
            this.getTrackedTvShows();
            this.getRemovedTvShows();
            this.refreshing = false;
        });
    }

    getTrackedTvShows() :void {
        this.tvshowsService.getTrackedTvShows().subscribe(tvshows => this.trackedShows = tvshows);
    }

    getRemovedTvShows() :void {
        this.tvshowsService.getRemovedTvShows().subscribe(tvshows => this.removedShows = tvshows);
    }

    getDownloadingEpisodes() :void {
        this.episodeService.getDownloadingEpisodes().subscribe(episodes => this.downloadingEpisodes = episodes);
    }

    stopDownload(episode :Episode) :void {
        episode.DownloadingItem.AbortPending = true;
        this.episodeService.stopDownload(episode.ID).subscribe(response => {
            this.episodeService.getDownloadingEpisodes().subscribe(episodes => {
                this.downloadingEpisodes = episodes;
                episode.DownloadingItem.AbortPending = true;
            });
        });
    }
}
