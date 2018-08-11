import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";
import { Subject }    from 'rxjs';

import { TvShow } from './tvshow';

@Injectable({
    providedIn: 'root'
})
export class TvshowsService {
    constructor(private http :HttpClient) {}

    private trackedShowsSource =  new Subject<TvShow[]>();
    private removedShowsSource =  new Subject<TvShow[]>();
    private refreshingSource =  new Subject<boolean>();

    trackedShows = this.trackedShowsSource.asObservable();
    removedShows = this.removedShowsSource.asObservable();
    refreshing = this.refreshingSource.asObservable();

    getShows() {
        this.refreshingSource.next(true);
        this.getTrackedTvShows();
        this.getRemovedTvShows();

        // Sleep for 1s to have a visual refresh feedback
        setTimeout(() => {
            this.refreshingSource.next(false);
        }, 3000);
    }

    refreshTvShows() {
        return this.http.post('/api/v1/modules/watchlists/refresh', {});
    }

    getTrackedTvShows() {
        this.http.get<TvShow[]>('/api/v1/tvshows/tracked').subscribe(shows => {
            this.trackedShowsSource.next(shows);
        });
    }

    getRemovedTvShows() {
        this.http.get<TvShow[]>('/api/v1/tvshows/removed').subscribe(shows => {
            this.removedShowsSource.next(shows);
        });
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

    useDefaultTitle(show :TvShow, useDefault :boolean) :Observable<TvShow>{
        if (show.ID != 0) {
            show.UseDefaultTitle = useDefault;
            return this.http.put<TvShow>('/api/v1/tvshows/details/' + show.ID, show);
        }
    }
}
