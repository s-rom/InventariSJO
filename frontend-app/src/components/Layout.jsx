import { useState } from 'react';
import { NavLink } from 'react-router-dom';

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

export default function Layout({ children, onLogout, role }) {
  const nav = role === 'admin' ? [...NAV, ...ADMIN_NAV] : NAV;
  // sidebarMode: 'full' | 'icons' | 'none'
  const [sidebarMode, setSidebarMode] = useState('full');

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
    </div>
  );
}

