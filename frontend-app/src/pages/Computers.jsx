import { useState, useEffect } from 'react';

// Utilidad para ordenar arrays de objetos por clave
function sortBy(arr, key, asc = true) {
  return [...arr].sort((a, b) => {
    const v1 = typeof key === 'function' ? key(a) : a[key];
    const v2 = typeof key === 'function' ? key(b) : b[key];
    if (v1 == null && v2 == null) return 0;
    if (v1 == null) return 1;
    if (v2 == null) return -1;
    if (v1 < v2) return asc ? -1 : 1;
    if (v1 > v2) return asc ? 1 : -1;
    return 0;
  });
}
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
  // Ordenado y paginación
  const [sortCol, setSortCol] = useState('hostname');
  const [sortAsc, setSortAsc] = useState(true);
  const [page, setPage] = useState(1);
  const PAGE_SIZE = 20;

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

  // Ordenar
  const sortKeyMap = {
    hostname: d => d.hostname?.toLowerCase() ?? '',
    aula: d => R.roomMap[d.room_id]?.toLowerCase() ?? '',
    model: d => d.desktop_model_id ? R.dmMap[d.desktop_model_id]?.toLowerCase() : '',
    cpu: d => d.cpu_id ? R.cpuMap[d.cpu_id]?.toLowerCase() : '',
    so: d => R.osMap[d.os_id]?.toLowerCase() ?? '',
    ram: d => d.ram_gb ?? 0,
    storage: d => d.storage_gb ?? 0,
    wifi: d => d.has_wifi_card ? 1 : 0,
    mac: d => d.mac_address?.toLowerCase() ?? '',
    user: d => d.equipment_user_id ? R.equipMap[d.equipment_user_id]?.toLowerCase() : '',
    obs: d => d.observations?.toLowerCase() ?? '',
  };
  const sortedDesktops = sortBy(filteredDesktops, sortKeyMap[sortCol] || sortKeyMap.hostname, sortAsc);
  // Paginado
  const totalPages = Math.ceil(sortedDesktops.length / PAGE_SIZE) || 1;
  const pagedDesktops = sortedDesktops.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);


  // Laptops: filtrado, orden y paginación
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

  const laptopSortKeyMap = {
    hostname: l => l.hostname?.toLowerCase() ?? '',
    aula: l => R.roomMap[l.room_id]?.toLowerCase() ?? '',
    model: l => l.laptop_model_id ? R.lmMap[l.laptop_model_id]?.toLowerCase() : '',
    so: l => R.osMap[l.os_id]?.toLowerCase() ?? '',
    ram: l => l.ram_gb ?? 0,
    storage: l => l.storage_gb ?? 0,
    mac: l => l.mac_address?.toLowerCase() ?? '',
    user: l => l.equipment_user_id ? R.equipMap[l.equipment_user_id]?.toLowerCase() : '',
    obs: l => l.observations?.toLowerCase() ?? '',
  };
  const sortedLaptops = sortBy(filteredLaptops, laptopSortKeyMap[sortCol] || laptopSortKeyMap.hostname, sortAsc);
  const totalLaptopPages = Math.ceil(sortedLaptops.length / PAGE_SIZE) || 1;
  const pagedLaptops = sortedLaptops.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

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
          <>
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('hostname'); setSortAsc(sortCol === 'hostname' ? !sortAsc : true); }}>
                      Hostname {sortCol === 'hostname' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('aula'); setSortAsc(sortCol === 'aula' ? !sortAsc : true); }}>
                      Aula {sortCol === 'aula' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('model'); setSortAsc(sortCol === 'model' ? !sortAsc : true); }}>
                      Model base {sortCol === 'model' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('cpu'); setSortAsc(sortCol === 'cpu' ? !sortAsc : true); }}>
                      CPU {sortCol === 'cpu' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('so'); setSortAsc(sortCol === 'so' ? !sortAsc : true); }}>
                      SO {sortCol === 'so' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('ram'); setSortAsc(sortCol === 'ram' ? !sortAsc : true); }}>
                      RAM {sortCol === 'ram' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('storage'); setSortAsc(sortCol === 'storage' ? !sortAsc : true); }}>
                      Emmagatzematge {sortCol === 'storage' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('wifi'); setSortAsc(sortCol === 'wifi' ? !sortAsc : true); }}>
                      WiFi {sortCol === 'wifi' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('mac'); setSortAsc(sortCol === 'mac' ? !sortAsc : true); }}>
                      MAC {sortCol === 'mac' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('user'); setSortAsc(sortCol === 'user' ? !sortAsc : true); }}>
                      Usuari equip {sortCol === 'user' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('obs'); setSortAsc(sortCol === 'obs' ? !sortAsc : true); }}>
                      Observacions {sortCol === 'obs' && (sortAsc ? '▲' : '▼')}
                    </th>
                    {canEdit && <th></th>}
                  </tr>
                </thead>
                <tbody>
                  {pagedDesktops.map(d => (
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
            {/* Paginación */}
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 12, marginTop: 12 }}>
              <button className="btn btn-sm" disabled={page === 1} onClick={() => setPage(page - 1)}>Anterior</button>
              <span>Pàgina {page} de {totalPages}</span>
              <button className="btn btn-sm" disabled={page === totalPages} onClick={() => setPage(page + 1)}>Següent</button>
            </div>
          </>
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
          <>
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('hostname'); setSortAsc(sortCol === 'hostname' ? !sortAsc : true); }}>
                      Hostname {sortCol === 'hostname' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('aula'); setSortAsc(sortCol === 'aula' ? !sortAsc : true); }}>
                      Aula {sortCol === 'aula' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('model'); setSortAsc(sortCol === 'model' ? !sortAsc : true); }}>
                      Model {sortCol === 'model' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('so'); setSortAsc(sortCol === 'so' ? !sortAsc : true); }}>
                      SO {sortCol === 'so' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('ram'); setSortAsc(sortCol === 'ram' ? !sortAsc : true); }}>
                      RAM {sortCol === 'ram' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('storage'); setSortAsc(sortCol === 'storage' ? !sortAsc : true); }}>
                      Emmagatzematge {sortCol === 'storage' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('mac'); setSortAsc(sortCol === 'mac' ? !sortAsc : true); }}>
                      MAC {sortCol === 'mac' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('user'); setSortAsc(sortCol === 'user' ? !sortAsc : true); }}>
                      Usuari equip {sortCol === 'user' && (sortAsc ? '▲' : '▼')}
                    </th>
                    <th style={{ cursor: 'pointer' }} onClick={() => { setSortCol('obs'); setSortAsc(sortCol === 'obs' ? !sortAsc : true); }}>
                      Observacions {sortCol === 'obs' && (sortAsc ? '▲' : '▼')}
                    </th>
                    {canEdit && <th></th>}
                  </tr>
                </thead>
                <tbody>
                  {pagedLaptops.map(l => (
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
            {/* Paginación */}
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 12, marginTop: 12 }}>
              <button className="btn btn-sm" disabled={page === 1} onClick={() => setPage(page - 1)}>Anterior</button>
              <span>Pàgina {page} de {totalLaptopPages}</span>
              <button className="btn btn-sm" disabled={page === totalLaptopPages} onClick={() => setPage(page + 1)}>Següent</button>
            </div>
          </>
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
