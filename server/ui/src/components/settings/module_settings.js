import React from "react";

import TraktAuth from "./trakt_auth";
import TelegramAuth from "./telegram_auth";

class ModuleSettings extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            module: null,
            type: null,
        };

        this.state.module = this.props.module;
        this.state.type = this.props.type;
    }

    componentWillReceiveProps(nextProps) {
        this.setState({
            module: nextProps.module,
            type: nextProps.type,
        });
    }

    render() {
        if (this.state.type === "watchlist" && this.state.module.Name === "trakt") {
            return (
                <>
                    {this.state.module.Name}
                    <TraktAuth/>
                </>
            );
        }

        if (this.state.type === "notifier" && this.state.module.Name === "telegram") {
            return (
                <>
                    {this.state.module.Name}
                    <TelegramAuth/>
                </>
            );
        }

        if (this.state.type === "indexer") {
            return (
                <>
                    {this.state.module.Config.name}
                    <span className="float-right text-muted">
                        <small>
                            <code>{this.state.module.Config.url}</code>
                        </small>
                    </span>
                </>
            );
        }

        return (
            <>
                {this.state.module.Name}
            </>
        );
    }
}

export default ModuleSettings;
