import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { GEO_HOST } from '../../config';
import { AdminPageContainer, Header, Title, Button } from '../BeaconsAdminPage/styles.js';
import { FormContainer, FormGroup, Label, Input, Select, CheckboxContainer } from './styles.js';

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
        const response = await fetch(`${GEO_HOST}/api/v1/geo/map`);
        const data = await response.json();
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
        const response = await fetch(`${GEO_HOST}/api/v1/geo/beacons/${beaconId}`);
        const data = await response.json();
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

    const url = isEditing
      ? `${GEO_HOST}/api/v1/geo/beacons` // Assuming update is PUT to the collection endpoint with ID in body
      : `${GEO_HOST}/api/v1/geo/beacons`;
      
    const method = isEditing ? 'PUT' : 'POST';

    try {
      const response = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(beacon),
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to save beacon');
      }
      navigate('/beacons-admin');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (loading && isEditing) return <AdminPageContainer>Loading...</AdminPageContainer>;

  return (
    <AdminPageContainer>
      <Header>
        <Title>{isEditing ? 'Edit Beacon' : 'Create Beacon'}</Title>
        <Button onClick={() => navigate('/beacons-admin')}>Back to List</Button>
      </Header>
      <FormContainer as="form" onSubmit={handleSubmit}>
        <FormGroup>
          <Label htmlFor="label">Label</Label>
          <Input type="text" name="label" id="label" value={beacon.label} onChange={handleChange} required />
        </FormGroup>
        <FormGroup>
          <Label htmlFor="udi">UDI</Label>
          <Input type="text" name="udi" id="udi" value={beacon.udi} onChange={handleChange} />
        </FormGroup>
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
        <FormGroup>
          <CheckboxContainer>
            <Input type="checkbox" name="enabled" id="enabled" checked={beacon.enabled} onChange={handleChange} style={{ width: 'auto' }} />
            <Label htmlFor="enabled" style={{ marginBottom: 0 }}>Enabled</Label>
          </CheckboxContainer>
        </FormGroup>
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <Button type="submit" disabled={loading}>{loading ? 'Saving...' : 'Save Beacon'}</Button>
      </FormContainer>
    </AdminPageContainer>
  );
}
