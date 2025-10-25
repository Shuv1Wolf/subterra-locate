import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { GEO_HOST } from '../../config';
import Header from '../../components/Header';
import {
  AdminPageContainer,
  Button,
  BeaconTable,
  ActionsContainer,
} from './styles.js';

export default function BeaconsAdminPage() {
  const navigate = useNavigate();
  const [beacons, setBeacons] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchBeacons = async () => {
    try {
      setLoading(true);
      const response = await fetch(`${GEO_HOST}/api/v1/geo/beacons`);
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      const data = await response.json();
      setBeacons(data.data || []);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (beaconId) => {
    if (window.confirm('Are you sure you want to delete this beacon?')) {
      try {
        const response = await fetch(`${GEO_HOST}/api/v1/geo/beacons/${beaconId}`, {
          method: 'DELETE',
        });
        if (!response.ok) {
          throw new Error('Failed to delete beacon');
        }
        fetchBeacons(); // Refresh the list
      } catch (e) {
        setError(e.message);
      }
    }
  };

  useEffect(() => {
    fetchBeacons();
  }, []);

  if (loading) return <AdminPageContainer>Loading beacons...</AdminPageContainer>;
  if (error) return <AdminPageContainer>Error: {error}</AdminPageContainer>;

  return (
    <>
      <Header variant="page" title="Beacons Admin" />
      <AdminPageContainer>
        <div style={{ marginBottom: '20px' }}>
          <Button onClick={() => navigate('/beacons-admin/new')}>Add New Beacon</Button>
        </div>
        <BeaconTable>
          <thead>
          <tr>
            <th>Label</th>
            <th>Map ID</th>
            <th>UDI</th>
            <th>Enabled</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {beacons.map((beacon) => (
            <tr key={beacon.id}>
              <td>{beacon.label}</td>
              <td>{beacon.map_id}</td>
              <td>{beacon.udi}</td>
              <td>{beacon.enabled ? 'Yes' : 'No'}</td>
              <td>
                <ActionsContainer>
                  <Button onClick={() => navigate(`/beacons-admin/edit/${beacon.id}`)}>Edit</Button>
                  <Button onClick={() => handleDelete(beacon.id)}>Delete</Button>
                </ActionsContainer>
              </td>
            </tr>
          ))}
        </tbody>
        </BeaconTable>
      </AdminPageContainer>
    </>
  );
}
