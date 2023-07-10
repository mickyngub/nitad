import React, { useEffect, useState } from "react";
import { Container } from "react-bootstrap";
import { useParams } from "react-router-dom";
import { getProjectById } from "../../../services/project/project";
import "../index.css";
import ProjectContent from "./ProjectContent";
import ProjectOwner from "./ProjectOwner";

const getSubcategoryTitle = (subcategory) => {
  let s = "";
  subcategory?.forEach(({ title }, index) => {
    if (index > 0) s = s + ", ";
    s = s + title;
  });
  return s;
};

const Project = () => {
  const { id } = useParams();
  const [projectDetail, setProjectDetail] = useState({});
  useEffect(() => {
    getProjectById(id).then((data) => {
      setProjectDetail(data);
    });
  }, [id]);

  return (
    projectDetail && (
      <div className="project-image-container">
        <div className="project-section-container">
          <Container fluid={"xl"}>
            <div className="project-section-container-inner">
              <ProjectOwner
                authors={projectDetail?.authors}
                authorEmail={projectDetail?.emails}
                projectInspiration={projectDetail?.inspiration}
              />
              <ProjectContent
                title={projectDetail?.title}
                videos={projectDetail?.videos}
                images={projectDetail?.images}
                abstract={projectDetail?.abstract}
                description={projectDetail?.description}
                authors={projectDetail?.authors}
                projectInspiration={projectDetail?.inspiration}
                reportDownloadLink={projectDetail?.report}
                virtualLink={projectDetail?.virtualLink}
                authorEmail={projectDetail?.emails}
                views={projectDetail.views}
                tags={projectDetail?.category}
                projectType={getSubcategoryTitle(projectDetail?.subcategory)}
              />
            </div>
          </Container>
        </div>
      </div>
    )
  );
};

export default Project;
