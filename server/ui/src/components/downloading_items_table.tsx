import React from "react";
import {Link} from "react-router-dom";

import Helpers from "../utils/helpers";
import Movie from "../types/movie";
import Episode from "../types/episode";

type Props = {
    list :Movie[] | Episode[],
    type :string,
    abortDownload(id: number): void,
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

        return list.map((item: Movie | Episode) => (
            <tr key={item.ID}>
                <td>
                    <Link to={this.getMediaLink(item)}>
                        {Helpers.getMediaTitle(item)}
                    </Link>
                </td>
                <td className="uk-text-muted">
                    <i>
                        {item.DownloadingItem.CurrentTorrent.Name}
                    </i>
                </td>
                <td className="uk-text-center uk-table-shrink uk-text-nowrap">
                    {(item.DownloadingItem.Downloading) && (
                        <>
                            <progress className="uk-progress uk-margin-small-bottom" value={item.DownloadingItem.CurrentTorrent.PercentDone * 100} max="100">{item.DownloadingItem.CurrentTorrent.PercentDone * 100}%</progress>
                            <span className="uk-text-success">Started at {Helpers.formatDate(item.DownloadingItem.CurrentTorrent.CreatedAt, 'HH:mm DD/MM/YYYY')}</span>
                        </>
                    )}

                    {(item.DownloadingItem.Pending) && (
                        <span className="uk-text-italic">Looking for torrents</span>
                    )}
                </td>
                <td className={`uk-text-center uk-text-bold ${(item.DownloadingItem.FailedTorrents.length > 0) ? "uk-text-danger" : "uk-text-success"}`}>
                    {item.DownloadingItem.FailedTorrents.length}
                </td>
                <td>
                    <button className="uk-icon uk-icon-link uk-text-danger abort-download-button" 
                        data-uk-icon="icon: close; ratio: 0.75"
                        onClick={() => this.props.abortDownload(item.ID)}>
                    </button>
                </td>
            </tr>
        ));
    }

    render() {
        return (
            <div className="uk-width-1-1 uk-overflow-auto">
                <table className="uk-table uk-table-small uk-table-divider uk-table-middle">
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th className="uk-text-nowrap">Current download</th>
                            <th className="uk-text-nowrap">Download status</th>
                            <th className="uk-text-center uk-table-shrink uk-text-nowrap">Failed torrents</th>
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
