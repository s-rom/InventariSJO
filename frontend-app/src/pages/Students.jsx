import { useState, useEffect, useCallback } from 'react';
import { api } from '../api';
import { useAuth } from '../App';
import Combobox from '../components/Combobox';

const TABS = [
  { id: 'students',    label: '👨‍🎓 Alumnes per classe' },
  { id: 'assignments', label: '💻 Assignacions portàtils' },
];

const SHIFTS_LABEL = { morning: 'Matí', afternoon: 'Tarda' };

function useConfirmDelete() {
  const [delId, setDelId] = useState(null);
  return { delId, askDelete: setDelId, cancelDelete: () => setDelId(null), isAsking: (id) => delId === id };
}

// ─────────────────────────────────────────────
// Tab 1 — Students CRUD per class
// ─────────────────────────────────────────────
function StudentsTab() {
  const { role } = useAuth();
  const isTutor  = role === 'tutor';

  const [cycles,     setCycles]     = useState([]);
  const [cycleId,    setCycleId]    = useState(null);
  const [classes,    setClasses]    = useState([]);
  const [classId,    setClassId]    = useState(null);
  const [students,   setStudents]   = useState([]);
  const [allClasses, setAllClasses] = useState([]); // for "move to class"
  const [name,       setName]       = useState('');
  const [editId,     setEditId]     = useState(null);
  const [editName,   setEditName]   = useState('');
  const [editClassId,setEditClassId]= useState(null);
  const [saving,     setSaving]     = useState(false);
  const [err,        setErr]        = useState('');
  const [csvResult,  setCsvResult]  = useState(null);
  const cd = useConfirmDelete();

  // Load initial class list
  useEffect(() => {
    if (isTutor) {
      // Tutors only see their assigned classes
      api.listMyClasses()
        .then(d => {
          const list = d ?? [];
          setClasses(list);
          setAllClasses(list);
          if (list.length) setClassId(list[0].class_id);
        })
        .catch(() => {});
    } else {
      api.listCycles()
        .then(d => { const list = d ?? []; setCycles(list); if (list.length) setCycleId(list[0].cycle_id); })
        .catch(() => {});
    }
  }, [isTutor]);

  // Load classes when cycle changes (admin/editor only)
  useEffect(() => {
    if (isTutor) return;
    setClassId(null); setClasses([]);
    if (!cycleId) return;
    api.listClassesByCycle(cycleId)
      .then(d => { const list = d ?? []; setClasses(list); if (list.length) setClassId(list[0].class_id); })
      .catch(() => {});
  }, [cycleId, isTutor]);

  // Load all classes across cycles for the "move" dropdown (admin/editor only)
  useEffect(() => {
    if (isTutor) return;
    async function loadAll() {
      const allCycles = await api.listCycles().catch(() => []);
      const nested = await Promise.all(
        (allCycles ?? []).map(c =>
          api.listClassesByCycle(c.cycle_id)
            .then(cls => (cls ?? []).map(cl => ({ ...cl, cycleName: c.name })))
            .catch(() => [])
        )
      );
      setAllClasses(nested.flat());
    }
    loadAll();
  }, [isTutor]);

  const loadStudents = useCallback(() => {
    if (!classId) return;
    api.listStudentsByClass(classId).then(d => setStudents(d ?? [])).catch(() => {});
  }, [classId]);

  useEffect(() => { setStudents([]); loadStudents(); }, [loadStudents]);

  async function create(e) {
    e.preventDefault(); if (!classId) return; setErr(''); setSaving(true);
    try { await api.createStudent(classId, { full_name: name }); setName(''); loadStudents(); }
    catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function saveEdit(id) {
    try {
      const body = { full_name: editName };
      if (editClassId && editClassId !== classId) body.class_id = editClassId;
      await api.updateStudent(id, body);
      setEditId(null);
      loadStudents();
    } catch (ex) { setErr(ex.message); }
  }

  async function del(id) {
    try { await api.deleteStudent(id); loadStudents(); }
    catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  async function handleImportCSV(e) {
    const file = e.target.files?.[0];
    e.target.value = '';
    if (!file) return;
    setErr(''); setCsvResult(null);
    try {
      const result = await api.importStudentsCSV(classId, file);
      setCsvResult(result);
      loadStudents();
    } catch (ex) { setErr(ex.message); }
  }

  const selectedClass = classes.find(c => c.class_id === classId);
  const classOpts = allClasses.map(c => ({
    value: c.class_id,
    label: c.cycleName
      ? `${c.cycleName} — ${c.course}r ${c.class_label} (${SHIFTS_LABEL[c.shift] ?? c.shift})`
      : `${c.cycle_name ?? ''} — ${c.course}r ${c.class_label} (${SHIFTS_LABEL[c.shift] ?? c.shift})`,
  }));

  return (
    <div>
      {/* Filters */}
      <div className="filter-bar">
        {!isTutor && (
          <div className="filter-item">
            <label>Cicle</label>
            <select value={cycleId ?? ''} onChange={e => setCycleId(Number(e.target.value))} style={{ width: 160 }}>
              {cycles.length === 0 && <option value="">— cap cicle —</option>}
              {cycles.map(c => <option key={c.cycle_id} value={c.cycle_id}>{c.name}</option>)}
            </select>
          </div>
        )}
        <div className="filter-item">
          <label>Classe</label>
          <select value={classId ?? ''} onChange={e => setClassId(Number(e.target.value))} style={{ width: 200 }}>
            {classes.length === 0 && <option value="">{isTutor ? '— sense cursos assignats —' : '— cap classe —'}</option>}
            {classes.map(c => (
              <option key={c.class_id} value={c.class_id}>
                {isTutor ? `${c.cycle_name ?? ''} · ` : ''}{c.course}r {c.class_label} — {SHIFTS_LABEL[c.shift] ?? c.shift}
              </option>
            ))}
          </select>
        </div>
        {selectedClass && (
          <div className="filter-badge">
            {isTutor
              ? `${selectedClass.cycle_name ?? ''} · `
              : `${cycles.find(c => c.cycle_id === cycleId)?.name ?? ''} · `
            }{selectedClass.course}r {selectedClass.class_label} · {SHIFTS_LABEL[selectedClass.shift] ?? selectedClass.shift} · <strong>{students.length}</strong> alumnes
          </div>
        )}
      </div>

      {classId && (
        <>
          {/* Add student */}
          <div className="card" style={{ marginBottom: 14 }}>
            <form onSubmit={create} className="form-panel">
              <div className="form-grid" style={{ gridTemplateColumns: '1fr auto' }}>
                <div className="form-group">
                  <label>Nom complet de l&apos;alumne *</label>
                  <input
                    type="text"
                    value={name}
                    onChange={e => setName(e.target.value)}
                    placeholder="Cognoms, Nom"
                    required
                  />
                </div>
                <div className="form-group" style={{ justifyContent: 'flex-end' }}>
                  <label style={{ visibility: 'hidden' }}>_</label>
                  <button type="submit" className="btn btn-primary" disabled={saving || !name}>
                    {saving ? '…' : 'Afegir alumne'}
                  </button>
                </div>
              </div>
              {!isTutor && (
                <div style={{ marginTop: 10, display: 'flex', alignItems: 'center', gap: 12, flexWrap: 'wrap' }}>
                  <label className="btn btn-ghost btn-sm" style={{ cursor: 'pointer', marginBottom: 0 }}>
                    📥 Importar CSV
                    <input type="file" accept=".csv,text/csv" style={{ display: 'none' }} onChange={handleImportCSV} />
                  </label>
                  <span style={{ fontSize: 12, color: 'var(--muted)' }}>Primera columna = nom de l&apos;alumne, primera fila = capçalera</span>
                  {csvResult && (
                    <span style={{ fontSize: 12 }}>
                      ✅ {csvResult.imported} importats
                      {csvResult.skipped?.length > 0 && <> · ⚠️ {csvResult.skipped.length} omesos: {csvResult.skipped.join(', ')}</>}
                    </span>
                  )}
                </div>
              )}
              {err && <div className="error-msg">{err}</div>}
            </form>
          </div>

          {/* Students table */}
          <div className="card">
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Nom</th>
                    <th>Classe</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  {students.length === 0 && (
                    <tr><td colSpan={3} style={{ textAlign: 'center', color: 'var(--muted)', padding: 20 }}>
                      Cap alumne a aquesta classe.
                    </td></tr>
                  )}
                  {students.map(s => (
                    <tr key={s.student_id}>
                      <td>
                        {editId === s.student_id
                          ? <input type="text" value={editName} onChange={e => setEditName(e.target.value)} style={{ width: 260 }} autoFocus />
                          : <strong>{s.full_name}</strong>}
                      </td>
                      <td>
                        {editId === s.student_id
                          ? <Combobox
                              options={classOpts}
                              value={editClassId}
                              onChange={v => setEditClassId(v)}
                              placeholder="Mantenir classe actual…"
                              nullable
                            />
                          : <span style={{ color: 'var(--muted)', fontSize: 12 }}>
                              {cycles.find(c => c.cycle_id === cycleId)?.name} · {selectedClass?.course}r {selectedClass?.class_label}
                            </span>}
                      </td>
                      <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                        {editId === s.student_id ? (
                          <>
                            <button className="btn btn-primary btn-sm" onClick={() => saveEdit(s.student_id)}>Guardar</button>
                            <button className="btn btn-ghost btn-sm" style={{ marginLeft: 4 }} onClick={() => setEditId(null)}>Cancel·lar</button>
                          </>
                        ) : cd.isAsking(s.student_id) ? (
                          <>
                            <span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Eliminar alumne?</span>
                            <button className="btn btn-danger btn-sm" onClick={() => del(s.student_id)}>Sí</button>
                            <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button>
                          </>
                        ) : (
                          <>
                            <button className="btn btn-ghost btn-sm" onClick={() => { setEditId(s.student_id); setEditName(s.full_name); setEditClassId(classId); }}>
                              Editar
                            </button>
                            <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => cd.askDelete(s.student_id)}>
                              Eliminar
                            </button>
                          </>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </>
      )}
    </div>
  );
}

// ─────────────────────────────────────────────
// Tab 2 — Laptop assignments (class + year centric)
// ─────────────────────────────────────────────
function AssignmentsTab() {
  const { role } = useAuth();
  const isTutor  = role === 'tutor';

  function currentAcademicYear() {
    const now = new Date();
    const y = now.getFullYear();
    return now.getMonth() >= 8 ? `${y}-${y + 1}` : `${y - 1}-${y}`;
  }

  const [academicYears, setAcademicYears] = useState([]);
  const [selectedYear,  setSelectedYear]  = useState(currentAcademicYear());
  const [cycles,        setCycles]        = useState([]);
  const [cycleId,       setCycleId]       = useState(null);
  const [classes,       setClasses]       = useState([]);
  const [classId,       setClassId]       = useState(null);
  const [students,      setStudents]      = useState([]);
  const [assignments,   setAssignments]   = useState([]);
  const [laptops,       setLaptops]       = useState([]);
  const [lmMap,         setLmMap]         = useState({});
  const [pendingId,     setPendingId]     = useState(null);
  const [err,           setErr]           = useState('');
  const cd = useConfirmDelete();

  // Load years, laptops, models, and initial class list
  useEffect(() => {
    async function load() {
      const cur = currentAcademicYear();
      const [years, lts, lm] = await Promise.all([
        api.listAcademicYears().catch(() => []),
        api.listLaptops().catch(() => []),
        api.listLaptopModels().catch(() => []),
      ]);
      const yearList = [...new Set([cur, ...(years ?? [])])].sort().reverse();
      setAcademicYears(yearList);
      setLaptops(lts ?? []);
      setLmMap(Object.fromEntries((lm ?? []).map(m => [m.laptop_model_id, `${m.brand_name} ${m.model_name}`])));

      if (isTutor) {
        const myClasses = await api.listMyClasses().catch(() => []);
        const list = myClasses ?? [];
        setClasses(list);
        if (list.length) setClassId(list[0].class_id);
      } else {
        const cyc = await api.listCycles().catch(() => []);
        setCycles(cyc ?? []);
        if ((cyc ?? []).length) setCycleId(cyc[0].cycle_id);
      }
    }
    load();
  }, [isTutor]);

  // Load classes when cycle changes (admin/editor only)
  useEffect(() => {
    if (isTutor || !cycleId) return;
    api.listClassesByCycle(cycleId)
      .then(d => {
        const list = d ?? [];
        setClasses(list);
        setClassId(list.length ? list[0].class_id : null);
      })
      .catch(() => {});
  }, [cycleId, isTutor]);

  // Load students + assignments when class or year changes
  const loadData = useCallback(async () => {
    if (!classId || !selectedYear) return;
    const [sts, asgns] = await Promise.all([
      api.listStudentsByClass(classId).catch(() => []),
      api.listAssignmentsByClassAndYear(classId, selectedYear).catch(() => []),
    ]);
    setStudents(sts ?? []);
    setAssignments(asgns ?? []);
  }, [classId, selectedYear]);

  useEffect(() => { loadData(); }, [loadData]);

  const assignmentByStudent = Object.fromEntries(
    (assignments ?? []).map(a => [a.student_id, a])
  );

  const laptopOpts = laptops.map(l => ({
    value: l.computer_id,
    label: l.laptop_model_id ? `${l.hostname}  (${lmMap[l.laptop_model_id] ?? '—'})` : l.hostname,
  }));

  async function handleAssign(studentId, newLaptopId) {
    setPendingId(studentId);
    setErr('');
    const existing = assignmentByStudent[studentId];
    try {
      if (existing) await api.deleteAssignment(existing.assignment_id);
      if (newLaptopId) {
        await api.createAssignment(newLaptopId, {
          student_id:    studentId,
          class_id:      classId,
          academic_year: selectedYear,
        });
      }
      await loadData();
    } catch (ex) {
      setErr(ex.message);
    } finally {
      setPendingId(null);
    }
  }

  async function del(assignmentId) {
    try { await api.deleteAssignment(assignmentId); await loadData(); }
    catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  const selectedClass = classes.find(c => c.class_id === classId);
  const assignedCount = students.filter(s => assignmentByStudent[s.student_id]).length;

  return (
    <div>
      <div className="filter-bar">
        <div className="filter-item">
          <label>Any acadèmic</label>
          <select value={selectedYear} onChange={e => setSelectedYear(e.target.value)} style={{ width: 130 }}>
            {academicYears.map(y => <option key={y} value={y}>{y}</option>)}
          </select>
        </div>
        {!isTutor && (
          <div className="filter-item">
            <label>Cicle</label>
            <select value={cycleId ?? ''} onChange={e => setCycleId(Number(e.target.value))} style={{ width: 140 }}>
              {cycles.length === 0 && <option value="">— cap cicle —</option>}
              {cycles.map(c => <option key={c.cycle_id} value={c.cycle_id}>{c.name}</option>)}
            </select>
          </div>
        )}
        <div className="filter-item">
          <label>Classe</label>
          <select value={classId ?? ''} onChange={e => setClassId(Number(e.target.value))} style={{ width: 200 }}>
            {classes.length === 0 && <option value="">{isTutor ? '— sense cursos assignats —' : '— cap classe —'}</option>}
            {classes.map(c => (
              <option key={c.class_id} value={c.class_id}>
                {isTutor ? `${c.cycle_name ?? ''} · ` : ''}{c.course}r {c.class_label} — {SHIFTS_LABEL[c.shift] ?? c.shift}
              </option>
            ))}
          </select>
        </div>
        {selectedClass && (
          <div className="filter-badge">
            {isTutor
              ? `${selectedClass.cycle_name ?? ''} · `
              : `${cycles.find(c => c.cycle_id === cycleId)?.name ?? ''} · `
            }{selectedClass.course}r {selectedClass.class_label}
            {' · '}<strong>{assignedCount}</strong>/{students.length} assignats
          </div>
        )}
      </div>

      {err && <div className="error-msg" style={{ marginBottom: 12 }}>{err}</div>}

      {classId && (
        <div className="card">
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Alumne</th>
                  <th>Portàtil assignat</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {students.length === 0 && (
                  <tr><td colSpan={3} style={{ textAlign: 'center', color: 'var(--muted)', padding: 20 }}>
                    Cap alumne a aquesta classe.
                  </td></tr>
                )}
                {students.map(s => {
                  const assignment = assignmentByStudent[s.student_id];
                  return (
                    <tr key={s.student_id}>
                      <td><strong>{s.full_name}</strong></td>
                      <td style={{ minWidth: 260 }}>
                        {pendingId === s.student_id
                          ? <span style={{ color: 'var(--muted)', fontSize: 12 }}>Guardant…</span>
                          : <Combobox
                              options={laptopOpts}
                              value={assignment?.computer_id ?? null}
                              onChange={v => handleAssign(s.student_id, v)}
                              placeholder="Sense assignació…"
                              nullable
                            />
                        }
                      </td>
                      <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                        {assignment && (
                          cd.isAsking(assignment.assignment_id) ? (
                            <>
                              <span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Eliminar?</span>
                              <button className="btn btn-danger btn-sm" onClick={() => del(assignment.assignment_id)}>Sí</button>
                              <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button>
                            </>
                          ) : (
                            <button className="btn btn-danger btn-sm" onClick={() => cd.askDelete(assignment.assignment_id)}>Eliminar</button>
                          )
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  );
}



// ─────────────────────────────────────────────
// Main page
// ─────────────────────────────────────────────
export default function Students() {
  const { role } = useAuth();
  const isTutor  = role === 'tutor';
  const [tab, setTab] = useState('students');

  if (isTutor) {
    return (
      <>
        <div className="page-header">
          <h1 className="page-title">💻 Assignacions de portàtils</h1>
        </div>
        <AssignmentsTab />
      </>
    );
  }

  return (
    <>
      <div className="page-header">
        <h1 className="page-title">👨‍🎓 Alumnes &amp; Assignacions</h1>
      </div>

      <div className="ref-tabs" style={{ marginBottom: 20 }}>
        {TABS.map(t => (
          <button
            key={t.id}
            className={`ref-tab${tab === t.id ? ' active' : ''}`}
            onClick={() => setTab(t.id)}
          >
            {t.label}
          </button>
        ))}
      </div>

      {tab === 'students'    && <StudentsTab />}
      {tab === 'assignments' && <AssignmentsTab />}
    </>
  );
}
