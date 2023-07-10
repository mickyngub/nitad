import { sortedType } from "./components/Filter/utils";

export const sortCard = (projects, sortedBy) => {
  if (sortedBy === sortedType[0]) {
    return projects.sort(function (a, b) {
      var keyA = a.title,
        keyB = b.title;
      if (keyA < keyB) return -1;
      if (keyA > keyB) return 1;
      return 0;
    });
  } else {
    return projects.sort(function (a, b) {
      var keyA = a.views,
        keyB = b.views;
      if (keyA < keyB) return 1;
      if (keyA > keyB) return -1;
      return 0;
    });
  }
};
