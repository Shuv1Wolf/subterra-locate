import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import Select from 'react-select';
import { TransformWrapper, TransformComponent } from 'react-zoom-pan-pinch';
import Draggable from 'react-draggable';
import UserIcon from '../../assets/user.svg';
import { GEO_HOST, SYSTEM_HOST } from '../../config';
import {
  MapsPageContainer,
  MapSelectorContainer,
  MapSelect,
  MapInfo,
  MapWrapper,
  BackArrow,
  PopupContainer,
  PopupHeader,
  CloseButton,
  FilterButton,
} from './styles.js';

// Styles for React-Select
const selectStyles = {
  container: (provided) => ({ ...provided, width: '250px', }),
  control: (provided) => ({ ...provided, backgroundColor: '#181c24', borderColor: '#222a36', color: '#e3eafc', }),
  menu: (provided) => ({ ...provided, backgroundColor: '#181c24', }),
  option: (provided, state) => ({ ...provided, backgroundColor: state.isSelected ? '#1976d2' : state.isFocused ? '#232936' : '#181c24', color: '#e3eafc', }),
  multiValue: (provided) => ({ ...provided, backgroundColor: '#232936', }),
  multiValueLabel: (provided) => ({ ...provided, color: '#e3eafc', }),
  input: (provided) => ({ ...provided, color: '#e3eafc', }),
};

// Device Info Popup Component
const DeviceInfoPopup = ({ device, onClose }) => (
  <Draggable handle=".popup-header">
    <PopupContainer>
      <PopupHeader className="popup-header">
        <h2>{device.name}</h2>
        <CloseButton onClick={onClose}>&times;</CloseButton>
      </PopupHeader>
      <div>
        <p><strong>ID:</strong> {device.id}</p>
        <p><strong>Type:</strong> {device.type}</p>
        <p><strong>Model:</strong> {device.model}</p>
        <p><strong>MAC Address:</strong> {device.mac_address}</p>
        <p><strong>IP Address:</strong> {device.ip_address}</p>
        <p><strong>Status:</strong> {device.enabled ? 'Enabled' : 'Disabled'}</p>
      </div>
    </PopupContainer>
  </Draggable>
);

export default function NewMapPage() {
  const navigate = useNavigate();
  const [maps, setMaps] = useState([]);
  const [selectedMapId, setSelectedMapId] = useState('');
  const [devices, setDevices] = useState([]);
  const [allDevices, setAllDevices] = useState([]);
  const [stagedFilters, setStagedFilters] = useState([]);
  const [activeFilters, setActiveFilters] = useState([]);
  const [selectedDevice, setSelectedDevice] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const ws = useRef(null);

  const handleDeviceClick = async (deviceId) => {
    try {
      const response = await fetch(`${SYSTEM_HOST}/api/v1/system/device/${deviceId}`);
      if (!response.ok) throw new Error('Failed to fetch device details');
      const data = await response.json();
      setSelectedDevice(data);
    } catch (err) {
      console.error(err);
    }
  };

  const applyFilters = () => {
    setDevices([]); // Clear devices before applying new filters
    setActiveFilters(stagedFilters);
  };

  const handleFilterChange = (selectedOptions) => {
    setStagedFilters(selectedOptions);
  };

  // Fetch initial data
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const [mapsRes, devicesRes] = await Promise.all([
          fetch(`${GEO_HOST}/api/v1/geo/map`),
          fetch(`${SYSTEM_HOST}/api/v1/system/devices`),
        ]);
        if (!mapsRes.ok || !devicesRes.ok) throw new Error('Network response was not ok');
        const [mapsData, devicesData] = await Promise.all([mapsRes.json(), devicesRes.json()]);
        setMaps(mapsData.data);
        setAllDevices(devicesData.data);
        if (mapsData.data.length > 0) setSelectedMapId(mapsData.data[0].id);
      } catch (e) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  // WebSocket connection
  useEffect(() => {
    if (!selectedMapId) return;
    if (ws.current) ws.current.close();

    const wsHost = GEO_HOST.replace('http://', 'ws://');
    let wsUrl = `${wsHost}/api/v1/geo/location/device/monitor?org_id=org$1&map_id=${selectedMapId}`;
    
    const deviceIds = activeFilters.map(d => d.value);
    if (deviceIds.length > 0) {
      wsUrl += `&device_ids=${deviceIds.join(',')}`;
    }

    ws.current = new WebSocket(wsUrl);
    ws.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.event) {
        setDevices(prev => {
          const map = new Map(prev.map(d => [d.device_id, d]));
          for (const device of message.event) {
            if (device.x == null || device.y == null) map.delete(device.device_id);
            else map.set(device.device_id, device);
          }
          return Array.from(map.values());
        });
      }
    };

    return () => {
      if (ws.current) ws.current.close();
    };
  }, [selectedMapId, activeFilters]);

  const deviceOptions = allDevices.map(d => ({ value: d.id, label: d.name }));
  const selectedMap = maps.find((map) => map.id === selectedMapId);

  if (loading) return <MapsPageContainer>Loading data...</MapsPageContainer>;
  if (error) return <MapsPageContainer>Error: {error}</MapsPageContainer>;

  return (
    <MapsPageContainer>
      {selectedDevice && <DeviceInfoPopup device={selectedDevice} onClose={() => setSelectedDevice(null)} />}
      <BackArrow onClick={() => navigate('/')}>&#x2190;</BackArrow>
      <MapSelectorContainer>
        <MapSelect value={selectedMapId} onChange={(e) => setSelectedMapId(e.target.value)}>
          {maps.map((map) => (<option key={map.id} value={map.id}>{map.name}</option>))}
        </MapSelect>
        <Select isMulti options={deviceOptions} value={stagedFilters} onChange={handleFilterChange} styles={selectStyles} placeholder="Filter by device..." />
        <FilterButton onClick={applyFilters}>Apply</FilterButton>
        {selectedMap && <MapInfo><span>Level: {selectedMap.level}</span><span>Size: {selectedMap.width}x{selectedMap.height}</span></MapInfo>}
      </MapSelectorContainer>
      <MapWrapper>
        {selectedMap && (
          <TransformWrapper>
            <TransformComponent>
              <div style={{ position: 'relative', width: selectedMap.width, height: selectedMap.height }}>
                <div style={{ width: '100%', height: '100%' }} dangerouslySetInnerHTML={{ __html: selectedMap.svg_content }} />
                {devices.map((device) => (
                  <UserIcon
                    key={device.device_id}
                    title={`${device.device_name} (x: ${device.x.toFixed(2)}, y: ${device.y.toFixed(2)})`}
                    onClick={() => handleDeviceClick(device.device_id)}
                    style={{ position: 'absolute', top: `${device.y}px`, left: `${device.x}px`, width: '24px', height: '24px', transform: 'translate(-50%, -50%)', cursor: 'pointer' }}
                  />
                ))}
              </div>
            </TransformComponent>
          </TransformWrapper>
        )}
      </MapWrapper>
    </MapsPageContainer>
  );
}
