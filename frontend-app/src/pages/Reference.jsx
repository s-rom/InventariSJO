import { useState, useEffect, useCallback } from 'react';
import { api } from '../api';
import { useAuth } from '../App';
import Combobox from '../components/Combobox';

const RAM_TYPES     = ['DDR3', 'DDR4', 'DDR5', 'None'];
const STORAGE_TYPES = ['HDD', 'SSD', 'NVMe', 'None'];

const TABS = [
  { id: 'cpus',         label: '💾 CPUs' },
  { id: 'os',           label: '🖥 SO' },
  { id: 'brands',       label: '🏷 Marques' },
  { id: 'lmodels',      label: '💻 Models portàtil' },
  { id: 'dmodels',      label: '🖥️ Models sobretaula' },
  { id: 'pmodels',      label: '🖨️ Models impressora' },
  { id: 'psupplies',    label: '🖨️ Consumibles impressora' },
  { id: 'equipusers',   label: '👤 Usuaris equip' },
  { id: 'centers',      label: '🏢 Centres' },
  { id: 'rooms',        label: '🚪 Aules' },
  { id: 'cycles',       label: '📚 Cicles' },
  { id: 'classes',      label: '🎓 Cursos' },
  { id: 'prjmodels',    label: '📽️ Models projector' },
];

// ─────────────────────────────────────────────
// Generic confirm-delete hook
function useConfirmDelete() {
  const [delId, setDelId] = useState(null);
  return {
    delId,
    askDelete: setDelId,
    cancelDelete: () => setDelId(null),
    isAsking: (id) => delId === id,
  };
}

// ─────────────────────────────────────────────
function Section({ title, children }) {
  return (
    <div style={{ marginBottom: 32 }}>
      <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 14 }}>{title}</h2>
      {children}
    </div>
  );
}

// ─────────────────────────────────────────────
function CpusTab() {
  const [list, setList]           = useState([]);
  const [name, setName]           = useState('');
  const [score, setScore]         = useState('');
  const [saving, setSaving]       = useState(false);
  const [err, setErr]             = useState('');
  const [editing, setEditing]     = useState(null); // { cpu_id, model_name, benchmark_score }
  const [editErr, setEditErr]     = useState('');
  const [editSaving, setEditSaving] = useState(false);

  const load = useCallback(() => api.listCpus().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => { load(); }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try {
      await api.createCpu({ model_name: name, benchmark_score: score ? parseInt(score, 10) : null });
      setName(''); setScore(''); load();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function saveEdit() {
    setEditErr(''); setEditSaving(true);
    try {
      await api.updateCpu(editing.cpu_id, {
        model_name: editing.model_name,
        benchmark_score: editing.benchmark_score !== '' ? parseInt(editing.benchmark_score, 10) : null,
      });
      setEditing(null); load();
    } catch (ex) { setEditErr(ex.message); }
    finally { setEditSaving(false); }
  }

  return (
    <>
      <Section title="CPUs">
        <div className="card" style={{ marginBottom: 14 }}>
          <form onSubmit={create} className="form-panel">
            <div className="form-grid" style={{ gridTemplateColumns: '1fr auto auto' }}>
              <div className="form-group">
                <label>Nom del model *</label>
                <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="Intel Core i5-4590 3.30GHz" required />
              </div>
              <div className="form-group">
                <label>Benchmark (Passmark)</label>
                <input type="number" min="0" value={score} onChange={e => setScore(e.target.value)} placeholder="5354" />
              </div>
              <div className="form-group" style={{ justifyContent: 'flex-end' }}>
                <label style={{ visibility: 'hidden' }}>_</label>
                <button type="submit" className="btn btn-primary" disabled={saving || !name}>{saving ? '…' : 'Afegir'}</button>
              </div>
            </div>
            {err && <div className="error-msg">{err}</div>}
          </form>
        </div>

        <div className="card">
          <div className="table-wrap">
            <table>
              <thead><tr><th>Model</th><th>Benchmark</th><th></th></tr></thead>
              <tbody>
                {list.length === 0 && <tr><td colSpan={3} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
                {list.map(c => (
                  <tr key={c.cpu_id}>
                    <td>{c.model_name}</td>
                    <td>{c.benchmark_score ?? '—'}</td>
                    <td style={{ textAlign: 'right' }}>
                      <button className="btn btn-ghost btn-sm" onClick={() => { setEditErr(''); setEditing({ ...c, benchmark_score: c.benchmark_score ?? '' }); }}>Editar</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </Section>

      {editing && (
        <div
          style={{ position: 'fixed', inset: 0, zIndex: 1000, background: 'rgba(0,0,0,0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '32px 16px' }}
          onClick={e => { if (e.target === e.currentTarget) setEditing(null); }}
        >
          <div className="card" style={{ width: '100%', maxWidth: 480, padding: 24 }}>
            <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>Editar CPU</h2>
            <div className="form-group" style={{ marginBottom: 12 }}>
              <label>Nom del model *</label>
              <input type="text" value={editing.model_name} onChange={e => setEditing(v => ({ ...v, model_name: e.target.value }))} required autoFocus />
            </div>
            <div className="form-group" style={{ marginBottom: 16 }}>
              <label>Benchmark (Passmark)</label>
              <input type="number" min="0" value={editing.benchmark_score} onChange={e => setEditing(v => ({ ...v, benchmark_score: e.target.value }))} placeholder="5354" />
            </div>
            {editErr && <div className="error-msg" style={{ marginBottom: 12 }}>{editErr}</div>}
            <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
              <button className="btn btn-ghost" onClick={() => setEditing(null)}>Cancel·lar</button>
              <button className="btn btn-primary" onClick={saveEdit} disabled={editSaving || !editing.model_name}>{editSaving ? '…' : 'Guardar'}</button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

// ─────────────────────────────────────────────
function OsTab() {
  const [list, setList]           = useState([]);
  const [name, setName]           = useState('');
  const [saving, setSaving]       = useState(false);
  const [err, setErr]             = useState('');
  const [editing, setEditing]     = useState(null); // { os_id, name }
  const [editErr, setEditErr]     = useState('');
  const [editSaving, setEditSaving] = useState(false);

  const load = useCallback(() => api.listOS().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => { load(); }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try { await api.createOS({ name }); setName(''); load(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function saveEdit() {
    setEditErr(''); setEditSaving(true);
    try {
      await api.updateOS(editing.os_id, { name: editing.name });
      setEditing(null); load();
    } catch (ex) { setEditErr(ex.message); }
    finally { setEditSaving(false); }
  }

  return (
    <>
      <Section title="Sistemes Operatius">
        <div className="card" style={{ marginBottom: 14 }}>
          <form onSubmit={create} className="form-panel">
            <div className="form-grid" style={{ gridTemplateColumns: '1fr auto' }}>
              <div className="form-group">
                <label>Nom *</label>
                <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="Windows 11" required />
              </div>
              <div className="form-group" style={{ justifyContent: 'flex-end' }}>
                <label style={{ visibility: 'hidden' }}>_</label>
                <button type="submit" className="btn btn-primary" disabled={saving || !name}>{saving ? '…' : 'Afegir'}</button>
              </div>
            </div>
            {err && <div className="error-msg">{err}</div>}
          </form>
        </div>

        <div className="card">
          <div className="table-wrap">
            <table>
              <thead><tr><th>Nom</th><th></th></tr></thead>
              <tbody>
                {list.length === 0 && <tr><td colSpan={2} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
                {list.map(o => (
                  <tr key={o.os_id}>
                    <td>{o.name}</td>
                    <td style={{ textAlign: 'right' }}>
                      <button className="btn btn-ghost btn-sm" onClick={() => { setEditErr(''); setEditing({ ...o }); }}>Editar</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </Section>

      {editing && (
        <div
          style={{ position: 'fixed', inset: 0, zIndex: 1000, background: 'rgba(0,0,0,0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '32px 16px' }}
          onClick={e => { if (e.target === e.currentTarget) setEditing(null); }}
        >
          <div className="card" style={{ width: '100%', maxWidth: 400, padding: 24 }}>
            <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>Editar SO</h2>
            <div className="form-group" style={{ marginBottom: 16 }}>
              <label>Nom *</label>
              <input type="text" value={editing.name} onChange={e => setEditing(v => ({ ...v, name: e.target.value }))} required autoFocus />
            </div>
            {editErr && <div className="error-msg" style={{ marginBottom: 12 }}>{editErr}</div>}
            <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
              <button className="btn btn-ghost" onClick={() => setEditing(null)}>Cancel·lar</button>
              <button className="btn btn-primary" onClick={saveEdit} disabled={editSaving || !editing.name}>{editSaving ? '…' : 'Guardar'}</button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

// ─────────────────────────────────────────────
function BrandsTab() {
  const [list, setList]     = useState([]);
  const [name, setName]     = useState('');
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');

  const load = useCallback(() => api.listBrands().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => { load(); }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try { await api.createBrand({ name }); setName(''); load(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  return (
    <Section title="Marques">
      <div className="card" style={{ marginBottom: 14 }}>
        <form onSubmit={create} className="form-panel">
          <div className="form-grid" style={{ gridTemplateColumns: '1fr auto' }}>
            <div className="form-group">
              <label>Nom *</label>
              <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="HP, Dell, Lenovo…" required />
            </div>
            <div className="form-group" style={{ justifyContent: 'flex-end' }}>
              <label style={{ visibility: 'hidden' }}>_</label>
              <button type="submit" className="btn btn-primary" disabled={saving || !name}>{saving ? '…' : 'Afegir'}</button>
            </div>
          </div>
          {err && <div className="error-msg">{err}</div>}
        </form>
      </div>

      <div className="card">
        <div className="table-wrap">
          <table>
            <thead><tr><th>Nom</th><th></th></tr></thead>
            <tbody>
              {list.length === 0 && <tr><td colSpan={2} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
              {list.map(b => (
                <tr key={b.brand_id}>
                  <td>{b.name}</td>
                  <td></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Section>
  );
}

// ─────────────────────────────────────────────
function LaptopModelsTab() {
  const [list, setList]     = useState([]);
  const [refs, setRefs]     = useState(null);
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');

  const EMPTY_LM = { brand_id: null, model_name: '', cpu_id: null, base_ram_gb: '', base_ram_type: 'DDR4', base_storage_gb: '', base_storage_type: 'SSD', base_os_id: null };
  const [form, setForm] = useState(EMPTY_LM);
  function setF(k, v) { setForm(f => ({ ...f, [k]: v })); }

  const load = useCallback(() => api.listLaptopModels().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => {
    load();
    Promise.all([api.listBrands(), api.listCpus(), api.listOS()])
      .then(([brands, cpus, os]) => setRefs({
        brandOpts: (brands ?? []).map(b => ({ value: b.brand_id, label: b.name })),
        cpuOpts:   (cpus   ?? []).map(c => ({ value: c.cpu_id,   label: c.model_name })),
        osOpts:    (os     ?? []).map(o => ({ value: o.os_id,    label: o.name })),
      }))
      .catch(() => {});
  }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try {
      await api.createLaptopModel({
        brand_id:          form.brand_id,
        model_name:        form.model_name,
        cpu_id:            form.cpu_id            || null,
        base_ram_gb:       parseInt(form.base_ram_gb, 10),
        base_ram_type:     form.base_ram_type,
        base_storage_gb:   parseInt(form.base_storage_gb, 10),
        base_storage_type: form.base_storage_type,
        base_os_id:        form.base_os_id        || null,
      });
      setForm(EMPTY_LM); load();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  if (!refs) return <div className="empty">Carregant…</div>;
  const R = refs;

  return (
    <Section title="Models de portàtil">
      <div className="card" style={{ marginBottom: 14 }}>
        <form onSubmit={create} className="form-panel">
          <div className="form-grid">
            <div className="form-group">
              <label>Marca *</label>
              <Combobox options={R.brandOpts} value={form.brand_id} onChange={v => setF('brand_id', v)} placeholder="Marca…" />
            </div>
            <div className="form-group">
              <label>Nom del model *</label>
              <input type="text" value={form.model_name} onChange={e => setF('model_name', e.target.value)} placeholder="EliteBook 840 G3" required />
            </div>
            <div className="form-group">
              <label>CPU</label>
              <Combobox options={R.cpuOpts} value={form.cpu_id} onChange={v => setF('cpu_id', v)} placeholder="Opcional…" nullable />
            </div>
            <div className="form-group">
              <label>RAM base (GB) *</label>
              <input type="number" min="0" value={form.base_ram_gb} onChange={e => setF('base_ram_gb', e.target.value)} required />
            </div>
            <div className="form-group">
              <label>Tipus RAM *</label>
              <select value={form.base_ram_type} onChange={e => setF('base_ram_type', e.target.value)}>
                {RAM_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>
            <div className="form-group">
              <label>Emmagatzematge base (GB) *</label>
              <input type="number" min="0" value={form.base_storage_gb} onChange={e => setF('base_storage_gb', e.target.value)} required />
            </div>
            <div className="form-group">
              <label>Tipus emmagatzematge *</label>
              <select value={form.base_storage_type} onChange={e => setF('base_storage_type', e.target.value)}>
                {STORAGE_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>
            <div className="form-group">
              <label>SO base</label>
              <Combobox options={R.osOpts} value={form.base_os_id} onChange={v => setF('base_os_id', v)} placeholder="Opcional…" nullable />
            </div>
          </div>
          {err && <div className="error-msg">{err}</div>}
          <div className="form-actions">
            <button type="submit" className="btn btn-primary"
              disabled={saving || !form.brand_id || !form.model_name || !form.base_ram_gb || !form.base_storage_gb}>
              {saving ? '…' : 'Afegir model'}
            </button>
          </div>
        </form>
      </div>

      <div className="card">
        <div className="table-wrap">
          <table>
            <thead><tr><th>Marca</th><th>Model</th><th>RAM</th><th>Emmagatzematge</th><th></th></tr></thead>
            <tbody>
              {list.length === 0 && <tr><td colSpan={5} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
              {list.map(m => (
                <tr key={m.laptop_model_id}>
                  <td>{m.brand_name}</td>
                  <td>{m.model_name}</td>
                  <td>{m.base_ram_gb} GB {m.base_ram_type}</td>
                  <td>{m.base_storage_gb} GB {m.base_storage_type}</td>
                  <td></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Section>
  );
}

// ─────────────────────────────────────────────
function DesktopModelsTab() {
  const [list, setList]     = useState([]);
  const [refs, setRefs]     = useState(null);
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');
  const cd                  = useConfirmDelete();

  const EMPTY_DM = { brand_id: null, model_name: '', cpu_id: null, base_ram_gb: '', base_ram_type: 'DDR4', base_storage_gb: '', base_storage_type: 'SSD', base_os_id: null };
  const [form, setForm] = useState(EMPTY_DM);
  function setF(k, v) { setForm(f => ({ ...f, [k]: v })); }

  const load = useCallback(() => api.listDesktopModels().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => {
    load();
    Promise.all([api.listBrands(), api.listCpus(), api.listOS()])
      .then(([brands, cpus, os]) => setRefs({
        brandOpts: (brands ?? []).map(b => ({ value: b.brand_id, label: b.name })),
        cpuOpts:   (cpus   ?? []).map(c => ({ value: c.cpu_id,   label: c.model_name })),
        osOpts:    (os     ?? []).map(o => ({ value: o.os_id,    label: o.name })),
      }))
      .catch(() => {});
  }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try {
      await api.createDesktopModel({
        brand_id:          form.brand_id,
        model_name:        form.model_name,
        cpu_id:            form.cpu_id            || null,
        base_ram_gb:       parseInt(form.base_ram_gb, 10),
        base_ram_type:     form.base_ram_type,
        base_storage_gb:   parseInt(form.base_storage_gb, 10),
        base_storage_type: form.base_storage_type,
        base_os_id:        form.base_os_id        || null,
      });
      setForm(EMPTY_DM); load();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function del(id) {
    try { await api.deleteDesktopModel(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  if (!refs) return <div className="empty">Carregant…</div>;
  const R = refs;

  return (
    <Section title="Models de sobretaula">
      <div className="card" style={{ marginBottom: 14 }}>
        <form onSubmit={create} className="form-panel">
          <div className="form-grid">
            <div className="form-group">
              <label>Marca *</label>
              <Combobox options={R.brandOpts} value={form.brand_id} onChange={v => setF('brand_id', v)} placeholder="Marca…" />
            </div>
            <div className="form-group">
              <label>Nom del model *</label>
              <input type="text" value={form.model_name} onChange={e => setF('model_name', e.target.value)} placeholder="ProDesk 400 G3" required />
            </div>
            <div className="form-group">
              <label>CPU</label>
              <Combobox options={R.cpuOpts} value={form.cpu_id} onChange={v => setF('cpu_id', v)} placeholder="Opcional…" nullable />
            </div>
            <div className="form-group">
              <label>RAM base (GB) *</label>
              <input type="number" min="0" value={form.base_ram_gb} onChange={e => setF('base_ram_gb', e.target.value)} required />
            </div>
            <div className="form-group">
              <label>Tipus RAM *</label>
              <select value={form.base_ram_type} onChange={e => setF('base_ram_type', e.target.value)}>
                {RAM_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>
            <div className="form-group">
              <label>Emmagatzematge base (GB) *</label>
              <input type="number" min="0" value={form.base_storage_gb} onChange={e => setF('base_storage_gb', e.target.value)} required />
            </div>
            <div className="form-group">
              <label>Tipus emmagatzematge *</label>
              <select value={form.base_storage_type} onChange={e => setF('base_storage_type', e.target.value)}>
                {STORAGE_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>
            <div className="form-group">
              <label>SO base</label>
              <Combobox options={R.osOpts} value={form.base_os_id} onChange={v => setF('base_os_id', v)} placeholder="Opcional…" nullable />
            </div>
          </div>
          {err && <div className="error-msg">{err}</div>}
          <div className="form-actions">
            <button type="submit" className="btn btn-primary"
              disabled={saving || !form.brand_id || !form.model_name || !form.base_ram_gb || !form.base_storage_gb}>
              {saving ? '…' : 'Afegir model'}
            </button>
          </div>
        </form>
      </div>

      <div className="card">
        <div className="table-wrap">
          <table>
            <thead><tr><th>Marca</th><th>Model</th><th>RAM</th><th>Emmagatzematge</th><th></th></tr></thead>
            <tbody>
              {list.length === 0 && <tr><td colSpan={5} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
              {list.map(m => (
                <tr key={m.desktop_model_id}>
                  <td>{m.brand_name}</td>
                  <td>{m.model_name}</td>
                  <td>{m.base_ram_gb} GB {m.base_ram_type}</td>
                  <td>{m.base_storage_gb} GB {m.base_storage_type}</td>
                  <td style={{ textAlign: 'right' }}>
                    {cd.isAsking(m.desktop_model_id)
                      ? <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                          <button className="btn btn-danger btn-sm" onClick={() => del(m.desktop_model_id)}>Sí</button>
                          <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                      : <button className="btn btn-danger btn-sm" onClick={() => cd.askDelete(m.desktop_model_id)}>Eliminar</button>
                    }
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Section>
  );
}

// ─────────────────────────────────────────────
function EquipUsersTab() {
  const [list, setList]       = useState([]);
  const [name, setName]       = useState('');
  const [editId, setEditId]   = useState(null);
  const [editName, setEditName] = useState('');
  const [saving, setSaving]   = useState(false);
  const [err, setErr]         = useState('');
  const cd                    = useConfirmDelete();

  const load = useCallback(() => api.listEquipmentUsers().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => { load(); }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try { await api.createEquipmentUser({ name }); setName(''); load(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function save(id) {
    try { await api.updateEquipmentUser(id, { name: editName }); setEditId(null); load(); }
    catch (ex) { setErr(ex.message); }
  }

  async function del(id) {
    try { await api.deleteEquipmentUser(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  return (
    <Section title="Usuaris d'equip">
      <div className="card" style={{ marginBottom: 14 }}>
        <form onSubmit={create} className="form-panel">
          <div className="form-grid" style={{ gridTemplateColumns: '1fr auto' }}>
            <div className="form-group">
              <label>Nom *</label>
              <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="ex: Sala servidors" required />
            </div>
            <div className="form-group" style={{ justifyContent: 'flex-end' }}>
              <label style={{ visibility: 'hidden' }}>_</label>
              <button type="submit" className="btn btn-primary" disabled={saving || !name}>{saving ? '…' : 'Afegir'}</button>
            </div>
          </div>
          {err && <div className="error-msg">{err}</div>}
        </form>
      </div>
      <div className="card">
        <div className="table-wrap">
          <table>
            <thead><tr><th>Nom</th><th></th></tr></thead>
            <tbody>
              {list.length === 0 && <tr><td colSpan={2} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
              {list.map(u => (
                <tr key={u.equipment_user_id}>
                  <td>
                    {editId === u.equipment_user_id
                      ? <input type="text" value={editName} onChange={e => setEditName(e.target.value)} style={{ width: 220 }} autoFocus />
                      : u.name}
                  </td>
                  <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                    {editId === u.equipment_user_id ? (
                      <>
                        <button className="btn btn-primary btn-sm" onClick={() => save(u.equipment_user_id)}>Guardar</button>
                        <button className="btn btn-ghost btn-sm" style={{ marginLeft: 4 }} onClick={() => setEditId(null)}>Cancel·lar</button>
                      </>
                    ) : cd.isAsking(u.equipment_user_id) ? (
                      <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                        <button className="btn btn-danger btn-sm" onClick={() => del(u.equipment_user_id)}>Sí</button>
                        <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                    ) : (
                      <>
                        <button className="btn btn-ghost btn-sm" onClick={() => { setEditId(u.equipment_user_id); setEditName(u.name); }}>Editar</button>
                        <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => cd.askDelete(u.equipment_user_id)}>Eliminar</button>
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Section>
  );
}

// ─────────────────────────────────────────────
function CentersTab() {
  const [list, setList]         = useState([]);
  const [name, setName]         = useState('');
  const [editId, setEditId]     = useState(null);
  const [editName, setEditName] = useState('');
  const [saving, setSaving]     = useState(false);
  const [err, setErr]           = useState('');
  const cd                      = useConfirmDelete();

  const load = useCallback(() => api.listCenters().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => { load(); }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try { await api.createCenter({ name }); setName(''); load(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function save(id) {
    try { await api.updateCenter(id, { name: editName }); setEditId(null); load(); }
    catch (ex) { setErr(ex.message); }
  }

  async function del(id) {
    try { await api.deleteCenter(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  return (
    <Section title="Centres">
      <div className="card" style={{ marginBottom: 14 }}>
        <form onSubmit={create} className="form-panel">
          <div className="form-grid" style={{ gridTemplateColumns: '1fr auto' }}>
            <div className="form-group">
              <label>Nom del centre *</label>
              <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="IES Exemple" required />
            </div>
            <div className="form-group" style={{ justifyContent: 'flex-end' }}>
              <label style={{ visibility: 'hidden' }}>_</label>
              <button type="submit" className="btn btn-primary" disabled={saving || !name}>{saving ? '…' : 'Afegir'}</button>
            </div>
          </div>
          {err && <div className="error-msg">{err}</div>}
        </form>
      </div>
      <div className="card">
        <div className="table-wrap">
          <table>
            <thead><tr><th>Nom</th><th></th></tr></thead>
            <tbody>
              {list.length === 0 && <tr><td colSpan={2} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
              {list.map(c => (
                <tr key={c.center_id}>
                  <td>
                    {editId === c.center_id
                      ? <input type="text" value={editName} onChange={e => setEditName(e.target.value)} style={{ width: 260 }} autoFocus />
                      : c.name}
                  </td>
                  <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                    {editId === c.center_id ? (
                      <>
                        <button className="btn btn-primary btn-sm" onClick={() => save(c.center_id)}>Guardar</button>
                        <button className="btn btn-ghost btn-sm" style={{ marginLeft: 4 }} onClick={() => setEditId(null)}>Cancel·lar</button>
                      </>
                    ) : cd.isAsking(c.center_id) ? (
                      <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                        <button className="btn btn-danger btn-sm" onClick={() => del(c.center_id)}>Sí</button>
                        <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                    ) : (
                      <>
                        <button className="btn btn-ghost btn-sm" onClick={() => { setEditId(c.center_id); setEditName(c.name); }}>Editar</button>
                        <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => cd.askDelete(c.center_id)}>Eliminar</button>
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Section>
  );
}

// ─────────────────────────────────────────────
function RoomsTab() {
  const [centers, setCenters]         = useState([]);
  const [centerId, setCenterId]       = useState(null);
  const [rooms, setRooms]             = useState([]);
  const [name, setName]               = useState('');
  const [editId, setEditId]           = useState(null);
  const [editName, setEditName]       = useState('');
  const [saving, setSaving]           = useState(false);
  const [err, setErr]                 = useState('');
  const cd                            = useConfirmDelete();

  useEffect(() => {
    api.listCenters()
      .then(d => { const list = d ?? []; setCenters(list); if (list.length > 0) setCenterId(list[0].center_id); })
      .catch(() => {});
  }, []);

  const loadRooms = useCallback(() => {
    if (!centerId) return;
    api.listRoomsByCenter(centerId).then(d => setRooms(d ?? [])).catch(() => {});
  }, [centerId]);

  useEffect(() => { loadRooms(); }, [loadRooms]);

  async function create(e) {
    e.preventDefault(); if (!centerId) return; setErr(''); setSaving(true);
    try { await api.createRoom(centerId, { name }); setName(''); loadRooms(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function save(id) {
    try { await api.updateRoom(id, { name: editName }); setEditId(null); loadRooms(); }
    catch (ex) { setErr(ex.message); }
  }

  async function del(id) {
    try { await api.deleteRoom(id); loadRooms(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  return (
    <Section title="Aules">
      {/* Center selector */}
      <div style={{ marginBottom: 16, display: 'flex', alignItems: 'center', gap: 10 }}>
        <label style={{ fontSize: 12, fontWeight: 500, color: 'var(--muted)', whiteSpace: 'nowrap' }}>Centre:</label>
        <select
          value={centerId ?? ''}
          onChange={e => setCenterId(Number(e.target.value))}
          style={{ width: 240 }}
        >
          {centers.length === 0 && <option value="">— Primer crea un centre —</option>}
          {centers.map(c => <option key={c.center_id} value={c.center_id}>{c.name}</option>)}
        </select>
      </div>

      {centerId && (
        <>
          <div className="card" style={{ marginBottom: 14 }}>
            <form onSubmit={create} className="form-panel">
              <div className="form-grid" style={{ gridTemplateColumns: '1fr auto' }}>
                <div className="form-group">
                  <label>Nom de l&apos;aula *</label>
                  <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="Aula 1, Lab informatica…" required />
                </div>
                <div className="form-group" style={{ justifyContent: 'flex-end' }}>
                  <label style={{ visibility: 'hidden' }}>_</label>
                  <button type="submit" className="btn btn-primary" disabled={saving || !name}>{saving ? '…' : 'Afegir'}</button>
                </div>
              </div>
              {err && <div className="error-msg">{err}</div>}
            </form>
          </div>
          <div className="card">
            <div className="table-wrap">
              <table>
                <thead><tr><th>Aula</th><th></th></tr></thead>
                <tbody>
                  {rooms.length === 0 && <tr><td colSpan={2} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense aules per aquest centre.</td></tr>}
                  {rooms.map(r => (
                    <tr key={r.room_id}>
                      <td>
                        {editId === r.room_id
                          ? <input type="text" value={editName} onChange={e => setEditName(e.target.value)} style={{ width: 220 }} autoFocus />
                          : r.name}
                      </td>
                      <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                        {editId === r.room_id ? (
                          <>
                            <button className="btn btn-primary btn-sm" onClick={() => save(r.room_id)}>Guardar</button>
                            <button className="btn btn-ghost btn-sm" style={{ marginLeft: 4 }} onClick={() => setEditId(null)}>Cancel·lar</button>
                          </>
                        ) : cd.isAsking(r.room_id) ? (
                          <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                            <button className="btn btn-danger btn-sm" onClick={() => del(r.room_id)}>Sí</button>
                            <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                        ) : (
                          <>
                            <button className="btn btn-ghost btn-sm" onClick={() => { setEditId(r.room_id); setEditName(r.name); }}>Editar</button>
                            <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => cd.askDelete(r.room_id)}>Eliminar</button>
                          </>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </>
      )}
    </Section>
  );
}

// ─────────────────────────────────────────────
function CyclesTab() {
  const [list, setList]   = useState([]);
  const [name, setName]   = useState('');
  const [saving, setSaving] = useState(false);
  const [err, setErr]     = useState('');

  const load = useCallback(() => api.listCycles().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => { load(); }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try { await api.createCycle({ name }); setName(''); load(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  return (
    <Section title="Cicles formatius">
      <div className="card" style={{ marginBottom: 14 }}>
        <form onSubmit={create} className="form-panel">
          <div className="form-grid" style={{ gridTemplateColumns: '1fr auto' }}>
            <div className="form-group">
              <label>Nom del cicle *</label>
              <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="SMX, ASIX, DAW…" required />
            </div>
            <div className="form-group" style={{ justifyContent: 'flex-end' }}>
              <label style={{ visibility: 'hidden' }}>_</label>
              <button type="submit" className="btn btn-primary" disabled={saving || !name}>{saving ? '…' : 'Afegir'}</button>
            </div>
          </div>
          {err && <div className="error-msg">{err}</div>}
        </form>
      </div>
      <div className="card">
        <div className="table-wrap">
          <table>
            <thead><tr><th>Nom</th></tr></thead>
            <tbody>
              {list.length === 0 && <tr><td style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
              {list.map(c => (
                <tr key={c.cycle_id}>
                  <td>{c.name}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Section>
  );
}

// ─────────────────────────────────────────────
const SHIFTS = [
  { value: 'morning',   label: 'Matí' },
  { value: 'afternoon', label: 'Tarda' },
];

function ClassesTab() {
  const { role } = useAuth();
  const isTutor  = role === 'tutor';

  const [cycles, setCycles]     = useState([]);
  const [cycleId, setCycleId]   = useState(null);
  const [classes, setClasses]   = useState([]);
  const [tutors, setTutors]     = useState([]);
  const [saving, setSaving]     = useState(false);
  const [err, setErr]           = useState('');
  const cd                      = useConfirmDelete();

  const EMPTY_CL = { course: '', class_label: '', shift: 'morning', tutor_app_user_id: '' };
  const [form, setForm]         = useState(EMPTY_CL);
  const [editId, setEditId]     = useState(null);
  const [editForm, setEditForm] = useState({});

  function setF(k, v)  { setForm(f => ({ ...f, [k]: v })); }
  function setEF(k, v) { setEditForm(f => ({ ...f, [k]: v })); }

  // Tutors: load only their assigned classes (flat, no cycle selector needed)
  // Admins/editors: load cycles + all classes per cycle
  useEffect(() => {
    if (isTutor) {
      api.listMyClasses().then(d => setClasses(d ?? [])).catch(() => {});
    } else {
      api.listCycles()
        .then(d => { const list = d ?? []; setCycles(list); if (list.length > 0) setCycleId(list[0].cycle_id); })
        .catch(() => {});
      api.listUsers()
        .then(d => setTutors((d ?? []).filter(u => u.role_id === 'tutor')))
        .catch(() => {});
    }
  }, [isTutor]);

  const loadClasses = useCallback(() => {
    if (isTutor) {
      api.listMyClasses().then(d => setClasses(d ?? [])).catch(() => {});
    } else {
      if (!cycleId) return;
      api.listClassesByCycle(cycleId).then(d => setClasses(d ?? [])).catch(() => {});
    }
  }, [isTutor, cycleId]);

  useEffect(() => { if (!isTutor) loadClasses(); }, [loadClasses, isTutor]);

  async function create(e) {
    e.preventDefault(); if (!cycleId) return; setErr(''); setSaving(true);
    try {
      const data = {
        course:      parseInt(form.course, 10),
        class_label: form.class_label,
        shift:       form.shift,
      };
      if (form.tutor_app_user_id !== '') data.tutor_app_user_id = parseInt(form.tutor_app_user_id, 10);
      await api.createClass(cycleId, data);
      setForm(EMPTY_CL); loadClasses();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function save(id) {
    try {
      const patch = { class_label: editForm.class_label, shift: editForm.shift };
      if (!isTutor && editForm.tutor_app_user_id !== '') patch.tutor_app_user_id = parseInt(editForm.tutor_app_user_id, 10);
      await api.updateClass(id, patch);
      setEditId(null); loadClasses();
    } catch (ex) { setErr(ex.message); }
  }

  async function del(id) {
    try { await api.deleteClass(id); loadClasses(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  // ── Tutor view: flat list of assigned classes, only edit shift ──────────────
  if (isTutor) {
    return (
      <Section title="Els meus cursos">
        {classes.length === 0 && (
          <div className="card" style={{ padding: 20, color: 'var(--muted)', textAlign: 'center' }}>
            Encara no tens cap curs assignat. Contacta amb un administrador.
          </div>
        )}
        {classes.length > 0 && (
          <div className="card">
            <div className="table-wrap">
              <table>
                <thead><tr><th>Cicle</th><th>Curs</th><th>Etiqueta</th><th>Torn</th><th></th></tr></thead>
                <tbody>
                  {classes.map(cl => (
                    <tr key={cl.class_id}>
                      <td>{cl.cycle_name}</td>
                      <td>{cl.course}r</td>
                      <td>{cl.class_label}</td>
                      <td>
                        {editId === cl.class_id
                          ? <select value={editForm.shift} onChange={e => setEF('shift', e.target.value)} style={{ width: 90 }}>
                              {SHIFTS.map(s => <option key={s.value} value={s.value}>{s.label}</option>)}
                            </select>
                          : SHIFTS.find(s => s.value === cl.shift)?.label ?? cl.shift}
                      </td>
                      <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                        {editId === cl.class_id ? (
                          <>
                            <button className="btn btn-primary btn-sm" onClick={() => save(cl.class_id)}>Guardar</button>
                            <button className="btn btn-ghost btn-sm" style={{ marginLeft: 4 }} onClick={() => setEditId(null)}>Cancel·lar</button>
                          </>
                        ) : (
                          <button className="btn btn-ghost btn-sm" onClick={() => { setEditId(cl.class_id); setEditForm({ shift: cl.shift }); }}>Editar</button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        )}
        {err && <div className="error-msg" style={{ marginTop: 8 }}>{err}</div>}
      </Section>
    );
  }

  // ── Admin / Editor view ───────────────────────────────────────────────────
  return (
    <Section title="Cursos / Classes">
      {/* Cycle selector */}
      <div style={{ marginBottom: 16, display: 'flex', alignItems: 'center', gap: 10 }}>
        <label style={{ fontSize: 12, fontWeight: 500, color: 'var(--muted)', whiteSpace: 'nowrap' }}>Cicle:</label>
        <select
          value={cycleId ?? ''}
          onChange={e => setCycleId(Number(e.target.value))}
          style={{ width: 240 }}
        >
          {cycles.length === 0 && <option value="">— Primer crea un cicle —</option>}
          {cycles.map(c => <option key={c.cycle_id} value={c.cycle_id}>{c.name}</option>)}
        </select>
      </div>

      {cycleId && (
        <>
          <div className="card" style={{ marginBottom: 14 }}>
            <form onSubmit={create} className="form-panel">
              <div className="form-grid" style={{ gridTemplateColumns: 'repeat(4, 1fr) auto' }}>
                <div className="form-group">
                  <label>Curs *</label>
                  <input
                    type="number" min="1" max="4"
                    value={form.course}
                    onChange={e => setF('course', e.target.value)}
                    placeholder="1, 2…"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Etiqueta *</label>
                  <input
                    type="text"
                    value={form.class_label}
                    onChange={e => setF('class_label', e.target.value)}
                    placeholder="A, B, Matí…"
                    required
                  />
                </div>
                <div className="form-group">
                  <label>Torn *</label>
                  <select value={form.shift} onChange={e => setF('shift', e.target.value)}>
                    {SHIFTS.map(s => <option key={s.value} value={s.value}>{s.label}</option>)}
                  </select>
                </div>
                <div className="form-group">
                  <label>Tutor</label>
                  <select value={form.tutor_app_user_id} onChange={e => setF('tutor_app_user_id', e.target.value)}>
                    <option value="">— Cap —</option>
                    {tutors.map(u => <option key={u.app_user_id} value={u.app_user_id}>{u.username}</option>)}
                  </select>
                </div>
                <div className="form-group" style={{ justifyContent: 'flex-end' }}>
                  <label style={{ visibility: 'hidden' }}>_</label>
                  <button type="submit" className="btn btn-primary" disabled={saving || !form.course || !form.class_label}>
                    {saving ? '…' : 'Afegir'}
                  </button>
                </div>
              </div>
              {err && <div className="error-msg">{err}</div>}
            </form>
          </div>
          <div className="card">
            <div className="table-wrap">
              <table>
                <thead><tr><th>Curs</th><th>Etiqueta</th><th>Torn</th><th>Tutor</th><th></th></tr></thead>
                <tbody>
                  {classes.length === 0 && <tr><td colSpan={5} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense cursos per aquest cicle.</td></tr>}
                  {classes.map(cl => (
                    <tr key={cl.class_id}>
                      <td>{cl.course}r</td>
                      <td>
                        {editId === cl.class_id
                          ? <input type="text" value={editForm.class_label} onChange={e => setEF('class_label', e.target.value)} style={{ width: 100 }} autoFocus />
                          : cl.class_label}
                      </td>
                      <td>
                        {editId === cl.class_id
                          ? <select value={editForm.shift} onChange={e => setEF('shift', e.target.value)} style={{ width: 90 }}>
                              {SHIFTS.map(s => <option key={s.value} value={s.value}>{s.label}</option>)}
                            </select>
                          : SHIFTS.find(s => s.value === cl.shift)?.label ?? cl.shift}
                      </td>
                      <td>
                        {editId === cl.class_id
                          ? <select value={editForm.tutor_app_user_id ?? ''} onChange={e => setEF('tutor_app_user_id', e.target.value)} style={{ width: 130 }}>
                              <option value="">— Cap —</option>
                              {tutors.map(u => <option key={u.app_user_id} value={u.app_user_id}>{u.username}</option>)}
                            </select>
                          : tutors.find(u => u.app_user_id === cl.tutor_app_user_id)?.username ?? '—'}
                      </td>
                      <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                        {editId === cl.class_id ? (
                          <>
                            <button className="btn btn-primary btn-sm" onClick={() => save(cl.class_id)}>Guardar</button>
                            <button className="btn btn-ghost btn-sm" style={{ marginLeft: 4 }} onClick={() => setEditId(null)}>Cancel·lar</button>
                          </>
                        ) : cd.isAsking(cl.class_id) ? (
                          <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                            <button className="btn btn-danger btn-sm" onClick={() => del(cl.class_id)}>Sí</button>
                            <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                        ) : (
                          <>
                            <button className="btn btn-ghost btn-sm" onClick={() => { setEditId(cl.class_id); setEditForm({ class_label: cl.class_label, shift: cl.shift, tutor_app_user_id: cl.tutor_app_user_id ?? '' }); }}>Editar</button>
                            <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => cd.askDelete(cl.class_id)}>Eliminar</button>
                          </>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </>
      )}
    </Section>
  );
}

// ─────────────────────────────────────────────
// Main Reference page
// ─────────────────────────────────────────────
function PrinterModelsTab() {
  const PRINTER_TYPES = ['toner', 'ink', 'managed'];

  const [list, setList]     = useState([]);
  const [refs, setRefs]     = useState(null);
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');
  const [editing, setEditing]   = useState(null);
  const [editErr, setEditErr]   = useState('');
  const [editSaving, setEditSaving] = useState(false);
  const [modelSupplies, setModelSupplies] = useState([]);
  const [addSupplyId, setAddSupplyId]     = useState(null);
  const cd = useConfirmDelete();

  const EMPTY = { brand_id: null, model_name: '', printer_type: 'toner', print_color: 'Color' };
  const [form, setForm] = useState(EMPTY);
  function setF(k, v) { setForm(f => ({ ...f, [k]: v })); }

  const load = useCallback(() => api.listPrinterModels().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => {
    load();
    Promise.all([api.listBrands(), api.listPrinterSupplies()])
      .then(([brands, supplies]) => setRefs({
        brandOpts:   (brands   ?? []).map(b => ({ value: b.brand_id,        label: b.name })),
        supplyOpts:  (supplies ?? []).map(s => ({ value: s.printer_supply_id, label: `${s.name} (${s.supply_type})` })),
      }))
      .catch(() => {});
  }, [load]);

  const loadModelSupplies = useCallback((modelId) => {
    api.listSuppliesByModel(modelId).then(d => setModelSupplies(d ?? [])).catch(() => {});
  }, []);

  useEffect(() => {
    if (editing) { setAddSupplyId(null); loadModelSupplies(editing.printer_model_id); }
    else { setModelSupplies([]); }
  }, [editing, loadModelSupplies]);

  async function addSupply() {
    if (!addSupplyId || !editing) return;
    try {
      await api.addSupplyToModel(editing.printer_model_id, { printer_supply_id: addSupplyId });
      setAddSupplyId(null);
      loadModelSupplies(editing.printer_model_id);
    } catch (ex) { setEditErr(ex.message); }
  }

  async function removeSupply(supplyId) {
    try {
      await api.removeSupplyFromModel(editing.printer_model_id, supplyId);
      loadModelSupplies(editing.printer_model_id);
    } catch (ex) { setEditErr(ex.message); }
  }

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try {
      await api.createPrinterModel({ brand_id: form.brand_id, model_name: form.model_name, printer_type: form.printer_type, print_color: form.print_color });
      setForm(EMPTY); load();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function saveEdit() {
    setEditErr(''); setEditSaving(true);
    try {
      await api.updatePrinterModel(editing.printer_model_id, {
        brand_id: editing.brand_id, model_name: editing.model_name,
        printer_type: editing.printer_type, print_color: editing.print_color,
      });
      setEditing(null); load();
    } catch (ex) { setEditErr(ex.message); }
    finally { setEditSaving(false); }
  }

  async function del(id) {
    try { await api.deletePrinterModel(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  if (!refs) return <div className="empty">Carregant…</div>;

  return (
    <>
      <Section title="Models d'impressora">
        <div className="card" style={{ marginBottom: 14 }}>
          <form onSubmit={create} className="form-panel">
            <div className="form-grid">
              <div className="form-group">
                <label>Marca *</label>
                <Combobox options={refs.brandOpts} value={form.brand_id} onChange={v => setF('brand_id', v)} placeholder="Marca…" />
              </div>
              <div className="form-group">
                <label>Nom del model *</label>
                <input type="text" value={form.model_name} onChange={e => setF('model_name', e.target.value)} placeholder="LaserJet Pro M404n" required />
              </div>
              <div className="form-group">
                <label>Tipus *</label>
                <select value={form.printer_type} onChange={e => setF('printer_type', e.target.value)}>
                  {PRINTER_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>Color</label>
                <select value={form.print_color} onChange={e => setF('print_color', e.target.value)}>
                  <option value="Color">Color</option>
                  <option value="BN">Blanc i negre</option>
                </select>
              </div>
            </div>
            {err && <div className="error-msg">{err}</div>}
            <div className="form-actions">
              <button type="submit" className="btn btn-primary" disabled={saving || !form.brand_id || !form.model_name}>
                {saving ? '…' : 'Afegir model'}
              </button>
            </div>
          </form>
        </div>

        <div className="card">
          <div className="table-wrap">
            <table>
              <thead><tr><th>Marca</th><th>Model</th><th>Tipus</th><th>Color</th><th></th></tr></thead>
              <tbody>
                {list.length === 0 && <tr><td colSpan={5} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
                {list.map(m => (
                  <tr key={m.printer_model_id}>
                    <td>{m.brand_name}</td>
                    <td>{m.model_name}</td>
                    <td>{m.printer_type}</td>
                    <td>{m.print_color}</td>
                    <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                      {cd.isAsking(m.printer_model_id)
                        ? <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                            <button className="btn btn-danger btn-sm" onClick={() => del(m.printer_model_id)}>Sí</button>
                            <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                        : <>
                            <button className="btn btn-ghost btn-sm" onClick={() => { setEditErr(''); setEditing({ ...m }); }}>Editar</button>
                            <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => cd.askDelete(m.printer_model_id)}>Eliminar</button>
                          </>
                      }
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </Section>

      {editing && (
        <div style={{ position: 'fixed', inset: 0, zIndex: 1000, background: 'rgba(0,0,0,0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '32px 16px' }}
          onClick={e => { if (e.target === e.currentTarget) setEditing(null); }}>
          <div className="card" style={{ width: '100%', maxWidth: 620, padding: 24 }}>
            <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>Editar model d'impressora</h2>
            <div className="form-grid">
              <div className="form-group">
                <label>Marca *</label>
                <Combobox options={refs.brandOpts} value={editing.brand_id} onChange={v => setEditing(e => ({ ...e, brand_id: v }))} placeholder="Marca…" />
              </div>
              <div className="form-group">
                <label>Nom del model *</label>
                <input type="text" value={editing.model_name} onChange={e => setEditing(v => ({ ...v, model_name: e.target.value }))} required autoFocus />
              </div>
              <div className="form-group">
                <label>Tipus *</label>
                <select value={editing.printer_type} onChange={e => setEditing(v => ({ ...v, printer_type: e.target.value }))}>
                  {PRINTER_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                </select>
              </div>
              <div className="form-group">
                <label>Color</label>
                <select value={editing.print_color} onChange={e => setEditing(v => ({ ...v, print_color: e.target.value }))}>
                  <option value="Color">Color</option>
                  <option value="BN">Blanc i negre</option>
                </select>
              </div>
            </div>

            <hr style={{ margin: '16px 0', borderColor: 'var(--border)' }} />
            <div style={{ marginBottom: 12 }}>
              <div style={{ fontSize: 13, fontWeight: 600, marginBottom: 8 }}>Consumibles associats</div>
              {modelSupplies.length === 0
                ? <div style={{ fontSize: 12, color: 'var(--muted)', marginBottom: 8 }}>Cap consumible associat</div>
                : <div style={{ display: 'flex', flexWrap: 'wrap', gap: 6, marginBottom: 8 }}>
                    {modelSupplies.map(s => (
                      <span key={s.printer_supply_id} style={{ display: 'inline-flex', alignItems: 'center', gap: 4, background: 'var(--bg-subtle, #f3f4f6)', borderRadius: 4, padding: '2px 8px', fontSize: 12 }}>
                        {s.name} <span style={{ color: 'var(--muted)', fontSize: 11 }}>({s.supply_type})</span>
                        <button
                          onClick={() => removeSupply(s.printer_supply_id)}
                          style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--danger, #dc2626)', padding: '0 2px', lineHeight: 1, fontSize: 14 }}
                          title="Desvincular"
                        >×</button>
                      </span>
                    ))}
                  </div>
              }
              <div style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                <div style={{ flex: 1 }}>
                  <Combobox
                    options={refs.supplyOpts.filter(o => !modelSupplies.some(s => s.printer_supply_id === o.value))}
                    value={addSupplyId}
                    onChange={setAddSupplyId}
                    placeholder="Afegir consumible…"
                  />
                </div>
                <button className="btn btn-ghost btn-sm" onClick={addSupply} disabled={!addSupplyId}>Afegir</button>
              </div>
            </div>

            {editErr && <div className="error-msg" style={{ marginBottom: 12 }}>{editErr}</div>}
            <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 16 }}>
              <button className="btn btn-ghost" onClick={() => setEditing(null)}>Tancar</button>
              <button className="btn btn-primary" onClick={saveEdit} disabled={editSaving || !editing.brand_id || !editing.model_name}>{editSaving ? '…' : 'Guardar'}</button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

// ─────────────────────────────────────────────
function PrinterSuppliesTab() {
  const SUPPLY_TYPES = ['toner', 'ink'];

  const [list, setList]     = useState([]);
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');
  const [editing, setEditing] = useState(null);
  const [editErr, setEditErr] = useState('');
  const [editSaving, setEditSaving] = useState(false);
  const cd = useConfirmDelete();

  const EMPTY = { name: '', supply_type: 'toner' };
  const [form, setForm] = useState(EMPTY);
  function setF(k, v) { setForm(f => ({ ...f, [k]: v })); }

  const load = useCallback(() => api.listPrinterSupplies().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => { load(); }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try { await api.createPrinterSupply({ name: form.name, supply_type: form.supply_type }); setForm(EMPTY); load(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function saveEdit() {
    setEditErr(''); setEditSaving(true);
    try {
      await api.updatePrinterSupply(editing.printer_supply_id, { name: editing.name, supply_type: editing.supply_type });
      setEditing(null); load();
    } catch (ex) { setEditErr(ex.message); }
    finally { setEditSaving(false); }
  }

  async function del(id) {
    try { await api.deletePrinterSupply(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  return (
    <>
      <Section title="Consumibles d'impressora">
        <div className="card" style={{ marginBottom: 14 }}>
          <form onSubmit={create} className="form-panel">
            <div className="form-grid" style={{ gridTemplateColumns: '1fr auto auto' }}>
              <div className="form-group">
                <label>Nom *</label>
                <input type="text" value={form.name} onChange={e => setF('name', e.target.value)} placeholder="HP 85A, Epson T501 Cyan…" required />
              </div>
              <div className="form-group">
                <label>Tipus *</label>
                <select value={form.supply_type} onChange={e => setF('supply_type', e.target.value)}>
                  {SUPPLY_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
                </select>
              </div>
              <div className="form-group" style={{ justifyContent: 'flex-end' }}>
                <label style={{ visibility: 'hidden' }}>_</label>
                <button type="submit" className="btn btn-primary" disabled={saving || !form.name}>{saving ? '…' : 'Afegir'}</button>
              </div>
            </div>
            {err && <div className="error-msg">{err}</div>}
          </form>
        </div>

        <div className="card">
          <div className="table-wrap">
            <table>
              <thead><tr><th>Nom</th><th>Tipus</th><th></th></tr></thead>
              <tbody>
                {list.length === 0 && <tr><td colSpan={3} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
                {list.map(s => (
                  <tr key={s.printer_supply_id}>
                    <td>{s.name}</td>
                    <td>{s.supply_type}</td>
                    <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                      {cd.isAsking(s.printer_supply_id)
                        ? <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                            <button className="btn btn-danger btn-sm" onClick={() => del(s.printer_supply_id)}>Sí</button>
                            <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                        : <>
                            <button className="btn btn-ghost btn-sm" onClick={() => { setEditErr(''); setEditing({ ...s }); }}>Editar</button>
                            <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => cd.askDelete(s.printer_supply_id)}>Eliminar</button>
                          </>
                      }
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </Section>

      {editing && (
        <div style={{ position: 'fixed', inset: 0, zIndex: 1000, background: 'rgba(0,0,0,0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '32px 16px' }}
          onClick={e => { if (e.target === e.currentTarget) setEditing(null); }}>
          <div className="card" style={{ width: '100%', maxWidth: 400, padding: 24 }}>
            <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>Editar consumible</h2>
            <div className="form-group" style={{ marginBottom: 12 }}>
              <label>Nom *</label>
              <input type="text" value={editing.name} onChange={e => setEditing(v => ({ ...v, name: e.target.value }))} required autoFocus />
            </div>
            <div className="form-group" style={{ marginBottom: 16 }}>
              <label>Tipus *</label>
              <select value={editing.supply_type} onChange={e => setEditing(v => ({ ...v, supply_type: e.target.value }))}>
                {SUPPLY_TYPES.map(t => <option key={t} value={t}>{t}</option>)}
              </select>
            </div>
            {editErr && <div className="error-msg" style={{ marginBottom: 12 }}>{editErr}</div>}
            <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
              <button className="btn btn-ghost" onClick={() => setEditing(null)}>Cancel·lar</button>
              <button className="btn btn-primary" onClick={saveEdit} disabled={editSaving || !editing.name}>{editSaving ? '…' : 'Guardar'}</button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

// ─────────────────────────────────────────────
function ProjectorModelsTab() {
  const [list, setList]       = useState([]);
  const [refs, setRefs]       = useState(null);
  const [saving, setSaving]   = useState(false);
  const [err, setErr]         = useState('');
  const [editing, setEditing] = useState(null);
  const [editErr, setEditErr] = useState('');
  const [editSaving, setEditSaving] = useState(false);
  const cd = useConfirmDelete();

  const EMPTY = { brand_id: null, model_name: '' };
  const [form, setForm] = useState(EMPTY);
  function setF(k, v) { setForm(f => ({ ...f, [k]: v })); }

  const load = useCallback(() => api.listProjectorModels().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => {
    load();
    api.listBrands()
      .then(brands => setRefs({ brandOpts: (brands ?? []).map(b => ({ value: b.brand_id, label: b.name })) }))
      .catch(() => {});
  }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try {
      await api.createProjectorModel({ brand_id: form.brand_id, model_name: form.model_name });
      setForm(EMPTY); load();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function saveEdit() {
    setEditErr(''); setEditSaving(true);
    try {
      await api.updateProjectorModel(editing.projector_model_id, {
        brand_id: editing.brand_id, model_name: editing.model_name,
      });
      setEditing(null); load();
    } catch (ex) { setEditErr(ex.message); }
    finally { setEditSaving(false); }
  }

  async function del(id) {
    try { await api.deleteProjectorModel(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  if (!refs) return <div className="empty">Carregant…</div>;

  return (
    <>
      <Section title="Models de projector">
        <div className="card" style={{ marginBottom: 14 }}>
          <form onSubmit={create} className="form-panel">
            <div className="form-grid">
              <div className="form-group">
                <label>Marca *</label>
                <Combobox options={refs.brandOpts} value={form.brand_id} onChange={v => setF('brand_id', v)} placeholder="Marca…" />
              </div>
              <div className="form-group">
                <label>Nom del model *</label>
                <input type="text" value={form.model_name} onChange={e => setF('model_name', e.target.value)} placeholder="EX3260" required />
              </div>
            </div>
            {err && <div className="error-msg">{err}</div>}
            <div className="form-actions">
              <button type="submit" className="btn btn-primary" disabled={saving || !form.brand_id || !form.model_name}>
                {saving ? '…' : 'Afegir model'}
              </button>
            </div>
          </form>
        </div>

        <div className="card">
          <div className="table-wrap">
            <table>
              <thead><tr><th>Marca</th><th>Model</th><th></th></tr></thead>
              <tbody>
                {list.length === 0 && <tr><td colSpan={3} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>Sense dades</td></tr>}
                {list.map(m => (
                  <tr key={m.projector_model_id}>
                    <td>{m.brand_name}</td>
                    <td>{m.model_name}</td>
                    <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                      {cd.isAsking(m.projector_model_id)
                        ? <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                            <button className="btn btn-danger btn-sm" onClick={() => del(m.projector_model_id)}>Sí</button>
                            <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                        : <>
                            <button className="btn btn-ghost btn-sm" onClick={() => { setEditErr(''); setEditing({ ...m }); }}>Editar</button>
                            <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => cd.askDelete(m.projector_model_id)}>Eliminar</button>
                          </>
                      }
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </Section>

      {editing && (
        <div style={{ position: 'fixed', inset: 0, zIndex: 1000, background: 'rgba(0,0,0,0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '32px 16px' }}
          onClick={e => { if (e.target === e.currentTarget) setEditing(null); }}>
          <div className="card" style={{ width: '100%', maxWidth: 440, padding: 24 }}>
            <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 16 }}>Editar model de projector</h2>
            <div className="form-grid">
              <div className="form-group">
                <label>Marca *</label>
                <Combobox options={refs.brandOpts} value={editing.brand_id} onChange={v => setEditing(e => ({ ...e, brand_id: v }))} placeholder="Marca…" />
              </div>
              <div className="form-group">
                <label>Nom del model *</label>
                <input type="text" value={editing.model_name} onChange={e => setEditing(v => ({ ...v, model_name: e.target.value }))} required autoFocus />
              </div>
            </div>
            {editErr && <div className="error-msg" style={{ marginBottom: 12 }}>{editErr}</div>}
            <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 16 }}>
              <button className="btn btn-ghost" onClick={() => setEditing(null)}>Cancel·lar</button>
              <button className="btn btn-primary" onClick={saveEdit} disabled={editSaving || !editing.brand_id || !editing.model_name}>{editSaving ? '…' : 'Guardar'}</button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}

export default function Reference() {
  const [tab, setTab] = useState('cpus');

  return (
    <>
      <div className="page-header">
        <h1 className="page-title">📋 Dades bàsiques</h1>
        <p> WARNING: Aquesta secció conté referències a dades de tota l'aplicació</p>
      </div>

      {/* Tab bar */}
      <div className="ref-tabs">
        {TABS.map(t => (
          <button
            key={t.id}
            className={`ref-tab${tab === t.id ? ' active' : ''}`}
            onClick={() => setTab(t.id)}
          >
            {t.label}
          </button>
        ))}
      </div>

      <div style={{ marginTop: 20 }}>
        {tab === 'cpus'       && <CpusTab />}
        {tab === 'os'         && <OsTab />}
        {tab === 'brands'     && <BrandsTab />}
        {tab === 'lmodels'    && <LaptopModelsTab />}
        {tab === 'dmodels'    && <DesktopModelsTab />}
        {tab === 'equipusers' && <EquipUsersTab />}
        {tab === 'centers'    && <CentersTab />}
        {tab === 'rooms'      && <RoomsTab />}
        {tab === 'cycles'     && <CyclesTab />}
        {tab === 'classes'    && <ClassesTab />}
        {tab === 'pmodels'    && <PrinterModelsTab />}
        {tab === 'psupplies'  && <PrinterSuppliesTab />}
        {tab === 'prjmodels'  && <ProjectorModelsTab />}
      </div>
    </>
  );
}
