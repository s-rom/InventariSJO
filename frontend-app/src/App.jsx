import { useState, useEffect, useCallback } from 'react';
import { api, getToken } from './api';
import Login from './components/Login';
import ComputerList from './components/ComputerList';
import ComputerForm from './components/ComputerForm';

export default function App() {
  const [authed,    setAuthed]    = useState(null); // null = checking, false = login, true = main
  const [computers, setComputers] = useState([]);
  const [refData,   setRefData]   = useState(null); // { cpus, rooms, equipUsers, osList, cpuMap, roomMap, equipUserMap }
  const [showForm,  setShowForm]  = useState(false);
  const [loadErr,   setLoadErr]   = useState('');

  // --- Initial auth probe ---
  useEffect(() => {
    if (!getToken()) { setAuthed(false); return; }
    api.listComputers()
      .then(data => { setComputers(data ?? []); setAuthed(true); })
      .catch(() => { api.logout(); setAuthed(false); });
  }, []);

  // --- Load reference data when authed ---
  const loadRefData = useCallback(async () => {
    try {
      const [cpusRaw, osRaw, equipRaw, centers] = await Promise.all([
        api.listCpus(),
        api.listOS(),
        api.listEquipmentUsers(),
        api.listCenters(),
      ]);

      // Fetch rooms for all centers in parallel
      const roomsNested = await Promise.all(
        (centers ?? []).map(c =>
          api.listRoomsByCenter(c.center_id)
            .then(rooms => (rooms ?? []).map(r => ({ ...r, centerName: c.name })))
            .catch(() => [])
        )
      );
      const allRooms = roomsNested.flat();

      // Build lookup maps
      const cpuMap       = Object.fromEntries((cpusRaw ?? []).map(c => [c.cpu_id, c.model_name ?? `CPU #${c.cpu_id}`]));
      const roomMap      = Object.fromEntries(allRooms.map(r => [r.room_id, { name: r.name, centerName: r.centerName }]));
      const equipUserMap = Object.fromEntries((equipRaw ?? []).map(e => [e.equipment_user_id, e.name]));

      // Build combobox option arrays
      const cpuOptions      = (cpusRaw ?? []).map(c => ({ value: c.cpu_id, label: c.model_name ?? `CPU #${c.cpu_id}` }));
      const roomOptions     = allRooms.map(r => ({ value: r.room_id, label: `${r.centerName} › ${r.name}` }));
      const equipUserOptions = (equipRaw ?? []).map(e => ({ value: e.equipment_user_id, label: e.name }));

      setRefData({ cpuOptions, roomOptions, equipUserOptions, osList: osRaw ?? [], cpuMap, roomMap, equipUserMap });
    } catch (err) {
      setLoadErr(err.message);
    }
  }, []);

  useEffect(() => {
    if (authed) loadRefData();
  }, [authed, loadRefData]);

  // --- Auth handlers ---
  function onLogin() {
    api.listComputers()
      .then(data => { setComputers(data ?? []); setAuthed(true); })
      .catch(() => setAuthed(true));
  }

  async function onLogout() {
    await api.logout();
    setAuthed(false);
    setComputers([]);
    setRefData(null);
  }

  // --- Computer handlers ---
  function onComputerCreated(computer) {
    setComputers(prev => [...prev, computer]);
    setShowForm(false);
  }

  function onComputerDeleted(id) {
    setComputers(prev => prev.filter(c => c.computer_id !== id));
  }

  // --- Render ---
  if (authed === null) {
    return <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', color: 'var(--muted)' }}>Carregant…</div>;
  }

  if (!authed) {
    return <Login onLogin={onLogin} />;
  }

  return (
    <div className="layout">
      <aside className="sidebar">
        <div className="sidebar-brand">Inventari</div>
        <nav style={{ flex: 1 }}>
          <div className="nav-item active">Equips</div>
        </nav>
        <div className="sidebar-bottom">
          <button className="btn btn-ghost btn-sm" style={{ width: '100%', justifyContent: 'center' }} onClick={onLogout}>
            Tancar sessió
          </button>
        </div>
      </aside>

      <main className="main">
        <div className="page-header">
          <h1 className="page-title">Equips informàtics</h1>
          {loadErr && <span style={{ color: 'var(--danger)', fontSize: 12 }}>Error carregant dades: {loadErr}</span>}
        </div>

        {refData ? (
          <>
            <ComputerList
              computers={computers}
              cpuMap={refData.cpuMap}
              roomMap={refData.roomMap}
              equipUserMap={refData.equipUserMap}
              onDelete={onComputerDeleted}
              onAddClick={() => setShowForm(f => !f)}
              showForm={showForm}
            />

            {showForm && (
              <ComputerForm
                cpus={refData.cpuOptions}
                rooms={refData.roomOptions}
                equipUsers={refData.equipUserOptions}
                osList={refData.osList}
                onCreated={onComputerCreated}
                onCancel={() => setShowForm(false)}
              />
            )}
          </>
        ) : (
          <div className="card" style={{ padding: 24, color: 'var(--muted)' }}>Carregant dades de referència…</div>
        )}
      </main>
    </div>
  );
}
