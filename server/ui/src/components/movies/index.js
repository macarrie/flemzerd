import React from "react";
import { Route } from "react-router-dom";

import API from "../../utils/api";

import MediaMiniature from "../media_miniature";
import MovieDetails from "./movie_details";
import DownloadingItemTable from "../downloading_items_table";
import LoadingButton from "../loading_button";

class Movies extends React.Component {
    state = {
        tracked: null,
        removed: null,
        downloaded: null,
        downloading: null,
    };

    constructor(props) {
        super(props);

        this.poll = this.poll.bind(this);
        this.refreshWatchlists = this.refreshWatchlists.bind(this);

        this.renderDownloadingList = this.renderDownloadingList.bind(this);
        this.renderMovieList = this.renderMovieList.bind(this);
        this.abortMovieDownload = this.abortMovieDownload.bind(this);
    }

    componentDidMount() {
        this.getMovies();
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

    abortMovieDownload(id) {
        API.Movies.abortDownload(id).then(response => {
            this.getMovies();
        }).catch(error => {
            console.log("Abort movie download error: ", error);
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
                    abortDownload={this.abortMovieDownload}/>
                </>
            );
        }
    }

    renderMovieList() {
        return (
            <div className="uk-container" data-uk-filter="target: .item-filter">
                {this.renderDownloadingList()}

                <div className="uk-grid" data-uk-grid>
                    <div className="uk-width-expand">
                        <span className="uk-h3">Movies</span>
                    </div>
                    <ul className="uk-subnav uk-subnav-pill">
                        <li className="uk-visible@s uk-active" data-uk-filter-control="">
                            <button className="uk-button uk-button-text"> All </button>
                        </li>
                        <li className="uk-visible@s" data-uk-filter-control="filter: .tracked-movie">
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
                </div>
                <hr />

                <div className="uk-grid uk-grid-small uk-child-width-1-2 uk-child-width-1-6@l uk-child-width-1-4@m uk-child-width-1-3@s item-filter" data-uk-grid>
                    {this.state.downloading ? (
                        this.state.downloading.map(movie => (
                            <MediaMiniature key={movie.ID}
                                item={movie}
                                type="movie" />
                        ))
                    ) : ""
                    }
                    {this.state.tracked ? (
                        this.state.tracked.map(movie => (
                            <MediaMiniature key={movie.ID}
                                item={movie}
                                type="movie" />
                        ))
                    ) : ""
                    }
                    {this.state.downloaded ? (
                        this.state.downloaded.map(movie => (
                            <MediaMiniature key={movie.ID}
                                item={movie}
                                type="movie" />
                        ))
                    ) : ""
                    }
                    {this.state.removed ? (
                        this.state.removed.map(movie => (
                            <MediaMiniature key={movie.ID}
                                item={movie}
                                type="movie" />
                        ))
                    ) : ""
                    }
                </div>
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
