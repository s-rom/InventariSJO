const PAGE_SIZES = [10, 20, 50, 100];

export default function Pagination({ page, totalPages, pageSize, onPage, onPageSize, pageSizeId }) {
  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 12, marginTop: 12, flexWrap: 'wrap' }}>
      <button className="btn btn-sm" disabled={page === 1} onClick={() => onPage(1)}>&lt;&lt;</button>
      <button className="btn btn-sm" disabled={page === 1} onClick={() => onPage(page - 1)}>Anterior</button>
      <span>Pàgina {page} de {totalPages}</span>
      <button className="btn btn-sm" disabled={page === totalPages} onClick={() => onPage(page + 1)}>Següent</button>
      <button className="btn btn-sm" disabled={page === totalPages} onClick={() => onPage(totalPages)}>&gt;&gt;</button>
      <label htmlFor={pageSizeId} style={{ marginLeft: 16 }}>Mostrar:</label>
      <select
        id={pageSizeId}
        value={pageSize}
        onChange={e => { onPageSize(Number(e.target.value)); onPage(1); }}
        style={{ fontSize: 15, borderRadius: 6, border: '1px solid var(--border, #d0d7de)', padding: '6px 12px', background: '#fff', minWidth: 70 }}
      >
        {PAGE_SIZES.map(n => <option key={n} value={n}>{n}</option>)}
      </select>
    </div>
  );
}
