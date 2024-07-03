import axios from "axios";
import { createContext, useContext, useEffect, useState } from "react";
import Cookies from "js-cookie";

const AuthContext = createContext();

export const useAuth = () => useContext(AuthContext);

export default function AuthContextProvider(props) {
  const [token, setToken_] = useState(Cookies.get("Cookies"));

  const setToken = (newToken) => {
    setToken_(newToken);
  };

  useEffect(() => {
    if (token) {
      axios.defaults.headers.common["Authorization"] = "Bearer " + token;
      Cookies.set("Cookies", token, { expires: 3, secure: true });
    } else {
      delete axios.defaults.headers.common["Authorization"];
      Cookies.remove("token");
    }
  });

  //memoized value of the authentication context

  // const contextValue = useMemo(() => {
  //   token, setToken;
  // }, [token]);

  const contextValue = {
    token,
    setToken,
  };

  return (
    <AuthContext.Provider value={contextValue}>
      {" "}
      {props.children}{" "}
    </AuthContext.Provider>
  );
}
