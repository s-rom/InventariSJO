export default function StatusBadge({ status }) {
  const ok = status === 'actiu';
  return (
    <span style={{
      display: 'inline-block',
      padding: '2px 8px',
      borderRadius: 999,
      fontSize: 11,
      fontWeight: 600,
      background: ok ? 'var(--success-bg, #d1fae5)' : 'var(--danger-bg, #fee2e2)',
      color: ok ? 'var(--success, #065f46)' : 'var(--danger, #991b1b)',
    }}>
      {ok ? 'Actiu' : 'Baixa'}
    </span>
  );
}
