import React from "react";

import {RiCloseLine} from "react-icons/ri";
import {RiCheckLine} from "react-icons/ri";

type State = {
    param :Boolean,
    label :string,
};

class SettingsBoolean extends React.Component<any, State> {
    state :State = {
        param: false,
        label: "No label defined",
    }

    constructor(props :any) {
        super(props);

        this.state.param = this.props.param;
        this.state.label = this.props.label;
    }

    componentDidMount() {
    }

    render() {
        return (
            <>
                <li>
                    <div className="columns is-gapless">
                        <div className="column">
                            { this.state.label}
                        </div>
                        {this.state.param ? (
                            <div className="column is-narrow">
                                <button className={"button is-small is-naked is-disabled"}>
                                    <span className={"icon has-text-success"}>
                                        <RiCheckLine/>
                                    </span>
                                </button>
                            </div>
                        ) : (
                            <div className="column is-narrow">
                                <button className={"button is-small is-naked"}>
                                    <span className={"icon has-text-danger"}>
                                        <RiCloseLine/>
                                    </span>
                                </button>
                            </div>
                        )}
                    </div>
                </li>
            </>
        );
    }
}

export default SettingsBoolean;
