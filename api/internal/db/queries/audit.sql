-- name: InsertAuditLog :exec
INSERT INTO audit_log (
    table_name, record_id, event_type,
    old_values, new_values,
    changed_by_app_user_id, changed_by_username
) VALUES (
    @table_name, @record_id, @event_type,
    @old_values, @new_values,
    @changed_by_app_user_id, @changed_by_username
);

-- name: GetAuditLog :many
SELECT *
FROM audit_log
WHERE table_name = $1 AND record_id = $2
ORDER BY changed_at DESC;
