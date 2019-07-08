import React from "react";

import Helpers from "../utils/helpers";
import DownloadingItem from "../types/downloading_item";

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

        if (item.Downloaded) {
            return (
                <div>
                    <div className="uk-text-center uk-margin-small-top uk-margin-small-bottom">
                        <span className="uk-h4 uk-text-success">
                            <i className="material-icons md-18">check</i> Download successful
                        </span>
                    </div>
                </div>
            );
        }

        return "";
    }

    printDownloadingInfo() {
        let item = this.state.item;

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
                                <span
                                    className="uk-text-success">Started at {Helpers.formatDate(item.CurrentTorrent.CreatedAt, 'HH:mm DD/MM/YYYY')} </span>
                            </td>
                        </tr>
                        <tr>
                            <td><b>Current download</b></td>
                            {item.CurrentTorrent.Name ? (
                                <td className="uk-font-italic uk-text-muted">
                                    {item.CurrentTorrent.Name}
                                </td>
                            ) : (
                                <td className="uk-font-italic uk-text-muted">
                                    Unknown
                                </td>
                            )}
                        </tr>
                        <tr>
                            <td><b>Download directory</b></td>
                            {item.CurrentTorrent.DownloadDir ? (
                                <td className="uk-text-muted uk-font-italic">
                                    <code>
                                        {item.CurrentTorrent.DownloadDir}
                                    </code>
                                </td>
                            ) : (
                                <td className="uk-text-muted uk-font-italic">
                                    Unknown
                                </td>
                            )}
                        </tr>
                        <tr className={item.FailedTorrents.length > 0 ? "uk-text-danger" : "uk-text-success"}>
                            <td><b>Failed torrents</b></td>
                            <td>
                                <b>{item.FailedTorrents.length}</b>
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
                                    <i className="material-icons uk-text-warning md-18">warning</i>
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
            console.log("Downloading item: ", this.state.item);
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
