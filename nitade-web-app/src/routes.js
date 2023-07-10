import React from "react";
import { Route, Routes as BasicRoutes } from "react-router-dom";
import Layout from "./components/Layout";
import Category from "./pages/category/index";
import Home from "./pages/home/index";
import Project from "./pages/project/index";
import AboutUs from "./pages/aboutUs/index";

const Routes = () => {
  return (
    <Layout
      children={
        <BasicRoutes>
          <Route path="/" element={<Home />} />
          <Route path="/category" element={<Category />} />
          <Route path="/project/:id" element={<Project />} />
          <Route path="/about" element={<AboutUs />} />
        </BasicRoutes>
      }
    />
  );
};

export default Routes;
