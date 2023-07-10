import React from "react";
import { Button } from "react-bootstrap";
import { IoIosCloseCircle } from "react-icons/io";
import "./index.css";

const Tag = ({
  categoryTitle,
  subcategoryTitle,
  variant = "disabled",
  onCloseTag,
}) => {
  const onClickTag = () => {
    window.location.href = `/category?${categoryTitle}=${subcategoryTitle}`;
  };

  return (
    <Button className="tag-link d-flex align-items-center">
      <span
        onClick={(e) => {
          e.preventDefault();
          onClickTag();
        }}
        style={{ padding: "5px" }}
      >
        {subcategoryTitle}
      </span>
      {variant === "closeable" && (
        <span
          onClick={() => {
            onCloseTag();
          }}
        >
          <IoIosCloseCircle className="mx-1" style={{ zIndex: "500" }} />
        </span>
      )}
    </Button>
  );
};

export default Tag;
