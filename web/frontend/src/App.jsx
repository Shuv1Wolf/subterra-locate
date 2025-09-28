import React from "react";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import MapsPage from "./pages/NewMapPage";
import BeaconsAdminPage from "./pages/BeaconsAdminPage/index.jsx";
import BeaconFormPage from "./pages/BeaconFormPage";
import DevicesAdminPage from "./pages/DevicesAdminPage";
import DeviceFormPage from "./pages/DeviceFormPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/maps" element={<MapsPage />} />
        <Route path="/beacons-admin" element={<BeaconsAdminPage />} />
        <Route path="/beacons-admin/new" element={<BeaconFormPage />} />
        <Route path="/beacons-admin/edit/:beaconId" element={<BeaconFormPage />} />
        <Route path="/devices-admin" element={<DevicesAdminPage />} />
        <Route path="/devices-admin/new" element={<DeviceFormPage />} />
        <Route path="/devices-admin/edit/:deviceId" element={<DeviceFormPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
