import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { GEO_HOST } from '../../config';
import { apiClient } from '../../utils/api';
import Header from '../../components/Header';
import { AdminPageContainer, Button } from '../DevicesAdminPage/styles';
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
} from '../DeviceFormPage/styles';

export default function BeaconFormPage() {
  const navigate = useNavigate();
  const { beaconId } = useParams();
  const [beacon, setBeacon] = useState({
    udi: '',
    label: '',
    x: 0,
    y: 0,
    z: 0,
    org_id: 'org$1',
    enabled: true,
    map_id: '',
  });
  const [maps, setMaps] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const isEditing = Boolean(beaconId);

  useEffect(() => {
    const fetchMaps = async () => {
      try {
        const data = await apiClient.get(`${GEO_HOST}/api/v1/geo/map`);
        setMaps(data.data || []);
        if (data.data.length > 0 && !isEditing) {
          setBeacon(b => ({ ...b, map_id: data.data[0].id }));
        }
      } catch (e) {
        setError('Failed to load maps');
      }
    };

    const fetchBeacon = async () => {
      if (!isEditing) return;
      setLoading(true);
      try {
        const data = await apiClient.get(`${GEO_HOST}/api/v1/geo/beacons/${beaconId}`);
        setBeacon(data);
      } catch (e) {
        setError('Failed to load beacon data');
      } finally {
        setLoading(false);
      }
    };

    fetchMaps();
    fetchBeacon();
  }, [beaconId, isEditing]);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setBeacon(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : type === 'number' ? parseFloat(value) : value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const url = `${GEO_HOST}/api/v1/geo/beacons`;
    const method = isEditing ? 'put' : 'post';

    try {
      await apiClient[method](url, beacon);
      navigate('/beacons-admin');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this beacon?')) {
      try {
        await apiClient.delete(`${GEO_HOST}/api/v1/geo/beacons/${beaconId}`);
        navigate('/beacons-admin');
      } catch (e) {
        setError(e.message);
      }
    }
  };

  if (loading && isEditing) return <AdminPageContainer>Loading...</AdminPageContainer>;

  return (
    <>
      <Header variant="page" title={isEditing ? 'Edit Beacon' : 'Create Beacon'} />
      <AdminPageContainer>
        <FormContainer as="form" onSubmit={handleSubmit}>
          <FormBlock>
            <BlockTitle>
              <span>General Information</span>
              {isEditing && <span style={{ opacity: 0.5, fontSize: '0.8rem' }}>{beaconId}</span>}
            </BlockTitle>
            <FormGroup>
              <Label htmlFor="label">Label</Label>
              <Input type="text" name="label" id="label" value={beacon.label} onChange={handleChange} required />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="udi">UDI</Label>
              <Input type="text" name="udi" id="udi" value={beacon.udi} onChange={handleChange} />
            </FormGroup>
            <FormGroup>
              <CheckboxContainer>
                <Input type="checkbox" name="enabled" id="enabled" checked={beacon.enabled} onChange={handleChange} style={{ width: 'auto' }} />
                <Label htmlFor="enabled" style={{ marginBottom: 0 }}>Enabled</Label>
              </CheckboxContainer>
            </FormGroup>
          </FormBlock>

          <FormBlock>
            <BlockTitle>Location</BlockTitle>
            <FormGroup>
              <Label htmlFor="map_id">Map</Label>
              <Select name="map_id" id="map_id" value={beacon.map_id} onChange={handleChange} required>
                {maps.map(map => <option key={map.id} value={map.id}>{map.name}</option>)}
              </Select>
            </FormGroup>
            <FormGroup>
              <Label htmlFor="x">X Coordinate</Label>
              <Input type="number" name="x" id="x" value={beacon.x} onChange={handleChange} step="0.1" required />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="y">Y Coordinate</Label>
              <Input type="number" name="y" id="y" value={beacon.y} onChange={handleChange} step="0.1" required />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="z">Z Coordinate</Label>
              <Input type="number" name="z" id="z" value={beacon.z} onChange={handleChange} step="0.1" required />
            </FormGroup>
          </FormBlock>

          {error && <p style={{ color: 'red' }}>{error}</p>}

          <Footer>
            <Button type="submit" disabled={loading}>{loading ? 'Saving...' : 'Save Beacon'}</Button>
            {isEditing && <Button type="button" onClick={handleDelete} style={{ background: '#d32f2f' }}>Delete</Button>}
          </Footer>
        </FormContainer>
      </AdminPageContainer>
    </>
  );
}
