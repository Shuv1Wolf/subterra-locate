import styled from 'styled-components';

export const FormContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
  margin: 0 auto;
`;

export const FormBlock = styled.div`
  background: #1f2430;
  padding: 20px;
  border-radius: 8px;
`;

export const BlockTitle = styled.h3`
  color: #e3eafc;
  margin-top: 0;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const FormGroup = styled.div`
  margin-bottom: 15px;
  display: flex;
  flex-direction: column;
`;

export const FormRow = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
`;

export const Label = styled.label`
  display: block;
  margin-bottom: 5px;
  color: #a9b3d5;
`;

export const Input = styled.input`
  width: 100%;
  padding: 10px;
  background: #232936;
  color: #e3eafc;
  border: 1px solid #222a36;
  border-radius: 6px;
  box-sizing: border-box;
`;

export const TextArea = styled.textarea`
  width: 100%;
  padding: 10px;
  background: #232936;
  color: #e3eafc;
  border: 1px solid #222a36;
  border-radius: 6px;
  resize: vertical;
  box-sizing: border-box;
`;

export const SvgPreview = styled.div`
  width: 100%;
  height: 300px;
  background: #232936;
  border: 1px solid #222a36;
  border-radius: 6px;
  margin-top: 15px;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;

  svg {
    max-width: 100%;
    max-height: 100%;
  }
`;

export const Footer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 20px;
`;

export const StepIndicator = styled.div`
  display: flex;
  gap: 10px;
`;

export const Step = styled.div`
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: ${({ active }) => (active ? '#1976d2' : '#232936')};
  color: #fff;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
`;
