import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import Select from 'react-select';
import { TransformWrapper, TransformComponent } from 'react-zoom-pan-pinch';
import Draggable from 'react-draggable';
import UserIcon from '../../assets/user.svg';
import BeaconIcon from '../../assets/bullseye-animated.gif';
import { GEO_HOST, SYSTEM_HOST } from '../../config';
import { apiClient } from '../../utils/api';
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

export default function NewMapPage() {
  const navigate = useNavigate();
  const [maps, setMaps] = useState([]);
  const [selectedMapId, setSelectedMapId] = useState('');
  const [devices, setDevices] = useState([]);
  const [beacons, setBeacons] = useState([]);
  const [allDevices, setAllDevices] = useState([]);
  const [stagedDeviceFilters, setStagedDeviceFilters] = useState([]);
  const [activeDeviceFilters, setActiveDeviceFilters] = useState([]);
  const [stagedBeaconFilters, setStagedBeaconFilters] = useState([]);
  const [activeBeaconFilters, setActiveBeaconFilters] = useState([]);
  const [showDevices, setShowDevices] = useState(true);
  const [showBeacons, setShowBeacons] = useState(true);
  const [isPanelVisible, setIsPanelVisible] = useState(true);
  const [selectedDevice, setSelectedDevice] = useState(null);
  const [selectedBeacon, setSelectedBeacon] = useState(null);
  const [allBeacons, setAllBeacons] = useState([]);
  const [isBeaconModalOpen, setIsBeaconModalOpen] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const ws = useRef(null);
  const [contextMenu, setContextMenu] = useState(null);
  const [isPanningDisabled, setIsPanningDisabled] = useState(false);
  const [transformState, setTransformState] = useState({ scale: 1, positionX: 0, positionY: 0 });
  const mapContentRef = useRef(null);

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

    setContextMenu({ x, y });
  };

  const handleCloseContextMenu = () => {
    setContextMenu(null);
  };

  const handlePlaceBeacon = () => {
    if (!contextMenu) return;
    setIsBeaconModalOpen(true);
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
        {selectedDevice && <DeviceInfoPopup device={selectedDevice} onClose={() => setSelectedDevice(null)} />}
        {selectedBeacon && <BeaconInfoPopup beacon={selectedBeacon} onClose={() => setSelectedBeacon(null)} />}
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
                {contextMenu && (
                  <ContextMenu style={{ top: contextMenu.y, left: contextMenu.x }} onClick={(e) => e.stopPropagation()}>
                    <ContextMenuItem onClick={handlePlaceBeacon}>Place a beacon</ContextMenuItem>
                  </ContextMenu>
                )}
                {showDevices && devices.map((device) => (
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
                {showBeacons && beacons.map((beacon) => (
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
      </MapWrapper>
      </MapsPageContainer>
    </>
  );
}
