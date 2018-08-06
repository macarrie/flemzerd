import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";
import { Subject }    from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class ModulesService {
    constructor(private http :HttpClient) {}

    private indexersSource = new Subject<any>();
    private providersSource = new Subject<any>();
    private downloadersSource = new Subject<any>();
    private notifiersSource = new Subject<any>();
    private watchlistsSource = new Subject<any>();
    private mediacentersSource = new Subject<any>();

    indexers = this.indexersSource.asObservable();
    providers = this.providersSource.asObservable();
    downloaders = this.downloadersSource.asObservable();
    notifiers = this.notifiersSource.asObservable();
    watchlists = this.watchlistsSource.asObservable();
    mediacenters = this.mediacentersSource.asObservable();

    getProviders() {
        this.http.get('/api/v1/modules/providers/status').subscribe(mods => {
            this.providersSource.next(mods);
        });
    }

    getIndexers() {
        this.http.get('/api/v1/modules/indexers/status').subscribe(mods => {
            this.indexersSource.next(mods);
        });
    }

    getDownloaders() {
        this.http.get('/api/v1/modules/downloaders/status').subscribe(mods => {
            this.downloadersSource.next(mods);
        });
    }

    getNotifiers() {
        this.http.get('/api/v1/modules/notifiers/status').subscribe(mods => {
            this.notifiersSource.next(mods);
        });
    }

    getWatchlists() {
        this.http.get('/api/v1/modules/watchlists/status').subscribe(mods => {
            this.watchlistsSource.next(mods);
        });
    }

    getMediaCenters() {
        this.http.get('/api/v1/modules/mediacenters/status').subscribe(mods => {
            this.mediacentersSource.next(mods);
        });
    }
}
