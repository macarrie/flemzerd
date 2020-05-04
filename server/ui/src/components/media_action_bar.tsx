import React from "react";
import {Link} from "react-router-dom";

import Helpers from "../utils/helpers";

import Movie from "../types/movie";
import TvShow from "../types/tvshow";
import Episode from "../types/episode";

import {RiArrowLeftLine} from "react-icons/ri";
import {RiCloseLine} from "react-icons/ri";
import {RiCheckLine} from "react-icons/ri";
import {RiDownloadLine} from "react-icons/ri";
import {RiEraserLine} from "react-icons/ri";
import {RiStopCircleLine} from "react-icons/ri";
import {RiArrowGoBackLine} from "react-icons/ri";
import {RiCheckboxMultipleBlankLine} from "react-icons/ri";
import {RiSkipForwardLine} from "react-icons/ri";

import {FaGlobeEurope} from "react-icons/fa";
import {FaGlobeAsia} from "react-icons/fa";

type Props = {
    item :Movie | TvShow | Episode,
    type :string,
    downloadItem?(): void,
    useOriginalTitle?(): void,
    useDefaultTitle?(): void,
    treatAsRegularShow?(): void,
    treatAsAnime?(): void,
    markNotDownloaded?(): void,
    markDownloaded?(): void,
    abortDownload?(): void,
    skipTorrent?(): void,
    restoreItem?(): void,
    deleteItem?(): void,
};

type State = {
    item :Movie | TvShow | Episode | null,
};

class MediaActionBar extends React.Component<Props, State> {
    state :State = {
        item: null,
    };

    componentWillReceiveProps(nextProps :Props) {
        this.setState({ item: nextProps.item });
    }

    getBackLink() {
        if (this.props.type === "movie") {
            return "/movies";
        }

        if (this.props.type === "tvshow") {
            return "/tvshows";
        }

        if (this.props.type === "episode") {
            let item = this.state.item as Episode;
            if (item == null || item.TvShow == null) {
                return "/tvshows/";
            }
            return "/tvshows/" + item.TvShow.ID;
        }

        return "";
    }

    getDownloadButton() {
        if (this.state.item == null) {
            return null;
        }

        let item :any;
        if (this.props.type === "movie") {
            item = this.state.item as Movie;
        } else if (this.props.type === "tvshow") {
            item = this.state.item as TvShow;
        } else if (this.props.type === "episode") {
            item = this.state.item as Episode;
        }

        let buttonsList :any = [];

        if (this.props.type === "movie" || this.props.type === "episode") {
            if (item.DownloadingItem == null) {
                return null;
            }
            if (!item.DeletedAt && !item.DownloadingItem.Pending && !item.DownloadingItem.Downloading && !item.DownloadingItem.Downloaded && !Helpers.dateIsInFuture(item.Date)) {
                buttonsList.push(
                    <button className="uk-button uk-button-link"
                        key="download"
                        onClick={this.props.downloadItem}>
                        <span className="icon is-small">
                            <RiDownloadLine />
                        </span>
                        Download
                    </button>
                );
            }
        }

        return buttonsList;
    }

    getTitleControlButtons() {
        if (this.state.item == null) {
            return null;
        }

        let item :any;
        if (this.props.type === "movie") {
            item = this.state.item as Movie;
        } else if (this.props.type === "tvshow") {
            item = this.state.item as TvShow;
        } else if (this.props.type === "episode") {
            item = this.state.item as Episode;
        }

        let buttonsList :any = [];

        if ((this.props.type !== "episode") && !item.DeletedAt) {
            if (item.Title !== item.OriginalTitle) {
                if (item.UseDefaultTitle) {
                    buttonsList.push(
                        <button className="uk-button uk-button-link"
                            key="useOriginalTitle"
                            onClick={this.props.useOriginalTitle}>
                            <span className="icon is-small">
                                <RiCheckboxMultipleBlankLine />
                            </span>
                            Use original title
                        </button>
                    );
                } else {
                    buttonsList.push(
                        <button className="uk-button uk-button-link"
                            key="useDefaultTitle"
                            onClick={this.props.useDefaultTitle}>
                            <span className="icon is-small">
                                <RiCheckboxMultipleBlankLine />
                            </span>
                            Use default title
                        </button>
                    );
                }
            }
        }

        return buttonsList;
    }

    getDownloadControlButtons() {
        if (this.state.item == null) {
            return null;
        }

        let item :any;
        if (this.props.type === "movie") {
            item = this.state.item as Movie;
        } else if (this.props.type === "tvshow") {
            item = this.state.item as TvShow;
        } else if (this.props.type === "episode") {
            item = this.state.item as Episode;
        }

        if (item.DeletedAt) {
            return null;
        }

        let buttonsList :any = [];

        if (this.props.type === "movie" || this.props.type === "episode") {
            if (item.DownloadingItem == null) {
                return null;
            }
            if (item.DownloadingItem.Downloaded) {
                buttonsList.push(
                    <button className=""
                        key="markNotDownloaded"
                        onClick={this.props.markNotDownloaded}>
                        <span className="icon is-small">
                            <RiEraserLine />
                        </span>
                        Unmark as downloaded
                    </button>
                );
            } else {
                if (!item.DownloadingItem.Pending && !item.DownloadingItem.Downloading) {
                    buttonsList.push(
                        <button className=""
                            key="markDownloaded"
                            onClick={this.props.markDownloaded}>
                            <span className="icon is-small">
                                <RiCheckLine />
                            </span>
                            Mark as downloaded
                        </button>
                    );
                }
            }

            if (item.DownloadingItem.Downloading) {
                buttonsList.push(
                    <button className=""
                        key="skipTorrent"
                        onClick={this.props.skipTorrent}>
                        <span className="icon is-small">
                            <RiSkipForwardLine />
                        </span>
                        Skip torrent
                    </button>
                );
            }

            if (item.DownloadingItem.Pending || item.DownloadingItem.Downloading) {
                buttonsList.push(
                    <button className=""
                        key="abortDownload"
                        onClick={this.props.abortDownload}>
                        <span className="icon is-small">
                            <RiStopCircleLine />
                        </span>
                        Abort download
                    </button>
                );
            }
        }

        return buttonsList;
    }

    getAnimeControlButtons() {
        if (this.state.item == null) {
            return null;
        }

        let item :any;
        if (this.props.type === "movie") {
            item = this.state.item as Movie;
        } else if (this.props.type === "tvshow") {
            item = this.state.item as TvShow;
        } else if (this.props.type === "episode") {
            item = this.state.item as Episode;
        }

        if (item.DeletedAt) {
            return null;
        }

        let buttonsList :any = [];

        if (item.DeletedAt || this.props.type !== "tvshow") {
            return null;
        }

            if (item.IsAnime) {
                buttonsList.push(
                    <button className=""
                        key="treatAsAnime"
                        onClick={this.props.treatAsRegularShow}>
                        <span className="icon is-small">
                            <FaGlobeEurope />
                        </span>
                        Treat as regular show
                    </button>
                );
            } else {
                buttonsList.push(
                    <button className=""
                        key="treatAsAnime"
                        onClick={this.props.treatAsAnime}>
                        <span className="icon is-small">
                            <FaGlobeAsia />
                        </span>
                        Treat as anime
                    </button>
                );
            }

        return buttonsList;
    }

    getTopbarClass() {
        if (this.state.item == null) {
            return "";
        }

        let item :any;
        if (this.props.type === "movie") {
            item = this.state.item as Movie;
        } else if (this.props.type === "tvshow") {
            item = this.state.item as TvShow;
        } else if (this.props.type === "episode") {
            item = this.state.item as Episode;
        }

        if (item.DeletedAt) {
            return "removed-element";
        }

        if (this.props.type === "movie" || this.props.type === "episode") {
            if (item.DownloadingItem == null) {
                return "";
            }
            if (item.DownloadingItem.Downloaded) {
                return "downloaded-element";
            }
        }

        return "";
    }

    render() {
        if (this.state.item == null) {
            return (
                <div className={`action-bar ${this.getTopbarClass()}`}>
                    <div className="container">
                        <div className="columns">
                            <div className="column">
                                <Link to={this.getBackLink()}>
                                <span className="icon is-small">
                                    <RiArrowLeftLine />
                                </span>
                                    Back
                                </Link>
                            </div>
                            <div className="column is-narrow">
                                <button className="">
                                    <i>
                                        Loading...
                                    </i>
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            );
        }

        let item = this.state.item;

        return (
            <div className={`action-bar ${this.getTopbarClass()}`}>
                <div className="container">
                    <div className="columns">
                        <div className="column">
                            <Link to={this.getBackLink()}>
                                <span className="icon is-small">
                                    <RiArrowLeftLine />
                                </span>
                                Back
                            </Link>
                        </div>
                        <div className="column is-narrow">
                            {this.getDownloadButton()}
                            {this.getTitleControlButtons()}
                            {this.getAnimeControlButtons()}

                            {item.DeletedAt && (
                                <button className=""
                                    onClick={this.props.restoreItem}>
                                    <span className="icon is-small">
                                        <RiArrowGoBackLine />
                                    </span>
                                    Restore
                                </button>
                            )}

                            {this.getDownloadControlButtons()}

                            {(!item.DeletedAt && (this.props.type === "movie" || this.props.type === "tvshow")) && (
                                <button className=""
                                    onClick={this.props.deleteItem}>
                                    <span className="icon is-small">
                                        <RiCloseLine />
                                    </span>
                                    Remove
                                </button>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

export default MediaActionBar;
