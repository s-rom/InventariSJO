import { useState, useEffect, useCallback } from 'react';
import { api } from '../api';

export default function Users() {
  const [users,   setUsers]   = useState([]);
  const [roles,   setRoles]   = useState([]);
  const [err,     setErr]     = useState('');
  const [saving,  setSaving]  = useState(false);

  // Create form state
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [roleId,   setRoleId]   = useState('');

  // Inline role-change state: { [userId]: pendingRoleId }
  const [pendingRole, setPendingRole] = useState({});

  // Confirm-delete state
  const [delId, setDelId] = useState(null);

  const load = useCallback(async () => {
    try {
      const [u, r] = await Promise.all([api.listUsers(), api.listRoles()]);
      setUsers(u ?? []);
      setRoles(r ?? []);
      if (!roleId && r?.length) setRoleId(r[0].role_id);
    } catch (ex) {
      setErr(ex.message);
    }
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => { load(); }, [load]);

  async function handleCreate(e) {
    e.preventDefault();
    setErr(''); setSaving(true);
    try {
      await api.createUser({ username, password, role_id: roleId });
      setUsername(''); setPassword('');
      load();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function handleRoleChange(userId) {
    const newRole = pendingRole[userId];
    if (!newRole) return;
    setErr('');
    try {
      await api.updateUser(userId, { role_id: newRole });
      setPendingRole(p => { const n = { ...p }; delete n[userId]; return n; });
      load();
    } catch (ex) { setErr(ex.message); }
  }

  async function handleDelete(id) {
    setErr('');
    try {
      await api.deleteUser(id);
      setDelId(null);
      load();
    } catch (ex) { setErr(ex.message); }
  }



  return (
    <div style={{ padding: '24px 32px', maxWidth: 780 }}>
      <h1 style={{ fontSize: 20, fontWeight: 700, marginBottom: 24 }}>Gestió d'usuaris</h1>

      {/* ── Create user form ─────────────────────────────────────── */}
      <div className="card" style={{ marginBottom: 24 }}>
        <h2 style={{ fontSize: 15, fontWeight: 600, marginBottom: 14 }}>Nou usuari</h2>
        <form onSubmit={handleCreate} className="form-panel">
          <div className="form-grid" style={{ gridTemplateColumns: '1fr 1fr 1fr auto' }}>
            <div className="form-group">
              <label>Usuari *</label>
              <input
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                placeholder="nom.cognom"
                required
              />
            </div>
            <div className="form-group">
              <label>Contrasenya *</label>
              <input
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label>Rol *</label>
              <select value={roleId} onChange={e => setRoleId(e.target.value)} required>
                {roles.map(r => (
                  <option key={r.role_id} value={r.role_id}>{r.role_id}</option>
                ))}
              </select>
            </div>
            <div className="form-group" style={{ justifyContent: 'flex-end' }}>
              <label style={{ visibility: 'hidden' }}>_</label>
              <button type="submit" className="btn btn-primary" disabled={saving || !username || !password || !roleId}>
                {saving ? '…' : 'Crear'}
              </button>
            </div>
          </div>
          {err && <div className="error-msg">{err}</div>}
        </form>
      </div>

      {/* ── Users table ──────────────────────────────────────────── */}
      <div className="card">
        <div className="table-wrap">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Usuari</th>
                <th>Rol actual</th>
                <th>Canviar rol</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {users.length === 0 && (
                <tr>
                  <td colSpan={5} style={{ textAlign: 'center', color: 'var(--muted)', padding: 16 }}>
                    Sense usuaris
                  </td>
                </tr>
              )}
              {users.map(u => {
                const pending = pendingRole[u.app_user_id] ?? u.role_id;
                const changed = pending !== u.role_id;
                return (
                  <tr key={u.app_user_id}>
                    <td style={{ color: 'var(--muted)', width: 48 }}>{u.app_user_id}</td>
                    <td>{u.username}</td>
                    <td>
                      <span className="badge">{u.role_id}</span>
                    </td>
                    <td style={{ width: 220 }}>
                      <div style={{ display: 'flex', gap: 6, alignItems: 'center' }}>
                        <select
                          value={pending}
                          onChange={e => setPendingRole(p => ({ ...p, [u.app_user_id]: e.target.value }))}
                          style={{ flex: 1 }}
                        >
                          {roles.map(r => (
                            <option key={r.role_id} value={r.role_id}>{r.role_id}</option>
                          ))}
                        </select>
                        {changed && (
                          <button
                            className="btn btn-sm btn-primary"
                            onClick={() => handleRoleChange(u.app_user_id)}
                          >
                            Desar
                          </button>
                        )}
                      </div>
                    </td>
                    <td style={{ width: 90, textAlign: 'right' }}>
                      {delId === u.app_user_id ? (
                        <span style={{ display: 'flex', gap: 4, justifyContent: 'flex-end' }}>
                          <button className="btn btn-sm btn-danger" onClick={() => handleDelete(u.app_user_id)}>Eliminar</button>
                          <button className="btn btn-sm btn-ghost" onClick={() => setDelId(null)}>Cancel·lar</button>
                        </span>
                      ) : (
                        <button className="btn btn-sm btn-ghost" onClick={() => setDelId(u.app_user_id)}>🗑</button>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>

    </div>
  );
}
