import { useState, useEffect } from 'react';
import { api } from '../api';
import Combobox from './Combobox';

const RAM_TYPES     = ['DDR3', 'DDR4', 'DDR5', 'None'];
const STORAGE_TYPES = ['HDD', 'SSD', 'NVMe', 'None'];

/**
 * EditComputerModal
 *
 * Props:
 *   item    : desktop | laptop object (pre-filled form values)
 *   type    : 'desktop' | 'laptop'
 *   refs    : { cpuOpts, osOpts, equipOpts, roomOpts, lmOpts, dmOpts }
 *   onSave  : (updatedItem) => void
 *   onClose : () => void
 */
export default function EditComputerModal({ item, type, refs, onSave, onClose }) {
  const [form, setForm]     = useState(() => {
    // network_connection can arrive as a plain string (new API) or as
    // {connection_type_enum, valid} (old binary without MarshalJSON override).
    // Normalize to plain string | null so the <select> always works.
    const nc = item.network_connection;
    const normalizedNc = nc == null ? null
      : typeof nc === 'string' ? nc
      : nc.connection_type_enum || null;
    return { ...item, network_connection: normalizedNc };
  });
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');

  // Prevent background scroll while modal is open
  useEffect(() => {
    document.body.style.overflow = 'hidden';
    return () => { document.body.style.overflow = ''; };
  }, []);

  function set(key, val) { setForm(f => ({ ...f, [key]: val })); }

  async function handleSubmit(e) {
    e.preventDefault();
    setSaving(true);
    setErr('');
    try {
      const base = {
        hostname:          form.hostname,
        room_id:           form.room_id,
        os_id:             form.os_id             || null,
        ram_gb:            form.ram_gb            ? parseInt(form.ram_gb, 10)     : null,
        ram_type:          form.ram_type          || null,
        storage_gb:        form.storage_gb        ? parseInt(form.storage_gb, 10) : null,
        storage_type:      form.storage_type      || null,
        mac_address:       form.mac_address       || null,
        equipment_user_id: form.equipment_user_id || null,
        observations:      form.observations      || null,
      };

      let updated;
      if (type === 'desktop') {
        updated = await api.updateDesktop(item.computer_id, {
          ...base,
          desktop_model_id:  form.desktop_model_id  || null,
          cpu_id:            form.cpu_id             || null,
          has_wifi_card:     form.has_wifi_card,
          network_connection: form.network_connection || null,
        });
      } else {
        updated = await api.updateLaptop(item.computer_id, {
          ...base,
          laptop_model_id: form.laptop_model_id || null,
          serial_number:   form.serial_number   || null,
        });
      }
      onSave(updated ?? { ...item, ...form });
    } catch (e) {
      setErr(e.message);
      setSaving(false);
    }
  }

  const R = refs;
  const title = type === 'desktop' ? '🖥️ Editar sobretaula' : '💻 Editar portàtil';

  return (
    <div
      style={{
        position: 'fixed', inset: 0, zIndex: 1000,
        background: 'rgba(0,0,0,0.45)',
        display: 'flex', alignItems: 'flex-start', justifyContent: 'center',
        overflowY: 'auto', padding: '32px 16px',
      }}
      onClick={e => { if (e.target === e.currentTarget) onClose(); }}
    >
      <div
        style={{
          background: 'var(--surface)',
          borderRadius: 10,
          boxShadow: '0 8px 32px rgba(0,0,0,0.25)',
          width: '100%',
          maxWidth: 760,
          padding: '28px 32px 32px',
        }}
      >
        {/* Header */}
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 24 }}>
          <h2 style={{ margin: 0, fontSize: 18, fontWeight: 700 }}>{title}</h2>
          <button
            onClick={onClose}
            style={{
              background: 'none', border: 'none', cursor: 'pointer',
              fontSize: 20, color: 'var(--muted)', lineHeight: 1, padding: '4px 8px',
            }}
          >✕</button>
        </div>

        {err && (
          <div style={{
            marginBottom: 16, padding: '10px 14px', borderRadius: 6,
            background: 'rgba(220,50,50,0.08)', color: 'var(--danger)',
            border: '1px solid rgba(220,50,50,0.2)', fontSize: 13,
          }}>
            {err}
          </div>
        )}

        <form onSubmit={handleSubmit} className="form-panel" style={{ padding: 0 }}>
          <div className="form-grid">

            {/* Hostname */}
            <div className="form-group">
              <label>Hostname *</label>
              <input
                type="text"
                value={form.hostname}
                onChange={e => set('hostname', e.target.value)}
                required
              />
            </div>

            {/* Aula */}
            <div className="form-group">
              <label>Aula *</label>
              <Combobox
                options={R.roomOpts}
                value={form.room_id}
                onChange={v => set('room_id', v)}
                placeholder="Selecciona aula…"
              />
            </div>

            {/* Model */}
            {type === 'desktop' ? (
              <div className="form-group">
                <label>Model base</label>
                <Combobox
                  options={R.dmOpts}
                  value={form.desktop_model_id}
                  onChange={v => set('desktop_model_id', v)}
                  placeholder="Sense model…"
                  nullable
                />
              </div>
            ) : (
              <div className="form-group">
                <label>Model *</label>
                <Combobox
                  options={R.lmOpts}
                  value={form.laptop_model_id}
                  onChange={v => set('laptop_model_id', v)}
                  placeholder="Selecciona model…"
                />
              </div>
            )}

            {/* CPU (desktop only) */}
            {type === 'desktop' && (
              <div className="form-group">
                <label>CPU</label>
                <Combobox
                  options={R.cpuOpts}
                  value={form.cpu_id}
                  onChange={v => set('cpu_id', v)}
                  placeholder="Sense CPU…"
                  nullable
                />
              </div>
            )}

            {/* Serial number (laptop only) */}
            {type === 'laptop' && (
              <div className="form-group">
                <label>Número de sèrie</label>
                <input
                  type="text"
                  value={form.serial_number ?? ''}
                  onChange={e => set('serial_number', e.target.value)}
                  placeholder="SN123456789"
                />
              </div>
            )}

            {/* OS */}
            <div className="form-group">
              <label>Sistema operatiu</label>
              <Combobox
                options={R.osOpts}
                value={form.os_id}
                onChange={v => set('os_id', v)}
                placeholder="Sense SO…"
                nullable
              />
            </div>

            {/* RAM */}
            <div className="form-group">
              <label>RAM (GB)</label>
              <input
                type="number"
                min="0"
                value={form.ram_gb ?? ''}
                onChange={e => set('ram_gb', e.target.value)}
                placeholder="ex: 8"
              />
            </div>

            <div className="form-group">
              <label>Tipus RAM</label>
              <select value={form.ram_type ?? ''} onChange={e => set('ram_type', e.target.value)}>
                <option value="">—</option>
                {RAM_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>

            {/* Storage */}
            <div className="form-group">
              <label>Emmagatzematge (GB)</label>
              <input
                type="number"
                min="0"
                value={form.storage_gb ?? ''}
                onChange={e => set('storage_gb', e.target.value)}
                placeholder="ex: 256"
              />
            </div>

            <div className="form-group">
              <label>Tipus emmagatzematge</label>
              <select value={form.storage_type ?? ''} onChange={e => set('storage_type', e.target.value)}>
                <option value="">—</option>
                {STORAGE_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>

            {/* WiFi (desktop only) */}
            {type === 'desktop' && (
              <div className="form-group" style={{ display: 'flex', alignItems: 'center', gap: 10, paddingTop: 24 }}>
                <input
                  id="edit-wifi"
                  type="checkbox"
                  checked={form.has_wifi_card ?? false}
                  onChange={e => set('has_wifi_card', e.target.checked)}
                />
                <label htmlFor="edit-wifi" style={{ margin: 0 }}>Té targeta WiFi</label>
              </div>
            )}

            {/* Connexió de xarxa (desktop only) */}
            {type === 'desktop' && (
              <div className="form-group">
                <label>Connexió de xarxa</label>
                <select
                  value={form.network_connection ?? ''}
                  onChange={e => set('network_connection', e.target.value || null)}
                >
                  <option value="">—</option>
                  <option value="ethernet">Ethernet (cablejat)</option>
                  <option value="wifi">WiFi (sense fil)</option>
                </select>
              </div>
            )}

            {/* MAC */}
            <div className="form-group">
              <label>Adreça MAC</label>
              <input
                type="text"
                value={form.mac_address ?? ''}
                onChange={e => set('mac_address', e.target.value)}
                placeholder="AA:BB:CC:DD:EE:FF"
              />
            </div>

            {/* Usuari equip */}
            <div className="form-group">
              <label>Usuari equip</label>
              <Combobox
                options={R.equipOpts}
                value={form.equipment_user_id}
                onChange={v => set('equipment_user_id', v)}
                placeholder="Sense usuari…"
                nullable
              />
            </div>

            {/* Observations */}
            <div className="form-group wide">
              <label>Observacions</label>
              <textarea
                value={form.observations ?? ''}
                onChange={e => set('observations', e.target.value)}
                placeholder="Opcional"
                rows={2}
              />
            </div>

          </div>

          <div style={{ display: 'flex', gap: 10, justifyContent: 'flex-end', marginTop: 24 }}>
            <button type="button" className="btn" onClick={onClose} disabled={saving}>
              Cancel·lar
            </button>
            <button type="submit" className="btn btn-primary" disabled={saving}>
              {saving ? 'Desant…' : 'Desar canvis'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
