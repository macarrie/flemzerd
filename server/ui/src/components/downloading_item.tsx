import React from "react";

import Helpers from "../utils/helpers";
import DownloadingItem from "../types/downloading_item";
import Torrent from "../types/torrent";
import Empty from "./empty";

import {RiCheckLine} from "react-icons/ri";
import {RiDownloadLine} from "react-icons/ri";
import {RiUploadLine} from "react-icons/ri";
import {RiAlertLine} from "react-icons/ri";

import Moment from "react-moment";

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
                <div className={"columns is-mobile is-multiline"}>
                    <div className="column is-full has-text-centered">
                        <div className={"columns is-mobile is-gapless is-centered is-vcentered"}>
                            <div className={"column is-narrow"}>
                                <button className={"button is-naked is-disabled"}>
                                    <span className={"icon has-text-success"}>
                                        <RiCheckLine />
                                    </span>
                                </button>
                            </div>
                            <div className={"column is-narrow"}>
                                <span className="title is-4 has-text-success">
                                    Download successful
                                </span>
                            </div>
                        </div>
                    </div>
                    <div className={"column"}>
                        <table className="table is-fullwidth">
                            <tbody>
                            <tr>
                                <td className={"is-narrow"}><b>Downloaded torrent</b></td>
                                {currentTorrent.Name ? (
                                    <td className="has-text-grey">
                                        <i>
                                            {currentTorrent.Name}
                                        </i>
                                    </td>
                                ) : (
                                    <td className="has-text-grey">
                                        <i>
                                            Unknown
                                        </i>
                                    </td>
                                )}
                            </tr>
                            <tr className={Helpers.count(failedTorrents) > 0 ? "has-text-danger" : "has-text-success"}>
                                <td><b>Failed torrents</b></td>
                                <td>
                                    <b>{Helpers.count(failedTorrents)}</b>
                                </td>
                            </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            );
        }

        return "";
    }

    getTorrentListItem(currentTorrent :Torrent, torrent :Torrent) {
        if (torrent.Failed) {
            return (
                <li className="has-text-danger">
                    <i>
                        {torrent.Name} (failed)
                    </i>
                </li>
            );
        };

        if (torrent.ID === currentTorrent.ID && currentTorrent.ID !== 0) {
            return (
                <li className="has-text-info">
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
                <div className={"columns is-mobile is-vcentered is-multiline"}>
                    <div className="column is-full has-text-centered">
                        <div className={"columns is-mobile is-gapless is-centered is-vcentered"}>
                            <div className={"column is-narrow"}>
                                <button className={"button is-naked is-disabled"}>
                                    <span className={"icon has-text-grey"}>
                                        <RiDownloadLine />
                                    </span>
                                </button>
                            </div>
                            <div className={"column is-narrow"}>
                                <span className="title is-4">
                                    Download in progress
                                </span>
                            </div>
                        </div>
                    </div>
                    <table className="table is-fullwidth">
                        <tbody>
                        <tr>
                            <td><b>Downloading</b></td>
                            <td>
                                <div className="columns is-mobile is-multiline is-gapless">
                                    <div className="column is-full">
                                        <div className="columns is-mobile is-multiline is-vcentered">
                                            <div className="column is-full-mobile has-text-success">Started at <Moment date={currentTorrent.CreatedAt} format={"HH:mm DD/MMM/YYYY"}/> </div>
                                            <div className={"column is-full-mobile"}> ETA: <Moment date={currentTorrent.CreatedAt} format={"HH[h]mm[mn]"}/> left </div>
                                            <div className="column is-half-mobile">
                                                <button className={"button is-small is-naked is-disabled"}>
                                                    <span className={"icon is-small"}>
                                                        <RiDownloadLine/>
                                                    </span>
                                                </button>
                                                {Math.round(currentTorrent.RateDownload / 1024)} ko/s
                                            </div>
                                            <div className="column is-half-mobile">
                                                <button className={"button is-small is-naked is-disabled"}>
                                                    <span className={"icon is-small"}>
                                                        <RiUploadLine/>
                                                    </span>
                                                </button>
                                                {Math.round(currentTorrent.RateUpload / 1024)} ko/s
                                            </div>
                                        </div>
                                    </div>
                                    <div className="column">
                                        <div className="columns is-mobile is-vcentered">
                                            <div className={"column is-narrow"}> {Math.round(currentTorrent.PercentDone * 100)}% done </div>
                                            <div className="column">
                                                <progress className="progress is-info is-small" value={currentTorrent.PercentDone * 100} max="100">{currentTorrent.PercentDone * 100}%</progress>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </td>
                        </tr>
                        <tr>
                            <td><b>Current download</b></td>
                            {currentTorrent.Name ? (
                                <td className="has-text-grey is-italic">
                                    {currentTorrent.Name}
                                </td>
                            ) : (
                                <td className="has-text-grey is-italic">
                                    Unknown
                                </td>
                            )}
                        </tr>
                        <tr>
                            <td><b>Download list</b></td>
                            {item.TorrentList == null ? (
                                <td className="has-text-grey is-italic">
                                    Empty torrent list
                                </td>
                            ) : (
                                <td className="has-text-grey is-italic">
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
                                <td className="has-text-grey is-italic">
                                    <code>
                                        {currentTorrent.DownloadDir}
                                    </code>
                                </td>
                            ) : (
                                <td className="has-text-grey is-italic">
                                    Unknown
                                </td>
                            )}
                        </tr>
                        <tr className={(failedTorrents || []).length > 0 ? "has-text-danger" : "has-text-success"}>
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
                    <div className={"columns is-mobile is-vcentered is-multiline"}>
                        <div className="column is-full has-text-centered">
                            <div className={"columns is-mobile is-gapless is-centered is-vcentered"}>
                                <div className={"column is-narrow"}>
                                    <button className={"button is-naked is-disabled"}>
                                    <span className={"icon has-text-warning"}>
                                        <RiAlertLine />
                                    </span>
                                    </button>
                                </div>
                                <div className={"column is-narrow"}>
                                    <span className="title is-4">
                                        No torrents found
                                    </span>
                                </div>
                            </div>
                            <div>
                                No torrents have been found during last download. <br/>
                                Try again later (in case of recent releases) or after adding new indexers.
                            </div>
                        </div>
                    </div>
                );
            } else {
                return (
                    <Empty label={"Not downloading"} />
                );
            }
        }

        return "";
    }

    printDownloadPending() {
        if (this.state.item.Pending) {
            return (
                <Empty label={"Looking for torrents"}/>
            );
        }

        return "";
    }

    render() {
        if (typeof this.state.item === "undefined") {
            return (
                <div className="columns is-mobile is-multiline">
                    <div className="column is-12">
                        <span className="title is-2 has-text-centered">
                            <i>
                                Loading
                            </i>
                        </span>
                    </div>
                </div>
            );
        }

        return (
            <div className="columns is-mobile is-multiline">
                <div className="column is-12">
                    <span className="title is-4">Download status</span>
                    <hr />
                </div>
                <div className="column is-12">
                    {this.printDownloadedInfo()}
                    {this.printDownloadingInfo()}
                    {this.printNotDownloadingInfo()}
                    {this.printDownloadPending()}
                </div>
            </div>
        );
    }
}

export default DownloadingItemComponent;
