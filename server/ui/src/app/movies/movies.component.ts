import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";

import { Movie } from '../movie';

import { MovieService } from '../movie.service';
import { UtilsService } from '../utils.service';
import { ConfigService } from '../config.service';

@Component({
    selector: 'app-movies',
    templateUrl: './movies.component.html',
    styleUrls: ['../media.miniatures.scss']
})
export class MoviesComponent implements OnInit {
    trackedMovies :Movie[];
    removedMovies :Movie[];
    downloadedMovies :Movie[];
    downloadingMovies :Movie[];

    refreshing :boolean;
    refreshing_poll :boolean;
    show_downloaded_movies = true;

    config :any;
    moviesRefresh :any;

    constructor(
        private http :HttpClient,
        private configService :ConfigService,
        private movieService :MovieService,
        private utils :UtilsService
    ) {
        movieService.trackedMovies.subscribe(movies => {
            this.trackedMovies = movies;
        });
        movieService.removedMovies.subscribe(movies => {
            this.removedMovies = movies;
        });
        movieService.downloadedMovies.subscribe(movies => {
            this.downloadedMovies = movies;
        });
        movieService.downloadingMovies.subscribe(movies => {
            this.downloadingMovies = movies;
        });
        movieService.refreshing.subscribe(bool => {
            this.refreshing = bool;
        });
        configService.config.subscribe(cfg => {
            this.config = cfg;
        });
    }

    ngOnInit() {
        this.configService.getConfig();
        this.movieService.getMovies();
        this.moviesRefresh = setInterval(() => {
            this.configService.getConfig();
            this.movieService.getMovies();
        }, 30000);
    }

    ngOnDestroy() {
        clearInterval(this.moviesRefresh);
    }

    refreshMovies() :void {
        this.movieService.getMovies();
    }

    executePollLoop() :void {
        this.refreshing_poll = true;
        this.http.post('/api/v1/actions/poll', {}).subscribe(data => {
            //this.getTrackedMovies();
            //this.getRemovedMovies();
            //this.getDownloadedMovies();
            //this.getDownloadingMovies();
            this.refreshing_poll = false;
        });
    }

    //getTrackedMovies() :void {
    //this.movieService.getTrackedMovies().subscribe(movies => this.trackedMovies = movies);
    //}

    //getRemovedMovies() :void {
    //this.movieService.getRemovedMovies().subscribe(movies => this.removedMovies = movies);
    //}

    //getDownloadedMovies() :void {
    //this.movieService.getDownloadedMovies().subscribe(movies => this.downloadedMovies = movies);
    //}

    //getDownloadingMovies() :void {
    //this.movieService.getDownloadingMovies().subscribe(movies => this.downloadingMovies = movies);
    //}

    stopDownload(movie :Movie) :void {
        movie.DownloadingItem.AbortPending = true;
        this.movieService.stopDownload(movie.ID).subscribe(response => {
            this.movieService.getDownloadingMovies();
            movie.DownloadingItem.AbortPending = false;
        });
    }
}
