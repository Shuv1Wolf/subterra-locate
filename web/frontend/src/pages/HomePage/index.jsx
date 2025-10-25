import React from "react";
import { useNavigate } from "react-router-dom";
import Header from "../../components/Header";
import {
  HomePageContainer,
  TilesContainer,
  Tile,
  TileIcon,
  TileTitle,
  TileDescription,
  TileButton,
} from "./styles.js";

export default function HomePage() {
  const navigate = useNavigate();

  return (
    <>
      <Header />
      <HomePageContainer>
        <TilesContainer>
          {/* Maps tile */}
          <Tile>
            <TileIcon>🗺️</TileIcon>
            <TileTitle>Maps</TileTitle>
            <TileDescription>2D map of the mine</TileDescription>
            <TileButton onClick={() => navigate("/maps")}>Open</TileButton>
          </Tile>
          {/* Devices tile */}
          <Tile>
            <TileIcon>💻</TileIcon>
            <TileTitle>Devices Admin</TileTitle>
            <TileDescription>Manage devices</TileDescription>
            <TileButton onClick={() => navigate("/devices-admin")}>
              Open
            </TileButton>
          </Tile>
          {/* Beacons Admin tile */}
          <Tile>
            <TileIcon>⚙️</TileIcon>
            <TileTitle>Beacons Admin</TileTitle>
            <TileDescription>Manage beacons</TileDescription>
            <TileButton onClick={() => navigate("/beacons-admin")}>
              Open
            </TileButton>
          </Tile>
        </TilesContainer>
      </HomePageContainer>
    </>
  );
}
