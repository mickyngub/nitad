import React from "react";
import { Col, Container, Row } from "react-bootstrap";
import facebookIcon from "../../assets/Facebook_icon.svg";
import logo from "../../assets/nitade-logo.png";
import twitterIcon from "../../assets/Twitter_icon.svg";
import youtubeIcon from "../../assets/Youtube_icon.png";
import tiktokIcon from "../../assets/Tiktok_icon.png";
import "./index.css";
import gracz from "../../assets/sponsors/gracz.svg";
import bangchak from "../../assets/sponsors/bangchak.svg";
import chula_eng from "../../assets/sponsors/chula_eng.svg";
import doss from "../../assets/sponsors/doss.svg";
import mfeg from "../../assets/sponsors/mfeg.svg";
import nstda from "../../assets/sponsors/nstda.svg";
import seedwebs from "../../assets/sponsors/seedwebs.svg";

const socialMedia = [
  {
    icon: facebookIcon,
    link: "https://www.facebook.com/NITAD18TH",
    title: "NITAD 18",
  },
  {
    icon: twitterIcon,
    link: "https://twitter.com/nitad18th",
    title: "NITAD 18",
  },
  {
    icon: tiktokIcon,
    link: "https://vt.tiktok.com/ZSeXJ9thL",
    title: "nitad18th",
  },
  {
    name: "youtube",
    icon: youtubeIcon,
    link: "https://www.youtube.com/channel/UCEEycf5ZuKZMfySVz07LmrQ?sub_confirmation=1",
    title: "NITAD 18",
  },
];

const Footer = ({ isAboutUs }) => {
  return (
    <div className={`${!isAboutUs ? "footer" : "footer-about-us"}`}>
      <Container fluied={"md"}>
        <Row>
          <Col xs={12} md={4} xl={5} className={"footer-col left"}>
            <div className="w-100 text-center text-md-start">
              <h2 className="mb-4">Sponsors</h2>
            </div>
            <div className="w-100">
              <div className="d-flex flex-column flex-md-row align-items-center">
                <img src={gracz} className="sponsor" alt={`sponsors-gracz`} />
              </div>
              <div className="d-flex flex-column flex-md-row align-items-center flex-wrap">
                <img src={mfeg} className="sponsor" alt={`sponsors-mfeg`} />
                <img
                  src={chula_eng}
                  className="sponsor"
                  alt={`sponsors-chula_eng`}
                />
                <img
                  src={bangchak}
                  className="sponsor"
                  alt={`sponsors-bangchak`}
                />
              </div>
              <div className="d-flex flex-column flex-md-row align-items-center flex-wrap">
                <img src={nstda} className="sponsor" alt={`sponsors-nstda`} />
                <img src={doss} className="sponsor" alt={`sponsors-doss`} />
                <img
                  src={seedwebs}
                  className="sponsor"
                  alt={`sponsors-seedwebs`}
                />
              </div>
            </div>
          </Col>
          <Col xs={12} md={4} xl={2} className={"footer-col"}>
            <h2 className="mb-4">Contact us</h2>
            {socialMedia.map((item, index) => (
              <a
                key={index}
                className={"social-media"}
                href={item.link}
                target="_blank"
              >
                <img
                  src={item.icon}
                  alt="icon"
                  className={`social-media-icon ${
                    item?.name === "youtube" ? "youtube" : "me-2"
                  } `}
                />
                <label>{item.title}</label>
              </a>
            ))}
          </Col>
          <Col xs={12} md={4} xl={5} className={"footer-col"}>
            <a href="/">
              <img src={logo} className="footer-logo" alt="logo" />
            </a>
          </Col>
        </Row>
      </Container>
    </div>
  );
};

export default Footer;
