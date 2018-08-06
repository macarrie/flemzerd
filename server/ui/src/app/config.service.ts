import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";
import { Subject }    from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class ConfigService {
    constructor(private http :HttpClient) {}

    private configSource = new Subject<any>();
    config = this.configSource.asObservable();

    getConfig() {
        this.http.get('/api/v1/config').subscribe(cfg => {
            this.configSource.next(cfg);
        });
    }
}
