import { api, version } from "../const.js";
import axios from "axios";

const ENDPOINT = {
  SEARCH: "search",
};

const getSearch = () => {
  return axios
    .get(
      `${process.env.REACT_APP_ENDPOINT ?? ""}/${api}/${version}/${
        ENDPOINT.SEARCH
      }`
    )
    .then((res) => res.data.result);
};

export { getSearch };
