import { Component, OnInit, OnDestroy } from '@angular/core';

import { TvShow } from '../tvshow';
import { Movie } from '../movie';
import { Episode } from '../episode';

import { TvshowsService } from '../tvshows.service';
import { EpisodeService } from '../episode.service';
import { MovieService } from '../movie.service';

@Component({
    selector: 'app-dashboard',
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
    trackedShows :TvShow[];
    downloadingEpisodes :Episode[];
    trackedMovies :Movie[];
    downloadedMovies :Movie[];
    downloadingMovies :Movie[];

    dashboardRefresh :any;

    constructor(
        private tvshowsService :TvshowsService,
        private episodeService :EpisodeService,
        private movieService :MovieService
    ) {
        movieService.trackedMovies.subscribe(movies => {
            this.trackedMovies = movies;
        });
        movieService.downloadedMovies.subscribe(movies => {
            this.downloadedMovies = movies;
        });
        movieService.downloadingMovies.subscribe(movies => {
            this.downloadingMovies = movies;
        });
        tvshowsService.trackedShows.subscribe(shows => {
            this.trackedShows = shows;
        });
        episodeService.downloadingEpisodes.subscribe(episodes => {
            this.downloadingEpisodes = episodes;
        });
    }

    ngOnInit() {
        this.movieService.getMovies();
        this.tvshowsService.getTrackedTvShows();
        this.episodeService.getDownloadingEpisodes();

        this.dashboardRefresh = setInterval(() => {
            this.movieService.getMovies();
            this.tvshowsService.getTrackedTvShows();
            this.episodeService.getDownloadingEpisodes();
        }, 30000);
    }

    ngOnDestroy() {
        clearInterval(this.dashboardRefresh);
    }
}
