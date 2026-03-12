import { NavLink } from 'react-router-dom';

const NAV = [
  { to: '/computers',             label: '🖥️  Equips' },
  { to: '/computers/new-desktop', label: '➕ Nou Sobretaula' },
  { to: '/computers/new-laptop',  label: '➕ Nou Portàtil' },
  { to: '/students',              label: '👨‍🎓 Alumnes' },
  { to: '/reference',             label: '📋 Dades bàsiques' },
];

export default function Layout({ children, onLogout }) {
  return (
    <div className="layout">
      <aside className="sidebar">
        <div className="sidebar-brand">Inventari</div>
        <nav style={{ flex: 1 }}>
          {NAV.map(({ to, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/computers'}
              className={({ isActive }) => `nav-item${isActive ? ' active' : ''}`}
              style={{ display: 'block', textDecoration: 'none' }}
            >
              {label}
            </NavLink>
          ))}
        </nav>
        <div className="sidebar-bottom">
          <button
            className="btn btn-ghost btn-sm"
            style={{ width: '100%', justifyContent: 'center' }}
            onClick={onLogout}
          >
            Tancar sessió
          </button>
        </div>
      </aside>

      <main className="main">{children}</main>
    </div>
  );
}
