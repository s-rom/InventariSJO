import { useState, useEffect, useCallback } from 'react';
import { api } from '../api';
import Combobox from '../components/Combobox';

const RAM_TYPES     = ['DDR3', 'DDR4', 'DDR5', 'None'];
const STORAGE_TYPES = ['HDD', 'SSD', 'NVMe', 'None'];

const TABS = [
  { id: 'cpus',     label: '💾 CPUs' },
  { id: 'os',       label: '🖥 SO' },
  { id: 'brands',   label: '🏷 Marques' },
  { id: 'lmodels',  label: '💻 Models portàtil' },
  { id: 'dmodels',  label: '🖥️ Models sobretaula' },
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
  const [list, setList]       = useState([]);
  const [name, setName]       = useState('');
  const [score, setScore]     = useState('');
  const [saving, setSaving]   = useState(false);
  const [err, setErr]         = useState('');
  const cd                    = useConfirmDelete();

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

  async function del(id) {
    try { await api.deleteCpu(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  return (
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
                    {cd.isAsking(c.cpu_id)
                      ? <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                          <button className="btn btn-danger btn-sm" onClick={() => del(c.cpu_id)}>Sí</button>
                          <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                      : <button className="btn btn-danger btn-sm" onClick={() => cd.askDelete(c.cpu_id)}>Eliminar</button>
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
function OsTab() {
  const [list, setList]     = useState([]);
  const [name, setName]     = useState('');
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');
  const cd                  = useConfirmDelete();

  const load = useCallback(() => api.listOS().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => { load(); }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try { await api.createOS({ name }); setName(''); load(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function del(id) {
    try { await api.deleteOS(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  return (
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
                    {cd.isAsking(o.os_id)
                      ? <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                          <button className="btn btn-danger btn-sm" onClick={() => del(o.os_id)}>Sí</button>
                          <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                      : <button className="btn btn-danger btn-sm" onClick={() => cd.askDelete(o.os_id)}>Eliminar</button>
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
function BrandsTab() {
  const [list, setList]     = useState([]);
  const [name, setName]     = useState('');
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');
  const cd                  = useConfirmDelete();

  const load = useCallback(() => api.listBrands().then(d => setList(d ?? [])).catch(() => {}), []);
  useEffect(() => { load(); }, [load]);

  async function create(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try { await api.createBrand({ name }); setName(''); load(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function del(id) {
    try { await api.deleteBrand(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
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
                  <td style={{ textAlign: 'right' }}>
                    {cd.isAsking(b.brand_id)
                      ? <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                          <button className="btn btn-danger btn-sm" onClick={() => del(b.brand_id)}>Sí</button>
                          <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                      : <button className="btn btn-danger btn-sm" onClick={() => cd.askDelete(b.brand_id)}>Eliminar</button>
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
function LaptopModelsTab() {
  const [list, setList]     = useState([]);
  const [refs, setRefs]     = useState(null);
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');
  const cd                  = useConfirmDelete();

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

  async function del(id) {
    try { await api.deleteLaptopModel(id); load(); } catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
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
                  <td style={{ textAlign: 'right' }}>
                    {cd.isAsking(m.laptop_model_id)
                      ? <><span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Segur?</span>
                          <button className="btn btn-danger btn-sm" onClick={() => del(m.laptop_model_id)}>Sí</button>
                          <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button></>
                      : <button className="btn btn-danger btn-sm" onClick={() => cd.askDelete(m.laptop_model_id)}>Eliminar</button>
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
// Main Reference page
export default function Reference() {
  const [tab, setTab] = useState('cpus');

  return (
    <>
      <div className="page-header">
        <h1 className="page-title">📋 Dades bàsiques</h1>
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
        {tab === 'cpus'    && <CpusTab />}
        {tab === 'os'      && <OsTab />}
        {tab === 'brands'  && <BrandsTab />}
        {tab === 'lmodels' && <LaptopModelsTab />}
        {tab === 'dmodels' && <DesktopModelsTab />}
      </div>
    </>
  );
}
