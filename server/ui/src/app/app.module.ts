import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { AppRoutingModule } from './/app-routing.module';
import { TvshowsComponent } from './tvshows/tvshows.component';
import { MoviesComponent } from './movies/movies.component';
import { StatusComponent } from './status/status.component';
import { SettingsComponent } from './settings/settings.component';
import { TvShowMiniatureComponent } from './tv-show-miniature/tv-show-miniature.component';
import { TvShowDetailsComponent } from './tv-show-details/tv-show-details.component';
import { EpisodeDetailsComponent } from './episode-details/episode-details.component';
import { DownloadingItemComponent } from './downloading-item/downloading-item.component';
import { MovieMiniatureComponent } from './movie-miniature/movie-miniature.component';
import { MovieDetailsComponent } from './movie-details/movie-details.component';

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    TvshowsComponent,
    MoviesComponent,
    StatusComponent,
    SettingsComponent,
    TvShowMiniatureComponent,
    TvShowDetailsComponent,
    EpisodeDetailsComponent,
    DownloadingItemComponent,
    MovieMiniatureComponent,
    MovieDetailsComponent
  ],
  imports: [
    BrowserModule,
      AppRoutingModule,
      HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
