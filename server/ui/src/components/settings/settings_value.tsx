import React from "react";

type State = {
    param :string,
    label :string,
    unit :string,
};

class SettingsValue extends React.Component<any, State> {
    state :State = {
        param: "",
        label: "No label defined",
        unit: "",
    }

    constructor(props :any) {
        super(props);

        this.state.param = this.props.param;
        this.state.label = this.props.label;
        this.state.unit = this.props.unit;
    }

    componentDidMount() {
    }

    render() {
        return (
            <>
                <li>
                    <div className={"columns is-mobile is-gapless"}>
                        <div className="column">
                            { this.state.label }
                        </div>
                        <div className={"column is-narrow"}>
                            <code>
                                { this.state.param }{ this.state.unit }
                            </code>
                        </div>
                    </div>
                </li>
            </>
        );
    }
}

export default SettingsValue;
