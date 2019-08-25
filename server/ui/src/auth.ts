import Helpers from "./utils/helpers";

class Auth {
    static IsLoggedIn = (): boolean => {
        return Auth.getToken() !== "";
    };

    static logout = () => {
        sessionStorage.removeItem("flemzerd_api_key");
        Helpers.eraseCookie("flemzerd_api_key");
        window.location.reload();
    };

    static setToken = (token: string, rememberme: boolean) => {
        sessionStorage.setItem('flemzerd_api_key', token);
        if (rememberme) {
            Helpers.createCookie("flemzerd_api_key", token, 14);
        }
    };

    static getToken = () :string => {
        let tokenSessionStorage = sessionStorage.getItem('flemzerd_api_key');
        if (tokenSessionStorage != null) {
            return tokenSessionStorage;
        }

        let tokenFromCookie = Helpers.readCookie('flemzerd_api_key');
        if (tokenFromCookie != null) {
            return tokenFromCookie;
        }

        return "";
    };
}

export default Auth;
