import React from "react";
import Banner from "./Banner";
import MostViewSection from "./MostViewSection";
import CategorySection from "./CategorySection";
import "../index.css";

const Home = () => {
  return (
    <div>
      <Banner />
      <MostViewSection />
      <CategorySection />
    </div>
  );
};

export default Home;
