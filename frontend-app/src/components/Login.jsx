import { useState } from 'react';
import { api } from '../api';

export default function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError]       = useState('');
  const [loading, setLoading]   = useState(false);

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
      </div>
    </div>
  );
}
