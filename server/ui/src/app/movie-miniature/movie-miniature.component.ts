import { Component, OnInit, Input } from '@angular/core';

import { Movie } from '../movie';
import { MovieService } from '../movie.service';

@Component({
    selector: 'movie-miniature',
    templateUrl: './movie-miniature.component.html',
    styleUrls: ['../media.miniatures.scss']
})
export class MovieMiniatureComponent implements OnInit {
    @Input() movie :Movie;

    constructor(private movieService :MovieService) {}

    ngOnInit() {}

    refreshMovie() :void {
        //this.movieService.getMovie(this.movie.ID).subscribe(movie => this.movie = movie);
        this.movieService.getMovies();
    }

    deleteMovie(movie :Movie) {
        if (movie.DownloadingItem.Downloading) {
            this.movieService.stopDownload(movie.ID).subscribe();
        }
        this.movieService.deleteMovie(movie.ID).subscribe(response => this.refreshMovie());
    }

    restoreMovie(id :number) {
        this.movieService.restoreMovie(id).subscribe(response => this.refreshMovie());
    }

    changeDownloadedState(m :Movie, downloaded :boolean) {
        this.movieService.changeDownloadedState(m, downloaded).subscribe(movie => {
            this.refreshMovie();
        });
    }
}
