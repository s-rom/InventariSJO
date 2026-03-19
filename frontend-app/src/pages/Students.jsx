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
// Tab 2 — Laptop assignments management
// ─────────────────────────────────────────────
function AssignmentsTab() {
  const { role } = useAuth();
  const isTutor  = role === 'tutor';

  const [laptops,       setLaptops]       = useState([]);
  const [laptopId,      setLaptopId]      = useState(null);
  const [assignments,   setAssignments]   = useState([]);
  const [allClasses,    setAllClasses]    = useState([]); // for dropdowns
  const [allStudents,   setAllStudents]   = useState([]); // all students flat
  const [lmMap,         setLmMap]         = useState({});
  const [saving,        setSaving]        = useState(false);
  const [showForm,      setShowForm]      = useState(false);
  const [err,           setErr]           = useState('');
  const cd = useConfirmDelete();

  const EMPTY_A = { student_id: null, class_id: null, academic_year: currentAcademicYear() };
  const [form, setForm]         = useState(EMPTY_A);
  const [editId, setEditId]     = useState(null);
  const [editForm, setEditForm] = useState({});

  function setF(k, v)  { setForm(f => ({ ...f, [k]: v })); }
  function setEF(k, v) { setEditForm(f => ({ ...f, [k]: v })); }

  function currentAcademicYear() {
    const now = new Date();
    const y = now.getFullYear();
    return now.getMonth() >= 8 ? `${y}-${y + 1}` : `${y - 1}-${y}`;
  }

  // Load laptops, models, classes + students on mount
  useEffect(() => {
    async function load() {
      const [lts, lm] = await Promise.all([api.listLaptops(), api.listLaptopModels()]).catch(() => [[], []]);
      const lmMapObj = Object.fromEntries((lm ?? []).map(m => [m.laptop_model_id, `${m.brand_name} ${m.model_name}`]));
      setLaptops(lts ?? []);
      setLmMap(lmMapObj);

      let flatClasses;
      if (isTutor) {
        // Tutors only see their own classes
        const myClasses = await api.listMyClasses().catch(() => []);
        flatClasses = (myClasses ?? []).map(cl => ({ ...cl, cycleName: cl.cycle_name ?? '' }));
      } else {
        // Admins/editors see all classes across all cycles
        const cycles = await api.listCycles().catch(() => []);
        const classesNested = await Promise.all(
          (cycles ?? []).map(c =>
            api.listClassesByCycle(c.cycle_id)
              .then(cls => (cls ?? []).map(cl => ({ ...cl, cycleName: c.name })))
              .catch(() => [])
          )
        );
        flatClasses = classesNested.flat();
      }
      setAllClasses(flatClasses);

      const studentsNested = await Promise.all(
        flatClasses.map(cl =>
          api.listStudentsByClass(cl.class_id)
            .then(sts => (sts ?? []).map(s => ({ ...s, className: `${cl.cycleName} ${cl.course}r${cl.class_label}` })))
            .catch(() => [])
        )
      );
      setAllStudents(studentsNested.flat());
    }
    load();
  }, [isTutor]);

  const loadAssignments = useCallback(() => {
    if (!laptopId) return;
    api.listAssignmentsByLaptop(laptopId).then(d => setAssignments(d ?? [])).catch(() => {});
  }, [laptopId]);

  useEffect(() => { setAssignments([]); setShowForm(false); setEditId(null); setForm(prev => ({ ...prev, student_id: null, class_id: null })); loadAssignments(); }, [loadAssignments]);

  async function create(e) {
    e.preventDefault(); if (!laptopId) return; setErr(''); setSaving(true);
    try {
      await api.createAssignment(laptopId, {
        student_id:    form.student_id,
        class_id:      form.class_id,
        academic_year: form.academic_year,
      });
      setShowForm(false);
      setForm(EMPTY_A);
      loadAssignments();
    } catch (ex) { setErr(ex.message); }
    finally { setSaving(false); }
  }

  async function saveEdit(id) {
    try {
      await api.updateAssignment(id, {
        student_id:    editForm.student_id    || undefined,
        class_id:      editForm.class_id      || undefined,
        academic_year: editForm.academic_year || undefined,
      });
      setEditId(null); loadAssignments();
    } catch (ex) { setErr(ex.message); }
  }

  async function del(id) {
    try { await api.deleteAssignment(id); loadAssignments(); }
    catch (ex) { setErr(ex.message); }
    cd.cancelDelete();
  }

  const laptopOpts = laptops.map(l => ({
    value: l.computer_id,
    label: l.laptop_model_id
      ? `${l.hostname}  (${lmMap[l.laptop_model_id] ?? '—'})`
      : l.hostname,
  }));

  const studentOpts = allStudents.map(s => ({
    value: s.student_id,
    label: `${s.full_name}  [${s.className}]`,
  }));

  const classOpts = allClasses.map(c => ({
    value: c.class_id,
    label: `${c.cycleName} · ${c.course}r ${c.class_label} (${SHIFTS_LABEL[c.shift] ?? c.shift})`,
  }));

  const selectedLaptop = laptops.find(l => l.computer_id === laptopId);

  return (
    <div>
      {/* Laptop selector */}
      <div className="filter-bar">
        <div className="filter-item" style={{ minWidth: 340 }}>
          <label>Portàtil</label>
          <Combobox
            options={laptopOpts}
            value={laptopId}
            onChange={v => setLaptopId(v)}
            placeholder="Cerca portàtil per hostname…"
          />
        </div>
        {selectedLaptop && (
          <div className="filter-badge">
            <strong>{selectedLaptop.hostname}</strong>
            {selectedLaptop.laptop_model_id && <> · {lmMap[selectedLaptop.laptop_model_id]}</>}
          </div>
        )}
      </div>

      {laptopId && (
        <>
          {/* Assignment list */}
          <div className="card" style={{ marginBottom: 14 }}>
            <div className="table-wrap">
              <table>
                <thead>
                  <tr>
                    <th>Alumne</th>
                    <th>Classe</th>
                    <th>Any acadèmic</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  {assignments.length === 0 && (
                    <tr><td colSpan={4} style={{ textAlign: 'center', color: 'var(--muted)', padding: 20 }}>
                      Cap assignació per a aquest portàtil.
                    </td></tr>
                  )}
                  {assignments.map(a => (
                    <tr key={a.assignment_id}>
                      <td>
                        {editId === a.assignment_id
                          ? <Combobox options={studentOpts} value={editForm.student_id} onChange={v => setEF('student_id', v)} placeholder="Alumne…" />
                          : <strong>{a.student_name}</strong>}
                      </td>
                      <td>
                        {editId === a.assignment_id
                          ? <Combobox options={classOpts} value={editForm.class_id} onChange={v => setEF('class_id', v)} placeholder="Classe…" />
                          : <span style={{ color: 'var(--muted)', fontSize: 12 }}>
                              {a.course}r{a.class_label} · {SHIFTS_LABEL[a.shift] ?? a.shift}
                            </span>}
                      </td>
                      <td>
                        {editId === a.assignment_id
                          ? <input type="text" value={editForm.academic_year} onChange={e => setEF('academic_year', e.target.value)} style={{ width: 100 }} />
                          : a.academic_year}
                      </td>
                      <td style={{ textAlign: 'right', whiteSpace: 'nowrap' }}>
                        {editId === a.assignment_id ? (
                          <>
                            <button className="btn btn-primary btn-sm" onClick={() => saveEdit(a.assignment_id)}>Guardar</button>
                            <button className="btn btn-ghost btn-sm" style={{ marginLeft: 4 }} onClick={() => setEditId(null)}>Cancel·lar</button>
                          </>
                        ) : cd.isAsking(a.assignment_id) ? (
                          <>
                            <span style={{ fontSize: 12, marginRight: 8, color: 'var(--muted)' }}>Eliminar?</span>
                            <button className="btn btn-danger btn-sm" onClick={() => del(a.assignment_id)}>Sí</button>
                            <button className="btn btn-ghost btn-sm" onClick={cd.cancelDelete}>No</button>
                          </>
                        ) : (
                          <>
                            <button className="btn btn-ghost btn-sm" onClick={() => {
                              setEditId(a.assignment_id);
                              setEditForm({ student_id: a.student_id, class_id: a.class_id, academic_year: a.academic_year });
                            }}>Editar</button>
                            <button className="btn btn-danger btn-sm" style={{ marginLeft: 4 }} onClick={() => cd.askDelete(a.assignment_id)}>Eliminar</button>
                          </>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Add assignment form */}
          {showForm ? (
            <div className="card">
              <form onSubmit={create} className="form-panel">
                <div style={{ fontSize: 14, fontWeight: 600, marginBottom: 14 }}>Nova assignació</div>
                <div className="form-grid">
                  <div className="form-group">
                    <label>Alumne *</label>
                    <Combobox
                      options={studentOpts}
                      value={form.student_id}
                      onChange={v => setF('student_id', v)}
                      placeholder="Cerca alumne…"
                    />
                  </div>
                  <div className="form-group">
                    <label>Classe *</label>
                    <Combobox
                      options={classOpts}
                      value={form.class_id}
                      onChange={v => setF('class_id', v)}
                      placeholder="Selecciona classe…"
                    />
                  </div>
                  <div className="form-group">
                    <label>Any acadèmic *</label>
                    <input
                      type="text"
                      value={form.academic_year}
                      onChange={e => setF('academic_year', e.target.value)}
                      placeholder="2025-2026"
                      required
                    />
                  </div>
                </div>
                {err && <div className="error-msg" style={{ marginTop: 10 }}>{err}</div>}
                <div className="form-actions">
                  <button type="submit" className="btn btn-primary" disabled={saving || !form.student_id || !form.class_id || !form.academic_year}>
                    {saving ? 'Guardant…' : 'Crear assignació'}
                  </button>
                  <button type="button" className="btn btn-ghost" onClick={() => { setShowForm(false); setErr(''); }}>
                    Cancel·lar
                  </button>
                </div>
              </form>
            </div>
          ) : (
            <button className="btn btn-primary" onClick={() => setShowForm(true)}>
              + Nova assignació
            </button>
          )}
        </>
      )}
    </div>
  );
}

// ─────────────────────────────────────────────
// Tutor view: student-centric assignment table
// ─────────────────────────────────────────────
function TutorView() {
  function currentAcademicYear() {
    const now = new Date();
    const y = now.getFullYear();
    return now.getMonth() >= 8 ? `${y}-${y + 1}` : `${y - 1}-${y}`;
  }

  const [classes,     setClasses]     = useState([]);
  const [classId,     setClassId]     = useState(null);
  const [students,    setStudents]    = useState([]);
  const [assignments, setAssignments] = useState([]);
  const [laptops,     setLaptops]     = useState([]);
  const [pendingId,   setPendingId]   = useState(null);
  const [err,         setErr]         = useState('');

  useEffect(() => {
    Promise.all([api.listMyClasses(), api.listLaptops()])
      .then(([cls, lts]) => {
        const list = cls ?? [];
        setClasses(list);
        if (list.length) setClassId(list[0].class_id);
        setLaptops(lts ?? []);
      })
      .catch(() => {});
  }, []);

  const loadData = useCallback(async () => {
    if (!classId) return;
    const [sts, asgns] = await Promise.all([
      api.listStudentsByClass(classId).catch(() => []),
      api.listAssignmentsByClass(classId).catch(() => []),
    ]);
    setStudents(sts ?? []);
    setAssignments(asgns ?? []);
  }, [classId]);

  useEffect(() => { loadData(); }, [loadData]);

  const assignmentByStudent = Object.fromEntries(
    (assignments ?? []).map(a => [a.student_id, a])
  );

  const laptopOpts = laptops.map(l => ({ value: l.computer_id, label: l.hostname }));

  async function handleAssign(studentId, newLaptopId) {
    setPendingId(studentId);
    setErr('');
    const existing = assignmentByStudent[studentId];
    try {
      if (existing) {
        await api.deleteAssignment(existing.assignment_id);
      }
      if (newLaptopId) {
        await api.createAssignment(newLaptopId, {
          student_id:    studentId,
          class_id:      classId,
          academic_year: currentAcademicYear(),
        });
      }
      await loadData();
    } catch (ex) {
      setErr(ex.message);
    } finally {
      setPendingId(null);
    }
  }

  const selectedClass = classes.find(c => c.class_id === classId);

  return (
    <div>
      <div className="filter-bar">
        <div className="filter-item">
          <label>Classe</label>
          <select value={classId ?? ''} onChange={e => setClassId(Number(e.target.value))} style={{ width: 240 }}>
            {classes.length === 0 && <option value="">— sense cursos assignats —</option>}
            {classes.map(c => (
              <option key={c.class_id} value={c.class_id}>
                {c.cycle_name ?? ''} · {c.course}r {c.class_label} — {SHIFTS_LABEL[c.shift] ?? c.shift}
              </option>
            ))}
          </select>
        </div>
        {selectedClass && (
          <div className="filter-badge">
            <strong>{students.length}</strong> alumnes · any {currentAcademicYear()}
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
                </tr>
              </thead>
              <tbody>
                {students.length === 0 && (
                  <tr><td colSpan={2} style={{ textAlign: 'center', color: 'var(--muted)', padding: 20 }}>
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
        <TutorView />
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
