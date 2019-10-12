import React from "react";

import API from "../../utils/api";
import Config from "../../types/config";
import MediaInfo from "../../types/media_info";

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
            let media_list :MediaInfo[] = response.data;

            this.setState({
                import_status: "scanned",
                media_list: response.data,
            });
        }).catch(error => {
            console.log("Getting movie scan list error: ", error);
        });
    }

    scanShows() {
        API.Modules.Scanner.scan_shows().then(response => {
            let media_list :MediaInfo[] = response.data;

            this.setState({
                import_status: "scanned",
                media_list: response.data,
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
                <span className="uk-text-success">
                    <small><i>scanned</i></small>
                </span>
            );
        } else if (this.state.import_status === "scanning") {
            return (
                <div className="uk-flex uk-flex-middle">
                    <div className="uk-margin-small-right"
                    data-uk-spinner="ratio: 0.5"></div>
                    <div>
                        <small><i>Scanning</i></small>
                    </div>
                </div>
            );
        } else if (this.state.import_status === "importing") {
            return (
                <div className="uk-flex uk-flex-middle">
                    <div className="uk-margin-small-right"
                    data-uk-spinner="ratio: 0.5"></div>
                    <div>
                        <small><i>Importing</i></small>
                    </div>
                </div>
            );
        } else if (this.state.import_status === "done") {
            return (
                <span className="uk-text-primary uk-text-bold">
                    <small><i>done</i></small>
                </span>
            );
        }

        return (
            <button className="uk-button uk-button-small uk-button-default uk-margin-small-left"
                    onClick={() => this.performScan()}>
                Scan
            </button>
        );
    }

    countSelectedItems() {
        return this.state.media_list.filter(item => item.selected).length;
    }

    countImportedItems() {
        return this.state.media_list.filter(item => item.import_status === "success").length;
    }

    countImportedErrorItems() {
        return this.state.media_list.filter(item => item.import_status === "error").length;
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
            <div>
                <hr className="uk-margin-small-top" />
                    <div className="uk-grid uk-flex uk-flex-middle" data-uk-grid>
                        <div className="uk-width-expand">
                            <span className="uk-h3">Select media to import ({this.countSelectedItems()}/{this.state.media_list.length})</span>
                        </div>
                        <div className="uk-button-group">
                            <button className="uk-button uk-button-default uk-button-small" onClick={this.selectAllMedia}>Select all</button>
                            <button className="uk-button uk-button-default uk-button-small" onClick={this.unselectAllMedia}>Clear selection</button>
                        </div>
                    </div>
                    <form className="uk-grid-small" 
                          data-uk-grid
                          onSubmit={this.importScannedMedia}>
                        <ul className="uk-list uk-list-striped no-stripes uk-width-1-1 uk-overflow-auto uk-height-max-large">
                        {this.state.media_list.map((media, index) =>
                            <li key={index}>
                                <label>
                                    <input className="uk-checkbox" type="checkbox" name="selected" value={index} checked={this.state.media_list[index].selected} onChange={this.selectMediaToImport} /> {media.title}
                                </label>
                            </li>
                        )}
                        </ul>

                        <div className="uk-width-1-2">
                            <input className="uk-button uk-button-danger uk-input"
                                   value="Cancel"
                                   onClick={this.abortImport}
                                   type="button"/>
                        </div>
                        <div className="uk-width-1-2">
                            <input className="uk-button uk-button-primary uk-input"
                                   value="Import"
                                   type="submit"/>
                        </div>
                    </form>
            </div>
        );
    }

    mediaImportStatusLabel(media :MediaInfo) {
        if (media.import_status === "error") {
            return (
                <div className="uk-text-danger"><small><i>Error</i></small></div>
            );
        }

        if (media.import_status === "success") {
            return (
                <div className="uk-text-success"><small><i>Imported</i></small></div>
            );
        }

        return (
            <div className="uk-flex uk-flex-middle">
                <div className="uk-margin-small-right" data-uk-spinner="ratio: 0.5"></div>
                <div>
                    <small><i>Pending</i></small>
                </div>
            </div>
        );
    }

    importRecapDialog() {
        return (
            <div>
                <hr className="uk-margin-small-top" />
                <div className="uk-grid uk-flex uk-flex-middle uk-margin-small-bottom" data-uk-grid>
                    <div className="uk-width-1-2">
                        <span className="uk-h3">Importing media</span>
                    </div>
                    <div className="uk-width-1-2 uk-grid-small" data-uk-grid> 
                        <div className="uk-width-1-4">
                            <b>
                                {this.state.media_list.length} total
                            </b>
                        </div>
                        <div className="uk-width-1-4 uk-text-muted">
                            <b>
                                {this.state.media_list.length - this.countImportedItems() - this.countImportedErrorItems()} pending
                            </b>
                        </div>
                        <div className="uk-width-1-4 uk-text-success">
                            <b>
                                {this.countImportedItems()} success
                            </b>
                        </div>
                        <div className="uk-width-1-4 uk-text-danger">
                            <b>
                                {this.countImportedErrorItems()} errors
                            </b>
                        </div>
                    </div>
                </div>
                <ul className="uk-list uk-list-striped no-stripes uk-padding-remove">
                {this.state.media_list.map((media, index) =>
                    <li key={index} className="uk-flex uk-flex-middle">
                        <div className="uk-width-expand">
                            {media.title}
                        </div>
                        {this.mediaImportStatusLabel(media)}
                    </li>
                )}
                </ul>

                <div className="uk-inline uk-width-1-1">
                    <button className="uk-button uk-button-primary uk-input"
                            onClick={this.abortImport}>
                        Close
                    </button>
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
                <div className="uk-flex uk-flex-middle">
                    <div className="uk-width-expand">
                        {label}
                    </div>
                    <div className="uk-flex uk-flex-middle">
                        <code className="uk-text-small">
                            <span className="uk-text-small">
                                {option}
                            </span>
                        </code>
                        <div className="uk-margin-small-left">
                            {this.scanButton()}
                        </div>
                    </div>
                </div>
                {this.importDialog()}
            </>
        );
    }
}

export default LibraryScan;
