import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";

@Injectable({
    providedIn: 'root'
})
export class ConfigService {
    constructor(private http :HttpClient) {}

    getConfig() :Observable<any> {
        return this.http.get('/api/v1/config');
    }
}
