import React from "react";

import {RiSearchLine} from "react-icons/ri";
import {RiCloseLine} from "react-icons/ri";

type Props = {
    updateSearch(value: string): void,
};

type State = {
    display_search_form :boolean,
};

class ItemSearchControls extends React.Component<Props, State> {
    constructor(props: Props) {
        super(props);
        this.state = {
            display_search_form: false,
        };

        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event: any) {
        this.props.updateSearch(event.target.value);
    }

    displaySearchForm() {
        this.setState({
            display_search_form: true,
        })
    }

    closeSearchForm() {
        this.setState({
            display_search_form: false,
        })
        this.props.updateSearch("");
    }

    render() {
        return (
            <div className="column is-half-mobile item-search">
                {this.state.display_search_form ? (
                    <div className="field">
                        <label className="label"></label>
                        <div className={"control has-icons-left has-icons-right"}>
                            <input className="input is-static"
                                   placeholder="Search"
                                   onChange={this.handleChange}
                                   autoFocus
                                   ref={input => input && input.focus()}
                                   type="text"/>
                            <span className="icon is-left">
                                <RiSearchLine />
                            </span>
                            <span className="icon is-right clickable"
                                  onClick={() => this.closeSearchForm()}>
                                <RiCloseLine />
                            </span>
                        </div>
                    </div>
                    ) : (
                    <div className="field">
                        <label className="label"></label>
                        <div className={"control has-icons-left has-icons-right"}
                             onClick={() => this.displaySearchForm()}>
                            <input className="input is-static is-disabled"
                                   readOnly={true}
                                   value={""}
                                   type="text"/>
                            <span className="icon is-left clickable">
                                <RiSearchLine />
                            </span>
                        </div>
                    </div>
                )}
            </div>
        );
    }
}

export default ItemSearchControls;
