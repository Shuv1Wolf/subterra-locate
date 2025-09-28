import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { SYSTEM_HOST } from '../../config';
import { AdminPageContainer, Header, Title, Button } from '../DevicesAdminPage/styles';
import { FormContainer, FormGroup, Label, Input, CheckboxContainer } from './styles';

export default function DeviceFormPage() {
  const navigate = useNavigate();
  const { deviceId } = useParams();
  const [device, setDevice] = useState({
    name: '',
    type: 'unknown',
    model: '',
    org_id: 'org$1',
    enabled: true,
    mac_address: '',
    ip_address: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const isEditing = Boolean(deviceId);

  useEffect(() => {
    const fetchDevice = async () => {
      if (!isEditing) return;
      setLoading(true);
      try {
        const response = await fetch(`${SYSTEM_HOST}/api/v1/system/device/${deviceId}`);
        const data = await response.json();
        setDevice(data);
      } catch (e) {
        setError('Failed to load device data');
      } finally {
        setLoading(false);
      }
    };
    fetchDevice();
  }, [deviceId, isEditing]);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setDevice(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    const url = isEditing
      ? `${SYSTEM_HOST}/api/v1/system/device`
      : `${SYSTEM_HOST}/api/v1/system/device`;
      
    const method = isEditing ? 'PUT' : 'POST';

    try {
      const response = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(device),
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Failed to save device');
      }
      navigate('/devices-admin');
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
        <Title>{isEditing ? 'Edit Device' : 'Create Device'}</Title>
        <Button onClick={() => navigate('/devices-admin')}>Back to List</Button>
      </Header>
      <FormContainer as="form" onSubmit={handleSubmit}>
        <FormGroup>
          <Label htmlFor="name">Name</Label>
          <Input type="text" name="name" id="name" value={device.name} onChange={handleChange} required />
        </FormGroup>
        <FormGroup>
          <Label htmlFor="type">Type</Label>
          <Input type="text" name="type" id="type" value={device.type} onChange={handleChange} />
        </FormGroup>
        <FormGroup>
          <Label htmlFor="model">Model</Label>
          <Input type="text" name="model" id="model" value={device.model} onChange={handleChange} />
        </FormGroup>
        <FormGroup>
          <Label htmlFor="mac_address">MAC Address</Label>
          <Input type="text" name="mac_address" id="mac_address" value={device.mac_address} onChange={handleChange} />
        </FormGroup>
        <FormGroup>
          <Label htmlFor="ip_address">IP Address</Label>
          <Input type="text" name="ip_address" id="ip_address" value={device.ip_address} onChange={handleChange} />
        </FormGroup>
        <FormGroup>
          <CheckboxContainer>
            <Input type="checkbox" name="enabled" id="enabled" checked={device.enabled} onChange={handleChange} style={{ width: 'auto' }} />
            <Label htmlFor="enabled" style={{ marginBottom: 0 }}>Enabled</Label>
          </CheckboxContainer>
        </FormGroup>
        {error && <p style={{ color: 'red' }}>{error}</p>}
        <Button type="submit" disabled={loading}>{loading ? 'Saving...' : 'Save Device'}</Button>
      </FormContainer>
    </AdminPageContainer>
  );
}
