import React from "react";
import {Route} from "react-router-dom";

import API from "../../utils/api";
import Const from "../../const";
import Movie from "../../types/movie";

import Loading from "../loading";
import MediaMiniature from "../media_miniature";
import MovieDetails from "./movie_details";
import DownloadingItemTable from "../downloading_items_table";
import LoadingButton from "../loading_button";

type State = {
    tracked :Movie[] | null,
    removed :Movie[] | null,
    downloaded :Movie[] | null,
    downloading :Movie[] | null,
    display_mode :number,
};

class Movies extends React.Component<any, State> {
    movies_refresh_interval :number;
    state :State = {
        tracked: null,
        removed: null,
        downloaded: null,
        downloading: null,
        display_mode: Const.DISPLAY_MINIATURES,
    };

    constructor(props: any) {
        super(props);

        this.movies_refresh_interval = 0;

        this.poll = this.poll.bind(this);
        this.refreshWatchlists = this.refreshWatchlists.bind(this);

        this.renderDownloadingList = this.renderDownloadingList.bind(this);
        this.renderMovieList = this.renderMovieList.bind(this);
        this.abortMovieDownload = this.abortMovieDownload.bind(this);
        this.skipTorrentDownload = this.skipTorrentDownload.bind(this);
    }

    componentDidMount() {
        this.getMovies();

        this.movies_refresh_interval = window.setInterval(this.getMovies.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.movies_refresh_interval);
    }

    getMovies() {
        API.Movies.tracked().then(response => {
            this.setState({tracked: response.data});
        }).catch(error => {
            console.log("Get tracked movies error: ", error);
        });

        API.Movies.downloading().then(response => {
            this.setState({downloading: response.data});
        }).catch(error => {
            console.log("Get downloading movies error: ", error);
        });

        API.Movies.downloaded().then(response => {
            this.setState({downloaded: response.data});
        }).catch(error => {
            console.log("Get downloaded movies error: ", error);
        });

        API.Movies.removed().then(response => {
            this.setState({removed: response.data});
        }).catch(error => {
            console.log("Get removed movies error: ", error);
        });
    }

    poll() {
        return API.Actions.poll();
    }

    refreshWatchlists() {
        return API.Modules.Watchlists.refresh();
    }

    skipTorrentDownload(id: number) {
        API.Movies.skipTorrent(id).then(response => {
            this.getMovies();
        }).catch(error => {
            console.log("Skip torrent download error: ", error);
        });
    }

    abortMovieDownload(id :number) {
        API.Movies.abortDownload(id).then(response => {
            this.getMovies();
        }).catch(error => {
            console.log("Abort movie download error: ", error);
        });
    }

    changeDisplayMode(display_mode :number) {
        this.setState({
            display_mode: display_mode,
        });
    }

    renderDownloadingList() {
        if (this.state.downloading) {
            return (
                <>
                <div className="uk-width-1-1">
                <span className="uk-h3">Downloading movies</span>
                <hr />
                </div>
                <DownloadingItemTable 
                list={this.state.downloading}
                type="movie"
                skipTorrent={this.skipTorrentDownload}
                abortDownload={this.abortMovieDownload}/>
                </>
            );
        }

        return "";
    }

    renderMovieListLayout() {
        let filterClass :string = "item-filter";
        if (this.state.downloading === null && this.state.tracked === null && this.state.downloaded === null && this.state.removed === null) {
            filterClass = "";
        }

        if (this.state.display_mode === Const.DISPLAY_MINIATURES) {
            return (
                <div className={`uk-grid uk-grid-small uk-child-width-1-2 uk-child-width-1-6@l uk-child-width-1-4@m uk-child-width-1-3@s ${filterClass}`} data-uk-grid>
                    {this.renderMovieListContent()}
                </div>
            );
        } else {
            return (
                <table className="uk-table uk-table-divider uk-table-small">
                    <thead>
                        <tr>
                            <th className="uk-width-expand">Title</th>
                            <th>State</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody className={`${filterClass}`}>
                        {this.renderMovieListContent()}
                    </tbody>
                </table>
            );
        }
    }

    renderMovieListContent() {
        if (this.state.downloading === null && this.state.tracked === null && this.state.downloaded === null && this.state.removed === null) {
            return null;
        }

        return (
            <>
                {this.state.downloading ? (
                    this.state.downloading.map(movie => (
                        <MediaMiniature key={movie.ID}
                        item={movie}
                        display_mode={this.state.display_mode}
                        type="movie" />
                    ))
                ) : ""
                }
                {this.state.tracked ? (
                    this.state.tracked.map(movie => (
                        <MediaMiniature key={movie.ID}
                        item={movie}
                        display_mode={this.state.display_mode}
                        type="movie" />
                    ))
                ) : ""
                }
                {this.state.downloaded ? (
                    this.state.downloaded.map(movie => (
                        <MediaMiniature key={movie.ID}
                        item={movie}
                        display_mode={this.state.display_mode}
                        type="movie" />
                    ))
                ) : ""
                }
                {this.state.removed ? (
                    this.state.removed.map(movie => (
                        <MediaMiniature key={movie.ID}
                        item={movie}
                        display_mode={this.state.display_mode}
                        type="movie" />
                    ))
                ) : ""
                }
            </>
        );
    }

    renderMovieList() {
        if (this.state.downloaded == null && this.state.downloading == null && this.state.removed == null && this.state.tracked == null) {
            return (
                <Loading/>
            );
        }

        return (
            <div className="uk-container" data-uk-filter="target: .item-filter">
                {this.renderDownloadingList()}

                <div className="uk-grid" data-uk-grid>
                    <div className="uk-width-expand">
                        <span className="uk-h3">Movies</span>
                    </div>
                    <ul className="uk-subnav uk-subnav-pill filter-container">
                        <li className="uk-visible@s" data-uk-filter-control="">
                            <button className="uk-button uk-button-text"> All </button>
                        </li>
                        <li className="uk-visible@s uk-active" data-uk-filter-control="filter: .tracked-movie">
                            <button className="uk-button uk-button-text"> Tracked </button>
                        </li>
                        <li className="uk-visible@s" data-uk-filter-control="filter: .future-movie">
                            <button className="uk-button uk-button-text"> Future </button>
                        </li>
                        <li className="uk-visible@s" data-uk-filter-control="filter: .downloaded-movie">
                            <button className="uk-button uk-button-text"> Downloaded </button>
                        </li>
                        <li className="uk-visible@s" data-uk-filter-control="filter: .removed-movie">
                            <button className="uk-button uk-button-text"> Removed </button>
                        </li>
                        <li>
                            <LoadingButton 
                                text="Check now"
                                loading_text="Checking"
                                action={this.poll}
                            />
                        </li>
                        <li>
                            <LoadingButton 
                                text="Refresh"
                                loading_text="Refreshing"
                                action={this.refreshWatchlists}
                            />
                        </li>
                    </ul>

                    <div className="uk-button-group">
                        <button className="uk-button uk-button-default uk-button-small" onClick={() => this.changeDisplayMode(Const.DISPLAY_MINIATURES)}>
                            <span uk-icon="icon: grid; ratio: 0.7"></span>
                        </button>
                        <button className="uk-button uk-button-default uk-button-small" onClick={() => this.changeDisplayMode(Const.DISPLAY_LIST)}>
                            <span uk-icon="icon: list; ratio: 0.7"></span>
                        </button>
                    </div>
                </div>
                <hr />

                {this.renderMovieListLayout()}
            </div>
        )
    }

    render() {
        const match = this.props.match;

        return (
            <>
            <Route path={`${match.path}/:id`} component={MovieDetails} />
            <Route exact
            path={match.path}
            render={this.renderMovieList} />
            </>
        );
    }
}

export default Movies;
