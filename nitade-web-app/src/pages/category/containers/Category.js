import React, { useEffect, useMemo, useState } from "react";
import { Container, Pagination } from "react-bootstrap";
import { useLocation, useNavigate } from "react-router-dom";
import Tag from "../../../components/Tag";
import { getAllCategory } from "../../../services/category/category";
import { getAllProjectBySubcategory } from "../../../services/project/project";
import CardGroup from "../components/CardGroup";
import Filter from "../components/Filter/Filter";
import "../index.css";
import { sortedType } from "../components/Filter/utils";

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

const getQueryString = (queryStringObject) => {
  let s = "";
  let entries = Object.entries(queryStringObject);
  entries.forEach(([key, value], index) => {
    if (index > 0) s = s + "&";
    s = s + key + "=" + value;
  });
  return s;
};

const Category = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [allCategory, setAllCategory] = useState([]);
  const [allSubcategory, setAllSubcategory] = useState([]);
  const [allProject, setAllProject] = useState([]);

  const [totalPage, setTotalPage] = useState(1);
  const queryStringObject = useMemo(() => {
    return getParams();
  }, [location.search]);

  const currentPage = useMemo(
    () =>
      isNaN(queryStringObject?.page) || Number(queryStringObject?.page) <= 0
        ? 1
        : Number(queryStringObject?.page),
    [queryStringObject?.page]
  );

  const limitPage = useMemo(
    () =>
      isNaN(queryStringObject?.limit) || Number(queryStringObject?.limit) <= 0
        ? 15
        : Number(queryStringObject?.limit),
    [queryStringObject?.limit]
  );

  const sortedBy = useMemo(
    () =>
      queryStringObject?.sort
        ? sortedType[queryStringObject?.sort]
        : sortedType[0],
    [queryStringObject]
  );

  const title = useMemo(() => {
    let s = "Results";
    let keys = Object.keys(queryStringObject).filter(
      (key) =>
        key.trim().toLocaleLowerCase() !== "sort" &&
        key.trim().toLocaleLowerCase() !== "page" &&
        key.trim().toLocaleLowerCase() !== "limit"
    );
    let entries = Object.entries(queryStringObject);

    if (keys.length === 1 && entries[0][1] !== "") {
      if (queryStringObject[keys[0]]?.includes(",")) {
        return s;
      } else return queryStringObject[keys[0]];
    }
    return s;
  }, [queryStringObject]);

  const handleSorted = (index) => {
    let temp = { ...queryStringObject };
    temp.sort = index;
    navigate(`?${getQueryString(temp)}`);
  };

  const handleChangePage = (pageNo) => {
    let temp = { ...queryStringObject };
    temp.page = pageNo;
    window.location.href = `category?${getQueryString(temp)}`;
  };
  const handleChangeFilter = ({ category, subcategory }) => {
    let temp = { ...queryStringObject };
    if (!temp[category]) {
      temp[category] = `${subcategory?.trim()}`;
    } else if (temp[category]?.includes(subcategory)) {
      let s = temp[category].split(",");
      const index = s.indexOf(subcategory);
      if (index > -1) {
        s.splice(index, 1);
      }
      if (s.length === 0) {
        delete temp[category];
      } else temp[category] = s.join(",");
    } else {
      temp[category] += `,${subcategory?.trim()}`;
    }
    delete temp.page;
    navigate(`?${getQueryString(temp)}`);
  };

  useEffect(() => {
    getAllCategory().then((data) => {
      setAllCategory(data);
      let temp = [];
      data?.forEach(({ subcategory }) => {
        subcategory?.forEach((item) => temp.push(item));
      });
      setAllSubcategory(temp);
    });
  }, []);

  const onCloseTag = ({ categoryTitle, subcategoryTitle }) => {
    handleChangeFilter({
      category: categoryTitle,
      subcategory: subcategoryTitle,
    });
  };

  useEffect(() => {
    if (allSubcategory.length > 0 && queryStringObject) {
      let value = Object.values(queryStringObject);
      let subcategoryList = [];
      value?.forEach((items) => {
        let s = items.split(",");
        s?.forEach((t) => {
          let temp = allSubcategory?.filter(
            ({ title }) => title.trim() === t.trim()
          );
          if (temp.length === 1) {
            subcategoryList?.push(temp[0].id);
          }
        });
      });

      let options = {
        page: currentPage,
        limit: limitPage,
      };
      getAllProjectBySubcategory({ subcategoryList, options })
        .then((data) => {
          if (data?.paginate?.totalPage >= 0) {
            if (data?.paginate?.totalPage !== totalPage)
              setTotalPage(
                (data.paginate.totalPage * data.paginate.limit) %
                  data.result.length ===
                  0
                  ? data.paginate.totalPage - 1
                  : data.paginate.totalPage
              );
          }
          setAllProject(data.result);
        })
        .catch((err) => {
          console.log(err);
          setAllProject([]);
          setTotalPage(0);
        });
    }
  }, [queryStringObject, allSubcategory, currentPage, limitPage, totalPage]);

  return (
    <div className="category-image-container">
      <div className="category-section-container">
        <Container fluid={"xl"}>
          <div className="category-section-container-inner">
            <div className="d-flex justify-content-between">
              <div
                className="w-25 d-none d-lg-block"
                style={{ marginRight: "32px" }}
              >
                <Filter
                  filters={allCategory}
                  queryStringObject={queryStringObject}
                  handleChangeFilter={handleChangeFilter}
                  handleSorted={handleSorted}
                />
              </div>
              <div className="w-100 category-group">
                <h1>{title}</h1>
                <div className="d-flex flex-wrap align-items-center mb-3">
                  <div className="breadcrumb-link d-none d-lg-block">
                    <a href="/">Home</a>
                  </div>
                  <span className="mx-2 d-none d-lg-block">/</span>
                  {Object.entries(queryStringObject)
                    .filter(
                      ([key]) =>
                        key.trim().toLocaleLowerCase() !== "sort" &&
                        key.trim().toLocaleLowerCase() !== "page" &&
                        key.trim().toLocaleLowerCase() !== "limit"
                    )
                    .map(([key, value], index) =>
                      value.split(",").map((subcategoryTitle) => (
                        <Tag
                          key={`${key}-${index}`}
                          subcategoryTitle={subcategoryTitle}
                          categoryTitle={key}
                          variant={"closeable"}
                          onCloseTag={() => {
                            onCloseTag({
                              categoryTitle: key,
                              subcategoryTitle: subcategoryTitle,
                            });
                          }}
                        />
                      ))
                    )}
                </div>
                <CardGroup cards={allProject} sortedBy={sortedBy} />
                {allProject.length > 0 && (
                  <Pagination
                    className="custom-pagination-wrapper mt-2 mt-lg-4"
                    size={"md"}
                  >
                    {currentPage >= 2 && (
                      <Pagination.First onClick={() => handleChangePage(1)} />
                    )}
                    {currentPage > 1 && (
                      <Pagination.Prev
                        onClick={() => handleChangePage(currentPage - 1)}
                      />
                    )}
                    {currentPage - 3 > 0 && (
                      <>
                        <Pagination.Item onClick={() => handleChangePage(1)}>
                          {1}
                        </Pagination.Item>
                        {currentPage - 3 > 1 && (
                          <Pagination.Ellipsis disabled />
                        )}
                      </>
                    )}
                    {currentPage - 2 > 0 && (
                      <Pagination.Item
                        onClick={() => handleChangePage(currentPage - 2)}
                      >
                        {currentPage - 2}
                      </Pagination.Item>
                    )}
                    {currentPage - 1 > 0 && (
                      <Pagination.Item
                        onClick={() => handleChangePage(currentPage - 1)}
                      >
                        {currentPage - 1}
                      </Pagination.Item>
                    )}
                    <Pagination.Item active>{currentPage}</Pagination.Item>
                    {currentPage + 1 <= totalPage && (
                      <Pagination.Item
                        onClick={() => handleChangePage(currentPage + 1)}
                      >
                        {currentPage + 1}
                      </Pagination.Item>
                    )}
                    {currentPage + 2 <= totalPage && (
                      <Pagination.Item
                        onClick={() => handleChangePage(currentPage + 2)}
                      >
                        {currentPage + 2}
                      </Pagination.Item>
                    )}
                    {currentPage + 3 <= totalPage && (
                      <>
                        {currentPage + 3 <= totalPage - 1 && (
                          <Pagination.Ellipsis disabled />
                        )}
                        <Pagination.Item
                          onClick={() => handleChangePage(totalPage)}
                        >
                          {totalPage}
                        </Pagination.Item>
                      </>
                    )}
                    {currentPage < totalPage && (
                      <Pagination.Next
                        onClick={() => handleChangePage(currentPage + 1)}
                      />
                    )}
                    {currentPage <= totalPage - 1 && (
                      <Pagination.Last
                        onClick={() => handleChangePage(totalPage)}
                      />
                    )}
                  </Pagination>
                )}
              </div>
            </div>
          </div>
        </Container>
      </div>
    </div>
  );
};

export default Category;
