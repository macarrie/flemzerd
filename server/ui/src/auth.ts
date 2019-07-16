class Auth {
    static IsLoggedIn = (): boolean => {
        return Auth.getToken() !== "";
    };

    static logout = () => {
        sessionStorage.removeItem("flemzerd_api_key");
        localStorage.removeItem("flemzerd_api_key");
        window.location.reload();
    };

    static setToken = (token: string, rememberme: boolean) => {
        sessionStorage.setItem('flemzerd_api_key', token);
        if (rememberme) {
            localStorage.setItem('flemzerd_api_key', token);
        }
    };

    static getToken = () :string => {
        let tokenSessionStorage = sessionStorage.getItem('flemzerd_api_key');
        if (tokenSessionStorage != null) {
            return tokenSessionStorage;
        }

        let tokenLocalStorage = localStorage.getItem('flemzerd_api_key');
        if (tokenLocalStorage != null) {
            return tokenLocalStorage;
        }

        return "";
    };
}

export default Auth;
