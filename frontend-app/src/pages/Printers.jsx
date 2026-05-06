import { useState, useEffect, useRef } from 'react';
import { api } from '../api';
import { useAuth } from '../App';
import Combobox from '../components/Combobox';

const PRINTER_TYPES  = ['toner', 'ink', 'managed'];
const DEVICE_STATUSES = ['actiu', 'baixa'];

const EMPTY = {
  printer_model_id:       null,
  status:                 'actiu',
  has_network_capability: false,
  uses_network:           false,
  ip_address:             '',
  room_id:                null,
  equipment_user_id:      null,
  observations:           '',
};

function StatusBadge({ status }) {
  const ok = status === 'actiu';
  return (
    <span style={{
      display: 'inline-block',
      padding: '2px 8px',
      borderRadius: 999,
      fontSize: 11,
      fontWeight: 600,
      background: ok ? 'var(--success-bg, #d1fae5)' : 'var(--danger-bg, #fee2e2)',
      color: ok ? 'var(--success, #065f46)' : 'var(--danger, #991b1b)',
    }}>
      {ok ? 'Actiu' : 'Baixa'}
    </span>
  );
}

function EditModal({ printer, refs, onClose, onSaved }) {
  const [form, setForm] = useState({
    printer_model_id:       printer.printer_model_id,
    status:                 printer.status,
    has_network_capability: printer.has_network_capability,
    uses_network:           printer.uses_network,
    ip_address:             printer.ip_address ?? '',
    room_id:                printer.room_id ?? null,
    equipment_user_id:      printer.equipment_user_id ?? null,
    observations:           printer.observations ?? '',
  });
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');

  function set(k, v) { setForm(f => ({ ...f, [k]: v })); }

  async function handleSubmit(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try {
      await api.updatePrinter(printer.printer_id, {
        printer_model_id:       form.printer_model_id ?? undefined,
        status:                 form.status,
        has_network_capability: form.has_network_capability,
        uses_network:           form.uses_network,
        ip_address:             form.ip_address || null,
        room_id:                form.room_id ?? null,
        equipment_user_id:      form.equipment_user_id ?? null,
        observations:           form.observations || null,
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
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 20 }}>Editar impressora #{printer.printer_id}</h2>
        <form onSubmit={handleSubmit} className="form-panel" style={{ padding: 0 }}>
          <div className="form-grid">
            <div className="form-group">
              <label>Model d'impressora *</label>
              <Combobox
                options={refs.modelOpts}
                value={form.printer_model_id}
                onChange={v => set('printer_model_id', v)}
                placeholder="Model…"
              />
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
          </div>

          <div className="form-group" style={{ marginTop: 4 }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer', userSelect: 'none' }}>
              <input type="checkbox" checked={form.has_network_capability} onChange={e => {
                const v = e.target.checked;
                set('has_network_capability', v);
                if (!v) { set('uses_network', false); set('ip_address', ''); }
              }} />
              Té capacitat de xarxa
            </label>
          </div>

          {form.has_network_capability && (
            <>
              <div className="form-group">
                <label style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer', userSelect: 'none' }}>
                  <input type="checkbox" checked={form.uses_network} onChange={e => set('uses_network', e.target.checked)} />
                  S'usa en xarxa
                </label>
              </div>
              <div className="form-group">
                <label>Adreça IP</label>
                <input type="text" value={form.ip_address} onChange={e => set('ip_address', e.target.value)} placeholder="192.168.1.100" />
              </div>
            </>
          )}

          <div className="form-group">
            <label>Observacions</label>
            <textarea value={form.observations} onChange={e => set('observations', e.target.value)} rows={3} style={{ resize: 'vertical' }} />
          </div>

          {err && <div className="error-msg" style={{ marginBottom: 12 }}>{err}</div>}
          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 8 }}>
            <button type="button" className="btn btn-ghost" onClick={onClose}>Cancel·lar</button>
            <button type="submit" className="btn btn-primary" disabled={saving || !form.printer_model_id}>{saving ? '…' : 'Guardar'}</button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default function Printers() {
  const { role } = useAuth();
  const canEdit  = role !== 'readonly';
  const isAdmin  = role === 'admin';

  const [printers,  setPrinters]  = useState([]);
  const [refs,      setRefs]      = useState(null);
  const [loading,   setLoading]   = useState(true);
  const [err,       setErr]       = useState('');
  const [query,     setQuery]     = useState('');
  const [editItem,  setEditItem]  = useState(null);
  const [showNew,   setShowNew]   = useState(false);
  const [delId,     setDelId]     = useState(null);
  const [toast,     setToast]     = useState(null);
  const toastTimer = useRef(null);

  function showToast(type, msg) {
    setToast({ type, msg });
    clearTimeout(toastTimer.current);
    toastTimer.current = setTimeout(() => setToast(null), 4000);
  }

  async function load() {
    try {
      const [data, models, equip, centers] = await Promise.all([
        api.listPrinters(),
        api.listPrinterModels(),
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
      const modelOpts = (models ?? []).map(m => ({ value: m.printer_model_id, label: `${m.brand_name} ${m.model_name} (${m.printer_type}, ${m.print_color})` }));
      const equipOpts = (equip  ?? []).map(e => ({ value: e.equipment_user_id, label: e.name }));
      const roomOpts  = allRooms.map(r => ({ value: r.room_id, label: `${r.centerName} › ${r.name}` }));
      setRefs({ modelOpts, equipOpts, roomOpts });
      setPrinters(data ?? []);
    } catch (e) {
      setErr(e.message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); }, []); // eslint-disable-line

  const filtered = printers.filter(p => {
    if (!query) return true;
    const q = query.toLowerCase();
    return (
      (p.model_name ?? '').toLowerCase().includes(q) ||
      (p.brand_name ?? '').toLowerCase().includes(q) ||
      (p.room_name  ?? '').toLowerCase().includes(q) ||
      (p.ip_address ?? '').toLowerCase().includes(q) ||
      (p.equipment_user_name ?? '').toLowerCase().includes(q)
    );
  });

  async function handleDelete(id) {
    try {
      await api.deletePrinter(id);
      showToast('ok', 'Impressora eliminada');
      setDelId(null);
      load();
    } catch (ex) { showToast('err', ex.message); setDelId(null); }
  }

  if (loading) return <div className="empty">Carregant…</div>;
  if (err)     return <div className="empty" style={{ color: 'var(--danger)' }}>{err}</div>;

  return (
    <>
      <div className="page-header">
        <h1 className="page-title">🖨️ Impressores</h1>
        {canEdit && (
          <button className="btn btn-primary" onClick={() => setShowNew(true)}>+ Nova impressora</button>
        )}
      </div>

      <div style={{ marginBottom: 14 }}>
        <input
          type="search"
          placeholder="Cercar per model, aula, IP…"
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
                <th>Model</th>
                <th>Tipus</th>
                <th>Estat</th>
                <th>Xarxa</th>
                <th>IP</th>
                <th>Aula</th>
                <th>Usuari</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {filtered.length === 0 && (
                <tr><td colSpan={9} style={{ textAlign: 'center', color: 'var(--muted)', padding: 20 }}>Sense impressores</td></tr>
              )}
              {filtered.map(p => (
                <tr key={p.printer_id}>
                  <td style={{ color: 'var(--muted)', fontSize: 12 }}>#{p.printer_id}</td>
                  <td><strong>{p.brand_name}</strong> {p.model_name}</td>
                  <td style={{ fontSize: 12 }}>{p.printer_type} · {p.print_color}</td>
                  <td><StatusBadge status={p.status} /></td>
                  <td style={{ fontSize: 12 }}>
                    {p.has_network_capability
                      ? (p.uses_network ? '✅ en xarxa' : '⬜ no en ús')
                      : '—'}
                  </td>
                  <td style={{ fontFamily: 'monospace', fontSize: 12 }}>{p.ip_address ?? '—'}</td>
                  <td>{p.room_name ?? '—'}</td>
                  <td>{p.equipment_user_name ?? '—'}</td>
                  <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                    {canEdit && (
                      <>
                        <button className="btn btn-ghost btn-sm" onClick={() => setEditItem(p)}>Editar</button>
                        {isAdmin && (
                          delId === p.printer_id
                            ? <><span style={{ fontSize: 12, marginLeft: 8, marginRight: 6, color: 'var(--muted)' }}>Segur?</span>
                                <button className="btn btn-danger btn-sm" onClick={() => handleDelete(p.printer_id)}>Sí</button>
                                <button className="btn btn-ghost btn-sm" onClick={() => setDelId(null)}>No</button></>
                            : <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => setDelId(p.printer_id)}>Eliminar</button>
                        )}
                      </>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {editItem && refs && (
        <EditModal
          printer={editItem}
          refs={refs}
          onClose={() => setEditItem(null)}
          onSaved={() => { setEditItem(null); showToast('ok', 'Impressora actualitzada'); load(); }}
        />
      )}

      {showNew && refs && (
        <NewPrinterModal
          refs={refs}
          onClose={() => setShowNew(false)}
          onSaved={() => { setShowNew(false); showToast('ok', 'Impressora creada'); load(); }}
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

function NewPrinterModal({ refs, onClose, onSaved }) {
  const [form, setForm] = useState({ ...EMPTY });
  const [saving, setSaving] = useState(false);
  const [err, setErr]       = useState('');

  function set(k, v) { setForm(f => ({ ...f, [k]: v })); }

  async function handleSubmit(e) {
    e.preventDefault(); setErr(''); setSaving(true);
    try {
      await api.createPrinter({
        printer_model_id:       form.printer_model_id,
        status:                 form.status,
        has_network_capability: form.has_network_capability,
        uses_network:           form.uses_network,
        ip_address:             form.ip_address || null,
        room_id:                form.room_id ?? null,
        equipment_user_id:      form.equipment_user_id ?? null,
        observations:           form.observations || null,
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
        <h2 style={{ fontSize: 16, fontWeight: 600, marginBottom: 20 }}>Nova impressora</h2>
        <form onSubmit={handleSubmit} className="form-panel" style={{ padding: 0 }}>
          <div className="form-grid">
            <div className="form-group">
              <label>Model d'impressora *</label>
              <Combobox
                options={refs.modelOpts}
                value={form.printer_model_id}
                onChange={v => set('printer_model_id', v)}
                placeholder="Model…"
              />
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
          </div>

          <div className="form-group" style={{ marginTop: 4 }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer', userSelect: 'none' }}>
              <input type="checkbox" checked={form.has_network_capability} onChange={e => {
                const v = e.target.checked;
                set('has_network_capability', v);
                if (!v) { set('uses_network', false); set('ip_address', ''); }
              }} />
              Té capacitat de xarxa
            </label>
          </div>

          {form.has_network_capability && (
            <>
              <div className="form-group">
                <label style={{ display: 'flex', alignItems: 'center', gap: 8, cursor: 'pointer', userSelect: 'none' }}>
                  <input type="checkbox" checked={form.uses_network} onChange={e => set('uses_network', e.target.checked)} />
                  S'usa en xarxa
                </label>
              </div>
              <div className="form-group">
                <label>Adreça IP</label>
                <input type="text" value={form.ip_address} onChange={e => set('ip_address', e.target.value)} placeholder="192.168.1.100" />
              </div>
            </>
          )}

          <div className="form-group">
            <label>Observacions</label>
            <textarea value={form.observations} onChange={e => set('observations', e.target.value)} rows={3} style={{ resize: 'vertical' }} />
          </div>

          {err && <div className="error-msg" style={{ marginBottom: 12 }}>{err}</div>}
          <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end', marginTop: 8 }}>
            <button type="button" className="btn btn-ghost" onClick={onClose}>Cancel·lar</button>
            <button type="submit" className="btn btn-primary" disabled={saving || !form.printer_model_id}>{saving ? '…' : 'Crear'}</button>
          </div>
        </form>
      </div>
    </div>
  );
}
