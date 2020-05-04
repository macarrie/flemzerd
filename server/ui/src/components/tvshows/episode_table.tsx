import React from "react";
import {Link} from "react-router-dom";
import Moment from 'react-moment';

import API from "../../utils/api";
import Helpers from "../../utils/helpers";
import Episode from "../../types/episode";

import {RiStopCircleLine} from "react-icons/ri";
import {RiDownloadLine} from "react-icons/ri";
import {RiCheckLine} from "react-icons/ri";
import {RiArrowUpDownLine} from "react-icons/ri";
import {RiEraserLine} from "react-icons/ri";
import {RiAlertLine} from "react-icons/ri";

type Props = {
    list: Episode[] | null,
};

type State = {
    list: Episode[] | null,
};

class EpisodeTable extends React.Component<Props, State> {
    state: State = {
        list: null,
    };

    constructor(props: Props) {
        super(props);

        this.state.list = this.props.list;
    }

    componentWillReceiveProps(nextProps: Props) {
        this.setState({list: nextProps.list});
    }

    render() {
        if (this.state.list == null || this.state.list.length === 0) {
            return (
                <div className="uk-text-muted uk-text-center">
                    <i>
                        No episodes for this season
                    </i>
                </div>
            );
        }

        return (
            <table className="table is-fullwidth is-hoverable">
                <thead>
                <tr>
                    <th>Title</th>
                    <th className="is-hidden-tablet">Release</th>
                    <th className="is-hidden-mobile">Release date</th>

                    <th className="is-hidden-tablet has-text-centered">Status</th>
                    <th className="is-hidden-mobile has-text-centered">Download status</th>

                    <th className="has-text-centered">Actions</th>
                </tr>
                </thead>
                <tbody>
                {this.state.list.map(episode => (
                    <EpisodeTableLine
                        key={episode.ID}
                        item={episode}/>
                ))}
                </tbody>
            </table>
        );
    }
}

type TableLineProps = {
    item: Episode,
};

type TableLineState = {
    item: Episode,
};

class EpisodeTableLine extends React.Component<TableLineProps, TableLineState> {
    state: TableLineState = {
        item: {} as Episode,
    };

    constructor(props: TableLineProps) {
        super(props);

        this.state.item = this.props.item;

        this.refreshEpisode = this.refreshEpisode.bind(this);

        this.downloadEpisode = this.downloadEpisode.bind(this);
        this.abortDownload = this.abortDownload.bind(this);
        this.markNotDownloaded = this.markNotDownloaded.bind(this);
        this.markDownloaded = this.markDownloaded.bind(this);
        this.changeDownloadedState = this.changeDownloadedState.bind(this);

        this.getDownloadButtons = this.getDownloadButtons.bind(this);
        this.getDownloadControlButtons = this.getDownloadControlButtons.bind(this);
        this.getEpisodeNumber = this.getEpisodeNumber.bind(this);
    }

    componentWillReceiveProps(nextProps: TableLineProps) {
        this.setState({item: nextProps.item});
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

    changeDownloadedState(downloaded_state: boolean) {
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
        if (this.state.item.ActionPending) {
            return (
                <button className={"button is-small is-loading is-naked"}></button>
            );
        }

        if (this.state.item.DownloadingItem.Downloaded) {
            return (
                <button className="button is-small is-naked has-text-success is-disabled"
                        data-tooltip={"Downloaded"}>
                    <span className="icon is-small">
                        <RiCheckLine/>
                    </span>
                </button>
            );
        }

        if (!(Helpers.dateIsInFuture(this.state.item.Date) || this.state.item.DownloadingItem.Downloaded || this.state.item.DownloadingItem.Downloading || this.state.item.DownloadingItem.Pending)) {
            return (
                <div>
                    <button className="button is-small is-naked"
                            key="downloadEpisode"
                            onClick={this.downloadEpisode}
                            data-tooltip="Download">
                        <span className="icon is-small">
                            <RiDownloadLine/>
                        </span>
                    </button>
                    {this.state.item.DownloadingItem.TorrentsNotFound && (
                        <button className="button is-small is-naked has-text-warning"
                                data-tooltip="No torrents found on last attempt">
                            <span className="icon is-small">
                                <RiAlertLine/>
                            </span>
                        </button>
                    )}
                </div>
            );
        }

        if (this.state.item.DownloadingItem.Downloading || this.state.item.DownloadingItem.Pending) {
            return (
                <button className={"button is-small is-naked is-disabled"}>
                    <span className="icon is-small">
                        <RiArrowUpDownLine/>
                    </span>
                </button>
        );
        }

        return "";
    }

    getDownloadControlButtons() {
        if (Helpers.dateIsInFuture(this.state.item.Date)) {
            return;
        }

        let buttonList: any[] = [];

        if (this.state.item.DownloadingItem.Downloaded) {
            buttonList.push(
                <button className="button is-small is-naked"
                        key="markNotDownloaded"
                        onClick={this.markNotDownloaded}
                        data-tooltip="Unmark as downloaded">
                        <span className="icon is-small">
                            <RiEraserLine/>
                        </span>
                </button>
            );
        }
        if (this.state.item.DownloadingItem.Downloading || this.state.item.DownloadingItem.Pending) {
            buttonList.push(
                <button className="button is-small is-naked"
                      key="abortDownload"
                      data-tooltip={"Abort download"}
                      onClick={this.abortDownload}>
                        <span className="icon is-small">
                            <RiStopCircleLine/>
                        </span>
                    </button>
            );
        }
        if (!this.state.item.DownloadingItem.Downloaded && !this.state.item.DownloadingItem.Pending && !this.state.item.DownloadingItem.Downloading) {
            buttonList.push(
                <button className="button is-small is-naked"
                        key="markDownloaded"
                        onClick={this.markDownloaded}
                        data-tooltip="Mark as downloaded">
                    <span className="icon is-small">
                        <RiCheckLine/>
                    </span>
                </button>
            );
        }

        return buttonList;
    }

    getEpisodeNumber() {
        if (this.state.item.TvShow.IsAnime) {
            return (
                <>
                    {Helpers.formatAbsoluteNumber(this.state.item.AbsoluteNumber)}
                </>
            );
        }

        return (
            <>
                {Helpers.formatNumber(this.state.item.Season)}x{Helpers.formatNumber(this.state.item.Number)}
            </>
        );
    }

    render() {
        return (
            <tr className={Helpers.dateIsInFuture(this.state.item.Date) ? "episode-list-future" : ""}>
                <td>
                    <span className="tag is-small episode-label">{this.getEpisodeNumber()}</span>
                    <Link to={`/tvshows/${this.state.item.TvShow.ID}/episodes/${this.state.item.ID}`}>
                        {this.state.item.Title}
                    </Link>
                </td>
                <td className="uk-table-shrink is-hidden-tablet uk-text-nowrap uk-text-center uk-hidden@s"><Moment date={this.state.item.Date} format={"DD MMM YYYY"}/></td>
                <td className="uk-table-shrink is-hidden-desktop is-hidden-mobile uk-text-nowrap uk-text-center uk-visible@s uk-hidden@m"><Moment date={this.state.item.Date} format={"ddd DD MMM YYYY"}/></td>
                <td className="uk-table-shrink is-hidden-mobile is-hidden-tablet-only uk-text-center uk-visible@m"><Moment date={this.state.item.Date} format={"dddd DD MMM YYYY"}/></td>
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
