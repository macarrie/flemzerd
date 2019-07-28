import React from "react";

import API from "../../utils/api";
import Const from "../../const";
import Module from "../../types/module";

import ModuleStatus from "./module_status";

type State = {
    watchlists: Module[],
    providers: Module[],
    notifiers: Module[],
    indexers: Module[],
    downloaders: Module[],
    mediacenters: Module[],
};

class Status extends React.Component<any, State> {
    status_refresh_interval :number;
    state :State = {
        watchlists: new Array<Module>(),
        providers: new Array<Module>(),
        notifiers: new Array<Module>(),
        indexers: new Array<Module>(),
        downloaders: new Array<Module>(),
        mediacenters: new Array<Module>(),
    };

    constructor(props: any) {
        super(props);

        this.status_refresh_interval = 0;
    }
    componentDidMount() {
        this.load();

        this.status_refresh_interval = window.setInterval(this.load.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.status_refresh_interval);
    }

    load() {
        this.getWatchlists();
        this.getProviders();
        this.getNotifiers();
        this.getIndexers();
        this.getDownloaders();
        this.getMediacenters();
    }

    getWatchlists() {
        API.Modules.Watchlists.status().then(response => {
            this.setState({watchlists: response.data});
        }).catch(error => {
            console.log("Get watchlists status error: ", error);
        });
    }

    getProviders() {
        API.Modules.Providers.status().then(response => {
            this.setState({providers: response.data});
        }).catch(error => {
            console.log("Get providers status error: ", error);
        });
    }

    getNotifiers() {
        API.Modules.Notifiers.status().then(response => {
            this.setState({notifiers: response.data});
        }).catch(error => {
            console.log("Get notifiers status error: ", error);
        });
    }

    getIndexers() {
        API.Modules.Indexers.status().then(response => {
            this.setState({indexers: response.data});
        }).catch(error => {
            console.log("Get indexers status error: ", error);
        });
    }

    getDownloaders() {
        API.Modules.Downloaders.status().then(response => {
            this.setState({downloaders: response.data});
        }).catch(error => {
            console.log("Get downloaders status error: ", error);
        });
    }

    getMediacenters() {
        API.Modules.Mediacenters.status().then(response => {
            this.setState({mediacenters: response.data});
        }).catch(error => {
            console.log("Get mediacenters status error: ", error);
        });
    }

    renderModuleList(list :Module[]) {
        if (list == null) {
            return (
                <li className="uk-text-center uk-text-muted">
                    <i>
                        None defined in configuration
                    </i>
                </li>
            );
        }

        if (list.length === 0) {
            return (
                <li className="uk-text-center">
                    <span data-uk-spinner="ratio: 0.8"></span>
                </li>
            );
        } else {
            return list.map((i :Module) => (
                <ModuleStatus key={i.Name} module={i} />
            ));
        }
    }

    render() {
        return (
            <div className="uk-container">
                <div className="uk-grid uk-child-width-1-1 uk-child-width-1-2@s" data-uk-grid="masonry: true">
                    <div>
                        <h3>Watchlists</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.renderModuleList(this.state.watchlists)}
                        </ul>
                    </div>
                    <div>
                        <h3>Providers</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.renderModuleList(this.state.providers)}
                        </ul>
                    </div>
                    <div>
                        <h3>Notifiers</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.renderModuleList(this.state.notifiers)}
                        </ul>
                    </div>
                    <div>
                        <h3>Indexers</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.renderModuleList(this.state.indexers)}
                        </ul>
                    </div>

                    <div>
                        <h3>Downloaders</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.renderModuleList(this.state.downloaders)}
                        </ul>
                    </div>
                    <div>
                        <h3>Media centers</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.renderModuleList(this.state.mediacenters)}
                        </ul>
                    </div>
                </div>
            </div>
        );
    }
}

export default Status;
