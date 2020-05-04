import React from "react";
import {Link} from "react-router-dom";

import Helpers from "../utils/helpers";
import Movie from "../types/movie";
import Episode from "../types/episode";
import Moment from "react-moment";

import {RiSkipForwardLine} from "react-icons/ri";
import {RiStopCircleLine} from "react-icons/ri";

type Props = {
    list :Movie[] | Episode[],
    type :string,
    abortDownload(id: number): void,
    skipTorrent(id: number): void,
};

type State = {
    list :Movie[] | Episode[],
};

class DownloadingItemTable extends React.Component<Props, State> {
    state :State = {
        list: [],
    };

    constructor(props :Props) {
        super(props);

        this.state.list = this.props.list;
    }

    componentWillReceiveProps(nextProps :Props) {
        this.setState({ list: nextProps.list });
    }

    getMediaLink(item :Movie | Episode) :string {
        if (this.props.type === "movie") {
            let typed :Movie = item as Movie;
            return `/movies/${typed.ID}`;
        } else if (this.props.type === "episode") {
            let typed :Episode = item as Episode;
            return `tvshows/${typed.TvShow.ID}/episodes/${typed.ID}`;
        } else {
            return "";
        }
    }

    renderTable() {
        let list :any;
        if (this.props.type === "movie") {
            list = this.state.list as Movie[];
        } else if (this.props.type === "episode") {
            list = this.state.list as Episode[];
        }


        return list.map((item: Movie | Episode) => {
            let currentTorrent = Helpers.getCurrentTorrent(item.DownloadingItem);
            let failedTorrents = Helpers.getFailedTorrents(item.DownloadingItem);

            return (
                <tr key={item.ID}>
                    <td>
                        <Link to={this.getMediaLink(item)}>
                            {Helpers.getMediaTitle(item)}
                        </Link>
                    </td>
                    <td className="has-text-grey-light">
                        <i>
                            {currentTorrent.Name}
                        </i>
                    </td>
                    <td className="is-narrow">
                        {(item.DownloadingItem.Downloading) && (
                            <>
                                <progress className="progress is-info is-small" value={currentTorrent.PercentDone * 100} max="100">{currentTorrent.PercentDone * 100}%</progress>
                                <span className="has-text-success">Started at <Moment date={currentTorrent.CreatedAt} format={"HH:mm DD/MMM/YYYY"}/></span>
                            </>
                        )}

                        {(item.DownloadingItem.Pending) && (
                            <span className="has-text-grey"><i> Looking for torrents </i></span>
                        )}
                    </td>
                    <td className={`has-text-centered ${((failedTorrents || []).length > 0) ? "has-text-danger" : "has-text-success"}`}>
                        <b>
                            {(failedTorrents || []).length}
                        </b>
                    </td>
                    <td>
                        <button className="button is-naked is-small"
                            data-tooltip="Skip current torrent download"
                            onClick={() => this.props.skipTorrent(item.ID)}>
                            <span className={"icon is-small"}>
                                <RiSkipForwardLine />
                            </span>
                        </button>
                        <button className="button is-naked is-small"
                            data-tooltip="Abort download"
                            onClick={() => this.props.abortDownload(item.ID)}>
                            <span className={"icon is-small"}>
                                <RiStopCircleLine />
                            </span>
                        </button>
                    </td>
                </tr>
            );
        });
    }

    render() {
        return (
            <div className="column">
                <table className="table is-fullwidth">
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th className="">Current download</th>
                            <th className="">Download status</th>
                            <th className="has-text-centered">Failed torrents</th>
                            <th></th>
                        </tr>
                    </thead>

                    <tbody>
                        {this.renderTable()}
                    </tbody>
                </table>
            </div>
        );
    }
}

export default DownloadingItemTable;
