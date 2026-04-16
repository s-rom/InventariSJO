import { useState, useEffect, createContext, useContext } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { api, getToken, getRole } from './api';
import Login    from './components/Login';
import Layout   from './components/Layout';
import Computers   from './pages/Computers';
import NewDesktop  from './pages/NewDesktop';
import NewLaptop   from './pages/NewLaptop';
import Students    from './pages/Students';
import Reference   from './pages/Reference';
import Users       from './pages/Users';

export const AuthContext = createContext({ role: null, username: null });
export function useAuth() { return useContext(AuthContext); }

function RequireAuth({ children }) {
  if (!getToken()) return <Navigate to="/login" replace />;
  return children;
}

function RequireAdmin({ role, children }) {
  if (role !== 'admin') return <Navigate to="/computers" replace />;
  return children;
}

function RequireNonTutor({ role, children }) {
  if (role === 'tutor') return <Navigate to="/students" replace />;
  return children;
}

export default function App() {
  const [authed,   setAuthed]   = useState(null); // null=checking, false=login, true=in
  const [role,     setRole]     = useState(getRole());
  const [username, setUsername] = useState(null);

  // Quick auth probe on mount
  useEffect(() => {
    if (!getToken()) { setAuthed(false); return; }
    api.me()
      .then((data) => { setRole(data.role_id); setUsername(data.username); setAuthed(true); })
      .catch(() => { api.logout(); setAuthed(false); });
  }, []);

  if (authed === null) {
    return (
      <div style={{ display:'flex', alignItems:'center', justifyContent:'center', height:'100vh', color:'var(--muted)' }}>
        Carregant…
      </div>
    );
  }

  function handleLogin(data) { setRole(data.role_id); setUsername(data.username); setAuthed(true); }
  async function handleLogout() { await api.logout(); setRole(null); setUsername(null); setAuthed(false); }

  return (
    <AuthContext.Provider value={{ role, username }}>
      <BrowserRouter>
        <Routes>
          <Route
            path="/login"
            element={authed ? <Navigate to={role === 'tutor' ? '/students' : '/computers'} replace /> : <Login onLogin={handleLogin} />}
          />
          <Route
            path="/*"
            element={
              <RequireAuth>
                <Layout onLogout={handleLogout} role={role} username={username}>
                  <Routes>
                    <Route index element={<Navigate to={role === 'tutor' ? '/students' : '/computers'} replace />} />
                    <Route path="computers"            element={<RequireNonTutor role={role}><Computers /></RequireNonTutor>} />
                    <Route path="computers/new-desktop" element={<RequireNonTutor role={role}><NewDesktop /></RequireNonTutor>} />
                    <Route path="computers/new-laptop"  element={<RequireNonTutor role={role}><NewLaptop /></RequireNonTutor>} />
                    <Route path="students"             element={<Students />} />
                    <Route path="reference"            element={<RequireNonTutor role={role}><Reference /></RequireNonTutor>} />
                    <Route path="admin/users"          element={<RequireAdmin role={role}><Users /></RequireAdmin>} />
                    <Route path="*"                    element={<Navigate to={role === 'tutor' ? '/students' : '/computers'} replace />} />
                  </Routes>
                </Layout>
              </RequireAuth>
            }
          />
          <Route path="/" element={<Navigate to={authed ? '/computers' : '/login'} replace />} />
        </Routes>
      </BrowserRouter>
    </AuthContext.Provider>
  );
}

