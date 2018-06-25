import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Routes } from '@angular/router';
import { DashboardComponent } from './dashboard/dashboard.component';
import { TvshowsComponent } from './tvshows/tvshows.component';
import { TvShowDetailsComponent } from './tv-show-details/tv-show-details.component';
import { EpisodeDetailsComponent } from './episode-details/episode-details.component';
import { MoviesComponent } from './movies/movies.component';
import { MovieDetailsComponent } from './movie-details/movie-details.component';
import { StatusComponent } from './status/status.component';
import { SettingsComponent } from './settings/settings.component';

const routes: Routes = [
    { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
    { path: 'dashboard', component: DashboardComponent },
    { path: 'tvshows', component: TvshowsComponent },
    { path: 'tvshows/:id', component: TvShowDetailsComponent },
    { path: 'episodes/:id', component: EpisodeDetailsComponent },
    { path: 'movies', component: MoviesComponent },
    { path: 'movies/:id', component: MovieDetailsComponent },
    { path: 'status', component: StatusComponent },
    { path: 'settings', component: SettingsComponent }
];

@NgModule({
    imports: [
        CommonModule,
        RouterModule.forRoot(routes)
    ],
    exports: [ RouterModule ],
    declarations: []
})
export class AppRoutingModule { }
