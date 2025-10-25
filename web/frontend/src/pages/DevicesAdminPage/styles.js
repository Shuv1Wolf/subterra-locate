import styled from 'styled-components';

export const AdminPageContainer = styled.div`
  padding: 24px;
  color: #e3eafc;
  min-height: calc(100vh - 50px);
`;

export const Button = styled.button`
  background: #1976d2;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: #1565c0;
  }
`;

export const BeaconTable = styled.table`
  width: 100%;
  border-collapse: collapse;
  
  th, td {
    border: 1px solid #222a36;
    padding: 12px;
    text-align: left;
  }

  th {
    background: #232936;
  }

  tr:nth-child(even) {
    background: #1f2430;
  }
`;

export const ActionsContainer = styled.div`
  display: flex;
  gap: 10px;
`;
