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

    componentWillReceiveProps(nextProps: Props) {
        this.setState({
            module: nextProps.module,
            type: nextProps.type,
        });
    }

    render() {
        if (this.state.type === "watchlist" && this.state.module.Name === "trakt") {
            return (
                <>
                    <div className={"columns is-mobile is-multiline is-gapless is-vcentered"}>
                        <div className={"column"}>
                            {this.state.module.Name}
                        </div>
                        <TraktAuth/>
                    </div>
                </>
            );
        }

        if (this.state.type === "notifier" && this.state.module.Name === "telegram") {
            return (
                <>
                    <div className={"columns is-mobile is-multiline is-gapless is-vcentered"}>
                        <div className={"column"}>
                            {this.state.module.Name}
                        </div>
                        <TelegramAuth/>
                    </div>
                </>
            );
        }

        if (this.state.type === "indexer") {
            return (
                <>
                    <div className={"columns is-multiline is-gapless is-vcentered"}>
                        <div className={"column is-narrow"}>
                            {this.state.module.Config.name}
                        </div>
                        <div className="column is-clipped">
                            <small>
                                <code>{this.state.module.Config.url}</code>
                            </small>
                        </div>
                    </div>
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
