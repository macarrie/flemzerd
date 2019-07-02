import React from "react";

import API from "../utils/api";
import Const from "../const";
import Loading from "./loading";

import ModuleStatus from "./module_status";

class Status extends React.Component {
    status_refresh_interval;
    state = {
        watchlists: null,
        providers: null,
        notifiers: null,
        indexers: null,
        downloaders: null,
        mediacenters: null,
    };

    componentDidMount() {
        this.load();

        this.status_refresh_interval = setInterval(this.load.bind(this), Const.DATA_REFRESH);
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

    render() {
        if (this.state.watchlists == null && this.state.providers == null && this.state.notifiers == null && this.state.indexers == null && this.state.downloaders == null && this.state.mediacenters == null) {
            return (
                <Loading/>
            );
        }

        return (
            <div className="uk-container">
                <div className="uk-grid uk-child-width-1-1 uk-child-width-1-2@s" data-uk-grid="masonry: true">
                    <div>
                        <h3>Watchlists</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.state.watchlists && (
                                this.state.watchlists.map(w => (
                                    <ModuleStatus key={w.Name}
                                                  module={w} />
                                ))
                            )}
                        </ul>
                    </div>
                    <div>
                        <h3>Providers</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.state.providers && (
                                this.state.providers.map(p => (
                                    <ModuleStatus key={p.Name}
                                                  module={p} />
                                ))
                            )}
                        </ul>
                    </div>
                    <div>
                        <h3>Notifiers</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.state.notifiers && (
                                this.state.notifiers.map(n => (
                                    <ModuleStatus key={n.Name}
                                                  module={n} />
                                ))
                            )}
                        </ul>
                    </div>
                    <div>
                        <h3>Indexers</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.state.indexers && (
                                this.state.indexers.map(i => (
                                    <ModuleStatus key={i.Name}
                                                  module={i} />
                                ))
                            )}
                        </ul>
                    </div>
                
                    <div>
                        <h3>Downloaders</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.state.downloaders && (
                                this.state.downloaders.map(d => (
                                    <ModuleStatus key={d.Name}
                                                  module={d} />
                                ))
                            )}
                        </ul>
                    </div>
                    <div>
                        <h3>Media centers</h3>
                        <ul className="uk-list uk-list-striped no-stripes">
                            {this.state.mediacenters && (
                                this.state.mediacenters.map(m => (
                                    <ModuleStatus key={m.Name}
                                                  module={m} />
                                ))
                            )}
                        </ul>
                    </div>
                </div>
            </div>
        );
    }
}

export default Status;
