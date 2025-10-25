import styled from 'styled-components';

export const FormContainer = styled.div`
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
`;

export const FormBlock = styled.div`
  background: #232936;
  border: 1px solid #222a36;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
`;

export const BlockTitle = styled.h2`
  color: #90caf9;
  margin-top: 0;
`;

export const FormGroup = styled.div`
  display: grid;
  grid-template-columns: 180px 1fr;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
`;

export const Label = styled.label`
  display: block;
  margin-bottom: 8px;
  color: #b0b8c1;
`;

export const Input = styled.input`
  width: 100%;
  padding: 10px;
  padding-right: 20px;
  box-sizing: border-box;
  background: #181c24;
  border: 1px solid #222a36;
  border-radius: 6px;
  color: #e3eafc;
  font-size: 16px;

  &:focus {
    outline: none;
    border-color: #90caf9;
  }
`;

export const Select = styled.select`
  width: 100%;
  padding: 10px;
  background: #181c24;
  border: 1px solid #222a36;
  border-radius: 6px;
  color: #e3eafc;
  font-size: 16px;

  &:focus {
    outline: none;
    border-color: #90caf9;
  }
`;

export const CheckboxContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 10px;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 24px;
`;
