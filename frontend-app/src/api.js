const BASE = '/api';
const TOKEN_KEY = 'inventari_token';
const ROLE_KEY   = 'inventari_role';

export function getToken() { return localStorage.getItem(TOKEN_KEY); }
function setToken(t) { localStorage.setItem(TOKEN_KEY, t); }
function clearToken() { localStorage.removeItem(TOKEN_KEY); }
export function getRole() { return localStorage.getItem(ROLE_KEY); }
function setRole(r) { localStorage.setItem(ROLE_KEY, r); }
function clearRole() { localStorage.removeItem(ROLE_KEY); }

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
    setRole(data.role_id);
    return data;
  },
  logout: async () => {
    await req('/auth/logout', { method: 'POST' }).catch(() => {});
    clearToken();
    clearRole();
  },
  me: () => req('/auth/me'),
  changePassword: (data) => req('/auth/change-password', { method: 'POST', body: JSON.stringify(data) }),

  // ── CPUs ──────────────────────────────────────────────────
  listCpus:   ()         => req('/cpus'),
  createCpu:  (data)     => req('/cpus', { method: 'POST',  body: JSON.stringify(data) }),
  updateCpu:  (id, data) => req(`/cpus/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteCpu:  (id)       => req(`/cpus/${id}`, { method: 'DELETE' }),

  // ── Operating Systems ─────────────────────────────────────
  listOS:    ()         => req('/os'),
  createOS:  (data)     => req('/os', { method: 'POST',  body: JSON.stringify(data) }),
  updateOS:  (id, data) => req(`/os/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteOS:  (id)       => req(`/os/${id}`, { method: 'DELETE' }),

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
  listEquipmentUsers:   ()         => req('/equipment-users'),
  createEquipmentUser:  (data)     => req('/equipment-users', { method: 'POST',  body: JSON.stringify(data) }),
  updateEquipmentUser:  (id, data) => req(`/equipment-users/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteEquipmentUser:  (id)       => req(`/equipment-users/${id}`, { method: 'DELETE' }),

  // ── Centers ───────────────────────────────────────────────
  listCenters:   ()         => req('/centers'),
  createCenter:  (data)     => req('/centers', { method: 'POST',  body: JSON.stringify(data) }),
  updateCenter:  (id, data) => req(`/centers/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteCenter:  (id)       => req(`/centers/${id}`, { method: 'DELETE' }),

  // ── Rooms ─────────────────────────────────────────────────
  listRoomsByCenter: (centerId)         => req(`/centers/${centerId}/rooms`),
  createRoom:        (centerId, data)   => req(`/centers/${centerId}/rooms`, { method: 'POST',  body: JSON.stringify(data) }),
  updateRoom:        (id, data)         => req(`/rooms/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteRoom:        (id)               => req(`/rooms/${id}`, { method: 'DELETE' }),

  // ── Cycles ────────────────────────────────────────────────
  listCycles:   ()         => req('/cycles'),
  createCycle:  (data)     => req('/cycles', { method: 'POST',  body: JSON.stringify(data) }),
  updateCycle:  (id, data) => req(`/cycles/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteCycle:  (id)       => req(`/cycles/${id}`, { method: 'DELETE' }),

  // ── Classes ───────────────────────────────────────────────
  listClassesByCycle: (cycleId)       => req(`/cycles/${cycleId}/classes`),
  listMyClasses:      ()              => req('/tutor/classes'),
  createClass:        (cycleId, data) => req(`/cycles/${cycleId}/classes`, { method: 'POST',  body: JSON.stringify(data) }),
  updateClass:        (id, data)      => req(`/classes/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteClass:        (id)            => req(`/classes/${id}`, { method: 'DELETE' }),

  // ── Students ──────────────────────────────────────────────
  listStudentsByClass: (classId)       => req(`/classes/${classId}/students`),
  createStudent:       (classId, data) => req(`/classes/${classId}/students`, { method: 'POST',  body: JSON.stringify(data) }),
  updateStudent:       (id, data)      => req(`/students/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteStudent:       (id)            => req(`/students/${id}`, { method: 'DELETE' }),
  importStudentsCSV: (classId, file) => {
    const token = getToken();
    const headers = {};
    if (token) headers['Authorization'] = `Bearer ${token}`;
    const form = new FormData();
    form.append('file', file);
    return fetch(`${BASE}/classes/${classId}/students/import`, { method: 'POST', headers, body: form })
      .then(async res => {
        const json = await res.json().catch(() => null);
        if (!res.ok) throw new Error(json?.error || `Error ${res.status}`);
        return json;
      });
  },

  // ── Laptop Assignments ────────────────────────────────────
  listAssignmentsByLaptop:      (laptopId)        => req(`/laptops/${laptopId}/assignments`),
  listAssignmentsByClass:       (classId)         => req(`/classes/${classId}/assignments`),
  listAssignmentsByClassAndYear:(classId, year)   => req(`/classes/${classId}/assignments?year=${encodeURIComponent(year)}`),
  listAssignmentsByYear:        (year)            => req(`/assignments?year=${encodeURIComponent(year)}`),
  listAcademicYears:            ()               => req('/academic-years'),
  createAssignment:        (laptopId, data) => req(`/laptops/${laptopId}/assignments`, { method: 'POST',  body: JSON.stringify(data) }),
  updateAssignment:        (id, data)       => req(`/assignments/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteAssignment:        (id)             => req(`/assignments/${id}`, { method: 'DELETE' }),

  // ── Computers (base) ────────────────────────────────────
  deleteComputer: (id)       => req(`/computers/${id}`, { method: 'DELETE' }),

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

  // ── Audit ─────────────────────────────────────────────
  getAuditLog: (table, recordId) => req(`/audit?table=${encodeURIComponent(table)}&record_id=${encodeURIComponent(recordId)}`),

  // ── Users (admin) ─────────────────────────────────────────
  listUsers:   ()             => req('/users'),
  createUser:  (data)         => req('/users', { method: 'POST',  body: JSON.stringify(data) }),
  updateUser:  (id, data)     => req(`/users/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteUser:  (id)           => req(`/users/${id}`, { method: 'DELETE' }),

  // ── Roles (admin) ─────────────────────────────────────────
  listRoles:   () => req('/roles'),
};

