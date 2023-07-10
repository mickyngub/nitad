import React from "react";

const Content = ({ title, content, button }) => {
  return (
    <>
      <div className="d-flex">
        <h2 className="project-content-title">{title}</h2>
        {button}
      </div>
      <p className="project-content-paragraph">&emsp; &emsp;{content}</p>
    </>
  );
};

export default Content;
