import React from "react";

import API from "../../utils/api";

import Notification from "./notification";

class Notifications extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            list: null,
        };

        this.getNotifications = this.getNotifications.bind(this);
        this.markRead = this.markRead.bind(this);
        this.markAllRead = this.markAllRead.bind(this);
        this.deleteAll = this.deleteAll.bind(this);
    }

    componentDidMount() {
        this.getNotifications();
    }

    getNotifications() {
        console.log("Get notifications");
        API.Notifications.all().then(response => {
            console.log("Notifications: ", response.data);
            this.setState({list: response.data});
        }).catch(error => {
            console.log("Get notifications error: ", error);
        });
    }

    markRead(id) {
        console.log("Mark notification read ID: ", id);
        API.Notifications.markRead(id).then(response => {
            this.getNotifications();
        }).catch(error => {
            console.log("Mark notification read error: ", error);
        });
    }

    markAllRead() {
        console.log("Mark all notifications read");
        API.Notifications.markAllRead().then(response => {
            this.getNotifications();
        }).catch(error => {
            console.log("Mark all notifications read error: ", error);
        });
    }

    deleteAll(id) {
        console.log("Delete all notifications");
        API.Notifications.deleteAll().then(response => {
            this.getNotifications();
        }).catch(error => {
            console.log("Delete all notifications error: ", error);
        });
    }

    render() {
        if (this.state.list == null) {
            return (
                <div>Loading</div>
            );
        }

        return (
            <div className="uk-container" data-uk-filter="target: .item-filter">
                <div className="uk-grid" data-uk-grid>
                    <div className="uk-width-expand">
                        <span className="uk-h3">Notifications</span>
                    </div>
                    <div>
                        <ul className="uk-subnav uk-subnav-pill">
                            <li data-uk-filter-control="" className="uk-active">
                                <button className="uk-button uk-button-text">
                                    All
                                </button>
                            </li>
                            <li data-uk-filter-control="filter: .read">
                                <button className="uk-button uk-button-text">
                                    Read
                                </button>
                            </li>
                            <li>
                                <button 
                                    className="uk-button uk-button-small uk-button-text"
                                    onClick={this.markAllRead}>
                                    <span className="uk-icon"
                                        data-uk-tooltip="delay: 500; title: Mark all notifications as read"
                                        data-uk-icon="icon: check; ratio: 0.75"></span>
                                    Mark all as read
                                </button>
                            </li>
                            <li>
                                <button 
                                    className="uk-button uk-button-small uk-button-text"
                                    onClick={this.deleteAll}>
                                    <span className="uk-icon"
                                        data-uk-tooltip="delay: 500; title: Delete all notifications"
                                        data-uk-icon="icon: trash; ratio: 0.75"></span>
                                    Delete all
                                </button>
                            </li>
                        </ul>
                    </div>
                </div>
                <hr />

                <div className="uk-grid uk-child-width-1-1">
                    <ul className="uk-list uk-list-striped no-stripes item-filter">
                        {this.state.list.map(notification => (
                            <Notification 
                                key={notification.ID}
                                markRead={this.markRead}
                                item={notification} />
                        ))}
                    </ul>
                </div>

            </div>
        );
    }
}

export default Notifications;
