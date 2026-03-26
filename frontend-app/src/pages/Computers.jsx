import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';
import { useAuth } from '../App';
import EditComputerModal from '../components/EditComputerModal';

function osEmoji(name = '') {
  const n = (name ?? '').toLowerCase();
  if (n.includes('chrome'))                    return '🟢 ' + name;
  if (n.includes('windows') || n.includes('win')) return '🪟 ' + name;
  if (n.includes('linux'))                     return '🐧 ' + name;
  return name || '—';
}

function ramLabel(gb, type) {
  if (!gb && !type) return '—';
  const parts = [];
  if (gb)   parts.push(`${gb} GB`);
  if (type) parts.push(type);
  return parts.join(' ');
}

function storageLabel(gb, type) {
  if (!gb && !type) return '—';
  const parts = [];
  if (gb)   parts.push(`${gb} GB`);
  if (type) parts.push(type);
  return parts.join(' ');
}

export default function Computers() {
  const navigate = useNavigate();
  const { role }  = useAuth();
  const canEdit   = role !== 'readonly';

  const [desktops,       setDesktops]       = useState([]);
  const [laptops,        setLaptops]        = useState([]);
  const [refs,           setRefs]           = useState(null);
  const [loading,        setLoading]        = useState(true);
  const [err,            setErr]            = useState('');
  const [query,          setQuery]          = useState('');
  const [editItem,       setEditItem]       = useState(null); // { item, type }

  useEffect(() => {
    async function load() {
      try {
        const [dt, lt, cpus, osList, equip, centers, lm, dm] = await Promise.all([
          api.listDesktops(),
          api.listLaptops(),
          api.listCpus(),
          api.listOS(),
          api.listEquipmentUsers(),
          api.listCenters(),
          api.listLaptopModels(),
          api.listDesktopModels(),
        ]);

        // Load rooms for all centers
        const roomsNested = await Promise.all(
          (centers ?? []).map(c =>
            api.listRoomsByCenter(c.center_id)
              .then(rs => (rs ?? []).map(r => ({ ...r, centerName: c.name })))
              .catch(() => [])
          )
        );
        const allRooms = roomsNested.flat();

        const cpuMap      = Object.fromEntries((cpus ?? []).map(c => [c.cpu_id,      c.model_name]));
        const osMap       = Object.fromEntries((osList ?? []).map(o => [o.os_id,      o.name]));
        const equipMap    = Object.fromEntries((equip ?? []).map(e => [e.equipment_user_id, e.name]));
        const roomMap     = Object.fromEntries(allRooms.map(r => [r.room_id, `${r.centerName} › ${r.name}`]));
        const lmMap       = Object.fromEntries((lm ?? []).map(m => [m.laptop_model_id,  `${m.brand_name} ${m.model_name}`]));
        const dmMap       = Object.fromEntries((dm ?? []).map(m => [m.desktop_model_id, `${m.brand_name} ${m.model_name}`]));

        const cpuOpts   = (cpus   ?? []).map(c => ({ value: c.cpu_id,              label: c.model_name }));
        const osOpts    = (osList ?? []).map(o => ({ value: o.os_id,               label: o.name }));
        const equipOpts = (equip  ?? []).map(e => ({ value: e.equipment_user_id,   label: e.name }));
        const roomOpts  = allRooms.map(r => ({ value: r.room_id, label: `${r.centerName} › ${r.name}` }));
        const lmOpts    = (lm     ?? []).map(m => ({ value: m.laptop_model_id,     label: `${m.brand_name} ${m.model_name}` }));
        const dmOpts    = (dm     ?? []).map(m => ({ value: m.desktop_model_id,    label: `${m.brand_name} ${m.model_name}` }));

        setRefs({ cpuMap, osMap, equipMap, roomMap, lmMap, dmMap, cpuOpts, osOpts, equipOpts, roomOpts, lmOpts, dmOpts });
        setDesktops(dt ?? []);
        setLaptops(lt ?? []);
      } catch (e) {
        setErr(e.message);
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  if (loading) return <div className="empty">Carregant…</div>;
  if (err)     return <div className="empty" style={{ color: 'var(--danger)' }}>Error: {err}</div>;

  const R = refs;

  const q = query.trim().toLowerCase();
  const filteredDesktops = q
    ? desktops.filter(d => [
        d.hostname,
        R.roomMap[d.room_id],
        d.desktop_model_id ? R.dmMap[d.desktop_model_id] : null,
        d.cpu_id ? R.cpuMap[d.cpu_id] : null,
        R.osMap[d.os_id],
        d.ram_gb != null ? `${d.ram_gb} GB` : null,
        d.ram_type,
        d.storage_gb != null ? `${d.storage_gb} GB` : null,
        d.storage_type,
        d.mac_address,
        d.equipment_user_id ? R.equipMap[d.equipment_user_id] : null,
        d.observations,
      ].filter(Boolean).join(' ').toLowerCase().includes(q))
    : desktops;

  const filteredLaptops = q
    ? laptops.filter(l => [
        l.hostname,
        R.roomMap[l.room_id],
        l.laptop_model_id ? R.lmMap[l.laptop_model_id] : null,
        R.osMap[l.os_id],
        l.ram_gb != null ? `${l.ram_gb} GB` : null,
        l.ram_type,
        l.storage_gb != null ? `${l.storage_gb} GB` : null,
        l.storage_type,
        l.mac_address,
        l.serial_number,
        l.equipment_user_id ? R.equipMap[l.equipment_user_id] : null,
        l.observations,
      ].filter(Boolean).join(' ').toLowerCase().includes(q))
    : laptops;

  async function handleDelete(id, hostname) {
    if (!confirm(`Eliminar l'equip "${hostname}"?`)) return;
    try {
      await api.deleteComputer(id);
      setDesktops(prev => prev.filter(d => d.computer_id !== id));
      setLaptops(prev  => prev.filter(l => l.computer_id !== id));
    } catch (err) {
      alert(err.message);
    }
  }

  function handleEditSave(updated, type) {
    if (type === 'desktop') {
      setDesktops(prev => prev.map(d => d.computer_id === updated.computer_id ? updated : d));
    } else {
      setLaptops(prev => prev.map(l => l.computer_id === updated.computer_id ? updated : l));
    }
    setEditItem(null);
  }

  return (
    <>
      {/* SEARCH BAR */}
      <div style={{ marginBottom: 24 }}>
        <input
          type="search"
          value={query}
          onChange={e => setQuery(e.target.value)}
          placeholder="Cercar per hostname, aula, model, CPU, MAC…"
          style={{ width: '100%', maxWidth: 480, padding: '8px 12px', borderRadius: 6, border: '1px solid var(--border)', fontSize: 14 }}
        />
        {q && (
          <span style={{ marginLeft: 12, fontSize: 13, color: 'var(--muted)' }}>
            {filteredDesktops.length + filteredLaptops.length} resultat{filteredDesktops.length + filteredLaptops.length !== 1 ? 's' : ''}
          </span>
        )}
      </div>

      {/* DESKTOPS */}
      <div className="page-header">
        <h1 className="page-title">🖥️ Sobretaules</h1>
        <button className="btn btn-primary" onClick={() => navigate('/computers/new-desktop')}>
          + Nou sobretaula
        </button>
      </div>

      <div className="card" style={{ marginBottom: 32 }}>
        {filteredDesktops.length === 0 ? (
          <div className="empty">{q ? 'Cap resultat.' : 'Cap sobretaula registrat.'}</div>
        ) : (
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Hostname</th>
                  <th>Aula</th>
                  <th>Model base</th>
                  <th>CPU</th>
                  <th>SO</th>
                  <th>RAM</th>
                  <th>Emmagatzematge</th>
                  <th>WiFi</th>
                  <th>MAC</th>
                  <th>Usuari equip</th>
                  <th>Observacions</th>
                  {canEdit && <th></th>}
                </tr>
              </thead>
              <tbody>
                {filteredDesktops.map(d => (
                  <tr key={d.computer_id}>
                    <td><strong>{d.hostname}</strong></td>
                    <td>{R.roomMap[d.room_id] ?? '—'}</td>
                    <td>{d.desktop_model_id ? R.dmMap[d.desktop_model_id] ?? '—' : <span style={{ color: 'var(--muted)' }}>sense model</span>}</td>
                    <td>{d.cpu_id ? R.cpuMap[d.cpu_id] ?? '—' : '—'}</td>
                    <td>{osEmoji(R.osMap[d.os_id])}</td>
                    <td>{ramLabel(d.ram_gb, d.ram_type)}</td>
                    <td>{storageLabel(d.storage_gb, d.storage_type)}</td>
                    <td style={{ textAlign: 'center' }}>{d.has_wifi_card ? '✔' : '—'}</td>
                    <td><code style={{ fontSize: 11 }}>{d.mac_address ?? '—'}</code></td>
                    <td>{d.equipment_user_id ? R.equipMap[d.equipment_user_id] ?? '—' : '—'}</td>
                    <td style={{ fontSize: 12, color: 'var(--muted)', maxWidth: 200, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{d.observations ?? '—'}</td>
                    {canEdit && (
                      <td>
                        <div style={{ display: 'flex', gap: 6 }}>
                          <button className="btn btn-sm" onClick={() => setEditItem({ item: d, type: 'desktop' })}>
                            Editar
                          </button>
                          <button className="btn btn-danger btn-sm" onClick={() => handleDelete(d.computer_id, d.hostname)}>
                            Eliminar
                          </button>
                        </div>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* LAPTOPS */}
      <div className="page-header">
        <h1 className="page-title">💻 Portàtils</h1>
        <button className="btn btn-primary" onClick={() => navigate('/computers/new-laptop')}>
          + Nou portàtil
        </button>
      </div>

      <div className="card">
        {filteredLaptops.length === 0 ? (
          <div className="empty">{q ? 'Cap resultat.' : 'Cap portàtil registrat.'}</div>
        ) : (
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Hostname</th>
                  <th>Aula</th>
                  <th>Model</th>
                  <th>SO</th>
                  <th>RAM</th>
                  <th>Emmagatzematge</th>
                  <th>MAC</th>
                  <th>Usuari equip</th>
                  <th>Observacions</th>
                  {canEdit && <th></th>}
                </tr>
              </thead>
              <tbody>
                {filteredLaptops.map(l => (
                  <tr key={l.computer_id}>
                    <td><strong>{l.hostname}</strong></td>
                    <td>{R.roomMap[l.room_id] ?? '—'}</td>
                    <td>{l.laptop_model_id ? R.lmMap[l.laptop_model_id] ?? '—' : '—'}</td>
                    <td>{osEmoji(R.osMap[l.os_id])}</td>
                    <td>{ramLabel(l.ram_gb, l.ram_type)}</td>
                    <td>{storageLabel(l.storage_gb, l.storage_type)}</td>
                    <td><code style={{ fontSize: 11 }}>{l.mac_address ?? '—'}</code></td>
                    <td>{l.equipment_user_id ? R.equipMap[l.equipment_user_id] ?? '—' : '—'}</td>
                    <td style={{ fontSize: 12, color: 'var(--muted)', maxWidth: 200, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>{l.observations ?? '—'}</td>
                    {canEdit && (
                      <td>
                        <div style={{ display: 'flex', gap: 6 }}>
                          <button className="btn btn-sm" onClick={() => setEditItem({ item: l, type: 'laptop' })}>
                            Editar
                          </button>
                          <button className="btn btn-danger btn-sm" onClick={() => handleDelete(l.computer_id, l.hostname)}>
                            Eliminar
                          </button>
                        </div>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* EDIT MODAL */}
      {editItem && (
        <EditComputerModal
          item={editItem.item}
          type={editItem.type}
          refs={R}
          onSave={updated => handleEditSave(updated, editItem.type)}
          onClose={() => setEditItem(null)}
        />
      )}
    </>
  );
}
