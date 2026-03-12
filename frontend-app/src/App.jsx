import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { api, getToken } from './api';
import Login    from './components/Login';
import Layout   from './components/Layout';
import Computers   from './pages/Computers';
import NewDesktop  from './pages/NewDesktop';
import NewLaptop   from './pages/NewLaptop';
import Students    from './pages/Students';
import Reference   from './pages/Reference';

function RequireAuth({ children }) {
  if (!getToken()) return <Navigate to="/login" replace />;
  return children;
}

export default function App() {
  const [authed, setAuthed] = useState(null); // null=checking, false=login, true=in

  // Quick auth probe on mount
  useEffect(() => {
    if (!getToken()) { setAuthed(false); return; }
    api.listDesktops()
      .then(() => setAuthed(true))
      .catch(() => { api.logout(); setAuthed(false); });
  }, []);

  if (authed === null) {
    return (
      <div style={{ display:'flex', alignItems:'center', justifyContent:'center', height:'100vh', color:'var(--muted)' }}>
        Carregant…
      </div>
    );
  }

  function handleLogin() { setAuthed(true); }
  async function handleLogout() { await api.logout(); setAuthed(false); }

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={authed ? <Navigate to="/computers" replace /> : <Login onLogin={handleLogin} />}
        />
        <Route
          path="/*"
          element={
            <RequireAuth>
              <Layout onLogout={handleLogout}>
                <Routes>
                  <Route index element={<Navigate to="/computers" replace />} />
                  <Route path="computers"            element={<Computers />} />
                  <Route path="computers/new-desktop" element={<NewDesktop />} />
                  <Route path="computers/new-laptop"  element={<NewLaptop />} />
                  <Route path="students"             element={<Students />} />
                  <Route path="reference"            element={<Reference />} />
                  <Route path="*"                    element={<Navigate to="/computers" replace />} />
                </Routes>
              </Layout>
            </RequireAuth>
          }
        />
        <Route path="/" element={<Navigate to={authed ? '/computers' : '/login'} replace />} />
      </Routes>
    </BrowserRouter>
  );
}

