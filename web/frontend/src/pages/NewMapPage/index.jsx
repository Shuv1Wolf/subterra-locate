import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import Select from 'react-select';
import { TransformWrapper, TransformComponent } from 'react-zoom-pan-pinch';
import Draggable from 'react-draggable';
import UserIcon from '../../assets/user.svg';
import BeaconIcon from '../../assets/bullseye-animated.gif';
import { GEO_HOST, SYSTEM_HOST } from '../../config';
import { apiClient, getDeviceHistory } from '../../utils/api';
import Header from '../../components/Header';
import {
  MapsPageContainer,
  MapSelectorContainer,
  MapSelect,
  MapInfo,
  MapWrapper,
  PopupContainer,
  PopupHeader,
  CloseButton,
  FilterButton,
  ContextMenu,
  ContextMenuItem,
  ModalBackdrop,
  ModalContent,
  BeaconList,
  BeaconListItem,
  FilterBlock,
  FilterHeader,
  CheckboxLabel,
  ToggleButton,
  TileButton,
  PaginationContainer,
  HistoryControlsContainer,
  ControlRow,
  StyledInput,
  HistoryHeader,
  PopupButtonContainer,
} from './styles.js';

// Styles for React-Select
const selectStyles = {
  container: (provided) => ({ ...provided, width: '250px' }),
  control: (provided) => ({ ...provided, backgroundColor: '#181c24', borderColor: '#222a36', color: '#e3eafc' }),
  menu: (provided) => ({ ...provided, backgroundColor: '#181c24' }),
  option: (provided, state) => ({ ...provided, backgroundColor: state.isSelected ? '#1976d2' : state.isFocused ? '#232936' : '#181c24', color: '#e3eafc' }),
  multiValue: (provided) => ({ ...provided, backgroundColor: '#232936' }),
  multiValueLabel: (provided) => ({ ...provided, color: '#e3eafc' }),
  input: (provided) => ({ ...provided, color: '#e3eafc' }),
};

const BeaconSelectorModal = ({ beacons, onSelect, onClose }) => (
  <ModalBackdrop onClick={onClose}>
    <ModalContent onClick={(e) => e.stopPropagation()}>
      <h3>Select the beacon to move</h3>
      <BeaconList>
        {beacons.map((beacon) => (
          <BeaconListItem key={beacon.id} onClick={() => onSelect(beacon)}>
            {beacon.label}
          </BeaconListItem>
        ))}
      </BeaconList>
    </ModalContent>
  </ModalBackdrop>
);

const ZoneSelectorModal = ({ zones, onSelect, onClose }) => (
  <ModalBackdrop onClick={onClose}>
    <ModalContent onClick={(e) => e.stopPropagation()}>
      <h3>Select the zone to edit</h3>
      <BeaconList>
        {zones.map((zone) => (
          <BeaconListItem key={zone.id} onClick={() => onSelect(zone)}>
            {zone.name}
          </BeaconListItem>
        ))}
      </BeaconList>
    </ModalContent>
  </ModalBackdrop>
);

// Device Info Popup Component
const DeviceInfoPopup = ({ device, onClose, onShowHistory }) => (
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
        <PopupButtonContainer>
          <TileButton onClick={() => onShowHistory(device)}>Show History</TileButton>
        </PopupButtonContainer>
      </div>
    </PopupContainer>
  </Draggable>
);

// Beacon Info Popup Component
const BeaconInfoPopup = ({ beacon, onClose }) => (
  <Draggable handle=".popup-header">
    <PopupContainer>
      <PopupHeader className="popup-header">
        <h2>{beacon.label || beacon.name}</h2>
        <CloseButton onClick={onClose}>&times;</CloseButton>
      </PopupHeader>
      <div>
        <p><strong>ID:</strong> {beacon.id}</p>
        <p><strong>Label:</strong> {beacon.label}</p>
        <p><strong>UDI:</strong> {beacon.udi}</p>
        <p><strong>Status:</strong> {beacon.enabled ? 'Enabled' : 'Disabled'}</p>
      </div>
    </PopupContainer>
  </Draggable>
);

// Zone Info Popup Component
const ZoneInfoPopup = ({ zone, onClose, onSave }) => {
  const [name, setName] = useState(zone.name);
  const [maxDevices, setMaxDevices] = useState(zone.max_device);
  const [color, setColor] = useState(zone.color);

  const handleSave = () => {
    onSave({ ...zone, name, max_device: maxDevices, color });
  };

  return (
    <Draggable handle=".popup-header">
      <PopupContainer>
        <PopupHeader className="popup-header">
          <h2>Edit Zone</h2>
          <CloseButton onClick={onClose}>&times;</CloseButton>
        </PopupHeader>
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: '120px auto',
            alignItems: 'center',
            gap: '10px',
            padding: '0 12px',
            width: '100%',
            maxWidth: '100%',
            boxSizing: 'border-box',
            wordBreak: 'break-word'
          }}
        >
          <p style={{ margin: 0, wordBreak: 'break-word' }}><strong>ID:</strong></p>
          <p style={{ margin: 0, wordBreak: 'break-word' }}>{zone.id}</p>
          
          <label htmlFor="zoneName" style={{ margin: 0 }}><strong>Name:</strong></label>
          <StyledInput
            id="zoneName"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            style={{ boxSizing: 'border-box', width: '100%', minWidth: 0, wordBreak: 'break-word' }}
          />

          <label htmlFor="maxDevices" style={{ margin: 0 }}><strong>Max Devices:</strong></label>
          <StyledInput
            id="maxDevices"
            type="number"
            value={maxDevices}
            onChange={(e) => setMaxDevices(parseInt(e.target.value, 10))}
            style={{ boxSizing: 'border-box', width: '100%', minWidth: 0, wordBreak: 'break-word' }}
          />

          <label htmlFor="zoneColor" style={{ margin: 0 }}><strong>Color:</strong></label>
          <input
            type="color"
            id="zoneColor"
            value={color}
            onChange={(e) => setColor(e.target.value)}
            style={{ padding: 0, border: 'none', background: 'transparent', height: '25px', width: '50px' }}
          />
          <div style={{ gridColumn: 'span 2', display: 'flex', justifyContent: 'center' }}>
            <PopupButtonContainer>
              <TileButton onClick={handleSave}>Save Changes</TileButton>
            </PopupButtonContainer>
          </div>
        </div>
      </PopupContainer>
    </Draggable>
  );
};

export default function NewMapPage() {
  const navigate = useNavigate();
  const [maps, setMaps] = useState([]);
  const [selectedMapId, setSelectedMapId] = useState('');
  const [devices, setDevices] = useState([]);
  const [beacons, setBeacons] = useState([]);
  const [zones, setZones] = useState([]);
  const [allDevices, setAllDevices] = useState([]);
  const [stagedDeviceFilters, setStagedDeviceFilters] = useState([]);
  const [activeDeviceFilters, setActiveDeviceFilters] = useState([]);
  const [stagedBeaconFilters, setStagedBeaconFilters] = useState([]);
  const [activeBeaconFilters, setActiveBeaconFilters] = useState([]);
  const [showDevices, setShowDevices] = useState(true);
  const [showBeacons, setShowBeacons] = useState(true);
  const [showZones, setShowZones] = useState(true);
  const [showZoneNames, setShowZoneNames] = useState(true);
  const [isPanelVisible, setIsPanelVisible] = useState(true);
  const [selectedDevice, setSelectedDevice] = useState(null);
  const [selectedBeacon, setSelectedBeacon] = useState(null);
  const [selectedZone, setSelectedZone] = useState(null);
  const [allBeacons, setAllBeacons] = useState([]);
  const [allZones, setAllZones] = useState([]);
  const [editingZone, setEditingZone] = useState(null);
  const [isZoneModalOpen, setIsZoneModalOpen] = useState(false);
  const [isBeaconModalOpen, setIsBeaconModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const ws = useRef(null);
  const [contextMenu, setContextMenu] = useState(null);
  const [isPanningDisabled, setIsPanningDisabled] = useState(false);
  const [transformState, setTransformState] = useState({ scale: 1, positionX: 0, positionY: 0 });
  const mapContentRef = useRef(null);
  const [deviceHistory, setDeviceHistory] = useState([]);
  const [historyTotal, setHistoryTotal] = useState(0);
  const [historyCurrentIndex, setHistoryCurrentIndex] = useState(0);
  const [historyLoading, setHistoryLoading] = useState(false);
  const [historyError, setHistoryError] = useState(null);
  const [deviceForHistory, setDeviceForHistory] = useState(null);
  const [historyFrom, setHistoryFrom] = useState('2021-09-01T00:00');
  const [historyTo, setHistoryTo] = useState('2025-12-01T23:59');
  const [clickedPoint, setClickedPoint] = useState(null);

  const handleZoneDragStop = (e, data) => {
    const { scale, positionX, positionY } = transformState;

    const worldX = (data.x - positionX) / scale;
    const worldY = (data.y - positionY) / scale;

    setEditingZone(prev => ({
      ...prev,
      position_x: worldX,
      position_y: worldY
    }));

    setIsPanningDisabled(false);
  };

  const handleMouseDown = (e) => {
    if (e.button === 2) { // Right mouse button
      setIsPanningDisabled(true);
    }
  };

  const handleMouseUp = (e) => {
    if (e.button === 2) { // Right mouse button
      setIsPanningDisabled(false);
    }
  };

  const handleContextMenu = (e) => {
    e.preventDefault();
    e.stopPropagation();
    if (!mapContentRef.current) return;

    const contentRect = mapContentRef.current.getBoundingClientRect();
    const { scale } = transformState;

    const x = (e.clientX - contentRect.left) / scale;
    const y = (e.clientY - contentRect.top) / scale;

    const clickedZone = zones.find(
      (zone) =>
        x >= zone.position_x &&
        x <= zone.position_x + zone.width &&
        y >= zone.position_y &&
        y <= zone.position_y + zone.height
    );

    setContextMenu({ x, y, zone: clickedZone });
  };

  const handleCloseContextMenu = () => {
    setContextMenu(null);
  };

  const handlePlaceBeacon = () => {
    if (!contextMenu) return;
    setIsBeaconModalOpen(true);
  };

  const handleEditZone = () => {
    if (!contextMenu) return;

    if (contextMenu.zone) {
      setEditingZone(contextMenu.zone);
      setContextMenu(null);
    } else {
      fetchZones();
      setIsZoneModalOpen(true);
    }
  };

  const handleZoneSelect = (zone) => {
    setEditingZone({
      ...zone,
      position_x: contextMenu.x,
      position_y: contextMenu.y,
    });
    setIsZoneModalOpen(false);
    setContextMenu(null);
  };

  const handleShowZoneInfo = () => {
    if (contextMenu?.zone) {
      setSelectedZone(contextMenu.zone);
      setContextMenu(null);
    }
  };

  const handleZoneUpdate = async (updatedZone) => {
    try {
      await apiClient.put(`${GEO_HOST}/api/v1/geo/zones`, updatedZone);
      setZones(prevZones => prevZones.map(z => z.id === updatedZone.id ? updatedZone : z));
      setSelectedZone(null);
    } catch (error) {
      console.error('Failed to update zone:', error);
      alert(`Error updating zone: ${error.message}`);
    }
  };

  const handleBeaconSelect = async (beacon) => {
    console.log('Selected beacon:', beacon);
    setIsBeaconModalOpen(false);
    console.log(contextMenu);
    if (!contextMenu) return;

    const updatedBeacon = {
      ...beacon,
      x: contextMenu.x,
      y: contextMenu.y,
      map_id: selectedMapId,
    };

    setContextMenu(null);

    try {
      await apiClient.put(`${GEO_HOST}/api/v1/geo/beacons`, updatedBeacon);
      alert('The beacon has been successfully relocated!');
    } catch (err) {
      console.error(err);
      alert(`Error: ${err.message}`);
    }
  };

  const handleSaveEditingZone = async () => {
    if (!editingZone) return;

    const zoneToUpdate = {
      ...editingZone,
      map_id: selectedMapId,
    };

    try {
      await apiClient.put(`${GEO_HOST}/api/v1/geo/zones`, zoneToUpdate);
      alert('Zone updated successfully!');
      setEditingZone(null);
    } catch (error) {
      console.error('Failed to update zone:', error);
      alert(`Error: ${error.message}`);
    }
  };

  const handleDeviceClick = async (deviceId) => {
    try {
      const data = await apiClient.get(`${SYSTEM_HOST}/api/v1/system/device/${deviceId}`);
      setSelectedDevice(data);
    } catch (err) {
      console.error(err);
    }
  };

  const handleBeaconClick = async (beaconId) => {
    try {
      const data = await apiClient.get(`${GEO_HOST}/api/v1/geo/beacons/${beaconId}`);
      setSelectedBeacon(data);
    } catch (err) {
      console.error(err);
    }
  };

  const fetchZones = async () => {
    const orgId = localStorage.getItem("selectedOrgId");
    if (!orgId) return;

    try {
      const response = await apiClient.get(`${GEO_HOST}/api/v1/geo/zones`, {
        params: {
          org_id: orgId,
        },
      });
      const filteredZones = response.data.filter(zone => !zone.map_id);
      setAllZones(filteredZones);
    } catch (error) {
      console.error("Failed to fetch zones:", error);
    }
  };

  const fetchHistoryBatch = async (device, skip = 0) => {
    setHistoryLoading(true);
    setHistoryError(null);
    try {
      const from = Math.floor(new Date(historyFrom).getTime() / 1000);
      const to = Math.floor(new Date(historyTo).getTime() / 1000);
      const historyData = await getDeviceHistory(selectedMapId, device.id, from, to, 100, skip);
      
      if (skip === 0) {
        setDeviceHistory(historyData.data);
        setHistoryTotal(historyData.total);
        setHistoryCurrentIndex(0);
      } else {
        setDeviceHistory(prev => [...prev, ...historyData.data]);
      }
    } catch (err) {
      setHistoryError(err.message);
    } finally {
      setHistoryLoading(false);
    }
  };

  const handleShowHistory = async (device) => {
    if (!device) return;
    setDeviceForHistory(device);
    fetchHistoryBatch(device, 0);
    setSelectedDevice(null); // Close the popup
  };

  const closeHistory = () => {
    setDeviceForHistory(null);
    setDeviceHistory([]);
    setHistoryTotal(0);
    setHistoryCurrentIndex(0);
  };

  // Effect for lazy loading history points
  useEffect(() => {
    const shouldLoadMore = historyCurrentIndex >= deviceHistory.length - 10 && deviceHistory.length < historyTotal;
    if (shouldLoadMore && !historyLoading) {
      fetchHistoryBatch(deviceForHistory, deviceHistory.length);
    }
  }, [historyCurrentIndex, deviceHistory, historyTotal, historyLoading, deviceForHistory]);

  const applyFilters = () => {
    setDevices([]);
    setBeacons([]);
    setActiveDeviceFilters(stagedDeviceFilters);
    setActiveBeaconFilters(stagedBeaconFilters);
  };

  const handleDeviceFilterChange = (selectedOptions) => {
    setStagedDeviceFilters(selectedOptions);
  };

  const handleBeaconFilterChange = (selectedOptions) => {
    setStagedBeaconFilters(selectedOptions);
  };

  // Fetch initial data
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const [mapsData, devicesData, beaconsData] = await Promise.all([
          apiClient.get(`${GEO_HOST}/api/v1/geo/map`),
          apiClient.get(`${SYSTEM_HOST}/api/v1/system/devices`),
          apiClient.get(`${GEO_HOST}/api/v1/geo/beacons`),
        ]);
        setMaps(mapsData.data);
        setAllDevices(devicesData.data);
        setAllBeacons(beaconsData.data);
        if (mapsData.data.length > 0) setSelectedMapId(mapsData.data[0].id);
      } catch (e) {
        setError(e.message);
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  // WebSocket connection for devices
  useEffect(() => {
    if (!selectedMapId || !showDevices) {
      if (ws.current) {
        ws.current.close();
      }
      setDevices([]);
      return;
    }

    const orgId = localStorage.getItem("selectedOrgId");
    const wsHost = GEO_HOST.replace('http://', 'ws://');
    let wsUrl = `${wsHost}/api/v1/geo/location/device/monitor?map_id=${selectedMapId}`;
    if (orgId) {
      wsUrl += `&org_id=${orgId}`;
    }
    
    const deviceIds = activeDeviceFilters.map(d => d.value);
    if (deviceIds.length > 0) {
      wsUrl += `&device_ids=${deviceIds.join(',')}`;
    }

    ws.current = new WebSocket(wsUrl);
    ws.current.onmessage = (event) => {
      const messages = event.data.split('\n');
      messages.forEach(messageStr => {
        if (messageStr.trim() === '') return;
        try {
          const message = JSON.parse(messageStr);
          if (message.event) {
            setDevices(prev => {
              const map = new Map(prev.map(d => [d.device_id, d]));
              for (const device of message.event) {
                if (device.x == null || device.y == null) {
                  map.delete(device.device_id);
                } else {
                  map.set(device.device_id, device);
                }
              }
              return Array.from(map.values());
            });
          }
        } catch (e) {
          console.error('Failed to parse WebSocket message:', e);
        }
      });
    };

    const currentWs = ws.current;
    return () => {
      if (currentWs) {
        currentWs.close();
      }
    };
  }, [selectedMapId, activeDeviceFilters, showDevices]);

  // WebSocket connection for beacons
  useEffect(() => {
    if (!selectedMapId || !showBeacons) {
      setBeacons([]);
      return;
    }

    const orgId = localStorage.getItem("selectedOrgId");
    const wsHost = GEO_HOST.replace('http://', 'ws://');
    let wsUrl = `${wsHost}/api/v1/geo/location/beacon/monitor?map_id=${selectedMapId}`;
    if (orgId) {
      wsUrl += `&org_id=${orgId}`;
    }
    
    const beaconIds = activeBeaconFilters.map(b => b.value);
    if (beaconIds.length > 0) {
      wsUrl += `&beacon_ids=${beaconIds.join(',')}`;
    }
    
    const beaconWs = new WebSocket(wsUrl);
    
    beaconWs.onmessage = (event) => {
      const messages = event.data.split('\n');
      messages.forEach(messageStr => {
        if (messageStr.trim() === '') return;
        try {
          const message = JSON.parse(messageStr);
          if (message.event) {
            setBeacons(prev => {
              const map = new Map(prev.map(b => [b.beacon_id, b]));
              for (const beacon of message.event) {
                 if (beacon.x == null || beacon.y == null) {
                  map.delete(beacon.beacon_id);
                } else {
                  map.set(beacon.beacon_id, beacon);
                }
              }
              return Array.from(map.values());
            });
          }
        } catch (e) {
          console.error('Failed to parse WebSocket message for beacons:', e);
        }
      });
    };

    return () => {
      beaconWs.close();
    };
  }, [selectedMapId, activeBeaconFilters, showBeacons]);

  // WebSocket connection for zones
  useEffect(() => {
    if (!selectedMapId || !showZones) {
      setZones([]);
      return;
    }

    const orgId = localStorage.getItem("selectedOrgId");
    const wsHost = GEO_HOST.replace('http://', 'ws://');
    let wsUrl = `${wsHost}/api/v1/geo/zones/monitor?map_id=${selectedMapId}`;
    if (orgId) {
      wsUrl += `&org_id=${orgId}`;
    }
    
    const zoneWs = new WebSocket(wsUrl);
    
    zoneWs.onmessage = (event) => {
      const messages = event.data.split('\n');
      messages.forEach(messageStr => {
        if (messageStr.trim() === '') return;
        try {
          const message = JSON.parse(messageStr);
          if (message.event) {
            setZones(prev => {
              const map = new Map(prev.map(z => [z.id, z]));
              for (const zone of message.event) {
                const normalizedZone = {
                  ...zone,
                  id: zone.zone_id,
                  name: zone.zone_name,
                  position_x: zone.positionX,
                  position_y: zone.positionY,
                };
                if (zone.deleted || zone.positionX == null || zone.positionY == null) {
                  map.delete(normalizedZone.id);
                } else {
                  map.set(normalizedZone.id, normalizedZone);
                }
              }
              return Array.from(map.values());
            });
          }
        } catch (e) {
          console.error('Failed to parse WebSocket message for zones:', e);
        }
      });
    };

    return () => {
      zoneWs.close();
    };
  }, [selectedMapId, showZones]);

  const deviceOptions = allDevices.map(d => ({ value: d.id, label: d.name }));
  const beaconOptions = allBeacons.map(b => ({ value: b.id, label: b.label }));
  const selectedMap = maps.find((map) => map.id === selectedMapId);

  if (loading) return <MapsPageContainer>Loading data...</MapsPageContainer>;
  if (error) return <MapsPageContainer>Error: {error}</MapsPageContainer>;

  return (
    <>
      <Header variant="page" title="Map" />
      <MapsPageContainer>
        {isBeaconModalOpen && (
          <BeaconSelectorModal
            beacons={allBeacons}
            onSelect={handleBeaconSelect}
            onClose={() => setIsBeaconModalOpen(false)}
          />
        )}
        {isZoneModalOpen && (
          <ZoneSelectorModal
            zones={allZones}
            onSelect={handleZoneSelect}
            onClose={() => setIsZoneModalOpen(false)}
          />
        )}
        {selectedDevice && <DeviceInfoPopup device={selectedDevice} onClose={() => setSelectedDevice(null)} onShowHistory={handleShowHistory} />}
        {selectedBeacon && <BeaconInfoPopup beacon={selectedBeacon} onClose={() => setSelectedBeacon(null)} />}
        {selectedZone && <ZoneInfoPopup zone={selectedZone} onClose={() => setSelectedZone(null)} onSave={handleZoneUpdate} />}
        <MapSelectorContainer className={isPanelVisible ? '' : 'hidden'}>
          <ToggleButton onClick={() => setIsPanelVisible(!isPanelVisible)}>
          {isPanelVisible ? 'Hide' : 'Show'}
        </ToggleButton>
        <MapSelect value={selectedMapId} onChange={(e) => setSelectedMapId(e.target.value)}>
          {maps.map((map) => (<option key={map.id} value={map.id}>{map.name}</option>))}
        </MapSelect>

        <FilterBlock>
          <FilterHeader>
            <CheckboxLabel>
              <input type="checkbox" checked={showDevices} onChange={() => setShowDevices(!showDevices)} />
              Devices
            </CheckboxLabel>
          </FilterHeader>
          <Select isMulti options={deviceOptions} value={stagedDeviceFilters} onChange={handleDeviceFilterChange} styles={selectStyles} placeholder="Filter by device..." />
        </FilterBlock>

        <FilterBlock>
          <FilterHeader>
            <CheckboxLabel>
              <input type="checkbox" checked={showBeacons} onChange={() => setShowBeacons(!showBeacons)} />
              Beacons
            </CheckboxLabel>
          </FilterHeader>
          <Select isMulti options={beaconOptions} value={stagedBeaconFilters} onChange={handleBeaconFilterChange} styles={selectStyles} placeholder="Filter by beacon..." />
        </FilterBlock>

        <FilterBlock>
          <FilterHeader>
            <CheckboxLabel>
              <input type="checkbox" checked={showZones} onChange={() => setShowZones(!showZones)} />
              Zones
            </CheckboxLabel>
            <CheckboxLabel>
              <input type="checkbox" checked={showZoneNames} onChange={() => setShowZoneNames(!showZoneNames)} />
              Show Names
            </CheckboxLabel>
          </FilterHeader>
        </FilterBlock>

        <FilterButton onClick={applyFilters}>Apply</FilterButton>

        {selectedMap && <MapInfo><span>Level: {selectedMap.level}</span><span>Size: {(selectedMap.width * selectedMap.scale_x).toFixed(2)}x{(selectedMap.height * selectedMap.scale_y).toFixed(2)} м</span></MapInfo>}
      </MapSelectorContainer>
      <MapWrapper onContextMenu={handleContextMenu} onMouseDown={handleMouseDown} onMouseUp={handleMouseUp} onClick={handleCloseContextMenu}>
        {selectedMap && (
          <TransformWrapper
            panning={{
              disabled: isPanningDisabled,
            }}
            wheel={{
              step: 0.2,
            }}
            pinch={{
              disabled: true,
            }}
            doubleClick={{
              disabled: true,
            }}
            onTransformed={(ref, state) => setTransformState(state)}
          >
            <TransformComponent>
              <div ref={mapContentRef} style={{ position: 'relative', width: selectedMap.width, height: selectedMap.height }}>
                <div style={{ width: '100%', height: '100%' }} dangerouslySetInnerHTML={{ __html: selectedMap.svg_content }} />
                {showZones && zones.map((zone) => (
                  <div
                    key={zone.id}
                    title={zone.name}
                    style={{
                      position: 'absolute',
                      left: `${zone.position_x}px`,
                      top: `${zone.position_y}px`,
                      width: `${zone.width}px`,
                      height: `${zone.height}px`,
                      backgroundColor: zone.color ? `rgba(${parseInt(zone.color.slice(1, 3), 16)}, ${parseInt(zone.color.slice(3, 5), 16)}, ${parseInt(zone.color.slice(5, 7), 16)}, 0.3)` : 'rgba(0, 255, 0, 0.3)',
                      border: `1px solid ${zone.color || 'rgba(0, 255, 0, 0.5)'}`,
                      pointerEvents: 'all',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                    }}
                  >
                    {showZoneNames && (
                      <span style={{ color: '#fff', fontSize: 12 / transformState.scale, textShadow: '1px 1px 2px black' }}>
                        {zone.name}
                      </span>
                    )}
                  </div>
                ))}
                {deviceHistory.length > 0 && (
                  <svg width={selectedMap.width} height={selectedMap.height} style={{ position: 'absolute', top: 0, left: 0, pointerEvents: 'none' }}>
                    <defs>
                      <marker
                        id="arrowhead"
                        markerWidth="5"
                        markerHeight="3.5"
                        refX="0"
                        refY="1.75"
                        orient="auto"
                      >
                        <polygon points="0 0, 5 1.75, 0 3.5" fill="#ff7f50" />
                      </marker>
                    </defs>
                    <g style={{ pointerEvents: 'all' }}>
                      <polyline
                        points={deviceHistory.slice(0, historyCurrentIndex + 1).map(p => `${p.x},${p.y}`).join(' ')}
                        fill="none"
                        stroke="#ff7f50"
                        strokeWidth={2 / transformState.scale}
                        markerMid="url(#arrowhead)"
                        markerEnd="url(#arrowhead)"
                      />
                      {deviceHistory.slice(0, historyCurrentIndex + 1).map((p, index) => (
                        <circle
                          key={p.id}
                          cx={p.x}
                          cy={p.y}
                          r={index === historyCurrentIndex ? 8 / transformState.scale : 5 / transformState.scale}
                          fill={index === historyCurrentIndex ? '#ff1111' : '#ff7f50'}
                          onClick={(e) => { e.stopPropagation(); setClickedPoint(p); }}
                          style={{ cursor: 'pointer', transition: 'r 0.3s ease' }}
                        />
                      ))}
                      {clickedPoint && (
                        <text
                          x={clickedPoint.x > selectedMap.width - 150 ? clickedPoint.x - 130 / transformState.scale : clickedPoint.x + 10 / transformState.scale}
                          y={clickedPoint.y}
                          fontSize={12 / transformState.scale}
                          fill="#ffffff"
                          onClick={(e) => { e.stopPropagation(); setClickedPoint(null); }}
                          style={{ cursor: 'pointer', textAnchor: clickedPoint.x > selectedMap.width - 150 ? 'end' : 'start' }}
                        >
                          {new Date(clickedPoint.timestamp).toLocaleString()}
                        </text>
                      )}
                    </g>
                  </svg>
                )}
                {contextMenu && (
                  <ContextMenu style={{ top: contextMenu.y, left: contextMenu.x }} onClick={(e) => e.stopPropagation()}>
                    <ContextMenuItem onClick={handlePlaceBeacon}>Place a beacon</ContextMenuItem>
                    <ContextMenuItem onClick={handleEditZone}>Edit Zone</ContextMenuItem>
                    {contextMenu.zone && <ContextMenuItem onClick={handleShowZoneInfo}>Show Info</ContextMenuItem>}
                  </ContextMenu>
                )}
                {editingZone && (
                  <Draggable
                    position={{ x: editingZone.position_x, y: editingZone.position_y }}
                    onStart={() => setIsPanningDisabled(true)}
                    onStop={(e, data) => handleZoneDragStop(e, data)}
                  >
                    <div
                      style={{
                        position: 'absolute',
                        width: editingZone.width,
                        height: editingZone.height,
                        border: '2px dashed #fff',
                        cursor: 'move',
                      }}
                    >
                      <div
                        style={{
                          position: 'absolute',
                          bottom: '5px',
                          left: '50%',
                          transform: 'translateX(-50%)',
                          display: 'flex',
                          gap: '10px',
                          zIndex: 1000,
                        }}
                      >
                        <button onClick={handleSaveEditingZone}>Save</button>
                        <button onClick={() => setEditingZone(null)}>Cancel</button>
                      </div>
                      <div
                        style={{
                          position: 'absolute',
                          bottom: '-5px',
                          right: '-5px',
                          width: '10px',
                          height: '10px',
                          backgroundColor: '#fff',
                          cursor: 'nwse-resize',
                        }}
                        onMouseDown={(e) => {
                          e.stopPropagation();
                          const startX = e.pageX;
                          const startY = e.pageY;
                          const startWidth = editingZone.width;
                          const startHeight = editingZone.height;

                          const doDrag = (e) => {
                            setEditingZone((prev) => ({
                              ...prev,
                              width: startWidth + e.pageX - startX,
                              height: startHeight + e.pageY - startY,
                            }));
                          };

                          const stopDrag = () => {
                            document.removeEventListener('mousemove', doDrag, false);
                            document.removeEventListener('mouseup', stopDrag, false);
                          };

                          document.addEventListener('mousemove', doDrag, false);
                          document.addEventListener('mouseup', stopDrag, false);
                        }}
                      />
                    </div>
                  </Draggable>
                )}
                {!deviceForHistory && showDevices && devices.map((device) => (
                  <UserIcon
                    key={device.device_id}
                    title={`${device.device_name} (x: ${device.x.toFixed(2)}, y: ${device.y.toFixed(2)})`}
                    onClick={() => handleDeviceClick(device.device_id)}
                    style={{
                      position: 'absolute',
                      top: `${device.y}px`,
                      left: `${device.x}px`,
                      width: '24px',
                      height: '24px',
                      transform: `translate(-50%, -50%) scale(${1 / transformState.scale})`,
                      cursor: 'pointer',
                    }}
                  />
                ))}
                {!deviceForHistory && showBeacons && beacons.map((beacon) => (
                  <div
                    key={beacon.beacon_id}
                    title={`${beacon.beacon_name} (x: ${beacon.x.toFixed(2)}, y: ${beacon.y.toFixed(2)})`}
                    onClick={() => handleBeaconClick(beacon.beacon_id)}
                    style={{
                      position: 'absolute',
                      top: `${beacon.y}px`,
                      left: `${beacon.x}px`,
                      width: '24px',
                      height: '24px',
                      transform: `translate(-50%, -50%) scale(${1 / transformState.scale})`,
                      cursor: 'pointer',
                    }}
                  >
                    <img
                      src={BeaconIcon}
                      alt={beacon.beacon_name}
                      style={{ width: '100%', height: '100%' }}
                    />
                  </div>
                ))}
              </div>
            </TransformComponent>
          </TransformWrapper>
        )}
        {deviceForHistory && (
          <Draggable handle=".history-header">
            <HistoryControlsContainer>
              <HistoryHeader className="history-header">History for {deviceForHistory.name}</HistoryHeader>
              <ControlRow>
                <label>From:</label>
                <StyledInput type="datetime-local" value={historyFrom} onChange={e => setHistoryFrom(e.target.value)} />
              </ControlRow>
              <ControlRow>
                <label>To:</label>
                <StyledInput type="datetime-local" value={historyTo} onChange={e => setHistoryTo(e.target.value)} />
              </ControlRow>
              <TileButton onClick={() => handleShowHistory(deviceForHistory)}>Apply</TileButton>
              {historyError && <p style={{ color: 'red' }}>{historyError}</p>}
              <PaginationContainer>
                <TileButton
                  onClick={() => setHistoryCurrentIndex(i => Math.max(0, i - 1))}
                  disabled={historyCurrentIndex === 0 || historyLoading}
                >
                  {'< Prev Point'}
                </TileButton>
                <span> {historyCurrentIndex + 1} / {historyTotal} </span>
                <TileButton
                  onClick={() => setHistoryCurrentIndex(i => Math.min(historyTotal - 1, i + 1))}
                  disabled={historyCurrentIndex >= historyTotal - 1 || historyLoading}
                >
                  {'Next Point >'}
                </TileButton>
              </PaginationContainer>
              {deviceHistory[historyCurrentIndex] && (
                <div style={{ fontSize: '12px', textAlign: 'center' }}>
                  <p>
                    Current: {new Date(deviceHistory[historyCurrentIndex].timestamp).toLocaleString()}
                  </p>
                </div>
              )}
              <TileButton onClick={closeHistory}>Close</TileButton>
            </HistoryControlsContainer>
          </Draggable>
        )}
      </MapWrapper>
      </MapsPageContainer>
    </>
  );
}
