import React from "react";
import { Link } from "react-router-dom";

import Helpers from "../utils/helpers";

class MediaActionBar extends React.Component {
    getBackLink() {
        if (this.props.type === "movie") {
            return "/movies";
        }

        if (this.props.type === "tvshow") {
            return "/tvshows";
        }

        if (this.props.type === "episode") {
            return "/tvshows/" + this.props.item.TvShow.ID;
        }

        return "";
    }

    getDownloadButton() {
        let item = this.props.item;
        let buttonsList = [];

        if (this.props.type === "movie" || this.props.type === "episode") {
            if (!item.DeletedAt && !item.DownloadingItem.Pending && !item.DownloadingItem.Downloading && !item.DownloadingItem.Downloaded && !Helpers.dateIsInFuture(item.Date)) {
                buttonsList.push(
                    <button className="uk-button uk-button-link"
                        key="download"
                        onClick={this.props.downloadItem}>
                        <span className="uk-icon uk-icon-link" data-uk-icon="icon: download; ratio: 0.75"></span>
                        Download
                    </button>
                );
            }
        }

        return buttonsList;
    }

    getTitleControlButtons() {
        let item = this.props.item;
        let buttonsList = [];

        if ((this.props.type !== "episode") && !item.DeletedAt) {
            if (item.Title !== item.OriginalTitle) {
                if (item.UseDefaultTitle) {
                    buttonsList.push(
                        <button className="uk-button uk-button-link"
                            key="useOriginalTitle"
                            onClick={this.props.useOriginalTitle}>
                            <span className="uk-icon uk-icon-link" data-uk-icon="icon: move; ratio: 0.75"></span>
                            Use original title
                        </button>
                    );
                } else {
                    buttonsList.push(
                        <button className="uk-button uk-button-link"
                            key="useDefaultTitle"
                            onClick={this.props.useDefaultTitle}>
                            <span className="uk-icon uk-icon-link" data-uk-icon="icon: move; ratio: 0.75"></span>
                            Use default title
                        </button>
                    );
                }
            }
        }

        return buttonsList;
    }

    getDownloadControlButtons() {
        let item = this.props.item;
        let buttonsList = [];

        if (item.DeletedAt) {
            return;
        }

        if (this.props.type === "movie" || this.props.type === "episode") {
            if (item.DownloadingItem.Downloaded) {
                buttonsList.push(
                    <button className="uk-button uk-button-link"
                        key="markNotDownloaded"
                        onClick={this.props.markNotDownloaded}>
                        <span className="uk-icon uk-icon-link" data-uk-icon="icon: push; ratio: 0.75"></span>
                        Unmark as downloaded
                    </button>
                );
            } else {
                if (!item.DownloadingItem.Pending && !item.DownloadingItem.Downloading) {
                    buttonsList.push(
                        <button className="uk-button uk-button-link"
                            key="markDownloaded"
                            onClick={this.props.markDownloaded}>
                            <span className="uk-icon uk-icon-link" data-uk-icon="icon: check; ratio: 0.75"></span>
                            Mark as downloaded
                        </button>
                    );
                }
            }

            if (item.DownloadingItem.Pending || item.DownloadingItem.Downloading) {
                buttonsList.push(
                    <button className="uk-button uk-button-link"
                        key="abortDownload"
                        onClick={this.props.abortDownload}>
                        <span className="uk-icon uk-icon-link" data-uk-icon="icon: close; ratio: 0.75"></span>
                        Abort download
                    </button>
                );
            }
        }

        return buttonsList;
    }

    render() {
        let item = this.props.item;

        //TODO: add "treat as anime" button
        return (
            <div className={`uk-light action-bar ${ item.DeletedAt ? "removed-element" : "" }`}>
                <div className="uk-container">
                    <div className="uk-grid uk-grid-collapse" data-uk-grid>
                        <div className="uk-width-expand">
                            <Link to={this.getBackLink()}>
                                <span className="uk-icon uk-icon-link" data-uk-icon="icon: arrow-left; ratio: 0.75"></span>
                                Back
                            </Link>
                        </div>
                        <div className="uk-width-auto">
                            {this.getDownloadButton()}
                            {this.getTitleControlButtons()}

                            {item.DeletedAt && (
                                <button className="uk-button uk-button-link"
                                    onClick={this.props.restoreItem}>
                                    <span className="uk-icon uk-icon-link" data-uk-icon="icon: reply; ratio: 0.75"></span>
                                    Restore
                                </button>
                            )}

                            {this.getDownloadControlButtons()}

                            {(!item.DeletedAt && (this.props.type === "movie" || this.props.type === "tvshow")) && (
                                <button className="uk-button uk-button-link"
                                    onClick={this.props.deleteItem}>
                                    <span className="uk-icon uk-icon-link" data-uk-icon="icon: close; ratio: 0.75"></span>
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
