import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";

@Injectable({
    providedIn: 'root'
})
export class ModulesService {
    constructor(private http :HttpClient) {}

    getProviders() :Observable<any> {
        return this.http.get('/api/v1/modules/providers/status');
    }

    getIndexers() :Observable<any> {
        return this.http.get('/api/v1/modules/indexers/status');
    }

    getDownloaders() :Observable<any> {
        return this.http.get('/api/v1/modules/downloaders/status');
    }

    getNotifiers() :Observable<any> {
        return this.http.get('/api/v1/modules/notifiers/status');
    }

    getWatchlists() :Observable<any> {
        return this.http.get('/api/v1/modules/watchlists/status');
    }

    getMediaCenters() :Observable<any> {
        return this.http.get('/api/v1/modules/mediacenters/status');
    }
}
