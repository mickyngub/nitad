import React from "react";
import { Button, Form } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import FilterCheckboxGroup from "../FilterCheckboxGroup/FilterCheckboxGroup";
import "./index.css";
import { sortedType } from "./utils";

const Filter = ({
  filters,
  queryStringObject,
  handleChangeFilter,
  handleSorted,
}) => {
  const navigate = useNavigate();
  return (
    <>
      <h3>Sorted by</h3>
      <div>
        <Form.Range
          min={0}
          max={1}
          onChange={(e) => handleSorted(e.target.value)}
          defaultValue={0}
          size="sm"
          className="sorted-range-bar"
        />
        <div className="d-flex justify-content-between">
          <Form.Label>{sortedType[0]}</Form.Label>
          <Form.Label>{sortedType[1]}</Form.Label>
        </div>
      </div>
      <Button
        variant="link mb-md-3 p-0 accordion-btn-link font-weight-normal"
        onClick={() => navigate("/category")}
      >
        Clear Filter
      </Button>
      {filters?.map(({ title, subcategory }, index) => {
        if (subcategory?.length === 0) return <></>;
        return (
          <FilterCheckboxGroup
            title={title?.trim()}
            subcategory={subcategory}
            key={`${title}-${index}`}
            index={index}
            selected={queryStringObject?.[title]}
            handleChangeFilter={handleChangeFilter}
          />
        );
      })}
    </>
  );
};

export default Filter;
