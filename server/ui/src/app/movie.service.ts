import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";
import { Subject }    from 'rxjs';

import { Movie } from './movie';

@Injectable({
    providedIn: 'root'
})
export class MovieService {
    constructor(private http :HttpClient) {}

    private trackedMoviesSource =  new Subject<Movie[]>();
    private removedMoviesSource =  new Subject<Movie[]>();
    private downloadedMoviesSource =  new Subject<Movie[]>();
    private downloadingMoviesSource =  new Subject<Movie[]>();
    private refreshingSource =  new Subject<boolean>();

    trackedMovies = this.trackedMoviesSource.asObservable();
    removedMovies = this.removedMoviesSource.asObservable();
    downloadedMovies = this.downloadedMoviesSource.asObservable();
    downloadingMovies = this.downloadingMoviesSource.asObservable();
    refreshing = this.refreshingSource.asObservable();

    getMovies() {
        this.refreshingSource.next(true);
        this.getTrackedMovies();
        this.getRemovedMovies();
        this.getDownloadedMovies();
        this.getDownloadingMovies();

        // Sleep for 1s to have a visual refresh feedback
        setTimeout(() => {
            this.refreshingSource.next(false);
        }, 3000);
    }

    getTrackedMovies() {
        this.http.get<Movie[]>('/api/v1/movies/tracked').subscribe(movies => {
            this.trackedMoviesSource.next(movies);
        });
    }

    getRemovedMovies() {
        this.http.get<Movie[]>('/api/v1/movies/removed').subscribe(movies => {
            this.removedMoviesSource.next(movies);
        });
    }

    getDownloadedMovies() {
        this.http.get<Movie[]>('/api/v1/movies/downloaded').subscribe(movies => {
            this.downloadedMoviesSource.next(movies);
        });
    }

    getDownloadingMovies() {
        this.http.get<Movie[]>('/api/v1/movies/downloading').subscribe(movies => {
            this.downloadingMoviesSource.next(movies);
        });
    }

    getMovie(id :number) :Observable<Movie> {
        return this.http.get<Movie>('/api/v1/movies/details/' + id);
    }

    changeDownloadedState(m :Movie, downloaded :boolean) :Observable<Movie>{
        if (m.ID != 0) {
            m.DownloadingItem.Downloaded = downloaded;
            return this.http.put<Movie>('/api/v1/movies/details/' + m.ID, m);
        }
    }

    useDefaultTitle(m :Movie, useDefault :boolean) :Observable<Movie>{
        if (m.ID != 0) {
            m.UseDefaultTitle = useDefault;
            return this.http.put<Movie>('/api/v1/movies/details/' + m.ID, m);
        }
    }

    deleteMovie(id :number) {
        if (id != 0) {
            return this.http.delete('/api/v1/movies/details/' + id);
        }
    }

    downloadMovie(id :number) {
        if (id != 0) {
            return this.http.post('/api/v1/movies/details/' + id + '/download', {});
        }
    }

    stopDownload(id :number) {
        if (id != 0) {
            return this.http.delete('/api/v1/movies/details/' + id + '/download');
        }
    }

    restoreMovie(id :number) {
        if (id != 0) {
            return this.http.post('/api/v1/movies/restore/' + id, {});
        }
    }
}
