import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';
import Combobox from '../components/Combobox';

const RAM_TYPES     = ['DDR3', 'DDR4', 'DDR5', 'None'];
const STORAGE_TYPES = ['HDD', 'SSD', 'NVMe', 'None'];

const EMPTY = {
  hostname:         '',
  room_id:          null,
  desktop_model_id: null,
  cpu_id:           null,
  os_id:            null,
  ram_gb:           '',
  ram_type:         '',
  storage_gb:       '',
  storage_type:     '',
  has_wifi_card:    false,
  mac_address:      '',
  equipment_user_id: null,
  observations:     '',
};

export default function NewDesktop() {
  const navigate = useNavigate();
  const [form, setForm]     = useState(EMPTY);
  const [refs, setRefs]     = useState(null);
  const [loading, setLoading] = useState(true);
  const [saving,  setSaving]  = useState(false);
  const [err,     setErr]     = useState('');

  useEffect(() => {
    async function load() {
      try {
        const [cpus, osList, equip, centers, dm] = await Promise.all([
          api.listCpus(),
          api.listOS(),
          api.listEquipmentUsers(),
          api.listCenters(),
          api.listDesktopModels(),
        ]);

        const roomsNested = await Promise.all(
          (centers ?? []).map(c =>
            api.listRoomsByCenter(c.center_id)
              .then(rs => (rs ?? []).map(r => ({ ...r, centerName: c.name })))
              .catch(() => [])
          )
        );
        const allRooms = roomsNested.flat();

        setRefs({
          cpuOpts:   (cpus   ?? []).map(c => ({ value: c.cpu_id,               label: c.model_name })),
          osOpts:    (osList ?? []).map(o => ({ value: o.os_id,                label: o.name })),
          equipOpts: (equip  ?? []).map(e => ({ value: e.equipment_user_id,    label: e.name })),
          roomOpts:  allRooms.map(r => ({ value: r.room_id, label: `${r.centerName} › ${r.name}` })),
          dmOpts:    (dm     ?? []).map(m => ({ value: m.desktop_model_id,     label: `${m.brand_name} ${m.model_name}` })),
        });
      } catch (e) {
        setErr(e.message);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  function set(key, val) { setForm(f => ({ ...f, [key]: val })); }

  async function handleSubmit(e) {
    e.preventDefault();
    setErr('');
    setSaving(true);
    try {
      const body = {
        hostname:          form.hostname,
        room_id:           form.room_id,
        desktop_model_id:  form.desktop_model_id || null,
        cpu_id:            form.cpu_id || null,
        os_id:             form.os_id || null,
        ram_gb:            form.ram_gb    ? parseInt(form.ram_gb,    10) : null,
        ram_type:          form.ram_type  || null,
        storage_gb:        form.storage_gb ? parseInt(form.storage_gb, 10) : null,
        storage_type:      form.storage_type || null,
        has_wifi_card:     form.has_wifi_card,
        mac_address:       form.mac_address || null,
        equipment_user_id: form.equipment_user_id || null,
        observations:      form.observations || null,
      };
      await api.createDesktop(body);
      navigate('/computers');
    } catch (e) {
      setErr(e.message);
    } finally {
      setSaving(false);
    }
  }

  if (loading) return <div className="empty">Carregant…</div>;

  const R = refs;

  return (
    <>
      <div className="page-header">
        <h1 className="page-title">🖥️ Nou sobretaula</h1>
      </div>

      <div className="card">
        <form onSubmit={handleSubmit} className="form-panel">
          <div className="form-grid">

            {/* Hostname */}
            <div className="form-group">
              <label>Hostname *</label>
              <input
                type="text"
                value={form.hostname}
                onChange={e => set('hostname', e.target.value)}
                placeholder="aula01-pc01"
                required
              />
            </div>

            {/* Room */}
            <div className="form-group">
              <label>Aula *</label>
              <Combobox
                options={R.roomOpts}
                value={form.room_id}
                onChange={v => set('room_id', v)}
                placeholder="Selecciona aula…"
              />
            </div>

            {/* Desktop model (optional) */}
            <div className="form-group">
              <label>Model base (opcional)</label>
              <Combobox
                options={R.dmOpts}
                value={form.desktop_model_id}
                onChange={v => set('desktop_model_id', v)}
                placeholder="Sense model base…"
                nullable
              />
            </div>

            {/* CPU */}
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
                value={form.ram_gb}
                onChange={e => set('ram_gb', e.target.value)}
                placeholder="ex: 8"
              />
            </div>

            <div className="form-group">
              <label>Tipus RAM</label>
              <select value={form.ram_type} onChange={e => set('ram_type', e.target.value)}>
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
                value={form.storage_gb}
                onChange={e => set('storage_gb', e.target.value)}
                placeholder="ex: 256"
              />
            </div>

            <div className="form-group">
              <label>Tipus emmagatzematge</label>
              <select value={form.storage_type} onChange={e => set('storage_type', e.target.value)}>
                <option value="">—</option>
                {STORAGE_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>

            {/* Wifi */}
            <div className="form-group" style={{ justifyContent: 'flex-end' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer' }}>
                <input
                  type="checkbox"
                  style={{ width: 'auto' }}
                  checked={form.has_wifi_card}
                  onChange={e => set('has_wifi_card', e.target.checked)}
                />
                Té targeta WiFi
              </label>
            </div>

            {/* MAC address */}
            <div className="form-group">
              <label>Adreça MAC {form.has_wifi_card ? '*' : ''}</label>
              <input
                type="text"
                value={form.mac_address}
                onChange={e => set('mac_address', e.target.value)}
                placeholder="AA:BB:CC:DD:EE:FF"
                required={form.has_wifi_card}
              />
            </div>

            {/* Equipment user */}
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
                value={form.observations}
                onChange={e => set('observations', e.target.value)}
                placeholder="Notes opcionals…"
              />
            </div>

          </div>

          {err && <div className="error-msg" style={{ marginTop: 12 }}>{err}</div>}

          <div className="form-actions">
            <button type="submit" className="btn btn-primary" disabled={saving || !form.hostname || !form.room_id}>
              {saving ? 'Guardant…' : 'Crear sobretaula'}
            </button>
            <button type="button" className="btn btn-ghost" onClick={() => navigate('/computers')}>
              Cancel·lar
            </button>
          </div>
        </form>
      </div>
    </>
  );
}
