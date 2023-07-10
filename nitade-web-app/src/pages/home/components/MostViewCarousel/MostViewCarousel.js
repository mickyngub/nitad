import React, { useEffect, useState } from "react";
import { Carousel, Col, Container, Row } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import { appConfig } from "../../../../config";
import Content from "../Content";
import "./index.css";

const { STORAGE_IMAGE_URL, BUCKET_NAME, COLLECTION_NAME } = appConfig;

const MostViewCarousel = ({ mostView }) => {
  const [currentProject, setCurrentProject] = useState({});
  const navigate = useNavigate();
  const handleClickSeemore = (id) => navigate(`/project/${id}`);
  useEffect(() => {
    setCurrentProject(mostView[0]);
  }, [mostView]);
  return (
    <div className="most-view">
      <Container fluid={"xl"}>
        <Row className="align-items-md-center most-view-row">
          <Col xs={12} lg={8} className="most-view-col">
            {mostView && (
              <Carousel
                fade
                onSlide={(e) => setCurrentProject(mostView[e])}
                interval={3000}
                controls={false}
              >
                {mostView?.map(({ title, images }, index) => (
                  <Carousel.Item key={title + index}>
                    {/* <div className="image-wrapper-fade left" /> */}
                    <div className="most-view-image-wrapper">
                      <img src={`${images[0]}`} alt={title} className="image" />
                    </div>
                    {/* <div
                      className="image-wrapper-fade right"
                      style={{ right: 0, top: 0 }}
                    /> */}
                  </Carousel.Item>
                ))}
              </Carousel>
            )}
          </Col>
          <Col xs={12} lg={4} className="most-view-col">
            <Content
              title={currentProject?.title}
              content={currentProject?.description}
              handleClickSeemore={() => handleClickSeemore(currentProject?.id)}
            />
          </Col>
        </Row>
      </Container>
    </div>
  );
};

export default MostViewCarousel;
