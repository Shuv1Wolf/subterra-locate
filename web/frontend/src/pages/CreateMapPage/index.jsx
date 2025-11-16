import React, { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { TransformWrapper, TransformComponent } from 'react-zoom-pan-pinch';
import { GEO_HOST } from '../../config';
import { createMap, updateMap, uploadMap } from '../../utils/api';
import Header from '../../components/Header';
import { AdminPageContainer, Button } from '../DevicesAdminPage/styles';
import {
  FormContainer,
  FormBlock,
  BlockTitle,
  FormGroup,
  FormRow,
  Label,
  Input,
  TextArea,
  Footer,
  SvgPreview,
  StepIndicator,
  Step,
} from './styles';

export default function CreateMapPage() {
  const navigate = useNavigate();
  const [step, setStep] = useState(1);
  const [map, setMap] = useState({
    name: '',
    description: '',
    svg_content: '',
    scale_x: 1,
    scale_y: 1,
    org_id: localStorage.getItem('selectedOrgId') || '',
    width: 0,
    height: 0,
    level: 0,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const fileInputRef = useRef(null);
  const [calibrationMode, setCalibrationMode] = useState('none'); // 'none', 'horizontal', 'vertical'
  const [firstPoint, setFirstPoint] = useState(null);
  const [secondPoint, setSecondPoint] = useState(null);

  const handleChange = (e) => {
    const { name, value, type } = e.target;
    setMap(prev => ({
      ...prev,
      [name]: type === 'number' ? parseFloat(value) : value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await updateMap(map);
      navigate('/maps-admin');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleFileUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    setLoading(true);
    setError(null);

    try {
      const reader = new FileReader();
      reader.onload = async (event) => {
        const svgContent = event.target.result;
        const viewBoxMatch = svgContent.match(/viewBox="(\d+)\s+(\d+)\s+(\d+)\s+(\d+)"/);
        let width = 0, height = 0;
        if (viewBoxMatch) {
          width = parseInt(viewBoxMatch[3], 10);
          height = parseInt(viewBoxMatch[4], 10);
        }

        const response = await uploadMap(file, map.id);
        setMap(prev => ({
          ...prev,
          svg_content: svgContent,
          width: width,
          height: height,
        }));
      };
      reader.readAsText(file);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handlePreviewClick = (e) => {
    if (calibrationMode === 'none') return;

    const svgElement = e.currentTarget.querySelector('svg');
    if (!svgElement) return;

    const rect = svgElement.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    if (!firstPoint) {
      setFirstPoint({ x, y });
    } else {
      if (calibrationMode === 'horizontal') {
        const pixelDistance = Math.abs(x - firstPoint.x);
        const newScale = pixelDistance > 0 ? parseFloat((1 / pixelDistance).toFixed(4)) : 0;
        setMap(prev => ({ ...prev, scale_x: newScale }));
      } else if (calibrationMode === 'vertical') {
        const pixelDistance = Math.abs(y - firstPoint.y);
        const newScale = pixelDistance > 0 ? parseFloat((1 / pixelDistance).toFixed(4)) : 0;
        setMap(prev => ({ ...prev, scale_y: newScale }));
      }
      setFirstPoint(null);
      setSecondPoint(null);
      setCalibrationMode('none');
    }
  };

  const handleMouseMove = (e) => {
    if (calibrationMode === 'none' || !firstPoint) return;

    const svgElement = e.currentTarget.querySelector('svg');
    if (!svgElement) return;

    const rect = svgElement.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;

    setSecondPoint({ x, y });
  };

  const nextStep = async () => {
    setLoading(true);
    setError(null);
    try {
      if (step === 1) {
        const newMap = await createMap({
          name: map.name,
          description: map.description,
          level: map.level,
          org_id: map.org_id,
        });
        setMap(newMap);
        setStep(2);
      } else if (step === 2) {
        await updateMap(map);
        setStep(3);
      }
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const prevStep = () => setStep(prev => Math.max(prev - 1, 1));

  return (
    <>
      <Header variant="page" title="Create Map" />
      <AdminPageContainer>
        <FormContainer as="form" onSubmit={handleSubmit}>
          {step === 1 && (
            <FormBlock>
              <BlockTitle>
                <span>Step 1: General Information</span>
              </BlockTitle>
              <FormGroup>
                <Label htmlFor="name">Name</Label>
                <Input type="text" name="name" id="name" value={map.name} onChange={handleChange} required />
              </FormGroup>
              <FormGroup>
                <Label htmlFor="level">Level</Label>
                <Input type="number" name="level" id="level" value={map.level} onChange={handleChange} />
              </FormGroup>
              <FormGroup>
                <Label htmlFor="description">Description</Label>
                <TextArea name="description" id="description" value={map.description} onChange={handleChange} rows="3" />
              </FormGroup>
            </FormBlock>
          )}

          {step === 2 && (
            <FormBlock>
              <BlockTitle>
                <span>Step 2: Upload Map</span>
                <Button type="button" onClick={() => fileInputRef.current.click()}>
                  Upload Map
                </Button>
                <input
                  type="file"
                  ref={fileInputRef}
                  onChange={handleFileUpload}
                  style={{ display: 'none' }}
                  accept=".svg"
                />
              </BlockTitle>
              <SvgPreview dangerouslySetInnerHTML={{ __html: map.svg_content }} />
            </FormBlock>
          )}

          {step === 3 && (
            <FormBlock>
              <BlockTitle>Step 3: Dimensions & Scale</BlockTitle>
              <p>Click a calibration button, then click two points on the map preview below to define a 1-meter segment.</p>
              <FormRow>
                <Button type="button" onClick={() => { setCalibrationMode('horizontal'); setFirstPoint(null); }}>Calibrate Horizontal (1m)</Button>
                <Button type="button" onClick={() => { setCalibrationMode('vertical'); setFirstPoint(null); }}>Calibrate Vertical (1m)</Button>
              </FormRow>
              <FormRow>
                <FormGroup>
                  <Label htmlFor="width">Width (px)</Label>
                  <Input type="number" name="width" id="width" value={map.width} onChange={handleChange} title="The width of the SVG in pixels." />
                </FormGroup>
                <FormGroup>
                  <Label htmlFor="height">Height (px)</Label>
                  <Input type="number" name="height" id="height" value={map.height} onChange={handleChange} title="The height of the SVG in pixels." />
                </FormGroup>
                <FormGroup>
                  <Label htmlFor="scale_x">Scale X (m/px)</Label>
                  <Input type="number" name="scale_x" id="scale_x" value={map.scale_x} onChange={handleChange} title="Calculated horizontal scale (meters per pixel)." />
                </FormGroup>
                <FormGroup>
                  <Label htmlFor="scale_y">Scale Y (m/px)</Label>
                  <Input type="number" name="scale_y" id="scale_y" value={map.scale_y} onChange={handleChange} title="Calculated vertical scale (meters per pixel)." />
                </FormGroup>
              </FormRow>
              <TransformWrapper>
                <TransformComponent>
                  <SvgPreview onMouseMove={handleMouseMove} onClick={handlePreviewClick} style={{ cursor: calibrationMode !== 'none' ? 'crosshair' : 'default' }}>
                    <div dangerouslySetInnerHTML={{ __html: map.svg_content }} />
                    {firstPoint && secondPoint && (
                      <svg width="100%" height="100%" style={{ position: 'absolute', top: 0, left: 0, pointerEvents: 'none' }}>
                        <line x1={firstPoint.x} y1={firstPoint.y} x2={secondPoint.x} y2={secondPoint.y} stroke="red" strokeWidth="2" strokeDasharray="5,5" />
                      </svg>
                    )}
                  </SvgPreview>
                </TransformComponent>
              </TransformWrapper>
            </FormBlock>
          )}

          {error && <p style={{ color: 'red' }}>{error}</p>}
          
          <Footer>
            <div>
              {step > 1 && <Button type="button" onClick={prevStep}>Previous</Button>}
            </div>
            <StepIndicator>
              <Step active={step === 1}>1</Step>
              <Step active={step === 2}>2</Step>
              <Step active={step === 3}>3</Step>
            </StepIndicator>
            <div>
              {step < 3 && <Button type="button" onClick={nextStep} disabled={loading}>{loading ? 'Saving...' : 'Next'}</Button>}
              {step === 3 && <Button type="submit" disabled={loading}>{loading ? 'Saving...' : 'Save Map'}</Button>}
            </div>
          </Footer>
        </FormContainer>
      </AdminPageContainer>
    </>
  );
}
