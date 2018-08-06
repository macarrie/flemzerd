import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";

import { TvShow } from '../tvshow';
import { Episode } from '../episode';

import { TvshowsService } from '../tvshows.service';
import { EpisodeService } from '../episode.service';
import { UtilsService } from '../utils.service';
import { ConfigService } from '../config.service';

@Component({
    selector: 'app-tvshows',
    templateUrl: './tvshows.component.html',
    styleUrls: ['../media.miniatures.scss']
})
export class TvshowsComponent implements OnInit {
    trackedShows :TvShow[];
    removedShows :TvShow[];
    downloadingEpisodes :Episode[];

    refreshing :boolean;
    refreshing_poll :boolean;

    config :any;
    showsRefresh :any;

    constructor(
        private http :HttpClient,
        private configService :ConfigService,
        private tvshowsService :TvshowsService,
        private episodeService :EpisodeService,
        private utils :UtilsService
    ) {
        tvshowsService.trackedShows.subscribe(shows => {
            this.trackedShows = shows;
        });
        tvshowsService.removedShows.subscribe(shows => {
            this.removedShows = shows;
        });
        episodeService.downloadingEpisodes.subscribe(episodes => {
            this.downloadingEpisodes = episodes;
        });
        tvshowsService.refreshing.subscribe(bool => {
            this.refreshing = bool;
        });
        configService.config.subscribe(cfg => {
            this.config = cfg;
        });
    }

    ngOnInit() {
        this.tvshowsService.getShows();
        this.episodeService.getDownloadingEpisodes();
        this.configService.getConfig();

        this.showsRefresh = setInterval(() => {
            this.configService.getConfig();
            this.tvshowsService.getShows();
            this.episodeService.getDownloadingEpisodes();
        }, 30000);
    }

    ngOnDestroy() {
        clearInterval(this.showsRefresh);
    }

    refreshTvShows() :void {
        this.tvshowsService.getShows();
    }

    executePollLoop() :void {
        this.refreshing_poll = true;
        this.http.post('/api/v1/actions/poll', {}).subscribe(data => {
            this.tvshowsService.getShows();
            this.refreshing_poll = false;
        });
    }

    stopDownload(episode :Episode) :void {
        episode.DownloadingItem.AbortPending = true;
        this.episodeService.stopDownload(episode.ID).subscribe(response => {
            this.episodeService.getDownloadingEpisodes();
            episode.DownloadingItem.AbortPending = false;
        });
    }
}
