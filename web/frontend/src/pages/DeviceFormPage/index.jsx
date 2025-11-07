import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { SYSTEM_HOST } from '../../config';
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
  CheckboxContainer,
  Footer,
} from './styles';

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
        const data = await apiClient.get(`${SYSTEM_HOST}/api/v1/system/device/${deviceId}`);
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

    const url = `${SYSTEM_HOST}/api/v1/system/device`;
    const method = isEditing ? 'put' : 'post';

    try {
      await apiClient[method](url, device);
      navigate('/devices-admin');
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this device?')) {
      try {
        await apiClient.delete(`${SYSTEM_HOST}/api/v1/system/device/${deviceId}`);
        navigate('/devices-admin');
      } catch (e) {
        setError(e.message);
      }
    }
  };

  if (loading && isEditing) return <AdminPageContainer>Loading...</AdminPageContainer>;

  return (
    <>
      <Header variant="page" title={isEditing ? 'Edit Device' : 'Create Device'} />
      <AdminPageContainer>
        <FormContainer as="form" onSubmit={handleSubmit}>
          <FormBlock>
            <BlockTitle>
              <span>General Information</span>
              {isEditing && <span style={{ opacity: 0.5, fontSize: '0.8rem' }}>{deviceId}</span>}
            </BlockTitle>
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
              <CheckboxContainer>
                <Input type="checkbox" name="enabled" id="enabled" checked={device.enabled} onChange={handleChange} style={{ width: 'auto' }} />
                <Label htmlFor="enabled" style={{ marginBottom: 0 }}>Enabled</Label>
              </CheckboxContainer>
            </FormGroup>
          </FormBlock>

          <FormBlock>
            <BlockTitle>Network Details</BlockTitle>
            <FormGroup>
              <Label htmlFor="mac_address">MAC Address</Label>
              <Input type="text" name="mac_address" id="mac_address" value={device.mac_address} onChange={handleChange} />
            </FormGroup>
            <FormGroup>
              <Label htmlFor="ip_address">IP Address</Label>
              <Input type="text" name="ip_address" id="ip_address" value={device.ip_address} onChange={handleChange} />
            </FormGroup>
          </FormBlock>

          {error && <p style={{ color: 'red' }}>{error}</p>}
          
          <Footer>
            <Button type="submit" disabled={loading}>{loading ? 'Saving...' : 'Save Device'}</Button>
            {isEditing && <Button type="button" onClick={handleDelete} style={{ background: '#d32f2f' }}>Delete</Button>}
          </Footer>
        </FormContainer>
      </AdminPageContainer>
    </>
  );
}
