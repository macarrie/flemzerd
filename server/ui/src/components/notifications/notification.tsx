import React from "react";
import {Link} from "react-router-dom";

import Helpers from "../../utils/helpers";
import Const from "../../const";
import Notification from "../../types/notification";

interface Props {
    item :Notification;
    markRead(id :number) :void;
};

type State = {
    notification: Notification,
};

class NotificationComponent extends React.Component<Props, State> {
    state: State = {
        notification: {} as Notification,
    };

    constructor(props: Props) {
        super(props);

        this.state.notification = this.props.item;

        this.getMediaLink = this.getMediaLink.bind(this);
        this.getMediaOverview = this.getMediaOverview.bind(this);
    }

    componentWillReceiveProps(nextProps: Props) {
        this.setState({ notification: nextProps.item });
    }

    getStatusClass() :string {
        let status :number = 0;

        if (this.state.notification == null) {
            return "";
        }

        switch (this.state.notification.Type) {
            case Const.NOTIFICATION_NEW_EPISODE:
                status = Const.INFO;
                break;
            case Const.NOTIFICATION_NEW_MOVIE:
                status = Const.INFO;
                break;
            case Const.NOTIFICATION_DOWNLOAD_START:
                status = Const.INFO;
                break;
            case Const.NOTIFICATION_DOWNLOAD_SUCCESS:
                status = Const.OK;
                break;
            case Const.NOTIFICATION_DOWNLOAD_FAILURE:
                status = Const.CRITICAL;
                break;
            case Const.NOTIFICATION_TEXT:
                status = Const.INFO;
                break;
            case Const.NOTIFICATION_NO_TORRENT:
                status = Const.WARNING;
                break;
            default:
                status = Const.UNKNOWN;
                break;
        }

        let classes :string = "notif-icon ";

        switch (status) {
            case Const.OK:
                classes += "ok";
                break;
            case Const.WARNING:
                classes += "warning";
                break;
            case Const.CRITICAL:
                classes += "critical";
                break;
            case Const.UNKNOWN:
                classes += "unknown";
                break;
            case Const.INFO:
                classes += "info";
                break;
            default:
                break;
        }

        return classes;
    }

    getMediaLink() {
        if (this.state.notification == null) {
            return;
        }

        if (this.state.notification.Movie.ID) {
            return (
                <Link to={`/movies/${this.state.notification.Movie.ID}`}>{Helpers.getMediaTitle(this.state.notification.Movie)}</Link>
            );
        }
        if (this.state.notification.Episode.ID) {
            return (
                <Link to={`/tvshows/${this.state.notification.Episode.TvShow.ID}/episodes/${this.state.notification.Episode.ID}`}>{Helpers.getMediaTitle(this.state.notification.Episode)}</Link>
            );
        }
    }

    getMediaOverview() {
        if (this.state.notification == null) {
            return;
        }

        if (this.state.notification.Movie.ID) {
            return (
                <blockquote className="uk-text-normal uk-padding-left-small">{this.state.notification.Movie.Overview}</blockquote>
            );
        }
        if (this.state.notification.Episode.ID) {
            return (
                <blockquote className="uk-text-normal uk-padding-left-small">{this.state.notification.Episode.Overview}</blockquote>
            );
        }
    }

    getNotificationTitle() {
        if (this.state.notification == null) {
            return;
        }

        if (this.state.notification.Type === Const.NOTIFICATION_NEW_EPISODE) {
            return (
                <span>
                    New episode aired: {this.getMediaLink()}
                </span>
            );
        }

        if (this.state.notification.Type === Const.NOTIFICATION_NEW_MOVIE) {
            return (
                <span>
                    New movie: {this.getMediaLink()}
                </span>
            );
        }

        if (this.state.notification.Type === Const.NOTIFICATION_DOWNLOAD_START) {
            return (
                <span>
                    Download started: {this.getMediaLink()}
                </span>
            );
        }

        if (this.state.notification.Type === Const.NOTIFICATION_DOWNLOAD_SUCCESS) {
            return (
                <span>
                    Download success: {this.getMediaLink()}
                </span>
            );
        }

        if (this.state.notification.Type === Const.NOTIFICATION_DOWNLOAD_FAILURE) {
            return (
                <span>
                    Download failure: {this.getMediaLink()}
                </span>
            );
        }

        if (this.state.notification.Type === Const.NOTIFICATION_TEXT) {
            return (
                <span>
                    {this.state.notification.Title}
                </span>
            );
        }

        if (this.state.notification.Type === Const.NOTIFICATION_NO_TORRENT) {
            return (
                <span>
                    No torrents found: {this.getMediaLink()}
                </span>
            );
        }
    }

    getNotificationContent() {
        if (this.state.notification == null) {
            return;
        }

        if (this.state.notification.Type === Const.NOTIFICATION_TEXT) {
            return (
                <span>
                    {this.state.notification.Content}
                </span>
            );
        }

        if (this.state.notification.Type === Const.NOTIFICATION_NO_TORRENT) {
            return (
                <span>
                    No torrents have been found for media. Try adding new indexers or wait for torrents to be available (for recently released movies or episodes)
                </span>
            );
        }

        return (
            <span>
                {this.getMediaOverview()}
            </span>
        );
    }

    render() {
        if (this.props.markRead == null) {
            console.log("markRead props is null");
            return;
        }

        if (this.state.notification == null) {
            return (
                <li>Loading</li>
            );
        }

        return (
            <li className={this.state.notification.Read ? "read" : ""}>
                <table
                    className={`uk-table-middle notification-table ${this.state.notification.Read ? "notification_read" : ""}`}>
                    <tbody>
                        <tr data-uk-toggle={`target: .notification-detailed-content-${this.state.notification.ID}`}>
                            <td className={`uk-width-auto uk-padding-medium-right notif-icon ${this.getStatusClass()}`}><div></div></td>
                            <td className="uk-width-expand">
                                {this.getNotificationTitle()}
                            </td>
                            <td className="uk-table-shrink uk-text-nowrap uk-text-muted uk-text-right">
                                <small>
                                    {Helpers.formatDate(this.state.notification.CreatedAt, 'ddd DD MMM hh:mm')}
                                </small>
                            </td>
                            {!this.state.notification.Read && (
                                <td className="notification_action_button uk-text-success"
                                    onClick={() => this.props.markRead(this.state.notification.ID)}>
                                    <span data-uk-icon="icon: check"></span>
                                </td>
                            )}
                        </tr>
                        <tr className={`notification-detailed-content-${this.state.notification.ID}`} hidden>
                            <td></td>
                            <td colSpan={2}>
                                {this.getNotificationContent()}
                            </td>
                            <td></td>
                        </tr>
                    </tbody>
                </table>
            </li>
        );
    }
}

export default NotificationComponent;
