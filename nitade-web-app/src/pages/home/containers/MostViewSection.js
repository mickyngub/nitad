import React, { useEffect, useState } from "react";
import { getMostViewProject } from "../../../services/project/project";
import MostViewCarousel from "../components/MostViewCarousel/MostViewCarousel";

const MostViewSection = () => {
  const [mostViewProject, setMostViewProject] = useState([]);

  useEffect(() => {
    getMostViewProject()
      .then((data) => setMostViewProject(data.result))
      .catch((err) => console.log(err));
  }, []);

  return (
    mostViewProject &&
    mostViewProject.length > 0 && (
      <MostViewCarousel mostView={mostViewProject} />
    )
  );
};

export default MostViewSection;
