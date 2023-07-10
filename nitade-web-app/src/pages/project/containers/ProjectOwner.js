import React from "react";
import { HiOutlineMail } from "react-icons/hi";
import { IconContext } from "react-icons";

const ProjectOwner = ({ authors, authorEmail, projectInspiration }) => {
  return (
    <div className="project-owner p-2 pe-5 mt-5">
      <div>
        <p className="project-owner-authors-paragraph">Authors</p>
        <p>
          {authors?.map((author, index) => (
            <span
              className="project-owner-authors-paragraph-span"
              key={`${author}-${index}`}
            >
              {author}
            </span>
          ))}
        </p>
        <div className="project-owner-authors-email">
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
      </div>
      <div className="project-owner-inspiration">
        <p>Project Inspiration</p>
        <p>"{projectInspiration}"</p>
      </div>
    </div>
  );
};

export default ProjectOwner;
