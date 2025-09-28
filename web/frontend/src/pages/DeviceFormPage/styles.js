import styled from 'styled-components';

export const FormContainer = styled.div`
  padding: 24px;
  max-width: 600px;
  margin: 0 auto;
`;

export const FormGroup = styled.div`
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