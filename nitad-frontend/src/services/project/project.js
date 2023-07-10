import { api, version } from "../const.js";
import axios from "axios";

const ENDPOINT = {
  PROJECT: "project",
};

// -1 for asc and 1 for desc (default)
/* option = {
byViews: ...,
byCreatedAt: ...,
byName: ...,
sort:..., (views/name/updatedAt/createAt)
by:...,
page:...,
limit:...,
}
*/

const getAllProjectBySubcategory = ({ subcategoryList, options }) => {
  let param = "";
  subcategoryList?.forEach((id) => {
    param = param + "subcategoryId=" + id + "&";
  }); // expected props []
  Object.entries(options).forEach(([key, value]) => {
    param = param + `${key}=${value}&`;
  });
  return axios
    .get(
      `${process.env.REACT_APP_ENDPOINT ?? ""}/${api}/${version}/${
        ENDPOINT.PROJECT
      }?${param}`
    )
    .then((res) => res.data);
};

const getAllProject = () => {
  return axios
    .get(
      `${process.env.REACT_APP_ENDPOINT ?? ""}/${api}/${version}/${
        ENDPOINT.PROJECT
      }?limit=-1`
    )
    .then((res) => res.data);
};

const getProjectById = (id) => {
  return axios
    .get(
      `${process.env.REACT_APP_ENDPOINT ?? ""}/${api}/${version}/${
        ENDPOINT.PROJECT
      }/${id}`
    )
    .then((res) => res.data.result);
};

const getMostViewProject = () => {
  let param = `sort=views&by=-1&limit=5`;
  return axios
    .get(
      `${process.env.REACT_APP_ENDPOINT ?? ""}/${api}/${version}/${
        ENDPOINT.PROJECT
      }?${param}`
    )
    .then((res) => res.data);
};

export { getAllProjectBySubcategory,getAllProject, getProjectById, getMostViewProject };
