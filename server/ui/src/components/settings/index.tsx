import React from "react";

import API from "../../utils/api";
import Auth from "../../auth";
import Const from "../../const";
import ModuleSettings from "../../types/module_settings";
import Config from "../../types/config";
import ConfigCheck from "../../types/config_check";

import Empty from "../empty";
import ModuleSettingsComponent from "./module_settings";
import ConfigCheckComponent from "./config_check";
import LibraryScan from "./scan_library";
import SettingsBoolean from "./settings_boolean";
import SettingsValue from "./settings_value";

type State = {
    config :Config | null,
    config_checks: ConfigCheck[] | null,
};

class Settings extends React.Component {
    settings_refresh_interval :number;
    state :State = {
        config: null,
        config_checks: null,
    };

    constructor(props: any) {
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
            if (Object.keys(response.data).length === 0) {
                this.setState({config_checks: []});
                return;
            }

            this.setState({config_checks: response.data});
        }).catch(error => {
            console.log("Check config error: ", error);
        });
    }

    renderModuleList(title: string | JSX.Element, collection: Array<ModuleSettings>, type: string) {
        let moduleList :Array<ModuleSettings> = [];

        if (!collection) {
            return (
                <div className={"column is-half"}>
                    <h4 className={"title is-4"}>{title}</h4>
                    <ul className="block-list">
                        <li className="has-text-grey-light has-text-centered"><i>No {type}s defined in configuration</i></li>
                    </ul>
                </div>
            );
        }

        Object.keys(collection).forEach((key: any) => {
            moduleList.push({
                Name: key,
                Config: collection[key]
            } as ModuleSettings);
        });

        return (
            <div className={"column is-half"}>
                <h4 className={"title is-4"}>{title}</h4>
                <ul className="block-list">
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
                <div className={"container"}>
                    <Empty label={"Loading"}/>
                </div>
            );
        }

        return (
            <>
                {this.state.config_checks.length > 0 && (
                    <div className="container">
                        <div className={"message is-danger configuration-errors"}>
                            <div className="message-header">
                                <p>Configuration errors</p>
                            </div>
                            <div className={"message-body"}>
                                <ul className={""}>
                                    {this.state.config_checks.map((check, index) => (
                                        <ConfigCheckComponent item={check} key={index}/>
                                    ))}
                                </ul>
                            </div>
                        </div>
                    </div>
                )}

                <div className="container">
                    <div className="columns is-multiline">
                        <div className="column is-full">
                            <h3 className="title is-3">Global settings</h3>
                        </div>

                        <div className={"column is-half"}>
                            <ul className="block-list">
                                <SettingsValue param={this.state.config.System.CheckInterval}
                                               label={"Polling interval"}
                                               unit={"mn"} />
                                <SettingsValue param={this.state.config.System.TorrentDownloadAttemptsLimit}
                                               label={"Torrent download attempts limit"} />
                            </ul>
                        </div>

                        <div className={"column is-half"}>
                            <ul className="block-list">
                                <SettingsBoolean param={this.state.config.Interface.Enabled}
                                    label={"Interface enabled"} />
                                <SettingsValue param={this.state.config.Interface.Port}
k                                              label={"Interface port"} />
                            </ul>
                        </div>

                        <div className={"column is-half"}>
                            <ul className="block-list">
                                <SettingsBoolean param={this.state.config.System.TrackShows}
                                                 label={"TV Shows tracking enabled"} />
                                <SettingsBoolean param={this.state.config.System.TrackMovies}
                                                 label={"Movie tracking enabled"} />
                                <SettingsBoolean param={this.state.config.System.AutomaticShowDownload}
                                                 label={"Automatic show download"} />
                                <SettingsBoolean param={this.state.config.System.AutomaticMovieDownload}
                                                 label={"Automatic movie download"} />
                                <SettingsValue param={this.state.config.System.PreferredMediaQuality ? this.state.config.System.PreferredMediaQuality : "None"}
                                               label={"Preferred media quality"} />
                            </ul>
                        </div>

                        <div className={"column is-half"}>
                            <ul className="block-list">
                                <SettingsValue param={this.state.config.Library.ShowPath}
                                               label={"TV Shows library path"} />
                                <SettingsValue param={this.state.config.Library.MoviePath}
                                               label={"Movie library path"} />
                            </ul>
                        </div>

                        <div className="column is-full">
                            <h3 className="title is-3">Notifications</h3>
                        </div>
                        <div className={"column is-half"}>
                            <ul className="block-list">
                                <SettingsBoolean param={this.state.config.Notifications.Enabled}
                                                 label={"Notifications enabled"} />
                                <SettingsBoolean param={this.state.config.Notifications.NotifyNewEpisode}
                                                 label={"Notify on new episode"} />
                                <SettingsBoolean param={this.state.config.Notifications.NotifyNewMovie}
                                                 label={"Notify on new movie"} />
                                <SettingsBoolean param={this.state.config.Notifications.NotifyDownloadComplete}
                                                 label={"Notify on on download complete"} />
                                <SettingsBoolean param={this.state.config.Notifications.NotifyFailure}
                                                 label={"Notify on on download failure"} />
                            </ul>
                        </div>

                        <div className="column is-full">
                            <h3 className="title is-3">Modules</h3>
                        </div>
                        {this.renderModuleList("Watchlists", this.state.config.Watchlists, "watchlist")}
                        {this.renderModuleList("Providers", this.state.config.Providers, "provider")}
                        {this.renderModuleList("Notifiers", this.state.config.Notifiers, 'notifier')}
                        {this.renderModuleList(<>Indexers <span
                            className="text-muted">(torznab)</span></>, this.state.config.Indexers.torznab, "indexer")}
                        {this.renderModuleList("Downloaders", this.state.config.Downloaders, "downloader")}
                        {this.renderModuleList("Media centers", this.state.config.MediaCenters, "mediacenter")}

                        <div className="column is-full">
                            <h4 className="title is-4">Library</h4>
                            <ul className="block-list">
                                <li>
                                    <LibraryScan scan_type="show" config={this.state.config} />
                                </li>
                                <li>
                                    <LibraryScan scan_type="movie" config={this.state.config} />
                                </li>
                            </ul>
                        </div>

                        <div className="column is-full">
                            <h4 className="title is-4">About flemzerd</h4>
                            <ul className="block-list">
                                <li>
                                    <div className={"columns is-gapless is-vcentered"}>
                                        <div className="column">
                                            Logged in as: <i>{this.state.config.Interface.Auth.Username}</i>
                                        </div>
                                        <div className={"column is-narrow"}>
                                            <button className="button is-danger is-small" onClick={() => Auth.logout()}>Log out</button>
                                        </div>
                                    </div>
                                </li>
                                <li>
                                    <div>
                                        Version: <code>{this.state.config.Version}</code>
                                    </div>
                                </li>
                                <li className="">
                                    <div className={"columns is-centered"}>
                                        <div className={"column is-narrow"}>
                                            flemzerd is an Open Source software. For feature ideas, notifying bugs or contributing, do not hesitate to head to the <a target="blank" href="https://github.com/macarrie/flemzerd">GitHub repo</a>
                                        </div>
                                    </div>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </>
        );
    }
}

export default Settings;
