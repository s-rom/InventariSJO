import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';
import Combobox from '../components/Combobox';
import MiniList from '../components/MiniList';

const RAM_TYPES     = ['DDR3', 'DDR4', 'DDR5', 'None'];
const STORAGE_TYPES = ['HDD', 'SSD', 'NVMe', 'None'];

const EMPTY = {
  hostname:          '',
  serial_number:     '',
  room_id:           null,
  laptop_model_id:   null,
  os_id:             null,
  ram_gb:            '',
  ram_type:          '',
  storage_gb:        '',
  storage_type:      '',
  mac_address:       '',
  equipment_user_id: null,
  observations:      '',
};

export default function NewLaptop() {
  const navigate = useNavigate();

  const [form,    setForm]    = useState(EMPTY);
  const [refs,    setRefs]    = useState(null);
  const [modelMap, setModelMap] = useState({});  // id -> full model object
  const [laptops, setLaptops] = useState([]);
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
        const [osList, equip, centers, lm, lt] = await Promise.all([
          api.listOS(),
          api.listEquipmentUsers(),
          api.listCenters(),
          api.listLaptopModels(),
          api.listLaptops(),
        ]);

        const roomsNested = await Promise.all(
          (centers ?? []).map(c =>
            api.listRoomsByCenter(c.center_id)
              .then(rs => (rs ?? []).map(r => ({ ...r, centerName: c.name })))
              .catch(() => [])
          )
        );
        const allRooms = roomsNested.flat();

        const osOpts   = (osList ?? []).map(o => ({ value: o.os_id,             label: o.name }));
        const roomOpts = allRooms.map(r => ({ value: r.room_id, label: `${r.centerName} › ${r.name}` }));
        const osMap    = Object.fromEntries(osOpts.map(o  => [o.value, o.label]));
        const roomMap  = Object.fromEntries(roomOpts.map(o => [o.value, o.label]));

        const lmModels = lm ?? [];
        const lmOpts   = lmModels.map(m => ({
          value: m.laptop_model_id,
          label: `${m.brand_name} ${m.model_name}`,
        }));
        const lmMapObj = Object.fromEntries(lmModels.map(m => [m.laptop_model_id, m]));
        const lmMap    = Object.fromEntries(lmOpts.map(o => [o.value, o.label]));

        setModelMap(lmMapObj);
        setRefs({
          osOpts,
          equipOpts: (equip  ?? []).map(e => ({ value: e.equipment_user_id, label: e.name })),
          roomOpts,
          lmOpts,
          osMap,
          roomMap,
          lmMap,
        });
        setLaptops(lt ?? []);
      } catch (e) {
        showToast('err', e.message);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  function set(key, val) { setForm(f => ({ ...f, [key]: val })); }

  // Derive selected model details for hint display
  const selectedModel = form.laptop_model_id ? modelMap[form.laptop_model_id] : null;

  async function handleSubmit(e) {
    e.preventDefault();
    setSaving(true);
    try {
      const body = {
        hostname:          form.hostname,
        serial_number:     form.serial_number     || null,
        room_id:           form.room_id,
        laptop_model_id:   form.laptop_model_id,
        os_id:             form.os_id             || null,
        ram_gb:            form.ram_gb            ? parseInt(form.ram_gb,    10) : null,
        ram_type:          form.ram_type          || null,
        storage_gb:        form.storage_gb        ? parseInt(form.storage_gb, 10) : null,
        storage_type:      form.storage_type      || null,
        mac_address:       form.mac_address       || null,
        equipment_user_id: form.equipment_user_id || null,
        observations:      form.observations      || null,
      };
      const created = await api.createLaptop(body);
      setLaptops(prev => [created, ...prev]);
      setForm(EMPTY);
      showToast('ok', `Portàtil “${body.hostname}” afegit correctament.`);
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
        <h1 className="page-title">💻 Nou portàtil</h1>
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
                placeholder="portàtil-001"
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

            {/* Serial number */}
            <div className="form-group">
              <label>Número de sèrie</label>
              <input
                type="text"
                value={form.serial_number}
                onChange={e => set('serial_number', e.target.value)}
                placeholder="ex: SN123456789"
              />
            </div>

            {/* Laptop model (required) */}
            <div className="form-group">
              <label>Model *</label>
              <Combobox
                options={R.lmOpts}
                value={form.laptop_model_id}
                onChange={v => set('laptop_model_id', v)}
                placeholder="Selecciona model…"
              />
            </div>

          </div>

          {/* Model hint */}
          {selectedModel && (
            <div className="model-hint">
              <span className="model-hint-label">Especificacions base del model:</span>
              <span>{selectedModel.brand_name} {selectedModel.model_name}</span>
              {selectedModel.base_ram_gb    && <span>· RAM: {selectedModel.base_ram_gb} GB {selectedModel.base_ram_type ?? ''}</span>}
              {selectedModel.base_storage_gb && <span>· Emmagatzematge: {selectedModel.base_storage_gb} GB {selectedModel.base_storage_type ?? ''}</span>}
              {selectedModel.base_os_id     && <span>· SO base: {R.osMap[selectedModel.base_os_id] ?? '—'}</span>}
              <span className="model-hint-note">Deixa buits els camps que no vulguis sobreescriure.</span>
            </div>
          )}

          <div className="form-section-title" style={{ marginTop: 20, marginBottom: 12 }}>
            Sobreescriptura d&apos;especificacions (opcional)
          </div>

          <div className="form-grid">

            {/* OS override */}
            <div className="form-group">
              <label>Sistema operatiu</label>
              <Combobox
                options={R.osOpts}
                value={form.os_id}
                onChange={v => set('os_id', v)}
                placeholder="Sense override…"
                nullable
              />
            </div>

            {/* RAM override */}
            <div className="form-group">
              <label>RAM (GB)</label>
              <input
                type="number"
                min="0"
                value={form.ram_gb}
                onChange={e => set('ram_gb', e.target.value)}
                placeholder="ex: 16"
              />
            </div>

            <div className="form-group">
              <label>Tipus RAM</label>
              <select value={form.ram_type} onChange={e => set('ram_type', e.target.value)}>
                <option value="">—</option>
                {RAM_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>

            {/* Storage override */}
            <div className="form-group">
              <label>Emmagatzematge (GB)</label>
              <input
                type="number"
                min="0"
                value={form.storage_gb}
                onChange={e => set('storage_gb', e.target.value)}
                placeholder="ex: 512"
              />
            </div>

            <div className="form-group">
              <label>Tipus emmagatzematge</label>
              <select value={form.storage_type} onChange={e => set('storage_type', e.target.value)}>
                <option value="">—</option>
                {STORAGE_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>

            {/* MAC */}
            <div className="form-group">
              <label>Adreça MAC</label>
              <input
                type="text"
                value={form.mac_address}
                onChange={e => set('mac_address', e.target.value)}
                placeholder="AA:BB:CC:DD:EE:FF"
              />
            </div>

            {/* Equipment user */}
            <div className="form-group">
              <label>Usuari equip <span style={{ color: 'var(--muted)', fontWeight: 400 }}>(buit = alumnes)</span></label>
              <Combobox
                options={R.equipOpts}
                value={form.equipment_user_id}
                onChange={v => set('equipment_user_id', v)}
                placeholder="Alumnes (per defecte)…"
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
            <button
              type="submit"
              className="btn btn-primary"
              disabled={saving || !form.hostname || !form.room_id || !form.laptop_model_id}
            >
              {saving ? 'Guardant…' : 'Afegir portàtil'}
            </button>
          </div>
        </form>
      </div>

      {/* PREVIEW LIST */}
      <div className="page-header" style={{ marginTop: 40 }}>
        <h2 className="page-title" style={{ fontSize: 18 }}>💻 Portàtils registrats ({laptops.length})</h2>
      </div>
      <div className="card">
        <MiniList items={laptops} type="laptop" refs={R} />
      </div>
    </>
  );
}
