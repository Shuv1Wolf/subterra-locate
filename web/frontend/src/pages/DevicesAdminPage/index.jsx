import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { SYSTEM_HOST } from '../../config';
import {
  AdminPageContainer,
  Header,
  Title,
  Button,
  BeaconTable as DeviceTable, // Renaming for clarity
  ActionsContainer,
} from './styles';

export default function DevicesAdminPage() {
  const navigate = useNavigate();
  const [devices, setDevices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchDevices = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${SYSTEM_HOST}/api/v1/system/devices`);
      if (!response.ok) throw new Error('Network response was not ok');
      const data = await response.json();
      setDevices(data.data || []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (deviceId) => {
    if (window.confirm('Are you sure you want to delete this device?')) {
      try {
        const response = await fetch(`${SYSTEM_HOST}/api/v1/system/device/${deviceId}`, {
          method: 'DELETE',
        });
        if (!response.ok) throw new Error('Failed to delete device');
        fetchDevices(); // Refresh the list
      } catch (e) {
        setError(e.message);
      }
    }
  };

  useEffect(() => {
    fetchDevices();
  }, []);

  if (loading) return <AdminPageContainer>Loading devices...</AdminPageContainer>;
  if (error) return <AdminPageContainer>Error: {error}</AdminPageContainer>;

  return (
    <AdminPageContainer>
      <Header>
        <Title>Devices Admin</Title>
        <Button onClick={() => navigate('/devices-admin/new')}>Add New Device</Button>
      </Header>
      <DeviceTable>
        <thead>
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Model</th>
            <th>MAC Address</th>
            <th>Enabled</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {devices.map((device) => (
            <tr key={device.id}>
              <td>{device.name}</td>
              <td>{device.type}</td>
              <td>{device.model}</td>
              <td>{device.mac_address}</td>
              <td>{device.enabled ? 'Yes' : 'No'}</td>
              <td>
                <ActionsContainer>
                  <Button onClick={() => navigate(`/devices-admin/edit/${device.id}`)}>Edit</Button>
                  <Button onClick={() => handleDelete(device.id)}>Delete</Button>
                </ActionsContainer>
              </td>
            </tr>
          ))}
        </tbody>
      </DeviceTable>
    </AdminPageContainer>
  );
}
