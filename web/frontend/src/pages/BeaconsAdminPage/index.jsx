import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { GEO_HOST } from '../../config';
import { apiClient } from '../../utils/api';
import Header from '../../components/Header';
import {
  AdminPageContainer,
  Button,
  BeaconTable,
  ActionsContainer,
  PaginationContainer,
} from './styles.js';

export default function BeaconsAdminPage() {
  const navigate = useNavigate();
  const [beacons, setBeacons] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [take] = useState(10); // Items per page
  const [total, setTotal] = useState(0);

  const fetchBeacons = async () => {
    try {
      setLoading(true);
      const skip = (page - 1) * take;
      const data = await apiClient.get(
        `${GEO_HOST}/api/v1/geo/beacons?total=true&skip=${skip}&take=${take}`
      );
      setBeacons(data.data || []);
      setTotal(data.total || 0);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (beaconId) => {
    if (window.confirm('Are you sure you want to delete this beacon?')) {
      try {
        await apiClient.delete(`${GEO_HOST}/api/v1/geo/beacons/${beaconId}`);
        fetchBeacons(); // Refresh the list
      } catch (e) {
        setError(e.message);
      }
    }
  };

  useEffect(() => {
    fetchBeacons();
  }, [page]);

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
