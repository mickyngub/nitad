import React from "react";
import { Container } from "react-bootstrap";

const Banner = () => {
  return (
    <div className="position-relative banner">
      <div className="banner-bg w-100 h-100 d-flex align-items-center position-absolute">
        <Container fluid={"xl"} className="">
          <div className="ps-5 banner-qoute-wrapper">
            <h1 className="banner-qoute">
              Bring yourself to
              <br /> a whole new level
              <br />
              <br /> Transform the world
              <br /> through technology
            </h1>
          </div>
        </Container>
      </div>
    </div>
  );
};

export default Banner;
