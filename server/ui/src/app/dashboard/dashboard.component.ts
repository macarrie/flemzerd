import { Component, OnInit } from '@angular/core';

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

    constructor(
        private tvshowsService :TvshowsService,
        private episodeService :EpisodeService,
        private movieService :MovieService
    ) {}

    ngOnInit() {
        this.getTrackedTvShows();
        this.getDownloadingEpisodes();
        this.getTrackedMovies();
        this.getDownloadedMovies();
        this.getDownloadingMovies();
    }

    getTrackedTvShows() :void {
        this.tvshowsService.getTrackedTvShows().subscribe(tvshows => this.trackedShows = tvshows);
    }

    getDownloadingEpisodes() :void {
        this.episodeService.getDownloadingEpisodes().subscribe(episodes => this.downloadingEpisodes = episodes);
    }

    getTrackedMovies() :void {
        this.movieService.getTrackedMovies().subscribe(movies => this.trackedMovies = movies);
    }

    getDownloadedMovies() :void {
        this.movieService.getDownloadedMovies().subscribe(movies => this.downloadedMovies = movies);
    }

    getDownloadingMovies() :void {
        this.movieService.getDownloadingMovies().subscribe(movies => this.downloadingMovies = movies);
    }
}
