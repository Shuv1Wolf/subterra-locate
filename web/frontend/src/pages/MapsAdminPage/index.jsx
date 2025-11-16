import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { GEO_HOST } from '../../config';
import { apiClient, getMaps, deleteMap } from '../../utils/api';
import Header from '../../components/Header';
import {
  AdminPageContainer,
  Button,
  BeaconTable as MapTable,
  ActionsContainer,
  PaginationContainer,
} from './styles';

export default function MapsAdminPage() {
  const navigate = useNavigate();
  const [maps, setMaps] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [take] = useState(10); // Items per page
  const [total, setTotal] = useState(0);

  const fetchMaps = async () => {
    try {
      setLoading(true);
      const selectedOrgId = localStorage.getItem("selectedOrgId");
      if (selectedOrgId) {
        const skip = (page - 1) * take;
        const data = await getMaps(selectedOrgId, skip, take);
        setMaps(data.data || []);
        setTotal(data.total || 0);
      }
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (mapId) => {
    if (window.confirm('Are you sure you want to delete this map?')) {
      try {
        await deleteMap(mapId);
        fetchMaps(); // Refresh the list
      } catch (e) {
        setError(e.message);
      }
    }
  };

  useEffect(() => {
    fetchMaps();
  }, [page]);

  if (loading) return <AdminPageContainer>Loading maps...</AdminPageContainer>;
  if (error) return <AdminPageContainer>Error: {error}</AdminPageContainer>;

  return (
    <>
      <Header variant="page" title="Maps Admin" />
      <AdminPageContainer>
        <div style={{ marginBottom: '20px' }}>
          <Button onClick={() => navigate('/new-map')}>Add New Map</Button>
        </div>
        <MapTable>
          <thead>
          <tr>
            <th>Name</th>
            <th>Description</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {maps.map((map) => (
            <tr key={map.id}>
              <td>{map.name}</td>
              <td>{map.description}</td>
              <td>
                <ActionsContainer>
                  <Button onClick={() => navigate(`/edit-map/${map.id}`)}>Edit</Button>
                  <Button onClick={() => handleDelete(map.id)}>Delete</Button>
                </ActionsContainer>
              </td>
            </tr>
          ))}
        </tbody>
        </MapTable>
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
