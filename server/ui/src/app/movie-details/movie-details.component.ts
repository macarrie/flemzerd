import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Location } from '@angular/common';

import { Movie } from '../movie';

import { MovieService } from '../movie.service';
import { FanartService } from '../fanart.service';
import { UtilsService } from '../utils.service';

@Component({
    selector: 'app-movie-details',
    templateUrl: './movie-details.component.html',
    styleUrls: ['../media.details.scss']
})
export class MovieDetailsComponent implements OnInit {
    movie :Movie;
    background_image :string;

    constructor(
        private movieService :MovieService,
        private fanartService :FanartService,
        private utils :UtilsService,
        private route :ActivatedRoute,
        private location :Location
    ) {}

    ngOnInit() {
        this.getMovie();
    }

    getMovie() :void {
        const id = +this.route.snapshot.paramMap.get('id');
        this.movieService.getMovie(id).subscribe(movie => {
            this.movie = movie;

            this.fanartService.getMovieFanart(movie).subscribe(fanartObj => {
                if (fanartObj["moviebackground"] && fanartObj["moviebackground"].length > 0) {
                    this.background_image = fanartObj["moviebackground"][0]["url"];
                }
            });
        });
    }

    deleteMovie(movie :Movie) :void {
        this.movieService.deleteMovie(movie.ID).subscribe(movie => {
            this.getMovie();
        });
    }

    restoreMovie(movie :Movie) :void {
        this.movieService.restoreMovie(movie.ID).subscribe(movie => {
            this.getMovie();
        });
    }

    changeDownloadedState(m :Movie, downloaded :boolean) {
        this.movieService.changeDownloadedState(m, downloaded).subscribe(movie => {
            this.getMovie();
        });
    }

    back() :void {
        this.location.back();
    }
}
