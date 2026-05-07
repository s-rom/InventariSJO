import { useState, useEffect, useRef } from 'react';
import { api } from '../api';
import { useAuth } from '../App';
import Combobox from '../components/Combobox';
import StatusBadge from '../components/StatusBadge';
import Pagination from '../components/Pagination';

const DEVICE_STATUSES = ['actiu', 'baixa'];

const EMPTY = {
  projector_model_id: null,
  serial_number:     '',
  status:            'actiu',
  room_id:           null,
  equipment_user_id: null,
  observations:      '',
};

function ProjectorFormFields({ form, set, refs }) {
  return (
    <div className="form-grid">
      <div className="form-group">
        <label>Model *</label>
        <Combobox options={refs.modelOpts} value={form.projector_model_id} onChange={v => set('projector_model_id', v)} placeholder="Selecciona model…" />
      </div>
      <div className="form-group">
        <label>Número de sèrie</label>
        <input type="text" value={form.serial_number} onChange={e => set('serial_number', e.target.value)} placeholder="Opcional" />
      </div>
      <div className="form-group">
        <label>Estat</label>
        <select value={form.status} onChange={e => set('status', e.target.value)}>
          {DEVICE_STATUSES.map(s => <option key={s} value={s}>{s}</option>)}
        </select>
      </div>
      <div className="form-group">
        <label>Aula</label>
        <Combobox options={refs.roomOpts} value={form.room_id} onChange={v => set('room_id', v)} placeholder="Opcional…" nullable />
      </div>
      <div className="form-group">
        <label>Usuari d'equip</label>
        <Combobox options={refs.equipOpts} value={form.equipment_user_id} onChange={v => set('equipment_user_id', v)} placeholder="Opcional…" nullable />
      </div>
      <div className="form-group" style={{ gridColumn: '1 / -1' }}>
        <label>Observacions</label>
        <textarea value={form.observations} onChange={e => set('observations', e.target.value)} rows={3} style={{ resize: 'vertical' }} />
      </div>
    </div>
  );
}

function EditModal({ projector, refs, onClose, onSaved }) {
  const [form, setForm] = useState({
    projector_model_id: projector.projector_model_id ?? null,
    serial_number:     projector.serial_number ?? '',
    status:            projector.status,
    room_id:           projector.room_id ?? null,
    equipment_user_id: projector.equipment_user_id ?? null,
    observations:      projector.observations ?? '',
  });
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');

  function set(k, v) { setForm(f => ({ ...f, [k]: v })); }

  async function handleSubmit(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try {
      await api.updateProjector(projector.projector_id, {
        projector_model_id: form.projector_model_id,
        serial_number:     form.serial_number || null,
        status:            form.status,
        room_id:           form.room_id ?? null,
        equipment_user_id: form.equipment_user_id ?? null,
        observations:      form.observations || null,
      });
      onSaved();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  return (
    <div
      style={{ position: 'fixed', inset: 0, zIndex: 1000, background: 'rgba(0,0,0,0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '32px 16px', overflowY: 'auto' }}
      onClick={e => { if (e.target === e.currentTarget) onClose(); }}
    >
      <div className="card" style={{ width: '100%', maxWidth: 560, padding: 24 }}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 20 }}>Editar projector #{projector.projector_id}</h2>
        <form onSubmit={handleSubmit} className="form-panel" style={{ padding: 0 }}>
          <ProjectorFormFields form={form} set={set} refs={refs} />
          {err && <div className="error-msg" style={{ marginBottom: 12 }}>{err}</div>}
          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 8 }}>
            <button type="button" className="btn btn-ghost" onClick={onClose}>Cancel·lar</button>
            <button type="submit" className="btn btn-primary" disabled={saving}>{saving ? '…' : 'Guardar'}</button>
          </div>
        </form>
      </div>
    </div>
  );
}

function NewModal({ refs, onClose, onSaved }) {
  const [form, setForm] = useState({ ...EMPTY });
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');

  function set(k, v) { setForm(f => ({ ...f, [k]: v })); }

  async function handleSubmit(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try {
      await api.createProjector({
        projector_model_id: form.projector_model_id,
        serial_number:     form.serial_number || null,
        status:            form.status,
        room_id:           form.room_id ?? null,
        equipment_user_id: form.equipment_user_id ?? null,
        observations:      form.observations || null,
      });
      onSaved();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  return (
    <div
      style={{ position: 'fixed', inset: 0, zIndex: 1000, background: 'rgba(0,0,0,0.45)', display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '32px 16px', overflowY: 'auto' }}
      onClick={e => { if (e.target === e.currentTarget) onClose(); }}
    >
      <div className="card" style={{ width: '100%', maxWidth: 560, padding: 24 }}>
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 20 }}>Nou projector</h2>
        <form onSubmit={handleSubmit} className="form-panel" style={{ padding: 0 }}>
          <ProjectorFormFields form={form} set={set} refs={refs} />
          {err && <div className="error-msg" style={{ marginBottom: 12 }}>{err}</div>}
          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 8 }}>
            <button type="button" className="btn btn-ghost" onClick={onClose}>Cancel·lar</button>
            <button type="submit" className="btn btn-primary" disabled={saving}>{saving ? '…' : 'Crear'}</button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default function Projectors() {
  const { role } = useAuth();
  const canEdit  = role !== 'readonly';
  const isAdmin  = role === 'admin';

  const [projectors, setProjectors] = useState([]);
  const [refs,       setRefs]       = useState(null);
  const [loading,    setLoading]    = useState(true);
  const [err,        setErr]        = useState('');
  const [query,      setQuery]      = useState('');
  const [editItem,   setEditItem]   = useState(null);
  const [showNew,    setShowNew]    = useState(false);
  const [delId,      setDelId]      = useState(null);
  const [toast,      setToast]      = useState(null);
  const [page,       setPage]       = useState(1);
  const [pageSize,   setPageSize]   = useState(20);
  const toastTimer = useRef(null);

  function showToast(type, msg) {
    setToast({ type, msg });
    clearTimeout(toastTimer.current);
    toastTimer.current = setTimeout(() => setToast(null), 4000);
  }

  async function load() {
    try {
      const [data, models, equip, centers] = await Promise.all([
        api.listProjectors(),
        api.listProjectorModels(),
        api.listEquipmentUsers(),
        api.listCenters(),
      ]);
      const roomsNested = await Promise.all(
        (centers ?? []).map(c =>
          api.listRoomsByCenter(c.center_id)
            .then(rs => (rs ?? []).map(r => ({ ...r, centerName: c.name })))
            .catch(() => [])
        )
      );
      const allRooms = roomsNested.flat();
      const modelOpts = (models ?? []).map(m => ({ value: m.projector_model_id, label: `${m.brand_name} ${m.model_name}` }));
      const equipOpts = (equip  ?? []).map(e => ({ value: e.equipment_user_id, label: e.name }));
      const roomOpts  = allRooms.map(r => ({ value: r.room_id, label: `${r.centerName} › ${r.name}` }));
      const roomMap   = Object.fromEntries(allRooms.map(r => [r.room_id, `${r.centerName} › ${r.name}`]));
      setRefs({ modelOpts, equipOpts, roomOpts, roomMap });
      setProjectors(data ?? []);
    } catch (e) {
      setErr(e.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); }, []); // eslint-disable-line

  const filtered = projectors.filter(p => {
    if (!query) return true;
    const q = query.toLowerCase();
    return (
      (p.model_name  ?? '').toLowerCase().includes(q) ||
      (p.brand_name  ?? '').toLowerCase().includes(q) ||
      (refs?.roomMap?.[p.room_id] ?? p.room_name ?? '').toLowerCase().includes(q) ||
      (p.serial_number ?? '').toLowerCase().includes(q) ||
      (p.equipment_user_name ?? '').toLowerCase().includes(q)
    );
  });

  const totalPages = Math.ceil(filtered.length / pageSize) || 1;
  const paged = filtered.slice((page - 1) * pageSize, page * pageSize);

  async function handleDelete(id) {
    try {
      await api.deleteProjector(id);
      showToast('ok', 'Projector eliminat');
      setDelId(null);
      load();
    } catch (ex) { showToast('err', ex.message); setDelId(null); }
  }

  if (loading) return <div className="empty">Carregant…</div>;
  if (err)     return <div className="empty" style={{ color: 'var(--danger)' }}>{err}</div>;

  return (
    <>
      <div className="page-header">
        <h1 className="page-title">📽️ Projectors</h1>
        {canEdit && (
          <button className="btn btn-primary" onClick={() => setShowNew(true)}>+ Nou projector</button>
        )}
      </div>

      <div style={{ marginBottom: 14 }}>
        <input
          type="search"
          placeholder="Cercar per marca, model, aula…"
          value={query}
          onChange={e => setQuery(e.target.value)}
          style={{ width: '100%', maxWidth: 360 }}
        />
      </div>

      <div className="card">
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Marca</th>
                <th>Model</th>
                <th>Núm. sèrie</th>
                <th>Estat</th>
                <th>Aula</th>
                <th>Usuari</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {paged.length === 0 && (
                <tr><td colSpan={8} style={{ textAlign: 'center', color: 'var(--muted)', padding: 20 }}>Sense projectors</td></tr>
              )}
              {paged.map(p => (
                <tr key={p.projector_id}>
                  <td style={{ color: 'var(--muted)', fontSize: 12 }}>#{p.projector_id}</td>
                  <td>{p.brand_name ?? '—'}</td>
                  <td>{p.model_name ?? '—'}</td>
                  <td style={{ fontFamily: 'monospace', fontSize: 12 }}>{p.serial_number ?? '—'}</td>
                  <td><StatusBadge status={p.status} /></td>
                  <td>{refs?.roomMap?.[p.room_id] ?? '—'}</td>
                  <td>{p.equipment_user_name ?? '—'}</td>
                  <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                    {canEdit && (
                      <>
                        <button className="btn btn-ghost btn-sm" onClick={() => setEditItem(p)}>Editar</button>
                        {isAdmin && (
                          delId === p.projector_id
                            ? <><span style={{ fontSize: 12, marginLeft: 8, marginRight: 6, color: 'var(--muted)' }}>Segur?</span>
                                <button className="btn btn-danger btn-sm" onClick={() => handleDelete(p.projector_id)}>Sí</button>
                                <button className="btn btn-ghost btn-sm" onClick={() => setDelId(null)}>No</button></>
                            : <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => setDelId(p.projector_id)}>Eliminar</button>
                        )}
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <Pagination
          page={page}
          totalPages={totalPages}
          pageSize={pageSize}
          onPage={setPage}
          onPageSize={setPageSize}
          pageSizeId="projectors-page-size"
        />
      </div>

      {editItem && refs && (
        <EditModal
          projector={editItem}
          refs={refs}
          onClose={() => setEditItem(null)}
          onSaved={() => { setEditItem(null); showToast('ok', 'Projector actualitzat'); load(); }}
        />
      )}

      {showNew && refs && (
        <NewModal
          refs={refs}
          onClose={() => setShowNew(false)}
          onSaved={() => { setShowNew(false); showToast('ok', 'Projector creat'); load(); }}
        />
      )}

      {toast && (
        <div style={{
          position: 'fixed', bottom: 24, right: 24, zIndex: 2000,
          padding: '12px 20px', borderRadius: 8, fontWeight: 500, fontSize: 14,
          background: toast.type === 'ok' ? 'var(--success, #065f46)' : 'var(--danger, #991b1b)',
          color: '#fff', boxShadow: '0 4px 16px rgba(0,0,0,0.18)',
        }}>
          {toast.type === 'ok' ? '✅' : '❌'} {toast.msg}
        </div>
      )}
    </>
  );
}
