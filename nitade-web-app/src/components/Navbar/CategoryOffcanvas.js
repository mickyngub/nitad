import Multiselect from "multiselect-react-dropdown";
import React, { useMemo, useState } from "react";
import { Accordion, Button, Form, Navbar, Offcanvas } from "react-bootstrap";
import { useLocation } from "react-router-dom";
import { sortedType } from "../../pages/category/components/Filter/utils";

function getParams(url = window.location) {
  let params = {};
  new URL(url).searchParams.forEach(function (val, key) {
    if (params[key] !== undefined) {
      if (!Array.isArray(params[key])) {
        params[key] = [params[key]];
      }
      params[key].push(val);
    } else {
      params[key] = val;
    }
  });

  return params;
}

const getDefaultSelectedCategoryTitle = ({ queryStringObject }) => {
  let entries = Object?.entries(queryStringObject);
  return entries.reduce(
    (acc, curr) => ({ ...acc, [curr[0]]: curr[1].split(",") }),
    {}
  );
};

const CategoryOffcanvas = ({
  allCategory,
  handleSelectCategoryTitle,
  handleClose,
  show,
  ...props
}) => {
  const location = useLocation();

  return (
    <Offcanvas
      id="offcanvasNavbar"
      aria-labelledby="offcanvasNavbarLabel"
      placement="start"
      show={show}
      onHide={handleClose}
      {...props}
    >
      <Offcanvas.Header closeButton>
        <Offcanvas.Title></Offcanvas.Title>
      </Offcanvas.Header>
      <Offcanvas.Body>
        <Button variant="link" href="/about" className="category-link-button">
          ABOUT US
        </Button>
        <Button
          variant="default"
          href="/category"
          className="category-link-button"
        >
          EXPLORE PROJECT
        </Button>
        <Accordion
          defaultActiveKey="0"
          flush
          className="category-accordion"
          alwaysOpen
        >
          {allCategory?.map(({ title: categoryTitle, subcategory }, index) => (
            <Accordion.Item eventKey={index} key={`${categoryTitle}-${index}`}>
              <Accordion.Header>{categoryTitle}</Accordion.Header>
              <Accordion.Body>
                {subcategory?.map(({ title: subcategoryTitle }, index) => (
                  <Button
                    variant={"link"}
                    size="large"
                    className="w-100 text-decoration-none"
                    href={`/category?${categoryTitle}=${subcategoryTitle}`}
                    key={`${subcategoryTitle}-${index}`}
                  >
                    {subcategoryTitle}
                  </Button>
                ))}
              </Accordion.Body>
            </Accordion.Item>
          ))}
        </Accordion>
      </Offcanvas.Body>
    </Offcanvas>
  );
};

export default CategoryOffcanvas;
