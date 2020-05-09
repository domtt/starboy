import constate from "constate";
import { useLocation, useHistory } from "react-router-dom";

export const [TokenProvider, useToken] = constate((): string | null => {
  const location = useLocation();
  const query = new URLSearchParams(location.search);
  const history = useHistory();
  // get token from localstorage
  const token = localStorage.getItem("token");
  if (token !== null && token.length > 0) {
    return token;
  }
  // if the token is in the URL
  const tokenFromQuery = query.get("token");
  if (tokenFromQuery !== null && tokenFromQuery.length > 0) {
    // put it in localstorage and remove it from the URL
    localStorage.setItem("token", tokenFromQuery);
    query.set("token", "");
    history.push("/");
  }
  return null
});