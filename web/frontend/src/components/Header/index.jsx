import React, { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import styled from "styled-components";
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

const OrgDropdownMenu = styled(DropdownMenu)`
  right: 0;
  left: auto;
`;

export default function Header({ variant = "home", title = "Subterra Locate" }) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isOrgMenuOpen, setIsOrgMenuOpen] = useState(false);
  const [organizations, setOrganizations] = useState([]);
  const [selectedOrg, setSelectedOrg] = useState(null);
  const navigate = useNavigate();
  const headerRef = useRef(null);

  useEffect(() => {
    fetch("/organizations.json")
      .then((response) => response.json())
      .then((data) => {
        setOrganizations(data);
        const storedOrgId = localStorage.getItem("selectedOrgId");
        if (storedOrgId) {
          const foundOrg = data.find((org) => org.orgId === storedOrgId);
          setSelectedOrg(foundOrg || data[0]);
        } else {
          setSelectedOrg(data[0]);
          localStorage.setItem("selectedOrgId", data[0].orgId);
        }
      });
  }, []);

  const handleNavigate = (path) => {
    navigate(path);
    setIsMenuOpen(false);
  };

  const handleOrgSelect = (org) => {
    setSelectedOrg(org);
    localStorage.setItem("selectedOrgId", org.orgId);
    setIsOrgMenuOpen(false);
  };

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (headerRef.current && !headerRef.current.contains(event.target)) {
        setIsMenuOpen(false);
        setIsOrgMenuOpen(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [headerRef]);

  return (
    <HeaderContainer ref={headerRef}>
      {variant === "home" ? (
        <MenuIcon onClick={() => setIsMenuOpen(!isMenuOpen)}>&#9776;</MenuIcon>
      ) : (
        <BackButton onClick={() => navigate(-1)}>&#8592;</BackButton>
      )}
      <Title>{title}</Title>
      {variant === "home" && <AlarmsButton>Alarms</AlarmsButton>}
      <OrgName onClick={() => setIsOrgMenuOpen(!isOrgMenuOpen)}>
        {selectedOrg ? selectedOrg.name : "Select Organization"}
      </OrgName>
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
      {isOrgMenuOpen && (
        <OrgDropdownMenu>
          {organizations.map((org) => (
            <MenuItem key={org.orgId} onClick={() => handleOrgSelect(org)}>
              {org.name}
            </MenuItem>
          ))}
        </OrgDropdownMenu>
      )}
    </HeaderContainer>
  );
}
