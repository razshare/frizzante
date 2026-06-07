-- name: FindSessionById :one
select * from sessions where id = :id;

-- name: AddSessionWithId :exec
insert into sessions(id, created_at, updated_at) values (:id, datetime(), datetime());

-- name: ModifySessionById :exec
update sessions set error = :error where id = :id;

-- name: FindTodosBySessionId :many
select * from todos where session_id = :session_id;

-- name: AddTodoWithIdAndSessionId :exec
insert into todos(id, session_id, description, checked) values(:id, :session_id, :description, :checked);

-- name: ToggleTodosByIdAndSessionId :exec
update todos set checked = :checked where id = :id and session_id = :session_id;

-- name: RemoveTodosByIdAndSessionId :exec
delete from todos where id = :id and session_id = :session_id;
