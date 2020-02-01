import constate from "constate";
import { useLocation, useHistory } from "react-router-dom";

export const [TokenProvider, useToken] = constate((): string | null => {
  const location = useLocation();
  const query = new URLSearchParams(location.search);
  const history = useHistory();
  // TODO: hide token from url and store it in localstorage
  // then search for it in localstorage as well
  const token = localStorage.getItem("token");
  if (token !== null && token.length > 0) {
    return token;
  }
  const tokenFromQuery = query.get("token");
  if (tokenFromQuery !== null && tokenFromQuery.length > 0) {
    localStorage.setItem("token", tokenFromQuery);
    query.set("token", "");
    history.push("/");
  }
  return null;
});
