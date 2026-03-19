import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api';
import { useAuth } from '../App';

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

        setRefs({ cpuMap, osMap, equipMap, roomMap, lmMap, dmMap });
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

  return (
    <>
      {/* DESKTOPS */}
      <div className="page-header">
        <h1 className="page-title">🖥️ Sobretaules</h1>
        <button className="btn btn-primary" onClick={() => navigate('/computers/new-desktop')}>
          + Nou sobretaula
        </button>
      </div>

      <div className="card" style={{ marginBottom: 32 }}>
        {desktops.length === 0 ? (
          <div className="empty">Cap sobretaula registrat.</div>
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
                  {canEdit && <th></th>}
                </tr>
              </thead>
              <tbody>
                {desktops.map(d => (
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
                    {canEdit && (
                      <td>
                        <button className="btn btn-danger btn-sm" onClick={() => handleDelete(d.computer_id, d.hostname)}>
                          Eliminar
                        </button>
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
        {laptops.length === 0 ? (
          <div className="empty">Cap portàtil registrat.</div>
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
                  {canEdit && <th></th>}
                </tr>
              </thead>
              <tbody>
                {laptops.map(l => (
                  <tr key={l.computer_id}>
                    <td><strong>{l.hostname}</strong></td>
                    <td>{R.roomMap[l.room_id] ?? '—'}</td>
                    <td>{l.laptop_model_id ? R.lmMap[l.laptop_model_id] ?? '—' : '—'}</td>
                    <td>{osEmoji(R.osMap[l.os_id])}</td>
                    <td>{ramLabel(l.ram_gb, l.ram_type)}</td>
                    <td>{storageLabel(l.storage_gb, l.storage_type)}</td>
                    <td><code style={{ fontSize: 11 }}>{l.mac_address ?? '—'}</code></td>
                    <td>{l.equipment_user_id ? R.equipMap[l.equipment_user_id] ?? '—' : '—'}</td>
                    {canEdit && (
                      <td>
                        <button className="btn btn-danger btn-sm" onClick={() => handleDelete(l.computer_id, l.hostname)}>
                          Eliminar
                        </button>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </>
  );
}
