import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';
import { Observable } from 'rxjs';
import { forkJoin } from "rxjs";

import { TvShow } from '../tvshow';
import { Episode } from '../episode';

import { TvshowsService } from '../tvshows.service';
import { EpisodeService } from '../episode.service';
import { FanartService } from '../fanart.service';
import { UtilsService } from '../utils.service';
import { ConfigService } from '../config.service';

@Component({
    selector: 'app-tv-show-details',
    templateUrl: './tv-show-details.component.html',
    styleUrls: ['../media.details.scss']
})
export class TvShowDetailsComponent implements OnInit {
    tvshow :TvShow;
    seasons :Episode[][];
    background_image :string;
    showStatusClasses = {};
    title :string;

    config :any;

    constructor(
        private tvshowsService :TvshowsService,
        private episodeService :EpisodeService,
        private fanartService :FanartService,
        private configService :ConfigService,
        private utils :UtilsService,
        private route :ActivatedRoute,
        private location :Location
    ) {
        configService.config.subscribe(cfg => {
            this.config = cfg;
        });
    }

    ngOnInit() {
        this.configService.getConfig();
        this.getShow();
    }

    getTitle(show :TvShow) :string {
        if (show.CustomName != "") {
            return show.CustomName;
        }

        if (show.UseDefaultTitle) {
            return show.Name;
        }

        return show.OriginalName;
    }

    getShow() :void {
        const id = +this.route.snapshot.paramMap.get('id');
        this.tvshowsService.getShow(id).subscribe(tvshow => {
            this.tvshow = tvshow ;
            this.showStatusClasses = {
                'btn-outline-success': this.tvshow.Status < 3,
                'btn-outline-danger': this.tvshow.Status == 3,
                'btn-outline-warning': this.tvshow.Status > 3,
            };
            this.title = this.getTitle(tvshow);
            if (this.title == "") {
                if (tvshow.UseDefaultTitle) {
                    this.title = tvshow.Name;
                } else {
                    this.title = tvshow.OriginalName;
                }
            }

            this.fanartService.getTvShowFanart(tvshow).subscribe(fanartObj => {
                if (fanartObj["showbackground"] && fanartObj["showbackground"].length > 0) {
                    this.background_image = fanartObj["showbackground"][0]["url"];
                }
            });

            if (tvshow.DeletedAt == null) {
                this.seasons = new Array<Episode[]>(tvshow.Seasons.length);
                for (var season in tvshow.Seasons) {
                    let n = tvshow.Seasons[season].SeasonNumber;
                    this.episodeService.getBySeason(tvshow.ID, n).subscribe(episodes => {
                        this.seasons[n] = episodes;
                    });
                }
            }
        });
    }

    useDefaultTitle(show :TvShow, defaultTitle :boolean) {
        this.tvshowsService.useDefaultTitle(show, defaultTitle).subscribe(show => {
            this.getShow();
        });
    }

    deleteShow() :void {
        this.tvshowsService.deleteShow(this.tvshow.ID).subscribe(response => this.getShow());
    }

    restoreShow() :void {
        this.tvshowsService.restoreShow(this.tvshow.ID).subscribe(response => this.getShow());
    }

    changeEpisodeDownloadedState(e :Episode, downloaded :boolean) {
        this.episodeService.changeEpisodeDownloadedState(e, downloaded).subscribe(episode => {
            this.seasons[e.Season][e.Number - 1] = episode;
        });
    }

    changeSeasonDownloadedState(season, downloaded :boolean) {
        let collapsed_state = season.collapse;
        let requests = new Array<Observable<Episode>>();
        for (let episode of season) {
            if (!episode.DownloadingItem.Pending && episode.DownloadingItem.Downloaded != downloaded && !episode.DownloadingItem.Downloading) {
                requests.push(this.episodeService.changeEpisodeDownloadedState(episode, downloaded));
            }
        }

        forkJoin(requests).subscribe(response => {
            this.episodeService.getBySeason(this.tvshow.ID, season[0].Season).subscribe(episodes => {
                this.seasons[season[0].Season] = episodes;
                this.seasons[season[0].Season]["collapse"] = collapsed_state;
            });
        });
    }

    back() :void {
        this.location.back();
    }

    getStatusText(show :TvShow) :string {
        switch (show["Status"]) {
            case 1:
                return "Continuing";
            case 2:
                return "Planned";
            case 3:
                return "Ended";
            default:
                return "Unknown";
        }
    }

    downloadEpisode(episode :Episode) :void {
        if (episode.DownloadingItem.Pending || episode.DownloadingItem.Downloading || episode.DownloadingItem.Downloaded) {
            return;
        }

        episode.DownloadingItem.Pending = true;
        this.episodeService.downloadEpisode(episode).subscribe(
            data =>  {
                // Refresh season
                this.episodeService.getBySeason(this.tvshow.ID, episode.Season).subscribe(episodes => {
                    this.seasons[episode.Season] = episodes;
                });
            },
            error => {
                episode.DownloadingItem.Pending = false;
            }
        );
    }

    downloadSeason(season) {
        let collapsed_state = season.collapse;
        let requests = new Array<Observable<Episode>>();
        for (let episode of season) {
            if (!episode.DownloadingItem.Pending && !episode.DownloadingItem.Downloading && !episode.DownloadingItem.Downloaded) {
                episode.DownloadingItem.Pending = true;
                requests.push(this.episodeService.downloadEpisode(episode));
            }
        }

        forkJoin(requests).subscribe(response => {
            this.episodeService.getBySeason(this.tvshow.ID, season[0].Season).subscribe(episodes => {
                this.seasons[season[0].Season] = episodes;
                this.seasons[season[0].Season]["collapse"] = collapsed_state;
            });
        });
    }

    setCustomTitle(title :string) {
        this.tvshow.CustomName = title;
        this.tvshowsService.updateShow(this.tvshow).subscribe(tvshow => {
            this.getShow();
        });
    }
}

