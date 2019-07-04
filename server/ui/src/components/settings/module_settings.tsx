import React from "react";

import ModuleSettings from "../../types/module_settings";
import TraktAuth from "./trakt_auth";
import TelegramAuth from "./telegram_auth";

type Props = {
    module :ModuleSettings,
    type :string,
};

type State = {
    module :ModuleSettings,
    type :string,
};

class ModuleSettingsComponent extends React.Component<Props, State> {
    state :State = {
        module: {} as ModuleSettings,
        type: "",
    };

    constructor(props :Props) {
        super(props);

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

export default ModuleSettingsComponent;
