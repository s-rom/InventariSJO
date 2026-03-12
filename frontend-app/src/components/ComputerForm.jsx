import { useState } from 'react';
import { api } from '../api';
import Combobox from './Combobox';

const RAM_TYPES     = ['None', 'DDR3', 'DDR4', 'DDR5'];
const STORAGE_TYPES = ['None', 'HDD', 'SSD', 'NVMe'];
const COMPUTER_TYPES = ['Desktop', 'Laptop'];

const EMPTY = {
  hostname:           '',
  cpu_id:             null,
  ram_gb:             0,
  ram_type:           'None',
  storage_gb:         0,
  storage_type:       'None',
  computer_type:      'Desktop',
  observations:       '',
  equipment_user_id:  null,
  room_id:            null,
  mac_address:        '',
  os_ids:             [],
};

/**
 * Props:
 *   cpus       : { value: number, label: string }[]
 *   rooms      : { value: number, label: string }[]
 *   equipUsers : { value: number, label: string }[]
 *   osList     : { os_id, name }[]
 *   onCreated  : (computer) => void
 *   onCancel   : () => void
 */
export default function ComputerForm({ cpus, rooms, equipUsers, osList, onCreated, onCancel }) {
  const [form, setForm]     = useState(EMPTY);
  const [error, setError]   = useState('');
  const [saving, setSaving] = useState(false);

  function set(field, value) {
    setForm(f => ({ ...f, [field]: value }));
  }

  function toggleOS(id) {
    setForm(f => ({
      ...f,
      os_ids: f.os_ids.includes(id)
        ? f.os_ids.filter(x => x !== id)
        : [...f.os_ids, id],
    }));
  }

  async function handleSubmit(e) {
    e.preventDefault();
    if (!form.hostname.trim()) { setError('El hostname és obligatori.'); return; }
    setError('');
    setSaving(true);
    try {
      const payload = {
        hostname:          form.hostname.trim(),
        cpu_id:            form.cpu_id,
        ram_gb:            Number(form.ram_gb),
        ram_type:          form.ram_type,
        storage_gb:        Number(form.storage_gb),
        storage_type:      form.storage_type,
        computer_type:     form.computer_type,
        observations:      form.observations.trim() || null,
        equipment_user_id: form.equipment_user_id,
        room_id:           form.room_id,
        mac_address:       form.mac_address.trim() || null,
        os_ids:            form.os_ids,
      };
      const computer = await api.createComputer(payload);
      setForm(EMPTY);
      onCreated(computer);
    } catch (err) {
      setError(err.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="card" style={{ marginTop: 16 }}>
      <form onSubmit={handleSubmit}>
        <div className="form-panel">
          <div className="form-title">Nou equip</div>
          <div className="form-grid">

            <div className="form-group">
              <label>Hostname *</label>
              <input type="text" value={form.hostname} onChange={e => set('hostname', e.target.value)} required />
            </div>

            <div className="form-group">
              <label>Tipus</label>
              <select value={form.computer_type} onChange={e => set('computer_type', e.target.value)}>
                {COMPUTER_TYPES.map(t => <option key={t}>{t}</option>)}
              </select>
            </div>

            <div className="form-group">
              <label>CPU</label>
              <Combobox
                options={cpus}
                value={form.cpu_id}
                onChange={v => set('cpu_id', v)}
                placeholder="Cerca CPU…"
                nullable
              />
            </div>

            <div className="form-group">
              <label>RAM (GB)</label>
              <input type="number" min="0" value={form.ram_gb} onChange={e => set('ram_gb', e.target.value)} />
            </div>

            <div className="form-group">
              <label>Tipus RAM</label>
              <select value={form.ram_type} onChange={e => set('ram_type', e.target.value)}>
                {RAM_TYPES.map(t => <option key={t}>{t}</option>)}
              </select>
            </div>

            <div className="form-group">
              <label>Emmagatzematge (GB)</label>
              <input type="number" min="0" value={form.storage_gb} onChange={e => set('storage_gb', e.target.value)} />
            </div>

            <div className="form-group">
              <label>Tipus emmagatzematge</label>
              <select value={form.storage_type} onChange={e => set('storage_type', e.target.value)}>
                {STORAGE_TYPES.map(t => <option key={t}>{t}</option>)}
              </select>
            </div>

            <div className="form-group">
              <label>Sala</label>
              <Combobox
                options={rooms}
                value={form.room_id}
                onChange={v => set('room_id', v)}
                placeholder="Cerca sala…"
                nullable
              />
            </div>

            <div className="form-group">
              <label>Usuari equip</label>
              <Combobox
                options={equipUsers}
                value={form.equipment_user_id}
                onChange={v => set('equipment_user_id', v)}
                placeholder="Cerca usuari…"
                nullable
              />
            </div>

            <div className="form-group">
              <label>Adreça MAC</label>
              <input type="text" value={form.mac_address} onChange={e => set('mac_address', e.target.value)} placeholder="Opcional" />
            </div>

            <div className="form-group wide">
              <label>Observacions</label>
              <textarea value={form.observations} onChange={e => set('observations', e.target.value)} placeholder="Opcional" />
            </div>

            <div className="form-group wide">
              <label>Sistema operatiu</label>
              <div className="os-grid">
                {osList.map(os => (
                  <label
                    key={os.os_id}
                    className={`os-check${form.os_ids.includes(os.os_id) ? ' selected' : ''}`}
                  >
                    <input
                      type="checkbox"
                      style={{ display: 'none' }}
                      checked={form.os_ids.includes(os.os_id)}
                      onChange={() => toggleOS(os.os_id)}
                    />
                    {os.name}
                  </label>
                ))}
                {osList.length === 0 && <span style={{ color: 'var(--muted)', fontSize: 12 }}>Cap SO registrat</span>}
              </div>
            </div>

          </div>

          {error && <div className="error-msg" style={{ marginTop: 12 }}>{error}</div>}

          <div className="form-actions">
            <button type="submit" className="btn btn-primary" disabled={saving}>
              {saving ? 'Desant…' : 'Desar equip'}
            </button>
            <button type="button" className="btn btn-ghost" onClick={onCancel}>
              Cancel·lar
            </button>
          </div>
        </div>
      </form>
    </div>
  );
}
