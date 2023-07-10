import React, { useMemo, useState } from "react";
import { Container } from "react-bootstrap";
import { getAllCategory } from "../../../services/category/category";
import Category from "../components/Catergory/Category";

const CategorySection = () => {
  const [allCategory, setAllCategory] = useState([]);
  const refs = useMemo(
    () => allCategory?.map((items) => ({ ref: React.createRef() })),
    [allCategory]
  );

  const handleClickSubcategory = ({ categoryTitle, subcategoryTitle }) =>
    (window.location.href = `/category?${categoryTitle}=${subcategoryTitle}`);

  React.useEffect(() => {
    getAllCategory().then((data) => setAllCategory(data));
  }, []);


  return allCategory.length > 0 ? (
    <div className="category-wrapper-section">
      <div className="category-wrapper-section-container">
        <Container fluid={"xl"}>
          {allCategory
            ?.slice(0, Math.min(5, allCategory.length))
            ?.map(
              ({ title, subcategory }, index) =>
                subcategory?.length > 0 && (
                  <Category
                    title={title}
                    subcategory={subcategory}
                    ref={refs?.[index]?.ref}
                    key={index}
                    index={index}
                    handleClickSubcategory={handleClickSubcategory}
                  />
                )
            )}
        </Container>
      </div>
    </div>
  ) : (
    <div className="text-center">
      <h1>There is no category</h1>
    </div>
  );
};

export default CategorySection;
