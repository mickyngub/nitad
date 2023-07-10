import Multiselect from "multiselect-react-dropdown";
import React, { useMemo, useState } from "react";
import { Button, Form, Offcanvas } from "react-bootstrap";
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

const FilterOffcanvas = ({
  allCategory,
  handleSelectCategoryTitle,
  show,
  handleClose,
  ...props
}) => {
  const location = useLocation();
  const queryStringObject = useMemo(() => {
    return getParams();
  }, [location.search]);

  const [sortedBy, setSortedBy] = useState(queryStringObject?.sort ?? 0);

  const [selectedCategoryTitle, setSelectedCategoryTitle] = useState(
    getDefaultSelectedCategoryTitle({ queryStringObject })
  );

  const getSelectedValue = ({ title }) => {
    if (queryStringObject?.[title]?.split(",").length > 0) {
      return queryStringObject?.[title]
        ?.split(",")
        .filter((title) => title !== "");
    } else return null;
  };
  return (
    <Offcanvas
      id="offcanvasNavbar"
      aria-labelledby="offcanvasNavbarLabel"
      onHide={handleClose}
      show={show}
      placement="start"
      {...props}
    >
      <Offcanvas.Header closeButton>
        <Offcanvas.Title></Offcanvas.Title>
      </Offcanvas.Header>
      <Offcanvas.Body className="d-flex flex-column">
        <div className="mb-auto">
          <h3>Sorted by</h3>
          <div>
            <Form.Range
              min={0}
              max={1}
              value={sortedBy}
              onChange={(e) => setSortedBy(e.target.value)}
            />
            <div className="d-flex justify-content-between ">
              <Form.Label>{sortedType[0]}</Form.Label>
              <Form.Label>{sortedType[1]}</Form.Label>
            </div>
          </div>
          <h6>Filter</h6>
          {allCategory?.map(({ title, subcategory }, index) => (
            <Multiselect
              key={`${title}-${index}`}
              isObject={false}
              placeholder={title}
              selectedValues={getSelectedValue({ title })}
              onKeyPressFn={function noRefCheck() {}}
              onRemove={(selectionList) => {
                setSelectedCategoryTitle({
                  ...selectedCategoryTitle,
                  [title]: selectionList,
                });
              }}
              onSearch={function noRefCheck() {}}
              onSelect={(selectionList) => {
                setSelectedCategoryTitle({
                  ...selectedCategoryTitle,
                  [title]: selectionList,
                });
              }}
              options={subcategory?.map(({ title }) => title.trim())}
              avoidHighlightFirstOption
              showArrow
              className="mb-3"
              style={{
                chips: {
                  backgroundColor: "rgba(0, 0, 0, 0.3)",
                  borderRadius: "3px",
                },
                multiselectContainer: {
                  color: "#828282",
                },
                searchBox: {
                  border: "none",
                  "border-bottom": "1px solid #828282",
                  "border-radius": "0px",
                },
              }}
            />
          ))}
        </div>

        <div className="zoom-btn-wrapper w-100">
          <div className="zoom-btn-wrapper-outer">
            <div className="zoom-btn-wrapper-inner">
              <Button
                onClick={() =>
                  handleSelectCategoryTitle({
                    selected: { ...selectedCategoryTitle, sort: sortedBy },
                  })
                }
                className="zoom-in-zoom-out"
              >
                Go
              </Button>
            </div>
          </div>
        </div>
      </Offcanvas.Body>
    </Offcanvas>
  );
};

export default FilterOffcanvas;
