import styled from "styled-components";

export const HeaderContainer = styled.header`
  display: flex;
  align-items: center;
  padding: 8px 16px;
  background-color: #1f232a;
  color: white;
  height: 50px;
  position: relative;
`;

export const MenuIcon = styled.div`
  font-size: 24px;
  margin-right: 16px;
  cursor: pointer;
`;

export const BackButton = styled.div`
  font-size: 24px;
  margin-right: 16px;
  cursor: pointer;
`;

export const Title = styled.h1`
  font-size: 20px;
  margin: 0;
  flex-grow: 1;
`;

export const AlarmsButton = styled.button`
  background-color: #4a4a4a;
  color: white;
  border: 1px solid #6a6a6a;
  padding: 6px 16px;
  margin-right: 16px;
  cursor: pointer;
  border-radius: 6px;
  font-size: 14px;
`;

export const OrgName = styled.div`
  font-size: 16px;
  font-weight: bold;
`;

export const DropdownMenu = styled.div`
  position: absolute;
  top: 50px;
  left: 0;
  background-color: #1f232a;
  border: 1px solid #6a6a6a;
  border-radius: 0 0 6px 6px;
  width: 200px;
  z-index: 10;
`;

export const MenuItem = styled.div`
  padding: 12px 16px;
  cursor: pointer;
  &:hover {
    background-color: #4a4a4a;
  }
`;
