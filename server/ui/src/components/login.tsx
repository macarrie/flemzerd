import React from "react";

class Login extends React.Component<any, any> {
    constructor(props: any) {
        super(props);
    }

    render() {
        return (
            <div className="uk-container uk-position-center">
                <div className="uk-flex uk-flex-middle uk-flex-center" data-uk-grid>
                    <div className="uk-width-1-1 uk-text-center">
                        <a className="uk-navbar-item uk-logo" href="/dashboard">
                            <span className="navbar-brand">flemzer</span>
                        </a>
                    </div>
                    <div className="uk-width-1-4">
                        <form className="uk-grid-small" data-uk-grid>
                            <div className="uk-inline uk-width-1-1">
                                <span className="uk-form-icon uk-form-icon-flip" uk-icon="icon: user"></span>
                                <input className="uk-input uk-width-1-1"
                                       placeholder="Username"
                                       type="text"/>
                            </div>

                            <div className="uk-inline uk-width-1-1">
                                <span className="uk-form-icon uk-form-icon-flip" uk-icon="icon: lock"></span>
                                <input className="uk-input uk-width-1-1"
                                       placeholder="Password"
                                       type="text"/>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        );
    }
}

export default Login;
