import React, { useEffect, useRef, useState } from "react";
import {
  Button,
  Card,
  Collapse,
  Container,
  Form,
  InputGroup,
  ListGroup,
  Nav,
  Navbar,
} from "react-bootstrap";
import { BsSearch } from "react-icons/bs";
import { IoArrowForward } from "react-icons/io5";
import { useMediaQuery } from "react-responsive";
import { useLocation, useNavigate } from "react-router-dom";
import logo from "../../assets/nitade-logo.png";
import { getAllCategory } from "../../services/category/category";
import { getSearch } from "../../services/search/search";
import { getSpatialLink } from "../../services/spatial/spatial";
import CategoryOffcanvas from "./CategoryOffcanvas";
import FilterOffcanvas from "./FilterOffcanvas";
import "./index.css";
import SearchOffcanvas from "./SearchOffcanvas";
import digitalMuseum from "../../assets/digital_museum.svg";

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

const useFocus = () => {
  const htmlElRef = useRef(null);
  const setFocus = () => {
    htmlElRef.current && htmlElRef.current.focus();
  };
  return [htmlElRef, setFocus];
};

const BasicNavbar = ({ isAboutUs }) => {
  const [open, setOpen] = useState(false);
  const [isMouseEnter, setIsMouseEnter] = useState(false);
  const [isClickSearch, setIsClickSearch] = useState(false);
  const [searchContent, setSearchContent] = useState("");
  const [allSearchTitleAble, setAllSearchTitleAble] = useState([]);
  const [openSearchOffcanvas, setOpenSearchOffcanvas] = useState(false);
  const [openCategoryOffcanvas, setOpenCategoryOffcanvas] = useState(false);
  const [openFilterOffcanvas, setOpenFilterOffcanvas] = useState(false);
  const [openNested, setOpenNested] = useState(false);
  const [openSearchBar, setOpenSearchBar] = useState(false);
  const [selectedTitle, setSelectedTitle] = useState("");
  const [allCategory, setAllCategory] = useState();
  const [spatialLink, setSpatialLink] = useState("");
  const [inputFocus, setInputFocus] = useFocus();
  const isMobile = useMediaQuery({ query: "(max-width: 992px)" });
  const isXXL = useMediaQuery({ query: "(max-width: 1400px)" });
  const navigate = useNavigate();
  const location = useLocation();

  const handleClsoeSearchOffcanvas = () => setOpenSearchOffcanvas(false);
  const handleOpenSearchOffcanvas = () => setOpenSearchOffcanvas(true);

  const handleClsoeCategoryOffcanvas = () => setOpenCategoryOffcanvas(false);
  const handleOpenCategoryOffcanvas = () => setOpenCategoryOffcanvas(true);

  const handleClsoeFilterOffcanvas = () => setOpenFilterOffcanvas(false);
  const handleOpenFilterOffcanvas = () => setOpenFilterOffcanvas(true);

  const handleSelectSearch = (item) => {
    if (item.type === "subcategory") {
      window.location.href = `/category?${item.value}`;
    } else window.location.href = `/project/${item.value}`;
    setOpenSearchBar(false);
    setIsClickSearch(false);
    setOpenSearchOffcanvas(false);
    setSearchContent("");
  };

  const renderToggleMenu = () => {
    if (location.pathname.includes("category")) {
      return (
        <>
          <FilterOffcanvas
            key={0}
            allCategory={allCategory}
            handleSelectCategoryTitle={handleSelectCategoryTitle}
            handleClose={handleClsoeFilterOffcanvas}
            show={openFilterOffcanvas}
          />
          <Navbar.Toggle onClick={handleOpenFilterOffcanvas} />
        </>
      );
    } else {
      return (
        <>
          <CategoryOffcanvas
            key={1}
            allCategory={allCategory}
            handleClose={handleClsoeCategoryOffcanvas}
            show={openCategoryOffcanvas}
            handleSelectCategoryTitle={handleSelectCategoryTitle}
          />
          <Navbar.Toggle onClick={handleOpenCategoryOffcanvas} />
        </>
      );
    }
  };

  const handleSelectedCategory = ({ category, subcategory }) => {
    navigate(`/category?${category}=${subcategory}`);
  };

  const handleSelectedTitle = ({ title }) => {
    if (title) {
      setSelectedTitle(title);
      setOpenNested(true);
    }
  };
  const handleMouseLeave = () => {
    setSelectedTitle("");
    setOpenNested(false);
    setOpen(false);
  };

  const handleOpenSearchBar = () => {
    setOpenSearchBar(true);
    setIsMouseEnter(true);
  };

  const handleCloseSearchBar = () => {
    setOpenSearchBar(false);
    setSearchContent("");
  };

  const handleSelectCategoryTitle = ({ selected }) => {
    let entries = Object.entries(selected);
    let str = "";
    entries.forEach(([key, value]) => {
      if (value.length > 0)
        str += `${key}=${
          key !== "sort" ? value.map((title) => title.trim()).join(",") : value
        }&`;
    });
    window.location.href = `category?${str}`;
  };

  useEffect(() => {
    if (openSearchBar) {
      setInputFocus(inputFocus);
    }
  }, [openSearchBar, inputFocus, setInputFocus]);

  useEffect(() => {
    let temp = [];
    getSearch().then((data) => {
      Object.entries(data).forEach(([key, value]) => {
        if (key.trim().toLowerCase() === "category") {
          data[key]?.forEach(({ subcategory, title: categoryTitle }) => {
            subcategory?.forEach(({ title }) => {
              let key = `${categoryTitle} / ${title.trim()}`;
              let value = `${categoryTitle}=${title.trim()}`;
              temp.push({
                type: "subcategory",
                key: key,
                value: value,
              });
            });
          });
        } else if (key.trim().toLowerCase() === "project") {
          data[key]?.forEach(({ id, title }) => {
            let key = title.trim();
            let value = id.trim();
            temp.push({
              type: "project",
              key: key,
              value: value,
            });
          });
        }
      });
      setAllSearchTitleAble(temp);
    });
    getAllCategory().then((data) => {
      setAllCategory(data);
    });
    getSpatialLink()
      .then((data) => {
        setSpatialLink(data.result.link);
      })
      .catch((err) => console(err));
  }, []);
  return (
    <Navbar
      variant="dark"
      expand="lg"
      className={`${
        !isAboutUs ? "nav-bg-dark" : "nav-bg-about-us"
      } custom-nav-bar`}
    >
      <Container fluid={"xl"} className="d-flex">
        <Navbar.Brand
          href="/"
          style={{ zIndex: openSearchBar ? 500 : 900, order: isMobile ? 1 : 0 }}
          className="navbar-logo"
        >
          <img src={logo} className="logo" alt="logo" />
        </Navbar.Brand>
        {!isMobile && (
          <Navbar.Collapse id="navbarScroll">
            <Nav
              className="my-2 my-lg-0 w-100"
              style={{ maxHeight: "100px", alignItems: "center" }}
              navbarScroll
            >
              <Nav.Link
                href="/about"
                className={`${!isAboutUs ? "" : "nav-link-about-us"}`}
                style={{ zIndex: 800 }}
              >
                ABOUT US
              </Nav.Link>
              <Nav.Link
                onMouseEnter={() => setOpen(true)}
                onMouseLeave={handleMouseLeave}
                className={`${!isAboutUs ? "" : "nav-link-about-us"}`}
                style={{ zIndex: 800 }}
              >
                EXPLORE PROJECT
                <div id="collapse-list">
                  <Collapse in={open}>
                    <Card
                      body
                      style={{ width: "250px" }}
                      className="card-collapse"
                    >
                      <ul className="title-category">
                        {allCategory?.map(({ title, subcategory }, index) => {
                          if (subcategory?.length === 0) return <></>;
                          return (
                            <div key={index}>
                              {index > 0 &&
                                allCategory?.[index - 1]?.subcategroy?.length >
                                  0 && <hr className="mt-0 mb-md-2"></hr>}
                              <li
                                onMouseEnter={() =>
                                  handleSelectedTitle({ title })
                                }
                              >
                                <span className="text-title">
                                  {title}
                                  <IoArrowForward className="ms-auto" />
                                </span>
                              </li>
                            </div>
                          );
                        })}
                      </ul>
                    </Card>
                  </Collapse>
                  <Collapse in={openNested} dimension="width">
                    <div id="nested-collapse-list">
                      <Card
                        body
                        style={{ width: "250px" }}
                        className="card-collapse"
                      >
                        <ul>
                          {allCategory?.[
                            allCategory?.findIndex(
                              ({ title }) => title === selectedTitle
                            )
                          ]?.subcategory?.map(({ title }, index) => (
                            <li
                              key={`${title}-${index}`}
                              onClick={() => {
                                handleSelectedCategory({
                                  category: selectedTitle,
                                  subcategory: title,
                                });
                                setOpenNested(false);
                                setOpen(false);
                              }}
                              className="subcategory-title"
                            >
                              {title}
                            </li>
                          ))}
                        </ul>
                      </Card>
                    </div>
                  </Collapse>
                </div>
              </Nav.Link>
              <Nav.Link
                href={spatialLink}
                target="_blank"
                className="ms-auto"
                style={{ zIndex: 800 }}
              >
                <img src={digitalMuseum} alt="Digital Museum" />
              </Nav.Link>
            </Nav>
          </Navbar.Collapse>
        )}
        {!isMobile && (
          <>
            <InputGroup
              className={`d-flex align-items-center flex-column search-form ${
                openSearchBar && "expand"
              }`}
              style={{ height: "0px" }}
            >
              <div
                style={{
                  width: "100%",
                  maxWidth: isXXL ? "600px" : "850px",
                  position: "absolute",
                  top: 0,
                  zIndex: openSearchBar ? 800 : 700,
                }}
              >
                <div
                  style={{ width: "100%", maxWidth: "1000px", display: "flex" }}
                >
                  <Form.Control
                    type="search"
                    placeholder="Search by Project Name, Search by Technology, Search by Major"
                    className={`search-bar ${openSearchBar && "expand"}`}
                    aria-label="Search"
                    ref={inputFocus}
                    value={searchContent}
                    onChange={(e) => setSearchContent(e.target.value)}
                    onBlur={() => {
                      // if (openSearchBar && !isMouseEnter && !isClickSearch)
                      //   handleCloseSearchBar();
                    }}
                  />
                  {!openSearchBar && <div className="block-form "></div>}
                  <Button
                    className={`search-btn ${openSearchBar && "clicked"}`}
                    style={{
                      borderBottomLeftRadius: 0,
                      borderTopLeftRadius: 0,
                    }}
                    variant={`${openSearchBar ? "primary" : "default"}`}
                    onClick={() => {
                      if (!openSearchBar) handleOpenSearchBar();
                      else handleCloseSearchBar();
                    }}
                    onMouseEnter={() => setIsMouseEnter(true)}
                    onMouseLeave={() => setIsMouseEnter(false)}
                  >
                    <BsSearch
                      className={`${
                        !isAboutUs ? "search-icon" : "search-icon-about-us"
                      }`}
                    />
                  </Button>
                </div>
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
                            onMouseDown={(e) => {
                              e.stopPropagation();
                              setIsClickSearch(true);
                              handleSelectSearch(item);
                            }}
                          >
                            <div className="search-results-item">
                              {getHighlightedText(item.key, searchContent)}
                            </div>
                          </ListGroup.Item>
                        ))}
                  </ListGroup>
                </Collapse>
              </div>
            </InputGroup>
          </>
        )}
        {isMobile && (
          <Button
            style={{ order: isMobile ? 3 : 0 }}
            variant="default"
            className="search-btn"
            onClick={handleOpenSearchOffcanvas}
          >
            <BsSearch className="search-icon" />
          </Button>
        )}
        <SearchOffcanvas
          key={2}
          placement={"start"}
          name={"start"}
          show={openSearchOffcanvas}
          handleClose={handleClsoeSearchOffcanvas}
          allSearchTitleAble={allSearchTitleAble}
          handleSelectSearch={handleSelectSearch}
        />
        {renderToggleMenu()}
      </Container>
    </Navbar>
  );
};

export default BasicNavbar;
