import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { GEO_HOST } from '../../config';
import { apiClient } from '../../utils/api';
import Header from '../../components/Header';
import { AdminPageContainer, Button } from '../ZonesAdminPage/styles';
import {
  FormContainer,
  FormBlock,
  BlockTitle,
  FormGroup,
  Label,
  Input,
  Select,
  CheckboxContainer,
  Footer,
} from './styles';

export default function ZoneFormPage() {
  const navigate = useNavigate();
  const { zoneId } = useParams();
  const [zone, setZone] = useState({
    name: '',
    map_id: '',
    org_id: 'org$1',
    position_x: 0,
    position_y: 0,
    width: 10,
    height: 10,
    type: 'rect',
    color: '#FF0000',
    max_device: 0,
  });
  const [maps, setMaps] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const isEditing = Boolean(zoneId);

  useEffect(() => {
    const fetchMaps = async () => {
      try {
        const data = await apiClient.get(`${GEO_HOST}/api/v1/geo/map`);
        setMaps(data.data || []);
        if (data.data.length > 0 && !isEditing) {
          setZone(z => ({ ...z, map_id: data.data[0].id }));
        }
      } catch (e) {
        setError('Failed to load maps');
      }
    };

    const fetchZone = async () => {
      if (!isEditing) return;
      setLoading(true);
      try {
        const data = await apiClient.get(`${GEO_HOST}/api/v1/geo/zones/${zoneId}`);
        setZone(data);
      } catch (e) {
        setError('Failed to load zone data');
      } finally {
        setLoading(false);
      }
    };

    fetchMaps();
    fetchZone();
  }, [zoneId, isEditing]);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setZone(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : type === 'number' ? parseFloat(value) : value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const url = `${GEO_HOST}/api/v1/geo/zones`;
    const method = isEditing ? 'put' : 'post';

    try {
      await apiClient[method](url, zone);
      navigate('/zones-admin');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this zone?')) {
      try {
        await apiClient.delete(`${GEO_HOST}/api/v1/geo/zones/${zoneId}`);
        navigate('/zones-admin');
      } catch (e) {
        setError(e.message);
      }
    }
  };

  if (loading && isEditing) return <AdminPageContainer>Loading...</AdminPageContainer>;

  const mapName = maps.find(m => m.id === zone.map_id)?.name;

  return (
    <>
      <Header variant="page" title={isEditing ? 'Edit Zone' : 'Create Zone'} />
      <AdminPageContainer>
        <FormContainer as="form" onSubmit={handleSubmit}>
          <FormBlock>
            <BlockTitle>
              <span>General Information</span>
              {isEditing && <span style={{ opacity: 0.5, fontSize: '0.8rem' }}>{zoneId}</span>}
            </BlockTitle>
            <FormGroup>
              <Label htmlFor="name">Name</Label>
              <Input type="text" name="name" id="name" value={zone.name} onChange={handleChange} required />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="color">Color</Label>
              <Input type="color" name="color" id="color" value={zone.color} onChange={handleChange} />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="max_device">Max Devices</Label>
              <Input type="number" name="max_device" id="max_device" value={zone.max_device} onChange={handleChange} min="0" />
            </FormGroup>
          </FormBlock>

          {isEditing && (
          <FormBlock>
            <BlockTitle>Location and Size</BlockTitle>
            <FormGroup>
              <Label htmlFor="map_id">Map</Label>
              <Input type="text" value={mapName || ''} disabled />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="position_x">X Coordinate</Label>
              <Input type="number" name="position_x" id="position_x" value={zone.position_x} onChange={handleChange} step="0.1" required />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="position_y">Y Coordinate</Label>
              <Input type="number" name="position_y" id="position_y" value={zone.position_y} onChange={handleChange} step="0.1" required />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="width">Width</Label>
              <Input type="number" name="width" id="width" value={zone.width} onChange={handleChange} step="0.1" required />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="height">Height</Label>
              <Input type="number" name="height" id="height" value={zone.height} onChange={handleChange} step="0.1" required />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="type">Type</Label>
              <Select name="type" id="type" value={zone.type} onChange={handleChange}>
                <option value="rect">Rectangle</option>
              </Select>
            </FormGroup>
          </FormBlock>
          )}

          {error && <p style={{ color: 'red' }}>{error}</p>}

          <Footer>
            <Button type="submit" disabled={loading}>{loading ? 'Saving...' : 'Save Zone'}</Button>
            {isEditing && <Button type="button" onClick={handleDelete} style={{ background: '#d32f2f' }}>Delete</Button>}
          </Footer>
        </FormContainer>
      </AdminPageContainer>
    </>
  );
}
