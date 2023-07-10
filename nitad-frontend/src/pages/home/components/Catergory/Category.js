import React, { useMemo, useState } from "react";
import { useMediaQuery } from "react-responsive";
import Slider from "react-slick";
import { appConfig } from "../../../../config";
import "./index.css";

const { STORAGE_IMAGE_URL, BUCKET_NAME, COLLECTION_NAME } = appConfig;

const Dot = (dots) => (
  <div>
    <ul className="dot-custom">{dots}</ul>
  </div>
);

const Pagination = () => <div className="d-flex custom-pagination" />;

const Category = React.forwardRef((props, ref) => {
  const { title, subcategory, index, handleClickSubcategory } = props;
  const [isDrag, setIsDrag] = useState(true);
  const isMobileSize = useMediaQuery({ query: `(max-width: 768px)` });
  const isMobileDevice = useMemo(() => {
    if (
      /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
        navigator.userAgent
      )
    ) {
      return true;
    }
    return false;
  }, [isMobileSize]);
  const settings = useMemo(() => {
    return {
      dots: true,
      infinite: false,
      speed: 1000,
      slidesToShow: 4,
      slidesToScroll: 4,
      arrows: false,
      appendDots: Dot,
      customPaging: Pagination,
      responsive: [
        {
          breakpoint: 992,
          settings: {
            slidesToShow: 4,
            slidesToScroll: 1,
            dots: true,
            appendDots: Dot,
            customPaging: Pagination,
          },
        },
        {
          breakpoint: 768,
          settings: {
            slidesToShow: 2,
            slidesToScroll: 1,
            dots: true,
            appendDots: Dot,
            customPaging: Pagination,
          },
        },
      ],
    };
  }, [isMobileSize]);

  const shouldRenderCustom = useMemo(() => {
    return isMobileDevice;
  }, [isMobileDevice]);

  const shouldRenderSlickSlider = useMemo(() => {
    return !isMobileDevice;
  }, [isMobileDevice]);

  return (
    <div className="category">
      <label className="title">{title}</label>
      <div className="category-slider">
        <div className={"category-slider-wrapper"}>
          {shouldRenderSlickSlider && (
            <Slider {...settings}>
              {subcategory?.map(({ title: subcategoryTitle, image }, index) => (
                <div
                  className="p-md-2 image-wrapper-outter"
                  key={title + index}
                  onMouseDown={() => setIsDrag(false)}
                  onMouseMove={() => !isDrag && setIsDrag(true)}
                  onClick={() => {
                    if (!isDrag) {
                      handleClickSubcategory({
                        categoryTitle: title,
                        subcategoryTitle: subcategoryTitle,
                      });
                    }
                  }}
                >
                  <div className="image-wrapper">
                    <h3 className="subcategory-card-title">
                      {subcategoryTitle}
                    </h3>
                    <img
                      src={`${image}`}
                      alt={title}
                      className="category-img"
                    />
                  </div>
                </div>
              ))}
            </Slider>
          )}

          {shouldRenderCustom && (
            <>
              <div
                className="custom-slider-wrapper"
                id={`category-wrapper-${index}`}
              >
                <div
                  className="custom-slider"
                  id={`category-${index}`}
                  style={{ width: `${subcategory?.length * 250 + 10}px` }}
                >
                  {subcategory?.map(
                    ({ title: subcategoryTitle, image }, index) => (
                      <div
                        className="p-md-2 image-wrapper-outter"
                        key={index}
                        onMouseDown={() => setIsDrag(false)}
                        onMouseMove={() => !isDrag && setIsDrag(true)}
                        onClick={() => {
                          if (!isDrag) {
                            handleClickSubcategory({
                              categoryTitle: title,
                              subcategoryTitle: subcategoryTitle,
                            });
                          }
                        }}
                      >
                        <div className="image-wrapper">
                          <h5 className="subcategory-card-title">
                            {subcategoryTitle}
                          </h5>
                          <img
                            src={`${image}`}
                            alt={subcategoryTitle}
                            className="category-img"
                          />
                        </div>
                      </div>
                    )
                  )}
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
});

export default Category;
