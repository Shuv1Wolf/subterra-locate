import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { SYSTEM_HOST } from '../../config';
import Header from '../../components/Header';
import {
  AdminPageContainer,
  Button,
  BeaconTable as DeviceTable, // Renaming for clarity
  ActionsContainer,
  PaginationContainer,
} from './styles';

export default function DevicesAdminPage() {
  const navigate = useNavigate();
  const [devices, setDevices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [take] = useState(10); // Items per page
  const [total, setTotal] = useState(0);

  const fetchDevices = async () => {
    try {
      setLoading(true);
      const skip = (page - 1) * take;
      const response = await fetch(
        `${SYSTEM_HOST}/api/v1/system/devices?total=true&skip=${skip}&take=${take}`
      );
      if (!response.ok) throw new Error('Network response was not ok');
      const data = await response.json();
      setDevices(data.data || []);
      setTotal(data.total || 0);
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
  }, [page]);

  if (loading) return <AdminPageContainer>Loading devices...</AdminPageContainer>;
  if (error) return <AdminPageContainer>Error: {error}</AdminPageContainer>;

  return (
    <>
      <Header variant="page" title="Devices Admin" />
      <AdminPageContainer>
        <div style={{ marginBottom: '20px' }}>
          <Button onClick={() => navigate('/devices-admin/new')}>Add New Device</Button>
        </div>
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
        <PaginationContainer>
          <Button onClick={() => setPage(page - 1)} disabled={page === 1}>
            Previous
          </Button>
          <span>
            Showing {(page - 1) * take + 1} - {Math.min(page * take, total)} of {total}
          </span>
          <Button onClick={() => setPage(page + 1)} disabled={page * take >= total}>
            Next
          </Button>
        </PaginationContainer>
      </AdminPageContainer>
    </>
  );
}
