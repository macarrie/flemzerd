import { Component, OnInit, OnDestroy } from '@angular/core';
import { Observable, of, Subject } from 'rxjs';
import { HttpClient } from "@angular/common/http";
import { takeUntil  } from 'rxjs/operators'; // for rxjs ^5.5.0 lettable operators

import { Movie } from '../movie';

import { MovieService } from '../movie.service';
import { UtilsService } from '../utils.service';
import { ConfigService } from '../config.service';

@Component({
    selector: 'app-movies',
    templateUrl: './movies.component.html',
    styleUrls: ['../media.miniatures.scss']
})
export class MoviesComponent implements OnInit, OnDestroy {
    private subs :Subject<any> = new Subject();

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
        movieService.trackedMovies.pipe(takeUntil(this.subs)).subscribe(movies => {
            this.trackedMovies = movies;
        });
        movieService.removedMovies.pipe(takeUntil(this.subs)).subscribe(movies => {
            this.removedMovies = movies;
        });
        movieService.downloadedMovies.pipe(takeUntil(this.subs)).subscribe(movies => {
            this.downloadedMovies = movies;
        });
        movieService.downloadingMovies.pipe(takeUntil(this.subs)).subscribe(movies => {
            this.downloadingMovies = movies;
        });
        movieService.refreshing.pipe(takeUntil(this.subs)).subscribe(bool => {
            this.refreshing = bool;
        });
        configService.config.pipe(takeUntil(this.subs)).subscribe(cfg => {
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
        this.subs.next();
        this.subs.complete();
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

    getTitle(movie :Movie) :string {
        if (movie.UseDefaultTitle) {
            return movie.Title;
        }

        return movie.OriginalTitle;
    }

    stopDownload(movie :Movie) :void {
        movie.DownloadingItem.AbortPending = true;
        this.movieService.stopDownload(movie.ID).subscribe(response => {
            this.movieService.getDownloadingMovies();
            movie.DownloadingItem.AbortPending = false;
        });
    }
}
