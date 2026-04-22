import { useState, useEffect } from 'react';
import { api } from '../api';

export default function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError]       = useState('');
  const [loading, setLoading]   = useState(false);

  // Handle redirect back from Google OAuth callback
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const token    = params.get('token');
    const role     = params.get('role');
    const uname    = params.get('username');
    const oauthErr = params.get('error');

    if (oauthErr) {
      const messages = {
        domain_not_allowed: 'El compte de Google no pertany al domini autoritzat.',
        invalid_state:      'Error de seguretat OAuth. Torna-ho a intentar.',
      };
      setError(messages[oauthErr] || `Error d'autenticació: ${oauthErr}`);
      window.history.replaceState({}, '', '/login');
      return;
    }

    if (token && role && uname) {
      api.storeSession(token, role);
      onLogin({ token, role_id: role, username: uname });
    }
  }, [onLogin]);

  async function handleSubmit(e) {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const data = await api.login(username, password);
      onLogin(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  function handleGoogleLogin() {
    window.location.href = '/api/auth/google';
  }

  return (
    <div className="login-wrap">
      <div className="login-card">
        <div className="login-title">Inventari</div>
        <div className="login-sub">Inicia sessió per continuar</div>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Usuari</label>
            <input
              type="text"
              value={username}
              onChange={e => setUsername(e.target.value)}
              autoFocus
              required
            />
          </div>
          <div className="form-group">
            <label>Contrasenya</label>
            <input
              type="password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              required
            />
          </div>
          {error && <div className="error-msg">{error}</div>}
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Entrant…' : 'Entrar'}
          </button>
        </form>
        <div className="login-divider">o</div>
        <button type="button" className="btn btn-google" onClick={handleGoogleLogin}>
          Inicia sessió amb Google
        </button>
      </div>
    </div>
  );
}
