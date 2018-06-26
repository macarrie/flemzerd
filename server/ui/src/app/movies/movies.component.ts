import { Component, OnInit } from '@angular/core';

import { Movie } from '../movie';

import { MovieService } from '../movie.service';
import { UtilsService } from '../utils.service';
import { ConfigService } from '../config.service';

@Component({
    selector: 'app-movies',
    templateUrl: './movies.component.html',
    styleUrls: ['./movies.component.scss']
})
export class MoviesComponent implements OnInit {
    trackedMovies :Movie[];
    removedMovies :Movie[];
    downloadedMovies :Movie[];
    downloadingMovies :Movie[];
    refreshing :boolean;
    show_downloaded_movies = true;
    config :any;

    constructor(
        private configService :ConfigService,
        private movieService :MovieService,
        private utils :UtilsService
    ) { }

    ngOnInit() {
        this.getConfig();
        this.getTrackedMovies();
        this.getRemovedMovies();
        this.getDownloadedMovies();
        this.getDownloadingMovies();
    }

    getConfig() {
        this.configService.getConfig().subscribe(config => {
            this.config = config;
        })
    }

    refreshMovies() :void {
        this.refreshing = true;
        this.movieService.refreshMovies().subscribe(data => {
            this.getTrackedMovies();
            this.getRemovedMovies();
            this.getDownloadedMovies();
            this.getDownloadingMovies();
            this.refreshing = false;
        });
    }

    getTrackedMovies() :void {
        this.movieService.getTrackedMovies().subscribe(movies => this.trackedMovies = movies);
    }

    getRemovedMovies() :void {
        this.movieService.getRemovedMovies().subscribe(movies => this.removedMovies = movies);
    }

    getDownloadedMovies() :void {
        this.movieService.getDownloadedMovies().subscribe(movies => this.downloadedMovies = movies);
    }

    getDownloadingMovies() :void {
        this.movieService.getDownloadingMovies().subscribe(movies => this.downloadingMovies = movies);
    }

    stopDownload(movie :Movie) :void {
        movie.DownloadingItem.AbortPending = true;
        this.movieService.stopDownload(movie.ID).subscribe(response => {
            this.movieService.getDownloadingMovies().subscribe(movies => {
                this.downloadingMovies = movies;
                movie.DownloadingItem.AbortPending = true;
            });
        });
    }
}
