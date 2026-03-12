import { useState, useRef, useEffect } from 'react';

/**
 * Combobox — searchable single-select.
 *
 * Props:
 *   options   : { value, label }[]
 *   value     : current value (or null)
 *   onChange  : (value) => void
 *   placeholder: string
 *   nullable  : bool — shows a "clear" option if true
 */
export default function Combobox({ options = [], value, onChange, placeholder = 'Cerca…', nullable = false }) {
  const [open, setOpen]   = useState(false);
  const [query, setQuery] = useState('');
  const wrapRef           = useRef(null);
  const inputRef          = useRef(null);

  const selected = options.find(o => o.value === value) ?? null;

  // Close on outside click
  useEffect(() => {
    function handler(e) {
      if (wrapRef.current && !wrapRef.current.contains(e.target)) {
        setOpen(false);
        setQuery('');
      }
    }
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, []);

  const filtered = query
    ? options.filter(o => o.label.toLowerCase().includes(query.toLowerCase()))
    : options;

  function select(val) {
    onChange(val);
    setOpen(false);
    setQuery('');
  }

  function handleFocus() {
    setOpen(true);
    setQuery('');
  }

  return (
    <div className="combobox" ref={wrapRef}>
      <input
        ref={inputRef}
        type="text"
        value={open ? query : (selected?.label ?? '')}
        placeholder={placeholder}
        onChange={e => setQuery(e.target.value)}
        onFocus={handleFocus}
        readOnly={!open}
        style={{ cursor: open ? 'text' : 'pointer' }}
      />
      {open && (
        <div className="combobox-dropdown">
          {nullable && (
            <div
              className={`combobox-option${value === null ? ' selected' : ''}`}
              onMouseDown={() => select(null)}
            >
              <em style={{ color: 'var(--muted)' }}>Cap / buit</em>
            </div>
          )}
          {filtered.length === 0 ? (
            <div className="combobox-empty">Sense resultats</div>
          ) : (
            filtered.map(o => (
              <div
                key={o.value}
                className={`combobox-option${o.value === value ? ' selected' : ''}`}
                onMouseDown={() => select(o.value)}
              >
                {o.label}
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}
