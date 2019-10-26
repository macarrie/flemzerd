import React from "react";

import Helpers from "../utils/helpers";
import DownloadingItem from "../types/downloading_item";
import Torrent from "../types/torrent";

type Props = {
    item: DownloadingItem,
};

type State = {
    item: DownloadingItem,
};

class DownloadingItemComponent extends React.Component<Props, State> {
    state: State = {
        item: {} as DownloadingItem,
    };

    constructor(props :Props) {
        super(props);

        this.state.item = this.props.item;
    }

    componentWillReceiveProps(nextProps :Props) {
        this.setState({item: nextProps.item});
    }

    printDownloadedInfo() {
        let item = this.state.item;
        let currentTorrent = Helpers.getCurrentTorrent(item);
        let failedTorrents = Helpers.getFailedTorrents(item);

        if (item.Downloaded) {
            return (
                <div>
                <div className="uk-text-center uk-margin-small-top uk-margin-small-bottom">
                <span className="uk-h4 uk-text-success">
                <i className="material-icons md-18">check</i> Download successful
                </span>
                </div>
                <table className="uk-table uk-table-small uk-table-divider">
                <tbody>
                <tr>
                <td><b>Downloaded torrent</b></td>
                {currentTorrent.Name ? (
                    <td className="uk-font-italic uk-text-muted">
                    {currentTorrent.Name}
                    </td>
                ) : (
                    <td className="uk-font-italic uk-text-muted">
                    Unknown
                    </td>
                )}
                </tr>
                <tr className={(failedTorrents || []).length > 0 ? "uk-text-danger" : "uk-text-success"}>
                <td><b>Failed torrents</b></td>
                <td>
                <b>{(failedTorrents || []).length}</b>
                </td>
                </tr>
                </tbody>
                </table>
                </div>
            );
        }

        return "";
    }

    getTorrentListItem(currentTorrent :Torrent, torrent :Torrent) {
        if (torrent.Failed) {
            return (
                <li className="uk-text-danger">
                    <i>
                        {torrent.Name} (failed)
                    </i>
                </li>
            );
        };

        if (torrent.ID === currentTorrent.ID && currentTorrent.ID !== 0) {
            return (
                <li className="uk-text-primary">
                    {torrent.Name} <i>(current)</i>
                </li>
            );
        };

        return (
            <li>
            {torrent.Name}
            </li>
        );
    }

    printDownloadingInfo() {
        let item = this.state.item;
        let currentTorrent = Helpers.getCurrentTorrent(item)
        let failedTorrents = Helpers.getFailedTorrents(item);

        if (item.Downloading && !item.Downloaded) {
            return (
                <div>
                    <div className="uk-text-center uk-margin-small-top uk-margin-small-bottom">
                        <span className="uk-h4">
                            <i className="material-icons md-18">swap_vert</i> Download in progress
                        </span>
                    </div>
                    <table className="uk-table uk-table-small uk-table-divider">
                        <tbody>
                        <tr>
                            <td><b>Downloading</b></td>
                            <td>
                                <div className="uk-grid-collapse" data-uk-grid="">
                                    <div className="uk-width-1-1">
                                        <div className="uk-child-1-4" data-uk-grid="">
                                            <div className="uk-text-success">Started at {Helpers.formatDate(currentTorrent.CreatedAt, 'HH:mm DD/MM/YYYY')} </div>
                                            <div> ETA: {Helpers.formatDate(currentTorrent.ETA, 'HH[h]mm[mn]')} left </div>
                                            <div className="uk-flex uk-flex-middle">
                                                <span className="uk-icon-link" data-uk-icon="download"></span>
                                                {Math.round(currentTorrent.RateDownload / 1024)} ko/s
                                            </div>
                                            <div className="uk-flex uk-flex-middle">
                                                <span className="uk-icon-link" data-uk-icon="upload"></span>
                                                {Math.round(currentTorrent.RateUpload / 1024)} ko/s
                                            </div>
                                        </div>
                                    </div>
                                    <div className="uk-width-1-1">
                                        <div className="uk-flex uk-flex-middle" data-uk-grid="">
                                            <div> {Math.round(currentTorrent.PercentDone * 100)}% done </div>
                                            <div className="uk-width-expand">
                                                <progress className="uk-progress" value={currentTorrent.PercentDone * 100} max="100">{currentTorrent.PercentDone * 100}%</progress>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </td>
                        </tr>
                        <tr>
                            <td><b>Current download</b></td>
                            {currentTorrent.Name ? (
                                <td className="uk-font-italic uk-text-muted">
                                    {currentTorrent.Name}
                                </td>
                            ) : (
                                <td className="uk-font-italic uk-text-muted">
                                    Unknown
                                </td>
                            )}
                        </tr>
                        <tr>
                            <td><b>Download list</b></td>
                            {item.TorrentList == null ? (
                                <td className="uk-font-italic uk-text-muted">
                                    Empty torrent list
                                </td>
                            ) : (
                                <td className="uk-font-italic uk-text-muted">
                                    <ul>
                                        {item.TorrentList.map(torrent => {
                                            return this.getTorrentListItem(currentTorrent, torrent)
                                        })}
                                    </ul>
                                </td>
                            )}
                        </tr>
                        <tr>
                            <td><b>Download directory</b></td>
                            {currentTorrent.DownloadDir ? (
                                <td className="uk-text-muted uk-font-italic">
                                    <code>
                                        {currentTorrent.DownloadDir}
                                    </code>
                                </td>
                            ) : (
                                <td className="uk-text-muted uk-font-italic">
                                    Unknown
                                </td>
                            )}
                        </tr>
                        <tr className={(failedTorrents || []).length > 0 ? "uk-text-danger" : "uk-text-success"}>
                            <td><b>Failed torrents</b></td>
                            <td>
                                <b>{(failedTorrents || []).length}</b>
                            </td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            );
        }

        return "";
    }

    printNotDownloadingInfo() {
        let item = this.state.item;

        if (!item.Downloading && !item.Downloaded && !item.Pending) {
            if (item.TorrentsNotFound) {
                return (
                    <div>
                        <div className="uk-text-center uk-margin-small-top uk-margin-small-bottom">
                            <span className="uk-text-muted">
                                <div className="uk-h4 uk-text-bold">
                                    <i className="material-icons uk-text-warning md-18 uk-margin-small-right">warning</i>
                                    No torrents found
                                </div>
                                No torrents have been found during last download. <br/>
                                Try again later (in case of recent releases) or after adding new indexers.
                            </span>
                        </div>
                    </div>
                );
            } else {
                return (
                    <div>
                        <div className="uk-text-center uk-margin-small-top uk-margin-small-bottom">
                            <span className="uk-h4 uk-text-muted">
                                Not downloading
                            </span>
                        </div>
                    </div>
                );
            }
        }

        return "";
    }

    printDownloadPending() {
        if (this.state.item.Pending) {
            return (
                <div>
                    <div className="uk-text-center uk-margin-small-top uk-margin-small-bottom">
                        <span className="uk-h4 uk-text-muted">
                            Looking for torrents
                        </span>
                    </div>
                </div>
            );
        }

        return "";
    }

    render() {
        if (typeof this.state.item === "undefined") {
            return (
                <div className="uk-background-default uk-border-rounded uk-padding-small">
                    <div className="uk-flex uk-flex-middle uk-flex-center uk-margin-large-top">
                        <span className="uk-margin-right" data-uk-spinner="ratio: 1"></span>
                        <div>
                                <span className="uk-h4 uk-text-muted">
                                    <i>
                                        Loading
                                    </i>
                                </span>
                        </div>
                    </div>
                </div>
            );
        }

        return (
            <div className="uk-background-default uk-border-rounded uk-padding-small">
                <div className="col-12 my-3">
                    <span className="uk-h4">Download status</span>
                </div>
                {this.printDownloadedInfo()}
                {this.printDownloadingInfo()}
                {this.printNotDownloadingInfo()}
                {this.printDownloadPending()}
            </div>
        );
    }
}

export default DownloadingItemComponent;
