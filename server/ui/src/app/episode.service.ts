import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { HttpClient } from "@angular/common/http";
import { Subject }    from 'rxjs';

import { Episode } from './episode';

@Injectable({
  providedIn: 'root'
})
export class EpisodeService {
    constructor(private http :HttpClient) {}

    private downloadingEpisodesSource =  new Subject<Episode[]>();

    downloadingEpisodes = this.downloadingEpisodesSource.asObservable();

    getEpisode(id :number): Observable<Episode> {
        return this.http.get<Episode>('/api/v1/tvshows/episodes/' +id);
    }

    getDownloadingEpisodes() {
        this.http.get<Episode[]>('/api/v1/tvshows/downloading').subscribe(episodes => {
            this.downloadingEpisodesSource.next(episodes);
        });
    }

    getBySeason(showID :number, seasonNumber :number) :Observable<Episode[]> {
        return this.http.get<Episode[]>('/api/v1/tvshows/details/' + showID + '/seasons/' + seasonNumber);
    }

    changeEpisodeDownloadedState(e :Episode, downloaded :boolean) :Observable<Episode>{
        if (e.ID != 0) {
            return this.http.put<Episode>('/api/v1/tvshows/episodes/' + e.ID, {
                "Downloaded": downloaded
            });
        }
    }

    downloadEpisode(e :Episode) :Observable<Episode>{
        if (e.ID != 0) {
            return this.http.post<Episode>('/api/v1/tvshows/episodes/' +e.ID+ "/download", {});
        }
    }

    deleteEpisode(e :Episode) {
        if (e.ID != 0) {
            return this.http.delete('/api/v1/tvshows/episodes/' +e.ID);
        }
    }

    stopDownload(id :number) {
        if (id != 0) {
            return this.http.delete('/api/v1/tvshows/episodes/' + id + '/download');
        }
    }
}
