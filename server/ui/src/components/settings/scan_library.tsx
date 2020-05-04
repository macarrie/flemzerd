import React from "react";

import API from "../../utils/api";
import Config from "../../types/config";
import MediaInfo from "../../types/media_info";
import Empty from "../empty";
import Helpers from "../../utils/helpers";

type Props = {
    scan_type :string,
    config :Config,
};

type State = {
    config :Config,
    media_list :MediaInfo[],
    import_status :string,
};

class LibraryScan extends React.Component<Props, State> {
    constructor(props :any) {
        super(props);

        this.state = {
            config: this.props.config,
            media_list: [],
            import_status: "",
        };

        this.scanButton = this.scanButton.bind(this);
        this.performScan = this.performScan.bind(this);
        this.selectMediaToImport = this.selectMediaToImport.bind(this);
        this.importScannedMedia = this.importScannedMedia.bind(this);
        this.changeMediaSelection = this.changeMediaSelection.bind(this);
        this.selectAllMedia = this.selectAllMedia.bind(this);
        this.unselectAllMedia = this.unselectAllMedia.bind(this);
        this.abortImport = this.abortImport.bind(this);
    }

    scanMovies() {
        API.Modules.Scanner.scan_movies().then(response => {
            let media_list = Array<MediaInfo>();
            if (response.data !== null) {
                media_list = response.data;
            }
            this.setState({
                import_status: "scanned",
                media_list: media_list,
            });
        }).catch(error => {
            console.log("Getting movie scan list error: ", error);
        });
    }

    scanShows() {
        API.Modules.Scanner.scan_shows().then(response => {
            let media_list :MediaInfo[] = new Array<MediaInfo>();
            for (let show in response.data) {
                media_list.push({
                    "title": show,
                    "selected": false,
                    "data": response.data[show],
                } as MediaInfo);
            }

            this.setState({
                import_status: "scanned",
                media_list: media_list,
            });
        }).catch(error => {
            console.log("Getting tvshow scan list error: ", error);
        });
    }

    performScan() {
        this.setState({import_status: "scanning"});
        if (this.props.scan_type === "movie") {
            this.scanMovies();
        } else {
            this.scanShows();
        }
    }

    selectMediaToImport(event :any) {
        event.persist();

        let media_info_list :MediaInfo[] = this.state.media_list;
        media_info_list[event.target.value].selected = event.target.checked;
        this.setState({
            media_list: media_info_list,
        });
    }

    updateStateMedia(id :string, status :string) {
        let media_list = this.state.media_list;
        for (let i = 0; i < media_list.length; i++) {
            if (media_list[i].Id === id) {
                media_list[i].import_status = status;
            }
        }

        this.setState({media_list: media_list});
    }

    importScannedMedia(event: any) {
        event.preventDefault();

        let selected_media :MediaInfo[] = this.state.media_list.filter(elt => elt.selected)
        this.setState({
            media_list: selected_media,
            import_status: "importing",
        }, () => {
            let batch_size :number = 10;
            let request_batches :any = new Array(Math.floor(selected_media.length / batch_size) + 1)
            for (let j = 0; j < request_batches.length; j++) {
                request_batches[j] = [];
            }

            let batch_nb :number = 0;
            for (let i = 0; i < selected_media.length; i++) {
                if (request_batches[batch_nb].length >= batch_size) {
                  batch_nb += 1;
                }

                if (this.props.scan_type === "movie") {
                    request_batches[batch_nb].push(API.Modules.Scanner.import_movie(selected_media[i]));
                } else {
                    request_batches[batch_nb].push(API.Modules.Scanner.import_show(selected_media[i]));
                }
            }

            for (let batch = 0; batch < request_batches.length; batch++) {
                Promise.all(request_batches[batch].map((promise :any) => promise.catch((err :any) => {
                    return err;
                }))).then(values => {
                    for (let req = 0; req < values.length; req++) {
                        let id = JSON.parse((values[req] as any).config.data)["Id"];
                        let retval = (values[req] as any).status;
                        let status :string = "";

                        if (retval !== 200) {
                            status = "error";
                        } else {
                            status = "success";
                        }
                        this.updateStateMedia(id, status);

                        // Check if import is done
                        if (this.state.media_list.length - this.countImportedItems() - this.countImportedErrorItems() <= 0) {
                            this.setState({
                                import_status: "done",
                            });
                        }
                    }
                });
            }
        });
    }

    scanButton() {
        if (this.state.import_status === "scanned") {
            return (
                <span className="has-text-success">
                    <small><i>scanned</i></small>
                </span>
            );
        } else if (this.state.import_status === "scanning") {
            return (
                <div className="columns is-vcentered">
                    <div className={"column"}>
                        <button className={"button is-small is-naked is-disabled is-loading"}>Scanning</button>
                    </div>
                    <div className={"column"}>
                        <small><i>Scanning</i></small>
                    </div>
                </div>
            );
        } else if (this.state.import_status === "importing") {
            return (
                <div className="columns is-vcentered">
                    <div className={"column"}>
                        <button className={"button is-small is-naked is-disabled is-loading"}>Importing</button>
                    </div>
                    <div className={"column"}>
                        <small><i>Importing</i></small>
                    </div>
                </div>
            );
        } else if (this.state.import_status === "done") {
            return (
                <span className="has-text-info has-text-bold">
                    <small><i>done</i></small>
                </span>
            );
        }

        return (
            <button className="button is-small is-info is-outlined"
                    onClick={() => this.performScan()}>
                Scan
            </button>
        );
    }

    countSelectedItems() {
        if (this.state.media_list === null)  {
            return 0;
        }

        return this.state.media_list.filter(item => item.selected).length;
    }

    countImportedItems() {
        if (this.state.media_list === null)  {
            return 0;
        }

        return this.state.media_list.filter(item => item.import_status === "success").length;
    }

    countImportedErrorItems() {
        if (this.state.media_list === null)  {
            return 0;
        }

        return this.state.media_list.filter(item => item.import_status === "error").length;
    }

    getSeasonCount(media :MediaInfo) {
        if (media.data == null) {
            return 0;
        }

        return Object.keys(media.data).length;
    }

    getEpisodeCount(media :MediaInfo) {
        if (media.data == null) {
            return 0;
        }

        let count = 0;
        for (let season in media.data) {
            count += Object.keys(media.data[season]).length;
        }

        return count;
    }

    changeMediaSelection(selected :boolean) {
        let list = this.state.media_list;
        for (let i = 0; i < list.length; i++) {
            list[i].selected = selected;
        }

        this.setState({media_list: list});
    }

    selectAllMedia() {
        this.changeMediaSelection(true);
    }

    unselectAllMedia() {
        this.changeMediaSelection(false);
    }

    abortImport() {
        this.setState({
            media_list: [],
            import_status: "",
        });
    }

    selectedMediaToImportDialog() {
        return (
            <div className={"column is-full"}>
                <hr />
                    {Helpers.count(this.state.media_list) > 0 ? (
                        <div className="columns is-gapless">
                            <div className="column">
                                <span className="title is-4">Select media to import ({this.countSelectedItems()}/{this.state.media_list.length})</span>
                            </div>
                            <div className="column is-narrow">
                                <div className={"buttons"}>
                                    <button className="button is-small is-info is-outlined" onClick={this.selectAllMedia}>Select all</button>
                                    <button className="button is-small is-info is-outlined" onClick={this.unselectAllMedia}>Clear selection</button>
                                </div>
                            </div>
                        </div>
                    ) : (
                        <Empty label={"No media detected in library"} />
                    )}
                    <form className="columns is-multiline"
                          onSubmit={this.importScannedMedia}>
                        <div className={"column is-full"}>
                            <ul className="">
                            {this.state.media_list.map((media, index) =>
                                <li key={index}>
                                    <div className={"field"}>
                                        <label className={"checkbox"}>
                                            <input type="checkbox" name="selected" value={index} checked={this.state.media_list[index].selected} onChange={this.selectMediaToImport} />
                                            <span></span>
                                            {media.title} &nbsp;
                                            {media.data != null && (
                                                <i className="has-text-grey">
                                                    ({this.getSeasonCount(media)} seasons, {this.getEpisodeCount(media)} episodes)
                                                </i>
                                            )}
                                        </label>
                                    </div>
                                </li>
                            )}
                            </ul>
                        </div>

                        <div className="column">
                            <input className="button is-danger is-fullwidth"
                                   value="Cancel"
                                   onClick={this.abortImport}
                                   type="button"/>
                        </div>
                        {Helpers.count(this.state.media_list) > 0 && (
                            <div className="column">
                                <input className="button is-success is-fullwidth"
                                       value="Import"
                                       type="submit"/>
                            </div>
                        )}
                    </form>
            </div>
        );
    }

    mediaImportStatusLabel(media :MediaInfo) {
        if (media.import_status === "error") {
            return (
                <div className="has-text-danger"><small><i>Error</i></small></div>
            );
        }

        if (media.import_status === "success") {
            return (
                <div className="has-text-success"><small><i>Imported</i></small></div>
            );
        }

        return (
            <div className="columns is-vcentered">
                <div className={"column"}>
                    <button className={"button is-small is-naked is-disabled is-loading"}>Pending</button>
                </div>
                <div className={"column"}>
                    <small><i>Pending</i></small>
                </div>
            </div>
        );
    }

    importRecapDialog() {
        return (
            <div className={"column is-full"}>
                <hr />
                <div className="columns is-vcentered is-multiline">
                    <div className="column is-half">
                        <span className="title is-4">Importing media</span>
                    </div>
                    <div className="columns column is-half">
                        <div className="column is-one-quarter">
                            <b>
                                {this.state.media_list.length} total
                            </b>
                        </div>
                        <div className="column is-one-quarter has-text-grey">
                            <b>
                                {this.state.media_list.length - this.countImportedItems() - this.countImportedErrorItems()} pending
                            </b>
                        </div>
                        <div className="column is-one-quarter has-text-success">
                            <b>
                                {this.countImportedItems()} success
                            </b>
                        </div>
                        <div className="column is-one-quarter has-text-danger">
                            <b>
                                {this.countImportedErrorItems()} errors
                            </b>
                        </div>
                    </div>
                    <div className={"column is-full"}>
                        <ul className="">
                        {this.state.media_list.map((media, index) =>
                            <li key={index} className="">
                                <div className={"columns"}>
                                    <div className="column">
                                        {media.title} &nbsp;
                                        {media.data != null && (
                                            <i className="has-text-grey">
                                                ({this.getSeasonCount(media)} seasons, {this.getEpisodeCount(media)} episodes)
                                            </i>
                                        )}
                                    </div>
                                    <div className={"column is-narrow"}>
                                        {this.mediaImportStatusLabel(media)}
                                    </div>
                                </div>
                            </li>
                        )}
                        </ul>
                    </div>

                    <div className="column is-full">
                        <button className="button is-danger is-fullwidth"
                                onClick={this.abortImport}>
                            Close
                        </button>
                    </div>
                </div>
            </div>
        );
    }

    importDialog() {
        if (this.state.import_status === "scanned") {
            return this.selectedMediaToImportDialog();
        } else if (this.state.import_status === "importing" || this.state.import_status === "done") {
            return this.importRecapDialog();
        } else {
            return null;
        }
    }

    render() {
        if (this.state.config == null) {
            return null;
        }

        let label = "";
        let option = "";
        if (this.props.scan_type === "movie") {
            label = "Movies library path";
            option = this.state.config.Library.MoviePath;
        } else {
            label = "TV shows library path";
            option = this.state.config.Library.ShowPath;
        }

        return (
            <>
                <div className="columns is-gapless is-multiline is-vcentered library-scan-dialog">
                    <div className="column">
                        {label}
                    </div>
                    <div className="column is-narrow">
                        <code>
                            <span className="">
                                {option}
                            </span>
                        </code>
                    </div>
                    <div className="column is-narrow">
                        {this.scanButton()}
                    </div>
                    {this.importDialog()}
                </div>
            </>
        );
    }
}

export default LibraryScan;
