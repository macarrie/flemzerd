import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';

import { Episode } from '../episode';

import { EpisodeService } from '../episode.service';
import { FanartService } from '../fanart.service';
import { UtilsService } from '../utils.service';

@Component({
    selector: 'app-episode-details',
    templateUrl: './episode-details.component.html',
    styleUrls: ['../media.details.scss']
})
export class EpisodeDetailsComponent implements OnInit {
    episode :Episode;
    background_image :string;

    constructor(
        private episodeService :EpisodeService,
        private fanartService :FanartService,
        private utils :UtilsService,
        private route :ActivatedRoute,
        private location :Location
    ) {}

    ngOnInit() {
        this.getEpisode();
    }

    getEpisode() :void {
        const id = +this.route.snapshot.paramMap.get('id');
        this.episodeService.getEpisode(id).subscribe(episode => {
            this.episode = episode;

            this.fanartService.getTvShowFanart(episode.TvShow).subscribe(fanartObj => {
                if (fanartObj["showbackground"] && fanartObj["showbackground"].length > 0) {
                    this.background_image = fanartObj["showbackground"][0]["url"];
                }
            });
        });
    }

    deleteEpisode(episode :Episode) :void {
        this.episodeService.deleteEpisode(episode).subscribe(response => {
            this.location.back();
        });
    }

    changeDownloadedState(e :Episode, downloaded :boolean) {
        this.episodeService.changeEpisodeDownloadedState(e, downloaded).subscribe(episode => {
            this.getEpisode();
        });
    }

    downloadEpisode(episode :Episode) :void {
        if (episode.DownloadingItem.Pending || episode.DownloadingItem.Downloading || episode.DownloadingItem.Downloaded) {
            return;
        }

        episode.DownloadingItem.Pending = true;
        this.episodeService.downloadEpisode(episode).subscribe(
            data =>  {
                this.getEpisode();
            },
            error => {
                episode.DownloadingItem.Pending = false;
            }
        );
    }

    back() :void {
        this.location.back();
    }
}
