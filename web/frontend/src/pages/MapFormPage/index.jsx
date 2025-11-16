import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { GEO_HOST } from '../../config';
import { getMap, createMap, updateMap, uploadMap } from '../../utils/api';
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
} from './styles';

export default function MapFormPage() {
  const navigate = useNavigate();
  const { mapId } = useParams();
  const [map, setMap] = useState({
    name: '',
    description: '',
    svg_content: '',
    scale_x: 1,
    scale_y: 1,
    org_id: 'org$1',
    width: 0,
    height: 0,
    level: 0,
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const fileInputRef = useRef(null);

  const isEditing = Boolean(mapId);

  useEffect(() => {
    const fetchMap = async () => {
      if (!isEditing) return;
      setLoading(true);
      try {
        const data = await getMap(mapId);
        setMap(data);
      } catch (e) {
        setError('Failed to load map data');
      } finally {
        setLoading(false);
      }
    };
    fetchMap();
  }, [mapId, isEditing]);

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
      if (isEditing) {
        await updateMap(map);
      } else {
        await createMap(map);
      }
      navigate('/maps-admin');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this map?')) {
      try {
        await apiClient.delete(`${GEO_HOST}/api/v1/geo/map/${mapId}`);
        navigate('/maps-admin');
      } catch (e) {
        setError(e.message);
      }
    }
  };

  const handleFileUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;

    setLoading(true);
    setError(null);

    try {
      const response = await uploadMap(file);
      setMap(prev => ({
        ...prev,
        svg_content: response.svg_content,
        width: response.width,
        height: response.height,
      }));
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading && isEditing) return <AdminPageContainer>Loading...</AdminPageContainer>;

  return (
    <>
      <Header variant="page" title={isEditing ? 'Edit Map' : 'Create Map'} />
      <AdminPageContainer>
        <FormContainer as="form" onSubmit={handleSubmit}>
          <FormBlock>
            <BlockTitle>
              <span>General Information</span>
              {isEditing && <span style={{ opacity: 0.5, fontSize: '0.8rem' }}>{mapId}</span>}
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

          <FormBlock>
            <BlockTitle>Dimensions</BlockTitle>
            <FormRow>
              <FormGroup>
                <Label htmlFor="width">Width</Label>
                <Input type="number" name="width" id="width" value={map.width} onChange={handleChange} title="The width of the map in meters." />
              </FormGroup>
              <FormGroup>
                <Label htmlFor="height">Height</Label>
                <Input type="number" name="height" id="height" value={map.height} onChange={handleChange} title="The height of the map in meters." />
              </FormGroup>
              <FormGroup>
                <Label htmlFor="scale_x">Scale X</Label>
                <Input type="number" name="scale_x" id="scale_x" value={map.scale_x} onChange={handleChange} title="The horizontal scale of the map (pixels per meter)." />
              </FormGroup>
              <FormGroup>
                <Label htmlFor="scale_y">Scale Y</Label>
                <Input type="number" name="scale_y" id="scale_y" value={map.scale_y} onChange={handleChange} title="The vertical scale of the map (pixels per meter)." />
              </FormGroup>
            </FormRow>
          </FormBlock>

          <FormBlock>
            <BlockTitle>
              <span>SVG Content</span>
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

          {error && <p style={{ color: 'red' }}>{error}</p>}
          
          <Footer>
            <Button type="submit" disabled={loading}>{loading ? 'Saving...' : 'Save Map'}</Button>
            {isEditing && <Button type="button" onClick={handleDelete} style={{ background: '#d32f2f' }}>Delete</Button>}
          </Footer>
        </FormContainer>
      </AdminPageContainer>
    </>
  );
}
