import React, { useState } from "react";
import { Card, Carousel } from "react-bootstrap";
import { IoIosArrowDown } from "react-icons/io";
import ReactPlayer from "react-player";
import { useNavigate } from "react-router-dom";
import { appConfig } from "../../../../config";
import "./index.css";
import useSound from "use-sound";

const { STORAGE_IMAGE_URL, BUCKET_NAME, COLLECTION_NAME } = appConfig;

const ProjectCard = ({ id, images, videos, title, description, tags }) => {
  const navigate = useNavigate();
  const [isSelected, setIsSelected] = useState(false);
  const toggle = () => setIsSelected(!isSelected);
  const [playActive] = useSound(
    `${process.env.PUBLIC_URL}/audio/click-01.mp3`,
    { id: "onhover" }
  );

  const handleRedirect = () => navigate(`/project/${id}`);
  return (
    <Card
      className={`zoom card-category`}
      onMouseEnter={() => {
        toggle();
        playActive();
      }}
      onMouseLeave={toggle}
    >
      <Carousel interval={null} indicators={false} controls={isSelected}>
        {Array.isArray(images) &&
          images.map((image, index) => (
            <Carousel.Item key={index}>
              <img
                className="w-100 h-100"
                alt={`title-no-${index + 1}`}
                src={`${image}`}
              />
              <Card.ImgOverlay
                onClick={handleRedirect}
                className={"project-card-wrapper"}
              >
                <Card.Title>{title}</Card.Title>
              </Card.ImgOverlay>
            </Carousel.Item>
          ))}
        {Array.isArray(videos) &&
          videos.map((video, index) => (
            <Carousel.Item key={index}>
              <ReactPlayer
                light={true}
                url={video}
                playing={true}
                controls={false}
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

      <Card.Body className={`seemore expand`}>
        <Card.Title className={"text-truncate text-truncate--2"}>
          {title}
        </Card.Title>
        <Card.Text className={"text-truncate text-truncate--3"}>
          {description}
        </Card.Text>
        <div
          className="seemore-link"
          onClick={(e) => {
            e.stopPropagation();
            handleRedirect();
          }}
        >
          See more <IoIosArrowDown />
        </div>
      </Card.Body>
    </Card>
  );
};

export default ProjectCard;
