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
    title :string;
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

    getTitle(movie :Movie) :string {
        if (movie.CustomTitle != "") {
            return movie.CustomTitle;
        }

        if (movie.UseDefaultTitle) {
            return movie.Title;
        }

        return movie.OriginalTitle;
    }

    getMovie() :void {
        const id = +this.route.snapshot.paramMap.get('id');
        this.movieService.getMovie(id).subscribe(movie => {
            this.movie = movie;
            this.title = this.getTitle(movie);
            if (this.title == "") {
                if (movie.UseDefaultTitle) {
                    this.title = movie.Title;
                } else {
                    this.title = movie.OriginalTitle;
                }
            }


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

    downloadMovie(id :number) {
        this.movie.DownloadingItem.Pending = true;
        this.movieService.downloadMovie(id).subscribe(response => this.getMovie());
    }

    changeDownloadedState(m :Movie, downloaded :boolean) {
        this.movieService.changeDownloadedState(m, downloaded).subscribe(movie => {
            this.getMovie();
        });
    }

    useDefaultTitle(m :Movie, defaultTitle :boolean) {
        this.movieService.useDefaultTitle(m, defaultTitle).subscribe(movie => {
            this.getMovie();
        });
    }

    setCustomTitle(title :string) {
        this.movie.CustomTitle = title;
        this.movieService.updateMovie(this.movie).subscribe(movie => {
            this.getMovie();
        });
    }

    back() :void {
        this.location.back();
    }
}
