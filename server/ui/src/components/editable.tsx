import React from "react";

type Props = {
    onFocus?(value: string): void,
    onFocusOut?(value: string): void,
    editingClassName :string,
    value :string,
};

type State = {
    isEditing :boolean,
    value :string,
};

class Editable extends React.Component<Props, State> {
    state :State = {
        isEditing: false,
        value: "",
    };

    constructor(props :Props) {
        super(props);

        this.state.value = this.props.value;

        this.handleFocus = this.handleFocus.bind(this);
        this.handleChange = this.handleChange.bind(this);
    }

    componentWillReceiveProps(nextProps: Props) {
        this.setState({ value: nextProps.value });
    }

    handleFocus(event :React.MouseEvent | React.FocusEvent) {
        if (this.state.isEditing) {
            if (typeof this.props.onFocusOut === "function") {
                this.props.onFocusOut(this.state.value);
            }
            this.setState({isEditing: false});
        } else {
            if (typeof this.props.onFocus === "function") {
                this.props.onFocus(this.state.value);
            }
            this.setState({isEditing: true});
        }
    }

    handleChange(event: any) {
        this.setState({
            value: event.target.value,
        });
    }

    render() {
        if (this.state.isEditing) {
            return (
                <input type="text"
                    value={this.state.value}
                    onBlur={this.handleFocus}
                    onChange={this.handleChange}
                    className={this.props.editingClassName}
                    autoFocus />
            )
        }

        return (
            <label
                onClick={this.handleFocus}>
                {this.state.value}
            </label>
        )
    }
}

export default Editable;
