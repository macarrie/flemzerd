import React from "react";
import { Link } from "react-router-dom";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";

class TvshowMiniature extends React.Component {
    constructor(props) {
        super(props);
        this.state = props.item;

        this.deleteMovie = this.deleteMovie.bind(this);
        this.restoreMovie = this.restoreMovie.bind(this);
        this.downloadMovie = this.downloadMovie.bind(this);

        this.changeDownloadedState = this.changeDownloadedState.bind(this);
        this.markDownloaded = this.markDownloaded.bind(this);
        this.markNotDownloaded = this.markNotDownloaded.bind(this);
    }

    getOverlayClassNames() {
        let overlay = "";

        if (Helpers.dateIsInFuture(this.state.Date)) {
            overlay = "overlay-yellow";
        }

        if (this.state.DeletedAt) {
            overlay = "overlay-red";
        }

        if (this.state.DownloadingItem.Downloaded) {
            overlay = "overlay-green";
        }

        if (this.state.DownloadingItem.Pending || this.state.DownloadingItem.Downloading) {
            overlay = "overlay-blue";
        }

        return overlay;
    }

    getDownloadStatusLabel() {
        let label = "";

        if (!this.state.DeletedAt && Helpers.dateIsInFuture(this.state.Date)) {
            label = "Future release";
        }

        if (!this.state.DeletedAt && (this.state.DownloadingItem.Downloading || this.state.DownloadingItem.Pending)) {
            if (this.state.DownloadingItem.Downloading) {
                label = "Downloading";

                if (this.state.DownloadingItem.Pending) {
                    label = "Pending";
                }
            }
        }

        if (this.state.DeletedAt) {
            label = "Removed";
        }
        if (this.state.DownloadingItem.Downloaded && !this.state.DeletedAt) {
            label = "Downloaded";
        }

        return (
            <span className="uk-position-center">
                {label}
            </span>
        );
    }

    getFilterClassName() {
        let className = "tracked-movie";

        if (this.state.DownloadingItem.Downloading || this.state.DownloadingItem.Pending) {
            className = "";
        }
        if (Helpers.dateIsInFuture(this.state.Date)) {
            className = "future-movie";
        }
        if (this.state.DownloadingItem.Downloaded) {
            className = "downloaded-movie";
        }
        if (this.state.DeletedAt) {
            className = "removed-movie";
        }

        return className;
    }

    refreshMovie() {
        API.Movies.get(this.state.ID).then(response => {
            this.setState(response.data);
        }).catch(error => {
            console.log("Refresh movie error: ", error);
        });
    }

    deleteMovie() {
        API.Movies.delete(this.state.ID).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Delete movie error: ", error);
        });
    }

    restoreMovie() {
        API.Movies.restore(this.state.ID).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Restore movie error: ", error);
        });
    }

    downloadMovie() {
        API.Movies.download(this.state.ID).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Download movie error: ", error);
        });
    }

    markDownloaded() {
        this.changeDownloadedState(true);
    }

    markNotDownloaded() {
        this.changeDownloadedState(false);
    }

    changeDownloadedState(downloaded_state) {
        API.Movies.changeDownloadedState(this.state.ID, downloaded_state).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Change movie downloaded state error: ", error);
        });
    }

    render() {
        return (
            <div className={this.getFilterClassName()}>
                <div className="uk-inline media-miniature">
                    <div>
                        <Link to={`/movies/${this.state.ID}`}>
                            <img src={ this.state.Poster } alt={ this.state.Title } className={`uk-border-rounded ${this.getOverlayClassNames() !== "" ? "black_and_white" : ""}`} data-uk-img />
                            <div className={`uk-overlay uk-border-rounded uk-position-cover ${this.getOverlayClassNames()}`}>
                                {this.getDownloadStatusLabel()}
                            </div>
                        </Link>
                    </div>
                
                    <div className="uk-icon uk-position-top-right">
                        { (this.state.DownloadingItem.Downloaded && (!this.state.DeletedAt)) ? (
                            <button className="uk-icon"
                                onClick={this.markNotDownloaded}
                                data-uk-tooltip="delay: 500; title: Unmark as downloaded"
                                data-uk-icon="icon: push; ratio: 0.75">
                            </button>
                        ) : ""}
                
                        { (!this.state.DownloadingItem.Downloaded && !this.state.DownloadingItem.Downloading && !this.state.DownloadingItem.Downloaded && !this.state.DeletedAt) ? (
                            <button className="uk-icon"
                                onClick={this.markDownloaded}
                                data-uk-tooltip="delay: 500; title: Mark as downloaded"
                                data-uk-icon="icon: check; ratio: 0.75">
                            </button>
                        ) : "" }
                
                        {/* TODO: Handle configuration */}
                        { (!this.state.DownloadingItem.Downloaded && !this.state.DownloadingItem.Downloading && !this.state.DownloadingItem.Pending && !this.state.DeletedAt && !Helpers.dateIsInFuture(this.state.Date)) ? (
                            <button className="uk-icon"
                                onClick={this.downloadMovie}
                                data-uk-tooltip="delay: 500; title: Download"
                                data-uk-icon="icon: download; ratio: 0.75">
                            </button>
                        ) : ""}
                
                        { (!this.state.DeletedAt && !this.state.DownloadingItem.Downloaded) ? (
                            <button className="uk-icon"
                                onClick={this.deleteMovie}
                                data-uk-tooltip="delay: 500; title: Remove"
                                data-uk-icon="icon: close; ratio: 0.75">
                            </button>
                        ) : ""}
                
                        { (this.state.DeletedAt && !this.state.DownloadingItem.Downloaded) ? (
                            <button className="uk-icon"
                                onClick={this.restoreMovie}
                                data-uk-tooltip="delay: 500; title: Restore"
                                data-uk-icon="icon: reply; ratio: 0.75">
                            </button>
                        ) : ""}
                    </div>
                </div>
            </div>
        );
    }
}

export default MovieMiniature;
