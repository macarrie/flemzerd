import React from "react";

class Editable extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            isEditing: false,
            value: "",
        }

        this.state.value = this.props.value;

        this.handleFocus = this.handleFocus.bind(this);
        this.handleChange = this.handleChange.bind(this);
    }

    componentWillReceiveProps(nextProps) {
        this.setState({ value: nextProps.value });
    }

    handleFocus() {
        //if (this.state.isEditing) {
            //console.log("Is Editing: exiting focus");
            //if (typeof this.props.onFocusOut === "function") {
                //this.props.onFocusOut(this.state.value);
            //}
            //this.setState({isEditing: false});
        //} else {
            //if (typeof this.props.onFocus === "function") {
                //this.props.onFocus(this.state.value);
            //}
            //this.setState({isEditing: true});
        //}
    }

    handleChange(event) {
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
            <>
            <label
                onClick={this.handleFocus}>
                {this.state.value}
            </label>
                <input type="text"
                    value={this.state.value}
                    onBlur={this.handleFocus}
                    onChange={this.handleChange}
                    className={this.props.editingClassName}
                    autoFocus />
            </>
        )
    }
}

export default Editable;
