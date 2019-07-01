import React from "react";
import { Link } from "react-router-dom";

import API from "../../utils/api";
import Helpers from "../../utils/helpers";

class EpisodeTable extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            list: null,
        };

        this.state.list = this.props.list;
    }

    componentWillReceiveProps(nextProps) {
        this.setState({ list: nextProps.list });
    }

    render() {
        return (
            <table className="uk-table uk-table-divider uk-table-small"> 
                <thead>
                    <tr>
                        <th>Title</th>
                        <th className="uk-hidden@m">Release</th>
                        <th className="uk-visible@m">Release date</th>

                        <th className="uk-table-shrink uk-text-nowrap uk-text-center uk-hidden@m">Status</th>
                        <th className="uk-table-shrink uk-text-nowrap uk-text-center uk-visible@m">Download status</th>

                        <th className="uk-table-shrink uk-text-nowrap uk-text-center">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {(this.state.list || []).map(episode => (
                        <EpisodeTableLine
                            key={episode.ID}
                            item={episode} />
                    ))}
                </tbody>
            </table>
        );
    }
}

class EpisodeTableLine extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            item: null,
        };

        this.state.item = this.props.item;

        this.refreshEpisode = this.refreshEpisode.bind(this);

        this.downloadEpisode = this.downloadEpisode.bind(this);
        this.abortDownload = this.abortDownload.bind(this);
        this.markNotDownloaded = this.markNotDownloaded.bind(this);
        this.markDownloaded = this.markDownloaded.bind(this);
        this.changeDownloadedState = this.changeDownloadedState.bind(this);

        this.getDownloadButtons = this.getDownloadButtons.bind(this);
        this.getDownloadControlButtons = this.getDownloadControlButtons.bind(this);
    }

    componentWillReceiveProps(nextProps) {
        this.setState({ item: nextProps.item });
    }

    refreshEpisode() {
        API.Episodes.get(this.state.item.ID).then(response => {
            this.setState({item: response.data});
        }).catch(error => {
            console.log("Change episode downloaded state error: ", error);
        });
    }

    downloadEpisode() {
        API.Episodes.download(this.state.item.ID).then(response => {
            this.refreshEpisode();
        }).catch(error => {
            console.log("Download episode error: ", error);
        });
    }

    markNotDownloaded() {
        this.changeDownloadedState(false);
    }

    markDownloaded() {
        this.changeDownloadedState(true);
    }

    changeDownloadedState(downloaded_state) {
        API.Episodes.changeDownloadedState(this.state.item.ID, downloaded_state).then(response => {
            this.refreshEpisode();
        }).catch(error => {
            console.log("Change episode downloaded state error: ", error);
        });
    }

    abortDownload() {
        API.Episodes.abortDownload(this.state.item.ID).then(response => {
            this.refreshEpisode();
        }).catch(error => {
            console.log("Abort episode download error: ", error);
        });
    }

    getDownloadButtons() {
        if (this.state.item.DownloadingItem.Downloaded) {
            return (
                <span className="uk-text-success"
                    data-uk-icon="icon: check; ratio: 0.75">
                </span>
            );
        }

        if (!(Helpers.dateIsInFuture(this.state.item.Date) || this.state.item.DownloadingItem.Downloaded || this.state.item.DownloadingItem.Downloading || this.state.item.DownloadingItem.Pending)) {
            return (
                <div>
                    <button className="uk-icon"
                        key="downloadEpisode"
                        onClick={this.downloadEpisode}
                        data-uk-tooltip="delay: 500; title: Download"
                        data-uk-icon="icon: download; ratio: 0.75">
                    </button>
                    {this.state.item.DownloadingItem.TorrentsNotFound && (
                        <div className="btn btn-sm">
                            <i className="material-icons uk-text-warning md-18"
                                title="No torrents found for episode on last download">warning</i>
                        </div>
                    )}
                </div>
            );
        }

        if (this.state.item.DownloadingItem.Downloading || this.state.item.DownloadingItem.Pending) {
            return (
                <span className="btn btn-sm p-0">
                    <i className="material-icons uk-text-small">swap_vert</i>
                </span>
            );
        }
    }

    getDownloadControlButtons() {
        if (Helpers.dateIsInFuture(this.state.item.Date)) {
            return;
        }

        let buttonList = [];

        if (this.state.item.DownloadingItem.Downloaded) {
            buttonList.push(
                <button className="uk-icon"
                    key="markNotDownloaded"
                    onClick={this.markNotDownloaded}
                    data-uk-tooltip="delay: 500; title: Unmark as downloaded"
                    data-uk-icon="icon: push; ratio: 0.75">
                </button>
            );
        }
        if (this.state.item.DownloadingItem.Downloading || this.state.item.DownloadingItem.Pending) {
            buttonList.push(
                <span className="btn btn-sm p-0"
                    key="abortDownload"
                    onClick={this.abortDownload}>
                    <span className="uk-icon uk-icon-link uk-text-danger" data-uk-icon="icon: close; ratio: 0.75"></span>
                </span>
            );
        }
        if (!this.state.item.DownloadingItem.Downloaded && !this.state.item.DownloadingItem.Pending && !this.state.item.DownloadingItem.Downloading) {
            buttonList.push(
                <button className="uk-icon"
                    key="markDownloaded"
                    onClick={this.markDownloaded}
                    data-uk-tooltip="delay: 500; title: Mark as downloaded"
                    data-uk-icon="icon: check; ratio: 0.75">
                </button>
            );
        }

        return buttonList;
    }

    render() {
        return (
            <tr className={Helpers.dateIsInFuture(this.state.item.Date) ? "episode-list-future" : ""}>
                <td>
                    <span className="uk-text-small uk-border-rounded episode-label">{Helpers.formatNumber(this.state.item.Season)}x{Helpers.formatNumber(this.state.item.Number)}</span>
                    <Link to={`/tvshows/${this.state.item.TvShow.ID}/episodes/${this.state.item.ID}`}>
                        {this.state.item.Title}
                    </Link>
                </td>
                <td className="uk-table-shrink uk-text-nowrap uk-text-center uk-hidden@s">{Helpers.formatDate(this.state.item.Date, 'DD MMM YYYY')}</td>
                <td className="uk-table-shrink uk-text-nowrap uk-text-center uk-visible@s uk-hidden@m">{Helpers.formatDate(this.state.item.Date, 'ddd DD MMM YYYY')}</td>
                <td className="uk-table-shrink uk-text-nowrap uk-text-center uk-visible@m">{Helpers.formatDate(this.state.item.Date, 'dddd DD MMM YYYY')}</td>
                <td className="uk-text-center">
                    {this.getDownloadButtons()}
                </td>
                <td className="uk-text-center">
                    {this.getDownloadControlButtons()}
                </td>
            </tr>
        );
    }
}

export default EpisodeTable;
