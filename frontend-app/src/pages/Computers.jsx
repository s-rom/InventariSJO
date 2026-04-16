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
  const isAdmin   = role === 'admin';

  const [desktops,       setDesktops]       = useState([]);
  const [laptops,        setLaptops]        = useState([]);
  const [refs,           setRefs]           = useState(null);
  const [loading,        setLoading]        = useState(true);
  const [err,            setErr]            = useState('');
  const [query,          setQuery]          = useState('');
  const [editItem,       setEditItem]       = useState(null); // { item, type }
  const [auditItem,      setAuditItem]      = useState(null); // { id, type, hostname }
  const [auditLog,       setAuditLog]       = useState([]);
  const [auditLoading,   setAuditLoading]   = useState(false);
  const [sortCol, setSortCol] = useState('hostname');
  const [sortAsc, setSortAsc] = useState(true);
  const [desktopPage, setDesktopPage] = useState(1);
  const [laptopPage, setLaptopPage] = useState(1);
  const [desktopPageSize, setDesktopPageSize] = useState(20);
  const [laptopPageSize, setLaptopPageSize] = useState(20);

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

        const cpuMap      = Object.fromEntries((cpus ?? []).map(c => [c.cpu_id, c.model_name]));
        const cpuScoreMap = Object.fromEntries((cpus ?? []).map(c => [c.cpu_id, c.benchmark_score]));
        const osMap       = Object.fromEntries((osList ?? []).map(o => [o.os_id,      o.name]));
        const equipMap    = Object.fromEntries((equip ?? []).map(e => [e.equipment_user_id, e.name]));
        const roomMap     = Object.fromEntries(allRooms.map(r => [r.room_id, `${r.centerName} › ${r.name}`]));
        const lmMap       = Object.fromEntries((lm ?? []).map(m => [m.laptop_model_id,  `${m.brand_name} ${m.model_name}`]));
        const dmMap       = Object.fromEntries((dm ?? []).map(m => [m.desktop_model_id, `${m.brand_name} ${m.model_name}`]));

        const cpuOpts   = (cpus   ?? []).map(c => ({ value: c.cpu_id, label: c.benchmark_score ? `${c.model_name} (${c.benchmark_score})` : c.model_name }));
        const osOpts    = (osList ?? []).map(o => ({ value: o.os_id,               label: o.name }));
        const equipOpts = (equip  ?? []).map(e => ({ value: e.equipment_user_id,   label: e.name }));
        const roomOpts  = allRooms.map(r => ({ value: r.room_id, label: `${r.centerName} › ${r.name}` }));
        const lmOpts    = (lm     ?? []).map(m => ({ value: m.laptop_model_id,     label: `${m.brand_name} ${m.model_name}` }));
        const dmOpts    = (dm     ?? []).map(m => ({ value: m.desktop_model_id,    label: `${m.brand_name} ${m.model_name}` }));

        setRefs({ cpuMap, cpuScoreMap, osMap, equipMap, roomMap, lmMap, dmMap, cpuOpts, osOpts, equipOpts, roomOpts, lmOpts, dmOpts });
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

  useEffect(() => {
    if (!auditItem) { setAuditLog([]); return; }
    setAuditLoading(true);
    setAuditLog([]);
    api.getAuditLog(auditItem.type, auditItem.id)
      .then(entries => setAuditLog(entries ?? []))
      .catch(() => setAuditLog([]))
      .finally(() => setAuditLoading(false));
  }, [auditItem]);

  if (loading) return (
    <div className="empty" style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 12 }}>
      <div className="spinner" style={{ width: 40, height: 40, border: '4px solid #eee', borderTop: '4px solid #007bff', borderRadius: '50%', animation: 'spin 1s linear infinite' }}></div>
      <span>Carregant…</span>
      <style>{`
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
  if (err)     return <div className="empty" style={{ color: 'var(--danger)' }}>Error: {err}</div>;

  const R = refs;

  const q = query.trim().toLowerCase();

  // Búsqueda texto sobremesas
  const filteredDesktops = !q ? desktops : desktops.filter(d => [
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
      ].filter(Boolean).join(' ').toLowerCase().includes(q));

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
  // Paginado independiente
  const totalDesktopPages = Math.ceil(sortedDesktops.length / desktopPageSize) || 1;
  const pagedDesktops = sortedDesktops.slice((desktopPage - 1) * desktopPageSize, desktopPage * desktopPageSize);


  // Laptops: filtrado, orden y paginación
  // Búsqueda texto portátiles
  const filteredLaptops = !q ? laptops : laptops.filter(l => [
        l.hostname,
        R.roomMap[l.room_id],
        l.laptop_model_id ? R.lmMap[l.laptop_model_id] : null,
        l.cpu_model_name,
        R.osMap[l.os_id],
        l.ram_gb != null ? `${l.ram_gb} GB` : null,
        l.ram_type,
        l.storage_gb != null ? `${l.storage_gb} GB` : null,
        l.storage_type,
        l.mac_address,
        l.serial_number,
        l.equipment_user_id ? R.equipMap[l.equipment_user_id] : null,
        l.observations,
      ].filter(Boolean).join(' ').toLowerCase().includes(q));

  const laptopSortKeyMap = {
    hostname: l => l.hostname?.toLowerCase() ?? '',
    aula: l => R.roomMap[l.room_id]?.toLowerCase() ?? '',
    model: l => l.laptop_model_id ? R.lmMap[l.laptop_model_id]?.toLowerCase() : '',
    cpu: l => l.cpu_benchmark_score ?? -1,
    so: l => R.osMap[l.os_id]?.toLowerCase() ?? '',
    ram: l => l.ram_gb ?? 0,
    storage: l => l.storage_gb ?? 0,
    mac: l => l.mac_address?.toLowerCase() ?? '',
    user: l => l.equipment_user_id ? R.equipMap[l.equipment_user_id]?.toLowerCase() : '',
    obs: l => l.observations?.toLowerCase() ?? '',
  };
  const sortedLaptops = sortBy(filteredLaptops, laptopSortKeyMap[sortCol] || laptopSortKeyMap.hostname, sortAsc);
  const totalLaptopPages = Math.ceil(sortedLaptops.length / laptopPageSize) || 1;
  const pagedLaptops = sortedLaptops.slice((laptopPage - 1) * laptopPageSize, laptopPage * laptopPageSize);

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

      {/* STATS */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(150px, 1fr))', gap: 12, marginBottom: 24 }}>
        {[
          { icon: '🖥️', value: desktops.length,                   label: 'Sobretaules',           color: '#3b82f6' },
          { icon: '💻', value: laptops.length,                    label: 'Portàtils',              color: '#8b5cf6' },
          { icon: '📦', value: desktops.length + laptops.length,  label: 'Total equips',           color: '#10b981' },
        ].map(({ icon, value, label, color }) => (
          <div key={label} style={{ background: '#fff', border: '1px solid var(--border)', borderRadius: 8, padding: '16px 20px', boxShadow: 'var(--shadow)', borderTop: `3px solid ${color}` }}>
            <div style={{ fontSize: 20, marginBottom: 6 }}>{icon}</div>
            <div style={{ fontSize: 30, fontWeight: 700, color, lineHeight: 1 }}>{value}</div>
            <div style={{ fontSize: 12, color: 'var(--muted)', marginTop: 6 }}>{label}</div>
          </div>
        ))}
      </div>

      {/* DESKTOPS */}
      <style>{`
        .filter-select {
          font-size: 15px;
          border-radius: 6px;
          border: 1px solid var(--border, #d0d7de);
          padding: 8px 18px;
          background: #fff;
          min-width: 140px;
          max-width: 180px;
          box-shadow: 0 1px 2px #0001;
          transition: border 0.2s;
        }
        .filter-select:focus {
          outline: 2px solid #007bff33;
          border-color: #007bff;
        }
      `}</style>
      <div className="page-header" style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 }}>
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
                      <td style={{ fontSize: 12 }}>{d.cpu_id
                          ? <>{R.cpuMap[d.cpu_id] ?? '—'}{R.cpuScoreMap[d.cpu_id] != null && <span style={{ color: 'var(--muted)' }}> ({R.cpuScoreMap[d.cpu_id]})</span>}</>
                          : '—'}</td>
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
                            {isAdmin && (
                              <button className="btn btn-sm" style={{ color: 'var(--muted)' }} onClick={() => setAuditItem({ id: d.computer_id, type: 'desktop', hostname: d.hostname })}>
                                Historial
                              </button>
                            )}
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
            {/* Paginación sobremesas */}
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 12, marginTop: 12 }}>
              <button className="btn btn-sm" disabled={desktopPage === 1} onClick={() => setDesktopPage(1)}>&lt;&lt;</button>
              <button className="btn btn-sm" disabled={desktopPage === 1} onClick={() => setDesktopPage(desktopPage - 1)}>Anterior</button>
              <span>Pàgina {desktopPage} de {totalDesktopPages}</span>
              <button className="btn btn-sm" disabled={desktopPage === totalDesktopPages} onClick={() => setDesktopPage(desktopPage + 1)}>Següent</button>
              <button className="btn btn-sm" disabled={desktopPage === totalDesktopPages} onClick={() => setDesktopPage(totalDesktopPages)}>&gt;&gt;</button>
              <label htmlFor="desktop-page-size-select" style={{ marginLeft: 16 }}>Mostrar:</label>
              <select
                id="desktop-page-size-select"
                value={desktopPageSize}
                onChange={e => {
                  setDesktopPageSize(Number(e.target.value));
                  setDesktopPage(1);
                }}
                className="filter-select"
                style={{ minWidth: 70 }}
              >
                {[10, 20, 50, 100].map(n => (
                  <option key={n} value={n}>{n}</option>
                ))}
              </select>
            </div>
          </>
        )}
      </div>

      {/* LAPTOPS */}
      <div className="page-header" style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 }}>
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
                      <td style={{ fontSize: 12 }}>
                        {l.cpu_model_name
                          ? <>{l.cpu_model_name}{l.cpu_benchmark_score != null && <span style={{ color: 'var(--muted)' }}> ({l.cpu_benchmark_score})</span>}</>
                          : '—'}
                      </td>
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
                            {isAdmin && (
                              <button className="btn btn-sm" style={{ color: 'var(--muted)' }} onClick={() => setAuditItem({ id: l.computer_id, type: 'laptop', hostname: l.hostname })}>
                                Historial
                              </button>
                            )}
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
            {/* Paginación portátiles */}
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 12, marginTop: 12 }}>
              <button className="btn btn-sm" disabled={laptopPage === 1} onClick={() => setLaptopPage(1)}>&lt;&lt;</button>
              <button className="btn btn-sm" disabled={laptopPage === 1} onClick={() => setLaptopPage(laptopPage - 1)}>Anterior</button>
              <span>Pàgina {laptopPage} de {totalLaptopPages}</span>
              <button className="btn btn-sm" disabled={laptopPage === totalLaptopPages} onClick={() => setLaptopPage(laptopPage + 1)}>Següent</button>
              <button className="btn btn-sm" disabled={laptopPage === totalLaptopPages} onClick={() => setLaptopPage(totalLaptopPages)}>&gt;&gt;</button>
              <label htmlFor="laptop-page-size-select" style={{ marginLeft: 16 }}>Mostrar:</label>
              <select
                id="laptop-page-size-select"
                value={laptopPageSize}
                onChange={e => {
                  setLaptopPageSize(Number(e.target.value));
                  setLaptopPage(1);
                }}
                className="filter-select"
                style={{ minWidth: 70 }}
              >
                {[10, 20, 50, 100].map(n => (
                  <option key={n} value={n}>{n}</option>
                ))}
              </select>
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

      {/* AUDIT MODAL */}
      {auditItem && (
        <AuditModal
          item={auditItem}
          log={auditLog}
          loading={auditLoading}
          onClose={() => setAuditItem(null)}
        />
      )}
    </>
  );
}

// ─── Audit helpers ────────────────────────────────────────────────────────────

function parseAuditJson(val) {
  if (!val) return null;
  if (typeof val === 'object') return val;
  try { return JSON.parse(val); } catch {}
  try { return JSON.parse(atob(val)); } catch {}
  return null;
}

const FIELD_LABELS = {
  hostname: 'Hostname', ram_gb: 'RAM (GB)', ram_type: 'Tipus RAM',
  storage_gb: 'Emmagatzematge (GB)', storage_type: 'Tipus emmagatzematge',
  has_wifi_card: 'Té WiFi', mac_address: 'MAC', observations: 'Observacions',
  cpu_id: 'CPU', os_id: 'SO', room_id: 'Sala', equipment_user_id: 'Usuari equip',
  desktop_model_id: 'Model sobretaula', laptop_model_id: 'Model portàtil',
  serial_number: 'Núm. sèrie',
};

const EVENT_META = {
  created: { label: 'Creat',     bg: '#dcfce7', color: '#16a34a' },
  updated: { label: 'Modificat', bg: '#dbeafe', color: '#2563eb' },
  deleted: { label: 'Eliminat',  bg: '#fee2e2', color: '#dc2626' },
};

function AuditDiff({ eventType, oldValues, newValues }) {
  const oldObj = parseAuditJson(oldValues) ?? {};
  const newObj = parseAuditJson(newValues) ?? {};

  if (eventType === 'created') {
    const rows = Object.entries(newObj).filter(([, v]) => v != null && v !== '');
    return (
      <table style={{ width: '100%', fontSize: 12, borderCollapse: 'collapse' }}>
        <thead><tr>
          <th style={{ textAlign: 'left', paddingBottom: 4, color: 'var(--muted)', fontWeight: 500, width: '40%' }}>Camp</th>
          <th style={{ textAlign: 'left', paddingBottom: 4, color: '#16a34a', fontWeight: 500 }}>Valor</th>
        </tr></thead>
        <tbody>
          {rows.map(([k, v]) => (
            <tr key={k}>
              <td style={{ paddingRight: 16, color: 'var(--muted)', paddingBottom: 2 }}>{FIELD_LABELS[k] || k}</td>
              <td style={{ color: '#16a34a' }}>{String(v)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  }

  if (eventType === 'deleted') {
    const rows = Object.entries(oldObj).filter(([, v]) => v != null && v !== '');
    return (
      <table style={{ width: '100%', fontSize: 12, borderCollapse: 'collapse' }}>
        <thead><tr>
          <th style={{ textAlign: 'left', paddingBottom: 4, color: 'var(--muted)', fontWeight: 500, width: '40%' }}>Camp</th>
          <th style={{ textAlign: 'left', paddingBottom: 4, color: '#dc2626', fontWeight: 500 }}>Valor eliminat</th>
        </tr></thead>
        <tbody>
          {rows.map(([k, v]) => (
            <tr key={k}>
              <td style={{ paddingRight: 16, color: 'var(--muted)', paddingBottom: 2 }}>{FIELD_LABELS[k] || k}</td>
              <td style={{ color: '#dc2626', textDecoration: 'line-through' }}>{String(v)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  }

  // updated — keys in new_values that actually differ from old_values
  const changed = Object.keys(newObj).filter(k => JSON.stringify(oldObj[k]) !== JSON.stringify(newObj[k]));
  if (changed.length === 0) {
    return <span style={{ fontSize: 12, color: 'var(--muted)' }}>Sense canvis de camp detectats.</span>;
  }
  return (
    <table style={{ width: '100%', fontSize: 12, borderCollapse: 'collapse' }}>
      <thead><tr>
        <th style={{ textAlign: 'left', paddingBottom: 4, color: 'var(--muted)', fontWeight: 500, width: '30%' }}>Camp</th>
        <th style={{ textAlign: 'left', paddingBottom: 4, color: '#dc2626', fontWeight: 500, width: '35%' }}>Anterior</th>
        <th style={{ textAlign: 'left', paddingBottom: 4, color: '#16a34a', fontWeight: 500 }}>Nou</th>
      </tr></thead>
      <tbody>
        {changed.map(k => (
          <tr key={k}>
            <td style={{ paddingRight: 16, color: 'var(--muted)', paddingBottom: 2 }}>{FIELD_LABELS[k] || k}</td>
            <td style={{ paddingRight: 16, color: '#dc2626', textDecoration: 'line-through' }}>{String(oldObj[k] ?? '—')}</td>
            <td style={{ color: '#16a34a' }}>{String(newObj[k] ?? '—')}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function AuditModal({ item, log, loading, onClose }) {
  return (
    <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,.45)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center', padding: 16 }}>
      <div style={{ background: '#fff', borderRadius: 10, width: '100%', maxWidth: 680, maxHeight: '80vh', display: 'flex', flexDirection: 'column', boxShadow: '0 20px 60px rgba(0,0,0,.25)' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '16px 20px', borderBottom: '1px solid var(--border)' }}>
          <div>
            <div style={{ fontWeight: 700, fontSize: 16 }}>Historial de canvis</div>
            <div style={{ fontSize: 13, color: 'var(--muted)' }}>
              {item.hostname} · {item.type === 'desktop' ? 'Sobretaula' : 'Portàtil'}
            </div>
          </div>
          <button className="btn btn-ghost btn-sm" onClick={onClose}>✕ Tancar</button>
        </div>
        <div style={{ flex: 1, overflowY: 'auto', padding: '16px 20px' }}>
          {loading ? (
            <div style={{ textAlign: 'center', color: 'var(--muted)', padding: 32 }}>Carregant historial…</div>
          ) : log.length === 0 ? (
            <div style={{ textAlign: 'center', color: 'var(--muted)', padding: 32 }}>
              Sense registres d'auditoria per a aquest equip.
            </div>
          ) : (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
              {log.map((entry, i) => {
                const meta = EVENT_META[entry.event_type] ?? { label: entry.event_type, bg: '#f3f4f6', color: 'var(--muted)' };
                const date = new Date(entry.changed_at);
                const dateStr = isNaN(date) ? entry.changed_at : date.toLocaleString('ca-ES', { dateStyle: 'short', timeStyle: 'short' });
                return (
                  <div key={i} style={{ border: '1px solid var(--border)', borderRadius: 8, overflow: 'hidden' }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 10, padding: '10px 14px', background: '#f9fafb' }}>
                      <span style={{ padding: '2px 8px', borderRadius: 99, fontSize: 11, fontWeight: 600, background: meta.bg, color: meta.color }}>
                        {meta.label}
                      </span>
                      <span style={{ fontSize: 12, color: 'var(--muted)' }}>{dateStr}</span>
                      <span style={{ fontSize: 12, marginLeft: 'auto' }}>👤 {entry.changed_by_username}</span>
                    </div>
                    <div style={{ padding: '10px 14px' }}>
                      <AuditDiff eventType={entry.event_type} oldValues={entry.old_values} newValues={entry.new_values} />
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
