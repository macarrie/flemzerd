import React from "react";
import {Link} from "react-router-dom";

import API from "../utils/api";
import Const from "../const";
import Config from "../types/config";
import Stats from "../types/stats";
import Episode from "../types/episode";

import Empty from "./empty";
import Movie from "../types/movie";
import Helpers from "../utils/helpers";
import DownloadingItemTable from "./downloading_items_table";

type State = {
    config :Config | null,
    stats :Stats | null,
    downloading_episodes :Episode[],
    downloading_movies :Movie[],
};

class Dashboard extends React.Component<any, State> {
    refresh_interval :number;

    state :State = {
        config: null,
        stats: null,
        downloading_episodes: Array<Episode>(),
        downloading_movies: Array<Movie>(),
    };

    constructor(props :any) {
        super(props);

        this.refresh_interval = 0;
    }

    componentDidMount() {
        this.load();

        this.refresh_interval = window.setInterval(this.load.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.refresh_interval);
    }

    load() {
        this.getConfig();
        this.getStats();
        this.getDownloadingEpisodes();
        this.getDownloadingMovies();
    }

    getConfig() {
        API.Config.get().then(response => {
            this.setState({config: response.data});
        }).catch(error => {
            console.log("Get config error: ", error);
        });
    }

    getStats() {
        API.Stats.get().then(response => {
            this.setState({stats: response.data.stats});
        }).catch(error => {
            console.log("Get stats error: ", error);
        });
    }

    getDownloadingEpisodes() {
        API.Episodes.downloading().then(response => {
            this.setState({downloading_episodes: response.data});
        }).catch(error => {
            console.log("Get downloading episodes error: ", error);
        });
    }

    getDownloadingMovies() {
        API.Movies.downloading().then(response => {
            this.setState({downloading_movies: response.data});
        }).catch(error => {
            console.log("Get downloading movies error: ", error);
        });
    }

    skipMovieTorrentDownload(id: number) {
        API.Movies.skipTorrent(id).then(response => {
            this.getDownloadingMovies();
        }).catch(error => {
            console.log("Skip torrent download error: ", error);
        });
    }

    abortMovieDownload(id :number) {
        API.Movies.abortDownload(id).then(response => {
            this.getDownloadingMovies();
        }).catch(error => {
            console.log("Abort movie download error: ", error);
        });
    }

    skipEpisodeTorrentDownload(id: number) {
        API.Episodes.skipTorrent(id).then(response => {
            this.getDownloadingEpisodes();
        }).catch(error => {
            console.log("Skip torrent download error: ", error);
        });
    }

    abortEpisodeDownload(id: number) {
        API.Episodes.abortDownload(id).then(response => {
            this.getDownloadingEpisodes();
        }).catch(error => {
            console.log("Abort episode download error: ", error);
        });
    }

    render() {
        if (this.state.config == null || this.state.stats == null) {
            return <Empty />;
        }

        return (
        <>
            <div className={"container"}>
                {(Helpers.count(this.state.downloading_episodes) + Helpers.count(this.state.downloading_movies) > 0) && (
                    <>
                        <div className="columns is-mobile is-vcentered is-multiline">
                            <div className="column is-narrow-tablet is-full-mobile">
                                <div className="title is-3">
                                    Currently downloading
                                </div>
                            </div>
                            <div className="column is-full-mobile">
                                <div className={"title is-5 has-text-grey subtitle-horizontal"}>
                                    ({Helpers.count(this.state.downloading_episodes)} episodes, {Helpers.count(this.state.downloading_movies)} movies)
                                </div>
                            </div>
                        </div>
                        <hr />
                        {Helpers.count(this.state.downloading_episodes) > 0 && (
                            <>
                                <DownloadingItemTable
                                    list={this.state.downloading_episodes}
                                    type="episode"
                                    skipTorrent={this.skipEpisodeTorrentDownload}
                                    abortDownload={this.abortEpisodeDownload}/>
                                <hr />
                            </>
                        )}
                        {Helpers.count(this.state.downloading_movies) > 0 && (
                            <>
                                <DownloadingItemTable
                                    list={this.state.downloading_movies}
                                    type="movie"
                                    skipTorrent={this.skipMovieTorrentDownload}
                                    abortDownload={this.abortMovieDownload}/>
                                <hr />
                            </>
                        )}
                    </>
                )}
            </div>
            <div className="container">
                <div className="columns has-text-centered is-multiline is-mobile">
                    {this.state.config.System.TrackShows && (
                            <>
                                <div className="column is-half-tablet is-full-mobile">
                                    <div className="title jumbo title-gradient">
                                        <Link to="/tvshows">
                                            {this.state.stats.Shows.Tracked}
                                        </Link>
                                    </div>
                                    <span className="title is-4">
                                        Tracked shows
                                    </span>
                                </div>
                                <div className="column is-half-tablet is-full-mobile">
                                    <div className="title jumbo title-gradient">
                                        <Link to="/tvshows">
                                            {this.state.stats.Episodes.Downloading}
                                        </Link>
                                    </div>
                                    <span className="title is-4">
                                        Downloading episodes
                                    </span>
                                </div>
                            </>
                        )}

                        {this.state.config.System.TrackMovies && (
                            <>
                                <div className="column is-one-third-tablet is-full-mobile">
                                    <div className="title jumbo title-gradient">
                                        <Link to="/movies">
                                            {this.state.stats.Movies.Tracked}
                                        </Link>
                                    </div>
                                    <span className="title is-4">
                                        New movies
                                    </span>
                                </div>
                                <div className="column is-one-third-tablet is-full-mobile">
                                    <div className="title jumbo is-1 title-gradient">
                                        <Link to="/movies">
                                            {this.state.stats.Movies.Downloading}
                                        </Link>
                                    </div>
                                    <span className="title is-4">
                                        Downloading movies
                                    </span>
                                </div>
                                <div className="column is-one-third-tablet is-full-mobile">
                                    <div className="title jumbo is-1 title-gradient">
                                        <Link to="/movies">
                                            {this.state.stats.Movies.Downloaded}
                                        </Link>
                                    </div>
                                    <span className="title is-4">
                                        Downloaded movies
                                    </span>
                                </div>
                            </>
                        )}
                    </div>
                </div>
            </>
        );
    }
}

export default Dashboard;
