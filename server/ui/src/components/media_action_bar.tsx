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
import {RiRefreshLine} from "react-icons/ri";

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
    refreshMetadata?(): void,
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
                    <button className=""
                        key="download"
                        onClick={this.props.downloadItem}>
                        <span className="icon is-small">
                            <RiDownloadLine />
                        </span>
                        <span className={"is-hidden-mobile"}>
                            Download
                        </span>
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
                        <button className=""
                            key="useOriginalTitle"
                            onClick={this.props.useOriginalTitle}>
                            <span className="icon is-small">
                                <RiCheckboxMultipleBlankLine />
                            </span>
                            <span className={"is-hidden-mobile"}>
                                Use original title
                            </span>
                        </button>
                    );
                } else {
                    buttonsList.push(
                        <button className=""
                            key="useDefaultTitle"
                            onClick={this.props.useDefaultTitle}>
                            <span className="icon is-small">
                                <RiCheckboxMultipleBlankLine />
                            </span>
                            <span className={"is-hidden-mobile"}>
                                Use default title
                            </span>
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
                        <span className={"is-hidden-mobile"}>
                            Unmark as downloaded
                        </span>
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
                            <span className={"is-hidden-mobile"}>
                                Mark as downloaded
                            </span>
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
                        <span className={"is-hidden-mobile"}>
                            Skip torrent
                        </span>
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
                        <span className={"is-hidden-mobile"}>
                            Abort download
                        </span>
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
                        <span className={"is-hidden-mobile"}>
                            Treat as regular show
                        </span>
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
                        <span className={"is-hidden-mobile"}>
                            Treat as anime
                        </span>
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
                        <div className="columns is-mobile">
                            <div className="column">
                                <Link to={this.getBackLink()}>
                                <span className="icon is-small">
                                    <RiArrowLeftLine />
                                </span>
                                    <span className={"is-hidden-mobile"}>
                                        Back
                                    </span>
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
                    <div className="columns is-mobile">
                        <div className="column">
                            <Link to={this.getBackLink()}>
                                <span className="icon is-small">
                                    <RiArrowLeftLine />
                                </span>
                                <span className={"is-hidden-mobile"}>
                                    Back
                                </span>
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
                                    <span className={"is-hidden-mobile"}>
                                        Restore
                                    </span>
                                </button>
                            )}

                            {this.getDownloadControlButtons()}

                            {(!item.DeletedAt && (this.props.type === "movie" || this.props.type === "tvshow")) && (
                                <>
                                <button className=""
                                        onClick={this.props.refreshMetadata}>
                                    <span className="icon is-small">
                                        <RiRefreshLine />
                                    </span>
                                    <span className={"is-hidden-mobile"}>
                                        Refresh metadata
                                    </span>
                                </button>

                                <button className=""
                                        onClick={this.props.deleteItem}>
                                    <span className="icon is-small">
                                        <RiCloseLine />
                                    </span>
                                    <span className={"is-hidden-mobile"}>
                                        Remove
                                    </span>
                                </button>
                                </>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

export default MediaActionBar;
