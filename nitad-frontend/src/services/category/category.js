import { api, version } from "../const.js";
import axios from "axios";

const ENDPOINT = {
  CATEGORY: "category",
};

const getAllCategory = () => {
  return axios
    .get(
      `${process.env.REACT_APP_ENDPOINT ?? ""}/${api}/${version}/${
        ENDPOINT.CATEGORY
      }`
    )
    .then((res) => res.data.result);
};

export { getAllCategory };
