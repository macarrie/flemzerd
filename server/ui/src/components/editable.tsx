import React from "react";

import TextareaAutosize from 'react-textarea-autosize';
import {RiEdit2Line} from "react-icons/ri";
import {RiEraserLine} from "react-icons/ri";

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
        this.emptyValue = this.emptyValue.bind(this);
    }

    componentWillReceiveProps(nextProps: Props) {
        // Updating state when editing overwrites user changes in edit field
        if (!this.state.isEditing) {
            this.setState({ value: nextProps.value });
        }
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

    emptyValue(event: any) {
        if (typeof this.props.onFocusOut === "function") {
            this.props.onFocusOut("");
        }
        this.setState({
            value: "",
            isEditing: false,
        });
    }

    render() {
        return (
            <div className="columns is-mobile is-vcentered inlineedit">
                <div className={"column"}>
                    <TextareaAutosize
                        minRows={1}
                        maxRows={3}
                        value={this.state.value}
                        onBlur={this.handleFocus}
                        onClick={this.handleFocus}
                        onChange={this.handleChange}
                        className={`textarea is-static has-fixed-size ${this.props.editingClassName}`}
                        autoFocus
                        readOnly={!this.state.isEditing} />
                </div>
                {this.state.isEditing ? (
                    <div className={"column is-narrow"}>
                        <button className={"button is-naked is-dark"}
                                onClick={this.emptyValue}>
                                <span className={"icon edit-button has-text-grey is-hidden"}>
                                    <RiEraserLine/>
                                </span>
                        </button>
                    </div>
                ) : (
                    <div className={"column is-narrow"}>
                        <button className={"button is-naked is-dark"}
                                onClick={this.handleFocus}>
                                <span className={"icon edit-button has-text-grey is-hidden"}>
                                    <RiEdit2Line/>
                                </span>
                        </button>
                    </div>
                )}
            </div>
        );
    }
}

export default Editable;
