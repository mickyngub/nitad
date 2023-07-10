import React, { useState } from "react";
import { Button, Collapse, Form } from "react-bootstrap";
import { IoIosArrowUp } from "react-icons/io";
import "./index.css";

const FilterCheckboxGroup = ({
  title,
  subcategory,
  selected,
  handleChangeFilter,
  index,
}) => {
  const [open, setOpen] = useState(true);
  return (
    <div className="mb-md-2">
      <Button
        variant="link"
        className="d-flex w-100 mb-md-2 align-items-center text-decoration-none accordion-btn-link font-weight-bold"
        onClick={() => setOpen(!open)}
      >
        <h4 variant="link" className="mb-0 category-title me-auto">
          {title}
        </h4>
        <IoIosArrowUp
          className={`ms-auto custom-arrow ${open ? "down" : "up"}`}
        />
      </Button>

      <Collapse in={open}>
        <Form>
          {subcategory?.map(({ title: subtitle }, index) => {
            return (
              <div key={`default-${title}-${subtitle}-${index}`}>
                <Form.Check
                  type={"checkbox"}
                  className={"custom-checkbox"}
                  id={`default-${title}-${subtitle}-${index}`}
                  label={subtitle?.trim()}
                  checked={selected?.includes(subtitle?.trim()) || false}
                  onChange={(e) => {
                    handleChangeFilter({
                      category: title?.trim(),
                      subcategory: subtitle?.trim(),
                    });
                  }}
                />
              </div>
            );
          })}
        </Form>
      </Collapse>
    </div>
  );
};

export default FilterCheckboxGroup;
