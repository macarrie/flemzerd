import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";

import { TvShow } from './tvshow';

@Injectable({
    providedIn: 'root'
})
export class TvshowsService {
    constructor(private http :HttpClient) {}

    refreshTvShows() {
        return this.http.post('/api/v1/modules/watchlists/refresh', {});
    }

    getTrackedTvShows(): Observable<TvShow[]> {
        return this.http.get<TvShow[]>('/api/v1/tvshows/tracked');
    }

    getRemovedTvShows(): Observable<TvShow[]> {
        return this.http.get<TvShow[]>('/api/v1/tvshows/removed');
    }

    getShow(id :number) :Observable<TvShow> {
        return this.http.get<TvShow>('/api/v1/tvshows/details/' + id);
    }

    deleteShow(id :number) {
        if (id != 0) {
            return this.http.delete('/api/v1/tvshows/details/' + id);
        }
    }

    restoreShow(id :number) {
        if (id != 0) {
            return this.http.post('/api/v1/tvshows/restore/' + id, {});
        }
    }
}
