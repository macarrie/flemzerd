import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";

import { Movie } from './movie';

@Injectable({
    providedIn: 'root'
})
export class MovieService {
    constructor(private http :HttpClient) {}

    refreshMovies() {
        return this.http.post('/api/v1/modules/watchlists/refresh', {});
    }

    getTrackedMovies(): Observable<Movie[]> {
        return this.http.get<Movie[]>('/api/v1/movies/tracked');
    }

    getRemovedMovies(): Observable<Movie[]> {
        return this.http.get<Movie[]>('/api/v1/movies/removed');
    }

    getDownloadedMovies(): Observable<Movie[]> {
        return this.http.get<Movie[]>('/api/v1/movies/downloaded');
    }

    getDownloadingMovies(): Observable<Movie[]> {
        return this.http.get<Movie[]>('/api/v1/movies/downloading');
    }

    getMovie(id :number) :Observable<Movie> {
        return this.http.get<Movie>('/api/v1/movies/details/' + id);
    }

    changeDownloadedState(m :Movie, downloaded :boolean) :Observable<Movie>{
        if (m.ID != 0) {
            return this.http.put<Movie>('/api/v1/movies/details/' + m.ID, {
                "Downloaded": downloaded
            });
        }
    }

    deleteMovie(id :number) {
        if (id != 0) {
            return this.http.delete('/api/v1/movies/details/' + id);
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
