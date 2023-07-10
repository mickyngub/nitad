import React, { useState } from "react";
import { Collapse, Form, ListGroup, Offcanvas } from "react-bootstrap";

const getHighlightedText = (text, highlight) => {
  const parts = text.split(new RegExp(`(${highlight})`, "gi"));
  return (
    <span>
      {" "}
      {parts.map((part, i) => (
        <span
          key={i}
          style={
            part.toLowerCase() === highlight.toLowerCase()
              ? { fontWeight: "bold", backgroundColor: "yellow" }
              : {}
          }
        >
          {part}
        </span>
      ))}{" "}
    </span>
  );
};

const SearchOffcanvas = ({
  allCategory,
  handleSelectCategoryTitle,
  handleClose,
  show,
  allSearchTitleAble,
  handleSelectSearch,
  ...props
}) => {
  const [openSearchBar, setOpenSearchBar] = useState(false);
  const [searchContent, setSearchContent] = useState("");
  return (
    <Offcanvas show={show} onHide={handleClose} {...props}>
      <Offcanvas.Header closeButton>
        <Offcanvas.Title></Offcanvas.Title>
      </Offcanvas.Header>
      <Offcanvas.Body>
        <p>
          Type in the project, the technology, or the major that you are looking
          for
        </p>
        <Form.Control
          id="search-form"
          type="search"
          placeholder="Enter search terms"
          aria-label="Search"
          onFocus={() => setOpenSearchBar(true)}
          onBlur={() => setOpenSearchBar(false)}
          onChange={(e) => setSearchContent(e.target.value)}
        />
        <Collapse in={openSearchBar} className="search-results">
          <ListGroup>
            {searchContent !== "" &&
              allSearchTitleAble
                .filter(({ key }) =>
                  key
                    .trim()
                    .toLowerCase()
                    .includes(searchContent.trim().toLocaleLowerCase())
                )
                .map((item, index) => (
                  <ListGroup.Item
                    action
                    key={index}
                    onClick={() => handleSelectSearch(item)}
                  >
                    <div className="search-results-item">
                      {getHighlightedText(item.key, searchContent)}
                    </div>
                  </ListGroup.Item>
                ))}
          </ListGroup>
        </Collapse>
      </Offcanvas.Body>
    </Offcanvas>
  );
};

export default SearchOffcanvas;
