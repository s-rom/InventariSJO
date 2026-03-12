import { api } from '../api';

/**
 * Props:
 *   computers     : Computer[]
 *   cpuMap        : { [id]: string }
 *   roomMap       : { [id]: { name, centerName } }
 *   equipUserMap  : { [id]: string }
 *   onDelete      : (id) => void
 *   onAddClick    : () => void
 *   showForm      : bool
 */
export default function ComputerList({ computers, cpuMap, roomMap, equipUserMap, onDelete, onAddClick, showForm }) {

  function formatRam(gb, type) {
    if (!gb && gb !== 0) return '—';
    return `${gb} GB ${type !== 'None' ? type : ''}`.trim();
  }

  function formatStorage(gb, type) {
    if (!gb && gb !== 0) return '—';
    return `${gb} GB ${type !== 'None' ? type : ''}`.trim();
  }

  async function handleDelete(id) {
    if (!confirm('Eliminar aquest equip?')) return;
    try {
      await api.deleteComputer(id);
      onDelete(id);
    } catch (err) {
      alert(err.message);
    }
  }

  return (
    <div className="card">
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', padding: '14px 16px', borderBottom: '1px solid var(--border)' }}>
        <span style={{ fontWeight: 600 }}>Equips ({computers.length})</span>
        <button className="btn btn-primary btn-sm" onClick={onAddClick}>
          {showForm ? 'Tancar formulari' : '+ Afegir equip'}
        </button>
      </div>

      {computers.length === 0 ? (
        <div className="empty">No hi ha equips registrats.</div>
      ) : (
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>Hostname</th>
                <th>Tipus</th>
                <th>CPU</th>
                <th>RAM</th>
                <th>Emmagatzematge</th>
                <th>Sala</th>
                <th>Usuari equip</th>
                <th>SO</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {computers.map(c => {
                const cpu      = c.cpu_id ? (cpuMap[c.cpu_id] ?? `#${c.cpu_id}`) : '—';
                const room     = c.room_id ? (roomMap[c.room_id] ? `${roomMap[c.room_id].centerName} › ${roomMap[c.room_id].name}` : `#${c.room_id}`) : '—';
                const eqUser   = c.equipment_user_id ? (equipUserMap[c.equipment_user_id] ?? `#${c.equipment_user_id}`) : '—';

                return (
                  <tr key={c.computer_id}>
                    <td style={{ fontWeight: 500 }}>{c.hostname}</td>
                    <td><span className="badge gray">{c.computer_type}</span></td>
                    <td style={{ fontSize: 12, color: 'var(--muted)' }}>{cpu}</td>
                    <td style={{ fontSize: 12 }}>{formatRam(c.ram_gb, c.ram_type)}</td>
                    <td style={{ fontSize: 12 }}>{formatStorage(c.storage_gb, c.storage_type)}</td>
                    <td style={{ fontSize: 12, color: 'var(--muted)' }}>{room}</td>
                    <td style={{ fontSize: 12 }}>{eqUser}</td>
                    <td>
                      {c.operating_systems?.map(os => (
                        <span key={os.os_id} className="badge">{os.name}</span>
                      ))}
                    </td>
                    <td>
                      <button className="btn btn-danger btn-sm" onClick={() => handleDelete(c.computer_id)}>
                        Eliminar
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
