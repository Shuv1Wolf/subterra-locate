import styled from 'styled-components';

export const MapsPageContainer = styled.div`
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  position: relative;
`;

export const MapSelectorContainer = styled.div`
  position: absolute;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(35, 41, 54, 0.9);
  border-radius: 12px;
  padding: 12px 24px;
  z-index: 10;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
  border: 1px solid #222a36;

  @media (max-width: 600px) {
    flex-direction: column;
    gap: 10px;
    padding: 10px;
  }
`;

export const MapSelect = styled.select`
  background: #181c24;
  color: #e3eafc;
  border: 1px solid #222a36;
  border-radius: 6px;
  padding: 8px 12px;
  font-size: 16px;
  cursor: pointer;

  &:focus {
    outline: none;
    border-color: #90caf9;
  }
`;

export const MapInfo = styled.div`
  color: #b0b8c1;
  font-size: 14px;
  display: flex;
  gap: 15px;
`;

export const MapWrapper = styled.div`
  width: 100%;
  height: 100%;
  
  .react-transform-wrapper {
    width: 100%;
    height: 100%;
  }

  .react-transform-component {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  svg {
    max-width: 100%;
    max-height: 100%;
  }
`;

export const BackArrow = styled.div`
  position: absolute;
  top: 20px;
  left: 20px;
  z-index: 10;
  cursor: pointer;
  font-size: 24px;
  color: #90caf9;
  transition: color 0.2s;

  &:hover {
    color: #fff;
  }
`;

export const FilterButton = styled.button`
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 8px 16px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: #1565c0;
  }
`;

export const PopupContainer = styled.div`
  position: absolute;
  top: 100px; /* Initial position, will be updated by Draggable */
  left: 100px; /* Initial position, will be updated by Draggable */
  background: #232936;
  border-radius: 12px;
  box-shadow: 0 5px 25px rgba(0, 0, 0, 0.5);
  width: 300px;
  color: #e3eafc;
  border: 1px solid #222a36;
  z-index: 100;

  p {
    margin: 8px 12px;
    line-height: 1.6;
    font-size: 14px;
  }

  strong {
    color: #b0b8c1;
  }
`;

export const PopupHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: #181c24;
  border-top-left-radius: 12px;
  border-top-right-radius: 12px;
  cursor: move;

  h2 {
    margin: 0;
    font-size: 16px;
    color: #90caf9;
  }
`;

export const CloseButton = styled.button`
  background: none;
  border: none;
  color: #b0b8c1;
  font-size: 20px;
  cursor: pointer;
  padding: 0;
  line-height: 1;

  &:hover {
    color: #fff;
  }
`;
