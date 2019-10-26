import React from "react";
import {Link} from "react-router-dom";

import API from "../utils/api";
import Const from "../const";
import Helpers from "../utils/helpers";
import Movie from "../types/movie";
import TvShow from "../types/tvshow";

type Props = {
    item :Movie | TvShow,
    type :string,
    display_mode :number,
};

type State = {
    item :Movie | TvShow;
    display_mode :number,
};

class MediaMiniature extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = {
            item: props.item,
            display_mode: props.display_mode,
        };

        this.deleteItem = this.deleteItem.bind(this);
        this.deleteMovie = this.deleteMovie.bind(this);
        this.deleteShow = this.deleteShow.bind(this);

        this.restoreItem = this.restoreItem.bind(this);
        this.restoreMovie = this.restoreMovie.bind(this);
        this.restoreShow = this.restoreShow.bind(this);

        this.downloadMovie = this.downloadMovie.bind(this);

        this.changeDownloadedState = this.changeDownloadedState.bind(this);
        this.markDownloaded = this.markDownloaded.bind(this);
        this.markNotDownloaded = this.markNotDownloaded.bind(this);
    }

    componentWillReceiveProps(nextProps: Props) {
        this.setState({
            item: nextProps.item,
            display_mode: nextProps.display_mode,
        });
    }

    getOverlayClassNames() {
        if (this.props.type === "movie") {
            return this.getMovieOverlayClassNames();
        } else if (this.props.type === "tvshow") {
            return this.getShowOverlayClassNames();
        }

        return;
    }

    getMovieOverlayClassNames() {
        let item = this.state.item as Movie;
        let overlay = "";

        if (Helpers.dateIsInFuture(item.Date)) {
            overlay = "overlay-yellow";
        }

        if (item.DeletedAt) {
            overlay = "overlay-red";
        }

        if (item.DownloadingItem.Downloaded) {
            overlay = "overlay-green";
        }

        if (item.DownloadingItem.Pending || item.DownloadingItem.Downloading) {
            overlay = "overlay-blue";
        }

        return overlay;
    }

    getShowOverlayClassNames() {
        let item = this.state.item as TvShow;
        let overlay = "";

        if (item.DeletedAt) {
            overlay = "overlay-red";
        }

        return overlay;
    }

    getDownloadStatusLabel() {
        let label = "";

        if (this.props.type === "movie") {
            let item = this.state.item as Movie;
            if (!item.DeletedAt && Helpers.dateIsInFuture(item.Date)) {
                label = "Future release";
            }

            if (!item.DeletedAt && (item.DownloadingItem.Downloading || item.DownloadingItem.Pending)) {
                if (item.DownloadingItem.Downloading) {
                    label = "Downloading";

                    if (item.DownloadingItem.Pending) {
                        label = "Pending";
                    }
                }
            }
            if (item.DownloadingItem.Downloaded && !item.DeletedAt) {
                label = "Downloaded";
            }
        }

        if (this.state.item.DeletedAt) {
            label = "Removed";
        }

        let classNames = "";
        if (this.state.display_mode === Const.DISPLAY_MINIATURES) {
            classNames = "uk-position-center";
        }

        return (
            <span className={classNames}>
                {label}
            </span>
        );
    }


    getFilterClassName() {
        if (this.props.type === "movie") {
            return this.getMovieFilterClassName();
        } else if (this.props.type === "tvshow") {
            return this.getTvshowFilterClassName();
        }

        return;
    }

    getTvshowFilterClassName() {
        let className = "tracked-show";

        if (this.state.item.DeletedAt) {
            className = "removed-show";
        }

        return className;
    }

    getMovieFilterClassName() {
        let item = this.state.item as Movie;
        let className = "tracked-movie";

        if (Helpers.dateIsInFuture(item.Date)) {
            className = "future-movie";
        }
        if (item.DownloadingItem.Downloaded) {
            className = "downloaded-movie";
        }
        if (item.DeletedAt) {
            className = "removed-movie";
        }

        return className;
    }


    getDownloadControlButtons() {
        if (this.props.type !== "movie") {
            return;
        }

        let item = this.state.item as Movie;
        let buttonList :any = [];

        if (item.DownloadingItem.Downloaded && (!item.DeletedAt)) {
            buttonList.push(
                <button className="uk-icon"
                        key="markNotDownloaded"
                        onClick={this.markNotDownloaded}
                        data-uk-tooltip="delay: 500; title: Unmark as downloaded"
                        data-uk-icon="icon: push; ratio: 0.75">
                </button>
            )
        }

        if (!item.DownloadingItem.Downloaded && !item.DownloadingItem.Downloading && !item.DownloadingItem.Downloaded && !item.DeletedAt) {
            buttonList.push(
                <button className="uk-icon"
                        key="markDownloaded"
                        onClick={this.markDownloaded}
                        data-uk-tooltip="delay: 500; title: Mark as downloaded"
                        data-uk-icon="icon: check; ratio: 0.75">
                </button>
            );
        }

        //TODO: Handle configuration
        if (!item.DownloadingItem.Downloaded && !item.DownloadingItem.Downloading && !item.DownloadingItem.Pending && !item.DeletedAt && !Helpers.dateIsInFuture(item.Date)) {
            buttonList.push(
                <button className="uk-icon"
                        key="downloadMovie"
                        onClick={this.downloadMovie}
                        data-uk-tooltip="delay: 500; title: Download"
                        data-uk-icon="icon: download; ratio: 0.75">
                </button>
            )
        }

        return buttonList;
    }


    refreshMovie() {
        API.Movies.get(this.state.item.ID).then(response => {
            this.setState(response.data);
        }).catch(error => {
            console.log("Refresh movie error: ", error);
        });
    }

    refreshShow() {
        API.Shows.get(this.state.item.ID).then(response => {
            this.setState(response.data);
        }).catch(error => {
            console.log("Refresh show error: ", error);
        });
    }


    deleteItem() {
        if (this.props.type === "movie") {
            this.deleteMovie();
        } else if (this.props.type === "tvshow") {
            return this.deleteShow();
        }
    }

    deleteMovie() {
        API.Movies.delete(this.state.item.ID).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Delete movie error: ", error);
        });
    }

    deleteShow() {
        API.Shows.delete(this.state.item.ID).then(response => {
            this.refreshShow();
        }).catch(error => {
            console.log("Delete movie error: ", error);
        });
    }


    restoreItem() {
        if (this.props.type === "movie") {
            this.restoreMovie();
        } else if (this.props.type === "tvshow") {
            return this.restoreShow();
        }
    }

    restoreMovie() {
        API.Movies.restore(this.state.item.ID).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Restore movie error: ", error);
        });
    }

    restoreShow() {
        API.Shows.restore(this.state.item.ID).then(response => {
            this.refreshShow();
        }).catch(error => {
            console.log("Restore show error: ", error);
        });
    }


    downloadMovie() {
        API.Movies.download(this.state.item.ID).then(response => {
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

    changeDownloadedState(downloaded_state: boolean) {
        API.Movies.changeDownloadedState(this.state.item.ID, downloaded_state).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Change movie downloaded state error: ", error);
        });
    }

    render() {
        if (this.state.display_mode === Const.DISPLAY_MINIATURES) {
            return (
                <div className={this.getFilterClassName()}>
                    <div className="uk-inline media-miniature">
                        <div>
                            <Link to={`/${this.props.type}s/${this.state.item.ID}`}>
                                <img data-src={this.state.item.Poster} alt={this.state.item.Title}
                                     className={`uk-border-rounded ${this.getOverlayClassNames() !== "" ? "black_and_white" : ""}`}
                                     data-uk-img/>
                                <div
                                    className={`uk-overlay uk-border-rounded uk-position-cover ${this.getOverlayClassNames()}`}>
                                    {this.getDownloadStatusLabel()}
                                </div>
                            </Link>
                        </div>

                        <div className="uk-icon uk-position-top-right">
                            {this.getDownloadControlButtons()}

                            {(!this.state.item.DeletedAt) ? (
                                <button className="uk-icon"
                                        onClick={this.deleteItem}
                                        data-uk-tooltip="delay: 500; title: Remove"
                                        data-uk-icon="icon: close; ratio: 0.75">
                                </button>
                            ) : ""}

                            {(this.state.item.DeletedAt) ? (
                                <button className="uk-icon"
                                        onClick={this.restoreItem}
                                        data-uk-tooltip="delay: 500; title: Restore"
                                        data-uk-icon="icon: reply; ratio: 0.75">
                                </button>
                            ) : ""}
                        </div>
                    </div>
                </div>
            );
        }

        if (this.state.display_mode === Const.DISPLAY_LIST) {
            return (
                <tr className={this.getFilterClassName()}>
                    <td>
                        <Link to={`/${this.props.type}s/${this.state.item.ID}`}>
                            {Helpers.getMediaTitle(this.state.item)}
                        </Link>
                    </td>
                    <td>
                        <small>
                            {this.getDownloadStatusLabel()}
                        </small>
                    </td>
                    <td>
                        {this.getDownloadControlButtons()}

                        {(!this.state.item.DeletedAt) ? (
                            <button className="uk-icon"
                                onClick={this.deleteItem}
                                data-uk-tooltip="delay: 500; title: Remove"
                                data-uk-icon="icon: close; ratio: 0.75">
                            </button>
                        ) : ""}

                        {(this.state.item.DeletedAt) ? (
                            <button className="uk-icon"
                                onClick={this.restoreItem}
                                data-uk-tooltip="delay: 500; title: Restore"
                                data-uk-icon="icon: reply; ratio: 0.75">
                            </button>
                        ) : ""}
                    </td>
                </tr>
            );
        }

        return null;
    }
}

export default MediaMiniature;
