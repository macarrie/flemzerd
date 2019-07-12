class Auth {
    static IsLoggedIn = (): boolean => {
        return Auth.getToken() !== "";
    };

    static logout = () => {
        sessionStorage.removeItem("flemzerd_api_key");
        window.location.reload();
    };

    static setToken = (token :string) => {
        console.log("Setting token: ", token);
        sessionStorage.setItem('flemzerd_api_key', token);
    };

    static getToken = () :string => {
        let token = sessionStorage.getItem('flemzerd_api_key');
        if (token == null) {
            return "";
        }

        return token;
    };
}

export default Auth;
