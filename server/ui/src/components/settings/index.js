import React from "react";

import API from "../../utils/api";
import ModuleSettings from "./module_settings";

class Settings extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            config: null,
        };

        this.renderModuleList = this.renderModuleList.bind(this);
    }

    componentDidMount() {
        this.getConfig();
    }

    getConfig() {
        API.Config.get().then(response => {
            console.log("Config: ", response.data);
            this.setState({config: response.data});
        }).catch(error => {
            console.log("Get config error: ", error);
        });
    }

    renderModuleSettings(module, type) {
    }

    renderModuleList(title, collection, type) {
        //TODO: Handle trakt auth
        //TODO: Handle telegram auth
        let moduleList = [];

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
            });
        });

        // TODO: Transform module settings to component
        return (
            <div>
                <h4>{title}</h4>
                <ul className="uk-list uk-list-striped no-stripes">
                    {moduleList.map(module =>
                        <li key={module.Name}>
                            <ModuleSettings
                                type={type}
                                module={module}/>
                        </li>
                    )}
                </ul>
            </div>
        );
    }

    render() {
        if (this.state.config == null) {
            return <div>Loading</div>;
        }

        return (
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
        );
    }
}

export default Settings;
