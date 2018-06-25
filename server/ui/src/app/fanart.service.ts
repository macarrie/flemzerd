import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";

import { TvShow } from './tvshow';
import { Movie } from './movie';

@Injectable({
    providedIn: 'root'
})
export class FanartService {
    FANART_TV_KEY :string = "648ff4214a1eea4416ad51417fc8a4e4";
    FANART_BASE_URL :string = "http://webservice.fanart.tv/v3/";

    params = {
        "params": {
            "api_key": this.FANART_TV_KEY
        }
    };

    constructor(private http :HttpClient) {}

    getTvShowFanart(show :TvShow) {
        return this.http.get(this.FANART_BASE_URL + "tv/" + show.MediaIds.Tvdb, this.params);
    }

    getMovieFanart(movie :Movie) {
        return this.http.get(this.FANART_BASE_URL + "movies/" + movie.MediaIds.Tmdb, this.params);
    }
}
