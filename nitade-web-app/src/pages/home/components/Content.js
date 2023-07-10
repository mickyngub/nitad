import { Button } from "react-bootstrap";
import React from "react";

const Content = ({ title, content, handleClickSeemore }) => {
  return (
    <div className="text-center most-view-content py-3 h-100">
      <h1 className="text-ellipsis-title">{title}</h1>
      <p className="text-ellipsis">{content}</p>
      <Button
        className="d-block w-100 seemore-btn mt-auto"
        style={{ backgroundColor: "#6166B3" }}
        onClick={handleClickSeemore}
      >
        See more
      </Button>
    </div>
  );
};

export default Content;
