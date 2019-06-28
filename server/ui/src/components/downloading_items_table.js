import React from "react";
import {Link} from "react-router-dom";

import Helpers from "../utils/helpers";

class DownloadingItemTable extends React.Component {
    constructor(props) {
        // TODO: Handle type
        super(props);

        this.state = {
            list: [],
        };

        console.log("Downloading item table list: ", this.props.list);
        this.state.list = this.props.list;
    }

    componentWillReceiveProps(nextProps) {
        console.log("Downloading item table list update: ", this.props.list);
        this.setState({ list: nextProps.list });
    }

    renderTable() {
        return this.state.list.map(item => (
            <tr key={item.ID}>
                <td>
                    <Link to={`/movies/${item.ID}`}>
                        {Helpers.getMovieTitle(item)}
                    </Link>
                </td>
                <td className="uk-text-muted">
                    <i>
                        {item.DownloadingItem.CurrentTorrent.Name}
                    </i>
                </td>
                <td className="uk-text-center uk-table-shrink uk-text-nowrap">
                    {(item.DownloadingItem.Downloading) && (
                        <span className="uk-text-success">Started at {Helpers.formatDate(item.DownloadingItem.CurrentTorrent.CreatedAt, 'HH:mm DD/MM/YYYY')}</span>
                    )}

                    {(item.DownloadingItem.Pending) && (
                        <span className="uk-text-italic">Looking for torrents</span>
                    )}
                </td>
                <td className={`uk-text-center uk-text-bold ${(item.DownloadingItem.FailedTorrents.length > 0) ? "uk-text-danger" : "uk-text-success"}`}>
                    {item.DownloadingItem.FailedTorrents.length}
                </td>
                <td>
                    TODO
                    <a className="uk-icon uk-icon-link uk-text-danger abort-download-button" data-uk-icon="icon: close; ratio: 0.75"
                        data-controller="movie"
                        data-action="click->movie#abortDownload"
                        data-movie-id="{{ .ID }}">
                    </a>
                </td>
            </tr>
        ));
    }

    render() {
        return (
            <div className="uk-width-1-1">
                <table className="uk-table uk-table-small uk-table-divider">
                    <thead>
                        <tr>
                            <th>Movie name</th>
                            <th className="uk-text-nowrap">current download</th>
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
