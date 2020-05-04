import React from "react";

import API from "../../utils/api";
import Const from "../../const";
import Notification from "../../types/notification";

import Empty from "../empty";
import NotificationComponent from "./notification";

import {RiCheckDoubleLine} from "react-icons/ri";
import {RiDeleteBin5Line} from "react-icons/ri";

type State = {
    list: Notification[] | null,
};

class Notifications extends React.Component<any, State> {
    notifications_refresh_interval: number;
    state: State = {
        list: null,
    };

    constructor(props: any) {
        super(props);

        this.notifications_refresh_interval = 0;

        this.getNotifications = this.getNotifications.bind(this);
        this.markRead = this.markRead.bind(this);
        this.markAllRead = this.markAllRead.bind(this);
        this.deleteAll = this.deleteAll.bind(this);
    }

    componentDidMount() {
        this.getNotifications();

        this.notifications_refresh_interval = window.setInterval(this.getNotifications.bind(this), Const.NOTIFICATIONS_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.notifications_refresh_interval);
    }

    getNotifications() {
        API.Notifications.all().then(response => {
            this.setState({list: response.data});
        }).catch(error => {
            console.log("Get notifications error: ", error);
        });
    }

    markRead(id: number) {
        API.Notifications.markRead(id).then(response => {
            this.getNotifications();
        }).catch(error => {
            console.log("Mark notification read error: ", error);
        });
    }

    markAllRead(event :React.MouseEvent) {
        API.Notifications.markAllRead().then(response => {
            this.getNotifications();
        }).catch(error => {
            console.log("Mark all notifications read error: ", error);
        });
    }

    deleteAll(event :React.MouseEvent) {
        API.Notifications.deleteAll().then(response => {
            this.getNotifications();
        }).catch(error => {
            console.log("Delete all notifications error: ", error);
        });
    }

    render() {
        if (this.state.list == null) {
            return (
                <div className={"container"}>
                    <Empty label={"Loading"} />
                </div>
            );
        }

        return (
            <div className="container">
                <div className="columns">
                    <div className="column">
                        <span className="title is-3">Notifications</span>
                    </div>
                    <div className={"column is-narrow buttons has-addons"}>
                        <button className="button is-naked is-underlined"
                            onClick={this.markAllRead}>
                            <span className={"icon icon-left is-small"}>
                                <RiCheckDoubleLine />
                            </span>
                            Mark all as read
                        </button>
                        <button className="button is-naked is-underlined"
                            onClick={this.deleteAll}>
                            <span className={"icon icon-left is-small"}>
                                <RiDeleteBin5Line />
                            </span>
                            Delete all
                        </button>
                    </div>
                </div>
                <hr />

                <div className="container">
                    <ul className="block-list">
                    {this.state.list.length === 0 ? 
                        (
                            <Empty label={"No notifications"} />
                        ) :
                        this.state.list.map(notification => (
                            <NotificationComponent
                            key={notification.ID}
                            markRead={this.markRead}
                            item={notification} />
                        ))
                    }
                    </ul>
                </div>

            </div>
        );
    }
}

export default Notifications;
