import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';
import Combobox from '../components/Combobox';
import MiniList from '../components/MiniList';

const RAM_TYPES     = ['DDR3', 'DDR4', 'DDR5', 'None'];
const STORAGE_TYPES = ['HDD', 'SSD', 'NVMe', 'None'];

const EMPTY = {
  hostname:          '',
  room_id:           null,
  desktop_model_id:  null,
  cpu_id:            null,
  os_id:             null,
  ram_gb:            '',
  ram_type:          '',
  storage_gb:        '',
  storage_type:      '',
  has_wifi_card:     false,
  mac_address:       '',
  network_connection: null,
  equipment_user_id: null,
  observations:      '',
};

export default function NewDesktop() {
  const navigate = useNavigate();
  const [form, setForm]     = useState(EMPTY);
  const [refs, setRefs]     = useState(null);
  const [desktops, setDesktops] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving,  setSaving]  = useState(false);
  const [toast,   setToast]   = useState(null); // { type: 'ok'|'err', msg }
  const toastTimer = useRef(null);

  function showToast(type, msg) {
    setToast({ type, msg });
    clearTimeout(toastTimer.current);
    toastTimer.current = setTimeout(() => setToast(null), 4000);
  }

  useEffect(() => {
    async function load() {
      try {
        const [cpus, osList, equip, centers, dm, dt] = await Promise.all([
          api.listCpus(),
          api.listOS(),
          api.listEquipmentUsers(),
          api.listCenters(),
          api.listDesktopModels(),
          api.listDesktops(),
        ]);

        const roomsNested = await Promise.all(
          (centers ?? []).map(c =>
            api.listRoomsByCenter(c.center_id)
              .then(rs => (rs ?? []).map(r => ({ ...r, centerName: c.name })))
              .catch(() => [])
          )
        );
        const allRooms = roomsNested.flat();

        const cpuOpts   = (cpus   ?? []).map(c => ({ value: c.cpu_id,               label: c.model_name }));
        const osOpts    = (osList ?? []).map(o => ({ value: o.os_id,                label: o.name }));
        const equipOpts = (equip  ?? []).map(e => ({ value: e.equipment_user_id,    label: e.name }));
        const roomOpts  = allRooms.map(r => ({ value: r.room_id, label: `${r.centerName} › ${r.name}` }));
        const dmOpts    = (dm     ?? []).map(m => ({ value: m.desktop_model_id,     label: `${m.brand_name} ${m.model_name}` }));

        const cpuMap   = Object.fromEntries(cpuOpts.map(o => [o.value, o.label]));
        const osMap    = Object.fromEntries(osOpts.map(o  => [o.value, o.label]));
        const roomMap  = Object.fromEntries(roomOpts.map(o => [o.value, o.label]));
        const dmMap    = Object.fromEntries(dmOpts.map(o => [o.value, o.label]));

        setRefs({ cpuOpts, osOpts, equipOpts, roomOpts, dmOpts, cpuMap, osMap, roomMap, dmMap });
        setDesktops(dt ?? []);
      } catch (e) {
        showToast('err', e.message);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  function set(key, val) { setForm(f => ({ ...f, [key]: val })); }

  async function handleSubmit(e) {
    e.preventDefault();
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
        network_connection: form.network_connection || null,
        equipment_user_id: form.equipment_user_id || null,
        observations:      form.observations || null,
      };
      const created = await api.createDesktop(body);
      setDesktops(prev => [created, ...prev]);
      setForm(EMPTY);
      showToast('ok', `Sobretaula “${body.hostname}” afegit correctament.`);
    } catch (e) {
      showToast('err', e.message);
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
        <button className="btn" onClick={() => navigate('/computers')}>← Tornar</button>
      </div>

      {/* TOAST */}
      {toast && (
        <div style={{
          marginBottom: 16, padding: '10px 16px', borderRadius: 6, fontSize: 14,
          background: toast.type === 'ok' ? 'rgba(34,197,94,0.1)' : 'rgba(220,50,50,0.08)',
          color:      toast.type === 'ok' ? '#16a34a'              : 'var(--danger)',
          border:     `1px solid ${toast.type === 'ok' ? 'rgba(34,197,94,0.3)' : 'rgba(220,50,50,0.2)'}`,
        }}>
          {toast.type === 'ok' ? '✔️ ' : '⚠️ '}{toast.msg}
        </div>
      )}

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
              <label>Adreça MAC</label>
              <input
                type="text"
                value={form.mac_address}
                onChange={e => set('mac_address', e.target.value)}
                placeholder="AA:BB:CC:DD:EE:FF"
              />
            </div>

            {/* Network connection */}
            <div className="form-group">
              <label>Connexió de xarxa</label>
              <select value={form.network_connection ?? ''} onChange={e => set('network_connection', e.target.value || null)}>
                <option value="">—</option>
                <option value="ethernet">Ethernet (cablejat)</option>
                <option value="wifi">WiFi (sense fil)</option>
              </select>
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

          <div className="form-actions">
            <button type="submit" className="btn btn-primary" disabled={saving || !form.hostname || !form.room_id}>
              {saving ? 'Guardant…' : 'Afegir sobretaula'}
            </button>
          </div>
        </form>
      </div>

      {/* PREVIEW LIST */}
      <div className="page-header" style={{ marginTop: 40 }}>
        <h2 className="page-title" style={{ fontSize: 18 }}>🖥️ Sobretaules registrats ({desktops.length})</h2>
      </div>
      <div className="card">
        <MiniList items={desktops} type="desktop" refs={R} />
      </div>
    </>
  );
}
