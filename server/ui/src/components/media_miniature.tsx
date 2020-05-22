import React from "react";
import {Link} from "react-router-dom";

import API from "../utils/api";
import Const from "../const";
import Helpers from "../utils/helpers";
import Movie from "../types/movie";
import TvShow from "../types/tvshow";
import {MediaMiniatureFilter} from "../types/media_miniature_filter";

import Fuse from "fuse.js";

import {RiCloseLine} from "react-icons/ri";
import {RiCheckLine} from "react-icons/ri";
import {RiDownloadLine} from "react-icons/ri";
import {RiEraserLine} from "react-icons/ri";
import {RiEyeLine} from "react-icons/ri";
import {RiArrowGoBackLine} from "react-icons/ri";

type Props = {
    item :Movie | TvShow,
    type :string,
    display_mode :number,
    filter :MediaMiniatureFilter,
    search :string,
};

type State = {
    item :Movie | TvShow;
    display_mode :number,
    filter :MediaMiniatureFilter,
    search :string,
};

class MediaMiniature extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = {
            item: props.item,
            display_mode: props.display_mode,
            filter: props.filter,
            search: props.search,
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
            filter: nextProps.filter,
            search: nextProps.search,
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
        let classNames = "";
        if (this.state.display_mode === Const.DISPLAY_MINIATURES) {
            classNames = "has-text-center";
        }

        if (this.props.type === "movie") {
            let item = this.state.item as Movie;
            if (!item.DeletedAt && Helpers.dateIsInFuture(item.Date)) {
                label = "Future release";
                classNames += "has-text-warning";
            }

            if (!item.DeletedAt && (item.DownloadingItem.Downloading || item.DownloadingItem.Pending)) {
                if (item.DownloadingItem.Pending) {
                    label = "Pending";
                }
                if (item.DownloadingItem.Downloading) {
                    label = "Downloading";
                }
                classNames += "has-text-info";
            }
            if (item.DownloadingItem.Downloaded && !item.DeletedAt) {
                label = "Downloaded";
                classNames += "has-text-success";
            }
        }

        if (this.state.item.DeletedAt) {
            label = "Removed";
            classNames += "has-text-danger";
        }

        return (
            <span className={classNames}>
                {label}
            </span>
        );
    }

    itemIsFiltered() :Boolean {
        if (this.state.filter === MediaMiniatureFilter.NONE) {
            return false;
        }

        if (this.props.type === "movie") {
            return this.movieItemIsFiltered();
        } else if (this.props.type === "tvshow") {
            return this.tvShowItemIsFiltered();
        }

        return false;
    }

    itemMatchesSearch() :Boolean {
        const options = {
            includeScore: true,
            shouldSort: true,
            threshold: 0.5,
            keys: ['name']
        }

        const fuse = new Fuse([{
            name: Helpers.getMediaTitle(this.state.item)
        }], options)

        const result = fuse.search(this.state.search)

        if (this.state.search === "") {
            return true;
        }
         if (result.length > 0) {
             console.log("SEARCH RESULT: ", result);
             return true;
         }

        return false;
    }

    tvShowItemIsFiltered() :Boolean {
        let item = this.state.item as TvShow;

        if (item.DeletedAt && this.state.filter === MediaMiniatureFilter.REMOVED) {
            return false;
        }

        if (this.state.filter === MediaMiniatureFilter.TRACKED) {
            if (item.DeletedAt) {
                return true;
            }

            return false;
        }

        return true;
    }

    movieItemIsFiltered() :Boolean {
        let item = this.state.item as Movie;

        if (Helpers.dateIsInFuture(item.Date) && !item.DeletedAt && this.state.filter === MediaMiniatureFilter.FUTURE) {
            return false;
        }
        if (item.DownloadingItem.Downloaded && this.state.filter === MediaMiniatureFilter.DOWNLOADED) {
            return false;
        }
        if (item.DeletedAt && this.state.filter === MediaMiniatureFilter.REMOVED) {
            return false;
        }

        if (this.state.filter === MediaMiniatureFilter.TRACKED) {
            if (Helpers.dateIsInFuture(item.Date) || item.DownloadingItem.Downloaded || item.DeletedAt) {
                return true;
            }

            return false;
        }

        return true;
    }

    getDownloadControlButtons() {
        if (this.props.type !== "movie") {
            return;
        }

        let item = this.state.item as Movie;
        let buttonList :any = [];

        if (item.DownloadingItem.Downloaded && (!item.DeletedAt)) {
            buttonList.push(
                <div className={"tile"}
                    key="markNotDownloaded">
                    <button className=""
                            onClick={this.markNotDownloaded}
                            data-tooltip="Unmark as downloaded">
                        <span className={"icon"}>
                            <RiEraserLine />
                        </span>
                    </button>
                </div>
            )
        }

        if (!item.DownloadingItem.Downloaded && !item.DownloadingItem.Downloading && !item.DownloadingItem.Downloaded && !item.DeletedAt) {
            buttonList.push(
                <div className={"tile"}
                     key="markDownloaded">
                    <button className=""
                            onClick={this.markDownloaded}
                            data-tooltip="Mark as downloaded">
                        <span className={"icon"}>
                            <RiCheckLine />
                        </span>
                    </button>
                </div>
            );
        }

        //TODO: Handle configuration
        if (!item.DownloadingItem.Downloaded && !item.DownloadingItem.Downloading && !item.DownloadingItem.Pending && !item.DeletedAt && !Helpers.dateIsInFuture(item.Date)) {
            buttonList.push(
                <div className={"tile"}
                    key="downloadMovie">
                    <button className=""
                            onClick={this.downloadMovie}
                            data-tooltip="Download">
                        <span className={"icon"}>
                            <RiDownloadLine />
                        </span>
                    </button>
                </div>
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
        // Update locally first to update display immediately
        let movie = this.state.item as Movie;
        movie.DeletedAt = new Date();
        this.setState({item: movie});

        API.Movies.delete(this.state.item.ID).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Delete movie error: ", error);
        });
    }

    deleteShow() {
        // Update locally first to update display immediately
        let show = this.state.item as TvShow;
        show.DeletedAt = new Date();
        this.setState({item: show});

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
        // Update locally first to update display immediately
        let movie = this.state.item as Movie;
        movie.DeletedAt = null;
        this.setState({item: movie});

        API.Movies.restore(this.state.item.ID).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Restore movie error: ", error);
        });
    }

    restoreShow() {
        // Update locally first to update display immediately
        let show = this.state.item as TvShow;
        show.DeletedAt = null;
        this.setState({item: show});

        API.Shows.restore(this.state.item.ID).then(response => {
            this.refreshShow();
        }).catch(error => {
            console.log("Restore show error: ", error);
        });
    }


    downloadMovie() {
        // Update locally first to update display immediately
        let movie = this.state.item as Movie;
        movie.DownloadingItem.Pending = true;
        this.setState({item: movie});

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
        // Update locally first to update display immediately
        let movie = this.state.item as Movie;
        movie.DownloadingItem.Downloaded = downloaded_state;
        this.setState({item: movie});

        API.Movies.changeDownloadedState(this.state.item.ID, downloaded_state).then(response => {
            this.refreshMovie();
        }).catch(error => {
            console.log("Change movie downloaded state error: ", error);
        });
    }

    render() {
        if (this.itemIsFiltered()) {
            return "";
        }
        if (!this.itemMatchesSearch()) {
            return "";
        }
        if (this.state.display_mode === Const.DISPLAY_MINIATURES) {
            return (
                <div className={`column is-half-mobile is-one-third-tablet is-2-desktop media-miniature is-relative`}>
                    <div className={"is-relative"}>
                        <Link to={`/${this.props.type}s/${this.state.item.ID}`}>
                            <div className={"thumbnail-container"}>
                                <span className={"thumbnail-alt"}>{this.state.item.Title}</span>
                                <img src={this.state.item.Poster} alt={""}
                                     className={`thumbnail ${this.getOverlayClassNames() !== "" ? "black_and_white" : ""}`} />
                            </div>
                            <div className={`thumbnail-overlay is-overlay ${this.getOverlayClassNames()}`}>
                                {this.getDownloadStatusLabel()}
                            </div>
                        </Link>
                    </div>

                    <div className="tile is-ancestor has-text-centered is-vertical media-miniature-controls is-overlay is-hidden-touch">
                        <div className={"tile"}>
                            <Link to={`/${this.props.type}s/${this.state.item.ID}`}>
                                <span data-tooltip="See details" className="icon">
                                    <RiEyeLine />
                                </span>
                            </Link>
                        </div>

                        {this.getDownloadControlButtons()}

                        {(!this.state.item.DeletedAt) ? (
                            <div className={"tile"}>
                                <button className="tooltip-danger"
                                        onClick={this.deleteItem}
                                        data-tooltip="Remove">
                                    <span className="icon">
                                        <RiCloseLine />
                                    </span>
                                </button>
                            </div>
                        ) : ""}

                        {(this.state.item.DeletedAt) ? (
                            <div className={"tile"}>
                                <button className=""
                                        onClick={this.restoreItem}
                                        data-tooltip="Restore">
                                    <span className="icon">
                                        <RiArrowGoBackLine />
                                    </span>
                                </button>
                            </div>
                        ) : ""}
                    </div>
                </div>
            );
        }

        if (this.state.display_mode === Const.DISPLAY_LIST) {
            return (
                <tr>
                    <td>
                        <Link to={`/${this.props.type}s/${this.state.item.ID}`}>
                            {Helpers.getMediaTitle(this.state.item)}
                        </Link>
                    </td>
                    <td className={"has-text-centered"}>
                        <div>
                            <small>
                                {this.getDownloadStatusLabel()}
                            </small>
                        </div>
                    </td>
                    <td className={"has-text-centered"}>
                        <div className={"action-buttons-container"}>
                            {this.getDownloadControlButtons()}

                            {(!this.state.item.DeletedAt) ? (
                                <div className={"tile"}>
                                    <button className=""
                                            onClick={this.deleteItem}
                                            data-tooltip="Remove">
                                    <span className={"icon"}>
                                        <RiCloseLine/>
                                    </span>
                                    </button>
                                </div>
                            ) : ""}

                            {(this.state.item.DeletedAt) ? (
                                <div className={"tile"}>
                                    <button className=""
                                            onClick={this.restoreItem}
                                            data-tooltip="Restore">
                                    <span className={"icon"}>
                                        <RiArrowGoBackLine/>
                                    </span>
                                    </button>
                                </div>
                            ) : ""}
                        </div>
                    </td>
                </tr>
            );
        }

        return null;
    }
}

export default MediaMiniature;
