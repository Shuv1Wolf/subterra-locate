import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import MapsPage from "./pages/NewMapPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/maps" element={<MapsPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
