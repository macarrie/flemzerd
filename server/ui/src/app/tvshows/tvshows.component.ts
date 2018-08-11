import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { HttpClient } from "@angular/common/http";
import { takeUntil  } from 'rxjs/operators'; // for rxjs ^5.5.0 lettable operators

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
export class TvshowsComponent implements OnInit, OnDestroy {
    private subs :Subject<any> = new Subject();

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
        tvshowsService.trackedShows.pipe(takeUntil(this.subs)).subscribe(shows => {
            this.trackedShows = shows;
        });
        tvshowsService.removedShows.pipe(takeUntil(this.subs)).subscribe(shows => {
            this.removedShows = shows;
        });
        episodeService.downloadingEpisodes.pipe(takeUntil(this.subs)).subscribe(episodes => {
            this.downloadingEpisodes = episodes;
        });
        tvshowsService.refreshing.pipe(takeUntil(this.subs)).subscribe(bool => {
            this.refreshing = bool;
        });
        configService.config.pipe(takeUntil(this.subs)).subscribe(cfg => {
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
        this.subs.next();
        this.subs.complete();
    }

    getTitle(show :TvShow) :string {
        if (show.UseDefaultTitle) {
            return show.Name;
        }

        return show.OriginalName;
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
