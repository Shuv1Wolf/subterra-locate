import React from "react";
import { useNavigate } from "react-router-dom";
import {
  HomePageContainer,
  Title,
  Description,
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
    <HomePageContainer>
      <Title>SubterraLocate</Title>
      <Description>
        SubterraLocate is a personnel positioning system for underground mines using BLE beacons.
        Supports motion tracking, historical data collection and route visualization for improved safety and operational control.
      </Description>
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
          <TileTitle>Devices</TileTitle>
          <TileDescription>Coming soon</TileDescription>
          <TileButton disabled>Open</TileButton>
        </Tile>
        {/* Beacons tile */}
        <Tile>
          <TileIcon>📡</TileIcon>
          <TileTitle>Beacons</TileTitle>
          <TileDescription>Coming soon</TileDescription>
          <TileButton disabled>Open</TileButton>
        </Tile>
      </TilesContainer>
    </HomePageContainer>
  );
}
