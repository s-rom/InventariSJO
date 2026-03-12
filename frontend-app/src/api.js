const BASE = '/api';
const TOKEN_KEY = 'inventari_token';

export function getToken() { return localStorage.getItem(TOKEN_KEY); }
function setToken(t) { localStorage.setItem(TOKEN_KEY, t); }
function clearToken() { localStorage.removeItem(TOKEN_KEY); }

async function req(path, opts = {}) {
  const token = getToken();
  const headers = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await fetch(BASE + path, { headers, ...opts });
  if (res.status === 204) return null;
  const json = await res.json().catch(() => null);
  if (!res.ok) throw new Error(json?.error || `Error ${res.status}`);
  return json;
}

export const api = {
  // Auth
  login: async (username, password) => {
    const data = await req('/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) });
    setToken(data.token);
    return data;
  },
  logout: async () => {
    await req('/auth/logout', { method: 'POST' }).catch(() => {});
    clearToken();
  },

  // Reference data
  listCpus:           () => req('/cpus'),
  listOS:             () => req('/os'),
  listEquipmentUsers: () => req('/equipment-users'),
  listCenters:        () => req('/centers'),
  listRoomsByCenter:  (centerId) => req(`/centers/${centerId}/rooms`),

  // Computers
  listComputers:   ()         => req('/computers'),
  createComputer:  (data)     => req('/computers', { method: 'POST',  body: JSON.stringify(data) }),
  updateComputer:  (id, data) => req(`/computers/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteComputer:  (id)       => req(`/computers/${id}`, { method: 'DELETE' }),
};
