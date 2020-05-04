import React from "react";
import {Link} from "react-router-dom";

import Helpers from "../../utils/helpers";
import Const from "../../const";
import Notification from "../../types/notification";
import Moment from "react-moment";

import {RiCheckLine} from "react-icons/ri";
import {RiArrowDropDownLine} from "react-icons/ri";
import {RiArrowDropUpLine} from "react-icons/ri";

interface Props {
    item :Notification;
    markRead(id :number) :void;
};

type State = {
    notification: Notification,
    display_detailed_contents: Boolean,
};

class NotificationComponent extends React.Component<Props, State> {
    state: State = {
        notification: {} as Notification,
        display_detailed_contents: false,
    };

    constructor(props: Props) {
        super(props);

        this.state.notification = this.props.item;

        this.getMediaLink = this.getMediaLink.bind(this);
        this.getMediaOverview = this.getMediaOverview.bind(this);
        this.toggleContent = this.toggleContent.bind(this);
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

        return "";
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

        return "";
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

        return "";
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

        if (this.state.notification.Type === Const.NOTIFICATION_DOWNLOAD_FAILURE) {
            return (
                <span>
                    {this.getMediaOverview()}
                    Maybe torrents could not be downloaded, or some essential modules for download were not available.
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

    toggleContent() {
        this.setState({ display_detailed_contents: !this.state.display_detailed_contents });
    }
    getContentToggleIcon() {
        return (
            <button className={"button is-naked is-small"}
                    onClick={this.toggleContent}>
                <span className={"icon"}>
                {this.state.display_detailed_contents ? (
                    <RiArrowDropUpLine />
                ) : (
                    <RiArrowDropDownLine />
                )}
                </span>
            </button>
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

        let read_class = "";
        if (this.state.notification.Read) {
            read_class = "read";
        }

        return (
            <li className={`notification_item ${read_class}`}>
                <div className={"columns is-gapless is-multiline is-vcentered"}>
                    <div className={`column is-narrow ${this.getStatusClass()}`}><div></div></div>
                    <div className="column"
                         onClick={this.toggleContent}>
                        {this.getNotificationTitle()}
                    </div>
                    <div className="column is-narrow">
                        {this.getContentToggleIcon()}
                    </div>
                    <div className="column is-narrow">
                        <span className={"has-text-grey"}>
                            <i>
                                <Moment date={this.state.notification.CreatedAt} format={"ddd DD MMM hh:mm"}/>
                            </i>
                        </span>
                    </div>
                    {!this.state.notification.Read && (
                        <div className="notification_action_button uk-text-success"
                            onClick={() => this.props.markRead(this.state.notification.ID)}>
                            <button className="button is-small is-naked"
                                    data-tooltip="Mark as read">
                                <span className={"icon has-text-success"}>
                                    <RiCheckLine />
                                </span>
                            </button>
                        </div>
                    )}
                    {this.state.display_detailed_contents && (
                        <div className={"column is-full notification-content"}>
                            <div>
                                <hr />
                                {this.getNotificationContent()}
                            </div>
                        </div>
                    )}
                </div>
            </li>
        );
    }
}

export default NotificationComponent;
