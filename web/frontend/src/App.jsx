import React, { useEffect } from "react";
import {
  BrowserRouter,
  Routes,
  Route,
  useNavigate,
  useLocation,
} from "react-router-dom";
import HomePage from "./pages/HomePage";
import MapsPage from "./pages/NewMapPage";
import CreateMapPage from "./pages/CreateMapPage";
import MapFormPage from "./pages/MapFormPage";
import BeaconsAdminPage from "./pages/BeaconsAdminPage/index.jsx";
import BeaconFormPage from "./pages/BeaconFormPage";
import DevicesAdminPage from "./pages/DevicesAdminPage";
import DeviceFormPage from "./pages/DeviceFormPage";
import MapsAdminPage from "./pages/MapsAdminPage";

const AppContent = () => {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/maps" element={<MapsPage />} />
      <Route path="/beacons-admin" element={<BeaconsAdminPage />} />
      <Route path="/beacons-admin/new" element={<BeaconFormPage />} />
      <Route
        path="/beacons-admin/edit/:beaconId"
        element={<BeaconFormPage />}
      />
      <Route path="/devices-admin" element={<DevicesAdminPage />} />
      <Route path="/devices-admin/new" element={<DeviceFormPage />} />
      <Route
        path="/devices-admin/edit/:deviceId"
        element={<DeviceFormPage />}
      />
      <Route path="/maps-admin" element={<MapsAdminPage />} />
      <Route path="/new-map" element={<CreateMapPage />} />
      <Route path="/edit-map/:mapId" element={<MapFormPage />} />
    </Routes>
  );
};

function App() {
  return (
    <BrowserRouter>
      <AppContent />
    </BrowserRouter>
  );
}

export default App;
