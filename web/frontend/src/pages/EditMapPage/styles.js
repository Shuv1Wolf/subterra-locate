import styled from 'styled-components';

export const MapsPageContainer = styled.div`
  display: flex;
  height: calc(100vh - 50px);
  color: #e3eafc;
`;

export const MapSelectorContainer = styled.div`
  width: 250px;
  padding: 20px;
  background: #181c24;
  overflow-y: auto;
  transition: width 0.3s;

  &.hidden {
    width: 40px;
  }
`;

export const MapSelect = styled.select`
  width: 100%;
  padding: 10px;
  margin-bottom: 20px;
  background: #232936;
  color: #e3eafc;
  border: 1px solid #222a36;
  border-radius: 6px;
`;

export const MapInfo = styled.div`
  margin-top: 20px;
  font-size: 14px;
  display: flex;
  justify-content: space-between;
`;

export const MapWrapper = styled.div`
  flex-grow: 1;
  background: #101418;
  position: relative;
  overflow: hidden;
`;

export const PopupContainer = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 300px;
  background: #181c24;
  border: 1px solid #222a36;
  border-radius: 8px;
  z-index: 1000;
  color: #e3eafc;
`;

export const PopupHeader = styled.div`
  padding: 10px 15px;
  background: #232936;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: move;
`;

export const CloseButton = styled.button`
  background: none;
  border: none;
  color: #e3eafc;
  font-size: 24px;
  cursor: pointer;
`;

export const FilterButton = styled.button`
  width: 100%;
  padding: 10px;
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  margin-top: 20px;

  &:hover {
    background: #1565c0;
  }
`;

export const ContextMenu = styled.div`
  position: absolute;
  background: #181c24;
  border: 1px solid #222a36;
  border-radius: 4px;
  z-index: 1001;
  padding: 5px 0;
`;

export const ContextMenuItem = styled.div`
  padding: 8px 12px;
  cursor: pointer;

  &:hover {
    background: #232936;
  }
`;

export const ModalBackdrop = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 2000;
`;

export const ModalContent = styled.div`
  background: #181c24;
  padding: 20px;
  border-radius: 8px;
  width: 400px;
`;

export const BeaconList = styled.ul`
  list-style: none;
  padding: 0;
  max-height: 300px;
  overflow-y: auto;
`;

export const BeaconListItem = styled.li`
  padding: 10px;
  cursor: pointer;

  &:hover {
    background: #232936;
  }
`;

export const FilterBlock = styled.div`
  margin-bottom: 15px;
`;

export const FilterHeader = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: 8px;
`;

export const CheckboxLabel = styled.label`
  display: flex;
  align-items: center;
  cursor: pointer;
`;

export const ToggleButton = styled.button`
  position: absolute;
  top: 10px;
  right: -30px;
  background: #1976d2;
  color: #fff;
  border: none;
  padding: 5px;
  cursor: pointer;
  border-radius: 0 4px 4px 0;
`;

export const TileButton = styled.button`
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 10px 15px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: #1565c0;
  }
`;

export const PaginationContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
`;

export const HistoryControlsContainer = styled.div`
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: #181c24;
  padding: 15px;
  border-radius: 8px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 350px;
`;

export const ControlRow = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

export const StyledInput = styled.input`
  background: #232936;
  color: #e3eafc;
  border: 1px solid #222a36;
  padding: 8px;
  border-radius: 4px;
`;

export const HistoryHeader = styled.h3`
  text-align: center;
  margin: 0 0 10px 0;
  cursor: move;
`;

export const PopupButtonContainer = styled.div`
  display: flex;
  justify-content: center;
  padding: 10px;
`;
