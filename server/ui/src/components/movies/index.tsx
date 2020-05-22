import React from "react";
import {Route} from "react-router-dom";

import API from "../../utils/api";
import Const from "../../const";
import Movie from "../../types/movie";
import {MediaMiniatureFilter} from "../../types/media_miniature_filter";

import Empty from "../empty";
import MediaMiniature from "../media_miniature";
import ItemFilterControls from "../item_filter_controls";
import MovieDetails from "./movie_details";
import DownloadingItemTable from "../downloading_items_table";
import LoadingButton from "../loading_button";
import Helpers from "../../utils/helpers";
import ItemSearchControls from "../item_search_controls";

import {RiLayoutRowLine} from "react-icons/ri";
import {RiLayoutGridLine} from "react-icons/ri";

type State = {
    tracked :Movie[] | null,
    removed :Movie[] | null,
    downloaded :Movie[] | null,
    downloading :Movie[] | null,
    display_mode :number,
    item_filter :MediaMiniatureFilter,
    search_filter :string,
};

class Movies extends React.Component<any, State> {
    movies_refresh_interval :number;
    state :State = {
        tracked: null,
        removed: null,
        downloaded: null,
        downloading: null,
        display_mode: Const.DISPLAY_MINIATURES,
        item_filter: MediaMiniatureFilter.TRACKED,
        search_filter: "",
    };

    constructor(props: any) {
        super(props);

        this.movies_refresh_interval = 0;

        this.setFilter = this.setFilter.bind(this);
        this.setSearch = this.setSearch.bind(this);
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

    setFilter(val :MediaMiniatureFilter) {
        this.setState({item_filter: val});
    }

    setSearch(val :string) {
        this.setState({search_filter: val});
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

    getActiveLayoutClass(targetLayout :number) :String {
        if (targetLayout === this.state.display_mode) {
            return "is-info is-outlined";
        }

        return "";
    }

    renderDownloadingList() {
        if (this.state.downloading) {
            return (
                <>
                    <div className="columns">
                        <div className="column">
                            <span className="title is-3">Downloading movies</span>
                        </div>
                    </div>
                    <hr />
                    <DownloadingItemTable
                        list={this.state.downloading}
                        type="movie"
                        skipTorrent={this.skipTorrentDownload}
                        abortDownload={this.abortMovieDownload}/>
                    <hr />
                </>
            );
        }

        return "";
    }

    renderMovieListLayout() {
        // Load finished but no movies are track by flemzerd
        let total_count = Helpers.count(this.state.tracked) + Helpers.count(this.state.downloaded) + Helpers.count(this.state.downloading) + Helpers.count(this.state.removed);
        if (this.state.item_filter === MediaMiniatureFilter.NONE && total_count === 0) {
            return <Empty label={"No movies"} />;
        }
        if (this.state.item_filter === MediaMiniatureFilter.TRACKED && Helpers.count(this.state.tracked) === 0) {
            return <Empty label={"No tracked movies"} />;
        }
        if (this.state.item_filter === MediaMiniatureFilter.DOWNLOADED && Helpers.count(this.state.downloaded) === 0) {
            return <Empty label={"No downloaded movies"} />;
        }
        if (this.state.item_filter === MediaMiniatureFilter.DOWNLOADING && Helpers.count(this.state.downloading) === 0) {
            return <Empty label={"No downloading movies"} />;
        }
        if (this.state.item_filter === MediaMiniatureFilter.REMOVED && Helpers.count(this.state.removed) === 0) {
            return <Empty label={"No removed movies"} />;
        }

        let futures :Array<Movie>;
        if (this.state.tracked !== null) {
            futures = this.state.tracked.filter(function(movie) { return Helpers.dateIsInFuture(movie.Date)})
        } else {
            futures = new Array<Movie>();
        }
        if (this.state.item_filter === MediaMiniatureFilter.FUTURE && Helpers.count(futures) === 0) {
            return <Empty label={"No future movies"} />;
        }

        if (this.state.display_mode === Const.DISPLAY_MINIATURES) {
            return (
                <div className={`columns is-multiline is-mobile`}>
                    {this.renderMovieListContent()}
                </div>
            );
        } else {
            return (
                <table className="table is-fullwidth media-miniature-table">
                    <thead>
                        <tr>
                            <th>Title</th>
                            <th className={"has-text-centered"}>State</th>
                            <th className={"has-text-centered is-narrow"}>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
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
                        filter={this.state.item_filter}
                        search={this.state.search_filter}
                        display_mode={this.state.display_mode}
                        type="movie" />
                    ))
                ) : ""
                }
                {this.state.tracked ? (
                    this.state.tracked.map(movie => (
                        <MediaMiniature key={movie.ID}
                        item={movie}
                        filter={this.state.item_filter}
                        search={this.state.search_filter}
                        display_mode={this.state.display_mode}
                        type="movie" />
                    ))
                ) : ""
                }
                {this.state.downloaded ? (
                    this.state.downloaded.map(movie => (
                        <MediaMiniature key={movie.ID}
                        item={movie}
                        filter={this.state.item_filter}
                        search={this.state.search_filter}
                        display_mode={this.state.display_mode}
                        type="movie" />
                    ))
                ) : ""
                }
                {this.state.removed ? (
                    this.state.removed.map(movie => (
                        <MediaMiniature key={movie.ID}
                        item={movie}
                        filter={this.state.item_filter}
                        search={this.state.search_filter}
                        display_mode={this.state.display_mode}
                        type="movie" />
                    ))
                ) : ""
                }
            </>
        );
    }

    renderMovieList() {
        return (
            <div className="container">
                {this.renderDownloadingList()}

                <div className="columns is-vcentered is-multiline is-mobile">
                    <div className="column is-narrow">
                        <span className="title is-3">Movies</span>
                    </div>
                    <ItemSearchControls
                        updateSearch={this.setSearch}
                    />
                    <ItemFilterControls
                        type={"movie"}
                        filterValue={this.state.item_filter}
                        updateFilter={this.setFilter}
                    />
                    <div className={"column is-narrow field is-marginless is-grouped"}>
                        <LoadingButton
                            text="Check now"
                            loading_text="Checking"
                            action={this.poll}
                        />
                        <LoadingButton
                            text="Refresh"
                            loading_text="Refreshing"
                            action={this.refreshWatchlists}
                        />
                    </div>
                    <div className="column is-narrow is-hidden-mobile">
                        <div className="buttons has-addons">
                            <button className={`button ${this.getActiveLayoutClass(Const.DISPLAY_MINIATURES)}`} onClick={() => this.changeDisplayMode(Const.DISPLAY_MINIATURES)}>
                                <span className="icon is-small">
                                    <RiLayoutGridLine />
                                </span>
                            </button>
                            <button className={`button ${this.getActiveLayoutClass(Const.DISPLAY_LIST)}`} onClick={() => this.changeDisplayMode(Const.DISPLAY_LIST)}>
                                <span className="icon is-small">
                                    <RiLayoutRowLine />
                                </span>
                            </button>
                        </div>
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
