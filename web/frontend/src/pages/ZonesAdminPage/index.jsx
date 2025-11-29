import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { GEO_HOST } from '../../config';
import { apiClient } from '../../utils/api';
import Header from '../../components/Header';
import {
  AdminPageContainer,
  Button,
  BeaconTable as ZoneTable,
  ActionsContainer,
  PaginationContainer,
} from './styles';

export default function ZonesAdminPage() {
  const navigate = useNavigate();
  const [zones, setZones] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [page, setPage] = useState(1);
  const [take] = useState(10); // Items per page
  const [total, setTotal] = useState(0);
  const [maps, setMaps] = useState([]);

  const fetchZones = async () => {
    try {
      setLoading(true);
      const skip = (page - 1) * take;
      const data = await apiClient.get(
        `${GEO_HOST}/api/v1/geo/zones?total=true&skip=${skip}&take=${take}`
      );
      setZones(data.data || []);
      setTotal(data.total || 0);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (zoneId) => {
    if (window.confirm('Are you sure you want to delete this zone?')) {
      try {
        await apiClient.delete(`${GEO_HOST}/api/v1/geo/zones/${zoneId}`);
        fetchZones(); // Refresh the list
      } catch (e) {
        setError(e.message);
      }
    }
  };

  const handleUnlinkFromMap = async (zone) => {
    if (window.confirm('Are you sure you want to unlink this zone from the map?')) {
      const updatedZone = {
        ...zone,
        map_id: '',
        position_x: 0,
        position_y: 0,
        width: 0,
        height: 0,
      };
      try {
        await apiClient.put(`${GEO_HOST}/api/v1/geo/zones`, updatedZone);
        fetchZones(); // Refresh the list
      } catch (e) {
        setError(e.message);
      }
    }
  };

  useEffect(() => {
    const fetchMaps = async () => {
      try {
        const mapsData = await apiClient.get(`${GEO_HOST}/api/v1/geo/map`);
        setMaps(mapsData.data || []);
      } catch (e) {
        console.error('Failed to load maps', e);
      }
    };

    fetchMaps();
    fetchZones();
  }, [page]);

  if (loading) return <AdminPageContainer>Loading zones...</AdminPageContainer>;
  if (error) return <AdminPageContainer>Error: {error}</AdminPageContainer>;

  return (
    <>
      <Header variant="page" title="Zones Admin" />
      <AdminPageContainer>
        <div style={{ marginBottom: '20px' }}>
          <Button onClick={() => navigate('/zones-admin/new')}>Add New Zone</Button>
        </div>
        <ZoneTable>
          <thead>
            <tr>
              <th>Name</th>
              <th>Map</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {zones.map((zone) => {
              const mapName = maps.find((map) => map.id === zone.map_id)?.name;
              return (
                <tr key={zone.id}>
                  <td>{zone.name}</td>
                  <td>{mapName || ''}</td>
                  <td>
                    <ActionsContainer>
                    <Button onClick={() => navigate(`/zones-admin/edit/${zone.id}`)}>Edit</Button>
                    <Button onClick={() => handleDelete(zone.id)}>Delete</Button>
                    {zone.map_id && (
                      <Button onClick={() => handleUnlinkFromMap(zone)}>Unlink from map</Button>
                    )}
                    </ActionsContainer>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </ZoneTable>
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
