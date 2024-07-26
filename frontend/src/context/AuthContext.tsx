import React, { createContext, useContext, useEffect, useState, ReactNode } from "react";
import axios from "axios";
import Cookies from "js-cookie";

// Define the shape of the authentication context
interface AuthContextType {
    token: string | undefined;
    setToken: (newToken: string | undefined) => void;
}


const AuthContext = createContext<AuthContextType | undefined>(undefined);


export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthContextProvider");
    }
    return context;
};


interface AuthContextProviderProps {
    children: ReactNode;
}

const AuthContextProvider: React.FC<AuthContextProviderProps> = ({ children }) => {
    const [token, setToken_] = useState<string | undefined>(Cookies.get("Cookies"));

    const setToken = (newToken: string | undefined) => {
        setToken_(newToken);
    };

    useEffect(() => {
        if (token) {
            axios.defaults.headers.common["Authorization"] = "Bearer " + token;
            Cookies.set("Cookies", token, { expires: 3, secure: true });
        } else {
            delete axios.defaults.headers.common["Authorization"];
            Cookies.remove("Cookies");
        }
    }, [token]); // Ensure useEffect has a dependency array

    const contextValue: AuthContextType = {
        token,
        setToken,
    };

    return (
        <AuthContext.Provider value={contextValue}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthContextProvider;
