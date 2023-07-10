import { api, version } from "../const.js";
import axios from "axios";

const ENDPOINT = {
  SPATIAL: "spatial",
};

const getSpatialLink = () => {
  return axios
    .get(
      `${process.env.REACT_APP_ENDPOINT ?? ""}/${api}/${version}/${
        ENDPOINT.SPATIAL
      }`
    )
    .then((res) => res.data);
};

export { getSpatialLink };
