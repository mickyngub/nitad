import React, { useEffect, useMemo, useRef, useState } from "react";
import { Button, Carousel } from "react-bootstrap";
import { HiOutlineMail } from "react-icons/hi";
import { IoMdPerson } from "react-icons/io";
import { IconContext } from "react-icons/lib";
import ReactPlayer from "react-player";
import Tag from "../../../components/Tag";
import { appConfig } from "../../../config";
import Content from "../components/Content/Content";
import spatialLogo from "../../../assets/spatial-logo.svg";
import pdfIcon from "../../../assets/pdf-icon.svg";
import spatialLogoDark from "../../../assets/spatial-logo-dark.svg";
import pdfIconDark from "../../../assets/pdf-icon-dark.svg";

const { STORAGE_IMAGE_URL, BUCKET_NAME, COLLECTION_NAME } = appConfig;

const DownloadButton = ({ title, src, alt, download, onAction, ...props }) => {
  return (
    <Button {...props} onClick={onAction}>
      <span>
        <img src={src} alt={alt} />
      </span>
      {title}
    </Button>
  );
};

const ProjectContent = ({
  title,
  videos,
  images,
  abstract,
  description,
  tags,
  projectType,
  views,
  authors,
  authorEmail,
  projectInspiration,
  reportDownloadLink,
  virtualLink,
  ...props
}) => {
  const [isReadmore, setIsReadmore] = useState(false);
  const [contentHeight, setContentHeight] = useState(0);
  const ref = useRef(null);
  const toggleReadmore = () => setIsReadmore(!isReadmore);
  const onClickDownloadLink = (link) => window.open(link);

  useEffect(() => {
    if (ref && ref?.current) {
      setContentHeight(ref?.current?.clientHeight);
    }
  }, [description, abstract]);

  return (
    <div className="d-flex flex-column project-content-container">
      <h1 className="project-title">{title}</h1>
      <div className="d-none flex-wrap align-items-center mb-3 d-lg-flex project-breadcrumb">
        <div className="breadcrumb-link d-none d-lg-block">
          <a href="/">Home</a>
        </div>
        <span className="mx-2 d-none d-lg-block">/</span>
        {tags?.map(({ subcategory, title: categoryTitle }, index) => {
          return subcategory?.map(({ title: subcategoryTitle }) => (
            <div key={`${categoryTitle}-${subcategoryTitle}-${index}`}>
              <Tag
                categoryTitle={categoryTitle}
                subcategoryTitle={subcategoryTitle}
              />
              {index % 3 === 0 && index > 0 && (
                <div className="break-item d-none d-lg-block"></div>
              )}
            </div>
          ));
        })}
      </div>
      <div className={"custom-carousel-wrapper project-carousel mb-4"}>
        <Carousel interval={null} indicators={false} controls={true}>
          {Array.isArray(images) &&
            images.map((image, index) => (
              <Carousel.Item key={index}>
                <img
                  className="h-100"
                  alt={`title-no-${index + 1}`}
                  src={`${image}`}
                />
              </Carousel.Item>
            ))}
          {Array.isArray(videos) &&
            videos.map((video, index) => (
              <Carousel.Item key={index}>
                <ReactPlayer
                  url={video}
                  controls={true}
                  width={"100%"}
                  height={"100%"}
                  config={{
                    youtube: {
                      playerVars: { origin: "https://www.youtube.com" },
                    },
                  }}
                />
              </Carousel.Item>
            ))}
        </Carousel>
      </div>
      <div className="d-flex justify-content-between project-wrapper-under-carousel">
        <div className="project-tag-and-view">
          <div className="mb-3">
            <label className="view-amount d-flex align-items-center">
              <IoMdPerson className="me-1" />
              {views} people are viewing this project
            </label>
          </div>
          <div className="d-flex" style={{ flexWrap: "wrap" }}>
            {tags?.map(({ subcategory, title: categoryTitle }, index) => {
              return subcategory?.map(({ title: subcategoryTitle }) => (
                <div key={`${categoryTitle}-${subcategoryTitle}-${index}`}>
                  <Tag
                    categoryTitle={categoryTitle}
                    subcategoryTitle={subcategoryTitle}
                  />
                  {index % 3 === 0 && index > 0 && (
                    <div className="break-item d-none d-lg-block"></div>
                  )}
                </div>
              ));
            })}
          </div>
        </div>
        <div className="d-flex d-lg-none w-100 project-download-on-mobile-device my-3">
          <div className="btn-download-slide">
            <DownloadButton
              title={"Virtual Presentation"}
              onAction={() => onClickDownloadLink(virtualLink)}
              variant="link"
              src={spatialLogoDark}
              alt="spatial-logo"
            />
          </div>
          <div className="btn-download-slide ">
            <DownloadButton
              title={"Download Full Report"}
              onAction={() => onClickDownloadLink(reportDownloadLink)}
              variant="link"
              src={pdfIconDark}
              alt="pdf-icon"
            />
          </div>
        </div>
        <div className="btn-download-slide d-none d-lg-block">
          <DownloadButton
            title={"Virtual Presentation"}
            onAction={() => onClickDownloadLink(virtualLink)}
            variant="link"
            src={spatialLogo}
            alt="spatial-logo"
          />
        </div>
      </div>
      <div className="d-lg-none order-4">
        <h3 className="project-owner-authors-paragraph">Authors</h3>
        <p className="mb-1">
          {authors?.map((author, index) => (
            <span
              className="project-owner-authors-paragraph-span"
              key={`${author}-${index}`}
            >
              {author}
            </span>
          ))}
        </p>
        <div className="project-owner-authors-email mb-2">
          <IconContext.Provider value={{ size: "2rem" }}>
            <HiOutlineMail />
          </IconContext.Provider>
          <p>
            {authorEmail?.map(
              (email, index) =>
                `${email}${index < authorEmail.length - 1 ? ", " : ""}`
            )}
          </p>
        </div>
        <div className="project-owner-inspiration text-align-left">
          <h3>Project Inspiration</h3>
          <p>"{projectInspiration}"</p>
        </div>
      </div>
      <div className={`readmore ${isReadmore && "expand"} project-content`}>
        <div ref={ref}>
          <Content title={"Description"} content={description} />
          <Content
            title={"Abstract"}
            content={abstract}
            button={
              <div className="btn-download-slide d-none d-lg-block ms-auto">
                <DownloadButton
                  title={"Download Full Report"}
                  onAction={() => onClickDownloadLink(reportDownloadLink)}
                  variant="link"
                  src={pdfIcon}
                  alt="pdf-logo"
                />
              </div>
            }
          />
        </div>
        {contentHeight > 250 && <span className="readmore-link"></span>}
      </div>
      <Button
        variant="link"
        className={"readmore-btn"}
        onClick={toggleReadmore}
      >
        {isReadmore ? "Read less" : "Read more"}
      </Button>
    </div>
  );
};

export default ProjectContent;
