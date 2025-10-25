import React, { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import {
  HeaderContainer,
  MenuIcon,
  AlarmsButton,
  DropdownMenu,
  MenuItem,
  BackButton,
  Title,
  OrgName,
} from "./styles";

export default function Header({ variant = 'home', title = 'Subterra Locate' }) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const navigate = useNavigate();
  const headerRef = useRef(null);

  const handleNavigate = (path) => {
    navigate(path);
    setIsMenuOpen(false);
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (headerRef.current && !headerRef.current.contains(event.target)) {
        setIsMenuOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [headerRef]);

  return (
    <HeaderContainer ref={headerRef}>
      {variant === 'home' ? (
        <MenuIcon onClick={() => setIsMenuOpen(!isMenuOpen)}>&#9776;</MenuIcon>
      ) : (
        <BackButton onClick={() => navigate(-1)}>&#8592;</BackButton>
      )}
      <Title>{title}</Title>
      {variant === 'home' && <AlarmsButton>Alarms</AlarmsButton>}
      <OrgName>Test organization</OrgName>
      {isMenuOpen && (
        <DropdownMenu>
          <MenuItem onClick={() => handleNavigate("/maps")}>Maps</MenuItem>
          <MenuItem onClick={() => handleNavigate("/devices-admin")}>
            Devices Admin
          </MenuItem>
          <MenuItem onClick={() => handleNavigate("/beacons-admin")}>
            Beacons Admin
          </MenuItem>
        </DropdownMenu>
      )}
    </HeaderContainer>
  );
}
