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
  // ── Auth ──────────────────────────────────────────────────
  login: async (username, password) => {
    const data = await req('/auth/login', { method: 'POST', body: JSON.stringify({ username, password }) });
    setToken(data.token);
    return data;
  },
  logout: async () => {
    await req('/auth/logout', { method: 'POST' }).catch(() => {});
    clearToken();
  },

  // ── CPUs ──────────────────────────────────────────────────
  listCpus:   ()         => req('/cpus'),
  createCpu:  (data)     => req('/cpus', { method: 'POST',  body: JSON.stringify(data) }),
  updateCpu:  (id, data) => req(`/cpus/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteCpu:  (id)       => req(`/cpus/${id}`, { method: 'DELETE' }),

  // ── Operating Systems ─────────────────────────────────────
  listOS:    ()     => req('/os'),
  createOS:  (data) => req('/os', { method: 'POST',   body: JSON.stringify(data) }),
  deleteOS:  (id)   => req(`/os/${id}`, { method: 'DELETE' }),

  // ── Brands ────────────────────────────────────────────────
  listBrands:   ()         => req('/brands'),
  createBrand:  (data)     => req('/brands', { method: 'POST',  body: JSON.stringify(data) }),
  updateBrand:  (id, data) => req(`/brands/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteBrand:  (id)       => req(`/brands/${id}`, { method: 'DELETE' }),

  // ── Laptop Models ─────────────────────────────────────────
  listLaptopModels:   ()         => req('/laptop-models'),
  getLaptopModel:     (id)       => req(`/laptop-models/${id}`),
  createLaptopModel:  (data)     => req('/laptop-models', { method: 'POST',  body: JSON.stringify(data) }),
  updateLaptopModel:  (id, data) => req(`/laptop-models/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteLaptopModel:  (id)       => req(`/laptop-models/${id}`, { method: 'DELETE' }),

  // ── Desktop Models ────────────────────────────────────────
  listDesktopModels:   ()         => req('/desktop-models'),
  getDesktopModel:     (id)       => req(`/desktop-models/${id}`),
  createDesktopModel:  (data)     => req('/desktop-models', { method: 'POST',  body: JSON.stringify(data) }),
  updateDesktopModel:  (id, data) => req(`/desktop-models/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteDesktopModel:  (id)       => req(`/desktop-models/${id}`, { method: 'DELETE' }),

  // ── Equipment Users ───────────────────────────────────────
  listEquipmentUsers: () => req('/equipment-users'),

  // ── Centers & Rooms ───────────────────────────────────────
  listCenters:       ()           => req('/centers'),
  listRoomsByCenter: (centerId)   => req(`/centers/${centerId}/rooms`),

  // ── Desktops ──────────────────────────────────────────────
  listDesktops:   ()         => req('/desktops'),
  getDesktop:     (id)       => req(`/desktops/${id}`),
  createDesktop:  (data)     => req('/desktops', { method: 'POST',  body: JSON.stringify(data) }),
  updateDesktop:  (id, data) => req(`/desktops/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),

  // ── Laptops ───────────────────────────────────────────────
  listLaptops:   ()         => req('/laptops'),
  getLaptop:     (id)       => req(`/laptops/${id}`),
  createLaptop:  (data)     => req('/laptops', { method: 'POST',  body: JSON.stringify(data) }),
  updateLaptop:  (id, data) => req(`/laptops/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
};

