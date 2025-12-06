import React, { useState } from 'react';
import Draggable from 'react-draggable';
import { PopupContainer, PopupHeader, CloseButton, TileButton, StyledInput } from './styles.js';
import { apiClient } from '../../utils/api';
import { GEO_HOST } from '../../config';

export default function ZoneDetailsPopup({ zone, onClose }) {
  if (!zone) return null;
  const info = zone.info || {};

  const [name, setName] = useState(zone.name);
  const [color, setColor] = useState(zone.color);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);

  // Синхронизация name/color при изменении zone (если не редактируется)
  React.useEffect(() => {
    if (!saving) {
      setName(zone.name);
      setColor(zone.color);
    }
  }, [zone.name, zone.color]);

  const handleSave = async () => {
    setSaving(true);
    setError(null);
    setSuccess(false);
    try {
      await apiClient.put(`${GEO_HOST}/api/v1/geo/zones`, {
        ...zone,
        name,
        color,
      });
      setSuccess(true);
      setTimeout(() => setSuccess(false), 1500);
    } catch (e) {
      setError(e.message || 'Failed to save');
    } finally {
      setSaving(false);
    }
  };

  return (
    <Draggable handle=".popup-header">
      <PopupContainer>
        <PopupHeader className="popup-header">
          <h2>Zone Information</h2>
          <CloseButton onClick={onClose}>&times;</CloseButton>
        </PopupHeader>
        <div
          style={{
            padding: '0 12px',
            wordBreak: 'break-word',
            display: 'grid',
            gridTemplateColumns: '120px auto',
            alignItems: 'center',
            gap: '10px',
            width: '100%',
            maxWidth: '100%',
            boxSizing: 'border-box',
          }}
        >
          <p style={{ margin: 0, wordBreak: 'break-word' }}><strong>ID:</strong></p>
          <p style={{ margin: 0, wordBreak: 'break-word' }}>{zone.id}</p>

          <label htmlFor="zoneName" style={{ margin: 0 }}><strong>Name:</strong></label>
          <StyledInput
            id="zoneName"
            type="text"
            value={name}
            onChange={e => setName(e.target.value)}
            style={{ boxSizing: 'border-box', width: '100%', minWidth: 0, wordBreak: 'break-word' }}
            disabled={saving}
          />

          <p style={{ margin: 0 }}><strong>Type:</strong></p>
          <p style={{ margin: 0 }}>{zone.type}</p>

          <p style={{ margin: 0 }}><strong>Size:</strong></p>
          <p style={{ margin: 0 }}>{zone.width} x {zone.height}</p>

          <label htmlFor="zoneColor" style={{ margin: 0 }}><strong>Color:</strong></label>
          <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
            <input
              type="color"
              id="zoneColor"
              value={color}
              onChange={e => setColor(e.target.value)}
              style={{ padding: 0, border: 'none', background: 'transparent', height: '25px', width: '50px' }}
              disabled={saving}
            />
            <span style={{ background: color, padding: '0 10px', borderRadius: 4, color: '#fff', fontSize: 12 }}>{color}</span>
          </div>

          <div style={{ gridColumn: 'span 2', margin: '10px 0' }}>
            <hr />
          </div>

          <p style={{ margin: 0 }}><strong>Device count:</strong></p>
          <p style={{ margin: 0 }}>{info.count || '—'}</p>
          <p style={{ margin: 0 }}><strong>Devices:</strong></p>
          <p style={{ margin: 0 }}>{info.devices || '—'}</p>
          <p style={{ margin: 0 }}><strong>Last entered:</strong></p>
          <p style={{ margin: 0 }}>{info.last_entered || '—'}</p>
          <p style={{ margin: 0 }}><strong>Last exited:</strong></p>
          <p style={{ margin: 0 }}>{info.last_exited || '—'}</p>

          {error && (
            <div style={{ gridColumn: 'span 2', color: 'red', marginTop: 8 }}>{error}</div>
          )}
          {success && (
            <div style={{ gridColumn: 'span 2', color: 'green', marginTop: 8 }}>Saved!</div>
          )}
          <div style={{ gridColumn: 'span 2', marginTop: 16, textAlign: 'center' }}>
            <TileButton as="button" onClick={handleSave} disabled={saving || (name === zone.name && color === zone.color)}>
              {saving ? 'Saving...' : 'Save'}
            </TileButton>
          </div>
        </div>
      </PopupContainer>
    </Draggable>
  );
}
