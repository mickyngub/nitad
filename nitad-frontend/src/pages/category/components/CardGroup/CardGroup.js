import React from "react";
import { Col, Row } from "react-bootstrap";
import { sortCard } from "../../utils";
import ProjectCard from "../ProjectCard";
import "./index.css";

const CardGroup = ({ cards, sortedBy }) => {
  return (
    <Row>
      {cards &&
        sortCard(cards, sortedBy)?.map(
          ({ id, title, videos, images, description, tags }, index) => (
            <Col
              xs={12}
              md={6}
              lg={4}
              className="mb-3 category-col d-flex justify-content-center"
              key={`${title}-${index}`}
            >
              <ProjectCard
                id={id}
                title={title}
                videos={videos}
                images={images}
                description={description}
                tags={tags}
              />
            </Col>
          )
        )}
    </Row>
  );
};

export default CardGroup;
