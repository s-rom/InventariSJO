import { useState } from 'react';
import { NavLink } from 'react-router-dom';
import { api } from '../api';

const NAV = [
  { to: '/computers',             label: '🖥️  Equips' },
  { to: '/computers/new-desktop', label: '➕ Nou Sobretaula' },
  { to: '/computers/new-laptop',  label: '➕ Nou Portàtil' },
  { to: '/students',              label: '👨‍🎓 Alumnes' },
  { to: '/reference',             label: '📋 Dades bàsiques' },
];

const ADMIN_NAV = [
  { to: '/admin/users', label: '👥 Usuaris' },
];

const ROLE_LABEL = { admin: 'Admin', editor: 'Editor', tutor: 'Tutor' };

function ChangePasswordModal({ onClose }) {
  const [current, setCurrent] = useState('');
  const [next,    setNext]    = useState('');
  const [confirm, setConfirm] = useState('');
  const [saving,  setSaving]  = useState(false);
  const [err,     setErr]     = useState('');
  const [ok,      setOk]      = useState(false);

  async function handleSubmit(e) {
    e.preventDefault();
    setErr('');
    if (next !== confirm) { setErr('Les contrasenyes no coincideixen'); return; }
    setSaving(true);
    try {
      await api.changePassword({ current_password: current, new_password: next });
      setOk(true);
    } catch (ex) {
      setErr(ex.message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <div
      style={{
        position: 'fixed', inset: 0, zIndex: 1000,
        background: 'rgba(0,0,0,0.35)',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
      }}
      onClick={e => { if (e.target === e.currentTarget) onClose(); }}
    >
      <div className="card" style={{ width: 340, padding: 24, position: 'relative' }}>
        <button
          onClick={onClose}
          style={{ position: 'absolute', top: 12, right: 14, background: 'none', border: 'none', fontSize: 18, cursor: 'pointer', color: 'var(--muted)' }}
        >×</button>
        <div style={{ fontSize: 15, fontWeight: 600, marginBottom: 18 }}>Canviar contrasenya</div>
        {ok ? (
          <div style={{ textAlign: 'center', padding: '12px 0' }}>
            <div style={{ fontSize: 28, marginBottom: 8 }}>✅</div>
            <div style={{ fontSize: 13, color: 'var(--muted)' }}>Contrasenya actualitzada correctament.</div>
            <button className="btn btn-primary" style={{ marginTop: 16, width: '100%' }} onClick={onClose}>Tancar</button>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="form-panel" style={{ padding: 0 }}>
            <div className="form-group" style={{ marginBottom: 12 }}>
              <label>Contrasenya actual *</label>
              <input type="password" value={current} onChange={e => setCurrent(e.target.value)} required autoFocus />
            </div>
            <div className="form-group" style={{ marginBottom: 12 }}>
              <label>Nova contrasenya *</label>
              <input type="password" value={next} onChange={e => setNext(e.target.value)} required minLength={8} />
            </div>
            <div className="form-group" style={{ marginBottom: 16 }}>
              <label>Confirmar nova contrasenya *</label>
              <input type="password" value={confirm} onChange={e => setConfirm(e.target.value)} required />
            </div>
            {err && <div className="error-msg" style={{ marginBottom: 10 }}>{err}</div>}
            <button type="submit" className="btn btn-primary" style={{ width: '100%' }} disabled={saving || !current || !next || !confirm}>
              {saving ? 'Guardant…' : 'Actualitzar'}
            </button>
          </form>
        )}
      </div>
    </div>
  );
}

export default function Layout({ children, onLogout, role, username }) {
  const nav = role === 'admin' ? [...NAV, ...ADMIN_NAV] : NAV;
  // sidebarMode: 'full' | 'icons' | 'none'
  const [sidebarMode,  setSidebarMode]  = useState('full');
  const [showChgPwd,   setShowChgPwd]   = useState(false);

  return (
    <div className="layout">

      <style>{`
        .sidebar {
          width: 240px;
          transition: width 0.3s ease;
          overflow: hidden;
        }

        .sidebar.sidebar-icons {
          width: 70px;
        }

        .sidebar-text {
          transition: opacity 0.2s ease;
          white-space: nowrap;
        }
      `}</style>

      {sidebarMode !== 'none' && (
        <aside className={`sidebar ${sidebarMode === 'icons' ? 'sidebar-icons' : ''}`}>
          <div className="sidebar-brand" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
            <span style={{ opacity: sidebarMode === 'icons' ? 0 : 1, transition: 'opacity .2s', width: sidebarMode === 'icons' ? 0 : 'auto', overflow: 'hidden' }}>Inventari</span>
            <div style={{ display: 'flex', gap: 4 }}>
              {sidebarMode === 'full' && (
                <button
                  className="sidebar-hide-btn"
                  title="Mostrar solo iconos"
                  onClick={() => setSidebarMode('icons')}
                  style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--muted)', fontSize: 18 }}
                >☰</button>
              )}
              {sidebarMode === 'icons' && (
                <>
                  <button
                    className="sidebar-hide-btn"
                    title="Ocultar barra lateral"
                    onClick={() => setSidebarMode('none')}
                    style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--muted)', fontSize: 18 }}
                  >⟨</button>
                  <button
                    className="sidebar-hide-btn"
                    title="Expandir barra lateral"
                    onClick={() => setSidebarMode('full')}
                    style={{ background: 'none', border: 'none', cursor: 'pointer', color: 'var(--muted)', fontSize: 18 }}
                  >⟩</button>
                </>
              )}
            </div>
          </div>
          <nav style={{ flex: 1 }}>
            {nav.map(({ to, label }) => (
              <NavLink
                key={to}
                to={to}
                end={to === '/computers'}
                className={({ isActive }) => `nav-item${isActive ? ' active' : ''}`}
                style={{ display: 'flex', alignItems: 'center', textDecoration: 'none', justifyContent: sidebarMode === 'icons' ? 'center' : 'flex-start' }}
              >
                <span style={{ fontSize: 20 }}>{label.split(' ')[0]}</span>
                {sidebarMode === 'full' && <span style={{ marginLeft: 10 }}>{label.slice(label.indexOf(' ')+1)}</span>}
              </NavLink>
            ))}
          </nav>
          {sidebarMode === 'full' && (
            <div className="sidebar-bottom">
              {username && (
                <div style={{ padding: '0 4px 10px', display: 'flex', flexDirection: 'column', gap: 2 }}>
                  <button
                    onClick={() => setShowChgPwd(true)}
                    style={{ background: 'none', border: 'none', cursor: 'pointer', padding: 0, textAlign: 'left' }}
                    title="Canviar contrasenya"
                  >
                    <span style={{ fontSize: 13, fontWeight: 600, color: 'var(--text)', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', display: 'block' }}>
                      {username} 🔒
                    </span>
                  </button>
                  <span style={{ fontSize: 11, color: 'var(--muted)' }}>
                    {ROLE_LABEL[role] ?? role}
                  </span>
                </div>
              )}
              <button
                className="btn btn-ghost btn-sm"
                style={{ width: '100%', justifyContent: 'center' }}
                onClick={onLogout}
              >
                Tancar sessió
              </button>
            </div>
          )}
        </aside>
      )}
      {sidebarMode === 'none' && (
        <button
          className="sidebar-show-btn"
          title="Mostrar barra lateral"
          onClick={() => setSidebarMode('icons')}
          style={{ position: 'absolute', left: 0, top: 16, zIndex: 20, background: '#fff', border: '1px solid #e0e3ea', borderRadius: '0 6px 6px 0', padding: '6px 10px', boxShadow: '0 1px 4px #0001', color: 'var(--muted)', cursor: 'pointer', fontSize: 18 }}
        >☰</button>
      )}
      <main className="main">{children}</main>
      {showChgPwd && <ChangePasswordModal onClose={() => setShowChgPwd(false)} />}
    </div>
  );
}

