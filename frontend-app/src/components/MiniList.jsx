/**
 * MiniList — compact preview table reused by NewDesktop and NewLaptop.
 *
 * Props:
 *   items  : desktop[] | laptop[]
 *   type   : 'desktop' | 'laptop'
 *   refs   : { cpuMap, osMap, equipMap, roomMap, lmMap, dmMap }
 */

function osLabel(name) {
  if (!name) return '—';
  const n = name.toLowerCase();
  if (n.includes('windows') || n.includes('win')) return '🪟 ' + name;
  if (n.includes('linux'))  return '🐧 ' + name;
  if (n.includes('chrome')) return '🟢 ' + name;
  return name;
}

const MAX = 10;

export default function MiniList({ items, type, refs }) {
  const R = refs ?? {};
  const recent = [...(items ?? [])]
    .sort((a, b) => new Date(b.created_at ?? 0) - new Date(a.created_at ?? 0))
    .slice(0, MAX);

  if (!items || items.length === 0) {
    return (
      <div className="empty" style={{ padding: '20px 0' }}>
        Cap {type === 'desktop' ? 'sobretaula' : 'portàtil'} registrat encara.
      </div>
    );
  }

  return (
    <div className="table-wrap">
      {items.length > MAX && (
        <div style={{ padding: '6px 12px', fontSize: 12, color: 'var(--muted)', borderBottom: '1px solid var(--border)' }}>
          Mostrant els últims {MAX} de {items.length}
        </div>
      )}
      <table>
        <thead>
          <tr>
            <th>Hostname</th>
            {type === 'laptop' && <th>Núm. sèrie</th>}
            <th>Aula</th>
            <th>Model</th>
            {type === 'desktop' && <th>CPU</th>}
            <th>SO</th>
            <th>RAM</th>
            <th>Emmagatzematge</th>
            {type === 'desktop' && <th>WiFi</th>}
            <th>Observacions</th>
          </tr>
        </thead>
        <tbody>
          {recent.map(item => (
            <tr key={item.computer_id}>
              <td><strong>{item.hostname}</strong></td>
              {type === 'laptop' && (
                <td><code style={{ fontSize: 11 }}>{item.serial_number ?? '—'}</code></td>
              )}
              <td>{R.roomMap?.[item.room_id] ?? '—'}</td>
              <td>
                {type === 'desktop'
                  ? (item.desktop_model_id ? (R.dmMap?.[item.desktop_model_id] ?? '—') : '—')
                  : (item.laptop_model_id  ? (R.lmMap?.[item.laptop_model_id]  ?? '—') : '—')}
              </td>
              {type === 'desktop' && (
                <td>{item.cpu_id ? (R.cpuMap?.[item.cpu_id] ?? '—') : '—'}</td>
              )}
              <td>{osLabel(R.osMap?.[item.os_id])}</td>
              <td>
                {item.ram_gb != null
                  ? `${item.ram_gb} GB${item.ram_type ? ' ' + item.ram_type : ''}`
                  : '—'}
              </td>
              <td>
                {item.storage_gb != null
                  ? `${item.storage_gb} GB${item.storage_type ? ' ' + item.storage_type : ''}`
                  : '—'}
              </td>
              {type === 'desktop' && (
                <td style={{ textAlign: 'center' }}>{item.has_wifi_card ? '✔' : '—'}</td>
              )}
              <td style={{ fontSize: 12, color: 'var(--muted)', maxWidth: 180, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                {item.observations ?? '—'}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
