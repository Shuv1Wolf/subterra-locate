import styled from 'styled-components';

export const HomePageContainer = styled.div`
  min-height: calc(100vh - 50px);
  min-width: 100vw;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  padding: 32px 0;

  @media (max-width: 700px) {
    padding: 16px 0;
  }
`;

export const TilesContainer = styled.div`
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 32px;
  width: 100%;
  justify-content: center;
  align-items: center;
  max-width: 1200px;

  @media (max-width: 700px) {
    flex-direction: column;
    gap: 20px;
  }
`;

export const Tile = styled.div`
  background: #232936;
  border: 1px solid #222a36;
  border-radius: 12px;
  padding: 32px 48px;
  box-shadow: 0 2px 16px rgba(0,0,0,0.18);
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 220px;
  width: 320px;
  max-width: 400px;
  margin: 0 auto;
  flex: 1 1 220px;

  @media (max-width: 700px) {
    padding: 24px 12vw;
    width: 90vw;
  }
`;

export const TileIcon = styled.span`
  font-size: 48px;
  color: #90caf9;
  margin-bottom: 12px;

  @media (max-width: 500px) {
    font-size: 36px;
  }
`;

export const TileTitle = styled.span`
  font-weight: 600;
  font-size: 20px;
  margin-bottom: 8px;
  color: #fff;

  @media (max-width: 500px) {
    font-size: 17px;
  }
`;

export const TileDescription = styled.span`
  color: #b0b8c1;
  font-size: 14px;
  margin-bottom: 16px;
  text-align: center;

  @media (max-width: 500px) {
    font-size: 13px;
  }
`;

export const TileButton = styled.button`
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 8px 24px;
  font-size: 16px;
  cursor: pointer;
  opacity: 1;
  transition: opacity 0.2s;
  width: 100%;
  max-width: 180px;

  &:disabled {
    cursor: not-allowed;
    opacity: 0.7;
  }

  @media (max-width: 500px) {
    padding: 8px 16px;
    font-size: 15px;
  }
`;
