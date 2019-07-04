import React from "react";

import API from "../../utils/api";
import Const from "../../const";
import ModuleSettings from "../../types/module_settings";
import Config from "../../types/config";

import Loading from "../loading";
import ModuleSettingsComponent from "./module_settings";
import ConfigCheck from "./config_check";

// TODO: fix any
type State = {
    config :Config | null,
    config_checks: any,
};

class Settings extends React.Component {
    settings_refresh_interval :number;
    state :State = {
        config: null,
        config_checks: null,
    };

    constructor(props) {
        super(props);

        this.settings_refresh_interval = 0;

        this.renderModuleList = this.renderModuleList.bind(this);
    }

    componentDidMount() {
        this.load();

        this.settings_refresh_interval = window.setInterval(this.load.bind(this), Const.DATA_REFRESH);
    }

    componentWillUnmount() {
        clearInterval(this.settings_refresh_interval);
    }

    load() {
        this.getConfig();
        this.checkConfig();
    }

    getConfig() {
        API.Config.get().then(response => {
            this.setState({config: response.data});
        }).catch(error => {
            console.log("Get config error: ", error);
        });
    }

    checkConfig() {
        API.Config.check().then(response => {
            this.setState({config_checks: response.data});
        }).catch(error => {
            console.log("Check config error: ", error);
        });
    }

    renderModuleList(title, collection, type) {
        let moduleList :Array<ModuleSettings> = [];

        if (!collection) {
            return (
                <div>
                    <h4>{title}</h4>
                    <ul className="uk-list uk-list-striped no-stripes">
                        <li className="uk-text-muted uk-text-center"><i>No {type}s defined in configuration</i></li>
                    </ul>
                </div>
            );
        }

        Object.keys(collection).forEach(key => {
            moduleList.push({
                Name: key,
                Config: collection[key]
            } as ModuleSettings);
        });

        // TODO: Transform module settings to component
        return (
            <div>
                <h4>{title}</h4>
                <ul className="uk-list uk-list-striped no-stripes">
                    {moduleList.map(module =>
                        <li key={module.Name}>
                            <ModuleSettingsComponent
                                type={type}
                                module={module}/>
                        </li>
                    )}
                </ul>
            </div>
        );
    }

    render() {
        if (this.state.config == null || this.state.config_checks == null) {
            return (
                <Loading/>
            );
        }

        return (
            <>
                {this.state.config_checks && (
                    <div className="uk-container">
                        <div className="uk-alert uk-alert-danger" data-uk-alert>
                            <div className="uk-flex uk-flex-middle">
                                <h3 className="uk-margin-remove">Configuration errors</h3>
                                <button className="uk-alert-close" data-uk-close></button>
                            </div>
                            <hr/>
                            <ul>
                                {this.state.config_checks.map((check, index) => (
                                    <ConfigCheck item={check} key={index}/>
                                ))}
                            </ul>
                        </div>
                    </div>
                )}

                <div className="uk-container">
                    <div className="uk-grid uk-grid-small uk-child-width-1-1 uk-child-width-1-2@m" data-uk-grid>
                        <div className="uk-width-1-1">
                            <h3 className="uk-heading-divider">Global settings</h3>
                        </div>

                        <div>
                            <ul className="uk-list uk-list-striped no-stripes">
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Polling interval
                                    </div>
                                    <div>
                                        <code>
                                            {this.state.config.System.CheckInterval}mn
                                        </code>
                                    </div>
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Torrent download attempts limit
                                    </div>
                                    <div>
                                        <code>
                                            {this.state.config.System.TorrentDownloadAttemptsLimit}
                                        </code>
                                    </div>
                                </li>
                            </ul>
                        </div>

                        <div>
                            <ul className="uk-list uk-list-striped no-stripes">
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Interface enabled
                                    </div>
                                    {this.state.config.Interface.Enabled ? (
                                        <div className="uk-text-success"><i
                                            className="material-icons float-right md-18">check_circle_outline</i></div>
                                    ) : (
                                        <div className="uk-text-danger"><i
                                            className="material-icons float-right md-18">clear</i></div>
                                    )}
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Interface port
                                    </div>
                                    <div>
                                        <code>
                                            {this.state.config.Interface.Port}
                                        </code>
                                    </div>
                                </li>
                            </ul>
                        </div>

                        <div>
                            <ul className="uk-list uk-list-striped no-stripes">
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        TV Shows tracking enabled
                                    </div>
                                    {this.state.config.System.TrackShows ? (
                                        <div className="uk-text-success"><i
                                            className="material-icons float-right md-18">check_circle_outline</i></div>
                                    ) : (
                                        <div className="uk-text-danger"><i
                                            className="material-icons float-right md-18">clear</i></div>
                                    )}
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Movie tracking enabled
                                    </div>
                                    {this.state.config.System.TrackMovies ? (
                                        <div className="uk-text-success"><i
                                            className="material-icons float-right md-18">check_circle_outline</i></div>
                                    ) : (
                                        <div className="uk-text-danger"><i
                                            className="material-icons float-right md-18">clear</i></div>
                                    )}
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Automatic show download
                                    </div>
                                    {this.state.config.System.AutomaticShowDownload ? (
                                        <div className="uk-text-success"><i
                                            className="material-icons float-right md-18">check_circle_outline</i></div>
                                    ) : (
                                        <div className="uk-text-danger"><i
                                            className="material-icons float-right md-18">clear</i></div>
                                    )}
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Automatic movie download
                                    </div>
                                    {this.state.config.System.AutomaticMovieDownload ? (
                                        <div className="uk-text-success"><i
                                            className="material-icons float-right md-18">check_circle_outline</i></div>
                                    ) : (
                                        <div className="uk-text-danger"><i
                                            className="material-icons float-right md-18">clear</i></div>
                                    )}
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Preferred media quality
                                    </div>
                                    {this.state.config.System.PreferredMediaQuality ? (
                                        <div>
                                            <code>
                                                {this.state.config.System.PreferredMediaQuality}
                                            </code>
                                        </div>
                                    ) : (
                                        <div>
                                            None
                                        </div>
                                    )}
                                </li>
                            </ul>
                        </div>

                        <div>
                            <ul className="uk-list uk-list-striped no-stripes">
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        TV Shows library path
                                    </div>
                                    <div className="uk-text-small">
                                        <code>
                                            {this.state.config.Library.ShowPath}
                                        </code>
                                    </div>
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Movies library path
                                    </div>
                                    <div>
                                        <code className="uk-text-small">
                                <span className="uk-text-small">
                                    {this.state.config.Library.MoviePath}
                                </span>
                                        </code>
                                    </div>
                                </li>
                            </ul>
                        </div>

                        <div className="uk-width-1-1">
                            <h3 className="uk-heading-divider">Notifications</h3>
                        </div>
                        <div>
                            <ul className="uk-list uk-list-striped no-stripes">
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Notifications enabled
                                    </div>
                                    {this.state.config.Notifications.Enabled ? (
                                        <div className="uk-text-success"><i
                                            className="material-icons float-right md-18">check_circle_outline</i></div>
                                    ) : (
                                        <div className="uk-text-danger"><i
                                            className="material-icons float-right md-18">clear</i></div>
                                    )}
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Notify on new episode
                                    </div>
                                    {this.state.config.Notifications.NotifyNewEpisode ? (
                                        <div className="uk-text-success"><i
                                            className="material-icons float-right md-18">check_circle_outline</i></div>
                                    ) : (
                                        <div className="uk-text-danger"><i
                                            className="material-icons float-right md-18">clear</i></div>
                                    )}
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Notify on download complete
                                    </div>
                                    {this.state.config.Notifications.NotifyDownloadComplete ? (
                                        <div className="uk-text-success"><i
                                            className="material-icons float-right md-18">check_circle_outline</i></div>
                                    ) : (
                                        <div className="uk-text-danger"><i
                                            className="material-icons float-right md-18">clear</i></div>
                                    )}
                                </li>
                                <li className="uk-flex uk-flex-middle">
                                    <div className="uk-width-expand">
                                        Notify on download failure
                                    </div>
                                    {this.state.config.Notifications.NotifyFailure ? (
                                        <div className="uk-text-success"><i
                                            className="material-icons float-right md-18">check_circle_outline</i></div>
                                    ) : (
                                        <div className="uk-text-danger"><i
                                            className="material-icons float-right md-18">clear</i></div>
                                    )}
                                </li>
                            </ul>
                        </div>

                        <div className="uk-width-1-1">
                            <h3 className="uk-heading-divider">Modules</h3>
                        </div>
                        {this.renderModuleList("Watchlists", this.state.config.Watchlists, "watchlist")}
                        {this.renderModuleList("Providers", this.state.config.Providers, "provider")}
                        {this.renderModuleList("Notifiers", this.state.config.Notifiers, 'notifier')}
                        {this.renderModuleList(<>Indexers <span
                            className="text-muted">(torznab)</span></>, this.state.config.Indexers.torznab, "indexer")}
                        {this.renderModuleList("Downloaders", this.state.config.Downloaders, "downloader")}
                        {this.renderModuleList("Media centers", this.state.config.MediaCenters, "mediacenter")}

                        <div className="uk-width-1-1">
                            <h4 className="uk-heading-divider">About flemzerd</h4>
                            <ul className="uk-list uk-list-striped no-stripes">
                                <li>Version: <code>{this.state.config.Version}</code></li>
                                <li className="uk-text-center uk-padding-small">flemzerd is an Open Source software. For
                                    feature ideas, notifying bugs or contributing, do not hesitate to head to the <a
                                        target="blank" href="https://github.com/macarrie/flemzerd">GitHub repo</a></li>
                            </ul>
                        </div>
                    </div>
                </div>
            </>
        );
    }
}

export default Settings;
