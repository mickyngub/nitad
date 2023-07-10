import React, { useMemo } from "react";
import { useMediaQuery } from "react-responsive";
import { useLocation } from "react-router-dom";
import useSound from "use-sound";
import Footer from "../Footer/Footer";
import BasicNavbar from "../Navbar/Navbar";
import "./index.css";

const Layout = ({ children }) => {
  const location = useLocation();
  const path = useMemo(() => location.pathname, [location.pathname]);
  const [playActive] = useSound(
    `${process.env.PUBLIC_URL}/audio/click-01.mp3`,
    { id: "onclick" }
  );

  const isMobileSize = useMediaQuery({ query: `(max-width: 768px)` });

  const renderBackgroundImageClasses = () => {
    if (isMobileSize) return "mobile-background";
    else if (path.includes("/project")) {
      return "project-background";
    } else if (path.includes("/category")) {
      return "category-background";
    } else if (path.includes("/about")) {
      return "about-background";
    } else {
      return "home-background";
    }
  };

  return (
    <div
      className={`${renderBackgroundImageClasses()} layout-wrapper`}
      onClick={playActive}
    >
      <BasicNavbar isAboutUs={path.includes("/about")} />
      {children}
      <Footer isAboutUs={path.includes("/about")} />
    </div>
  );
};

export default Layout;
