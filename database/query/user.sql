
-- name: CreateUser :one
INSERT INTO users (
  name,
  mobile,
  email,
  hashed_password,
  password_changed_at,
  address,
  created_by,
  last_updated_by,
  ip_from,
  user_agent
) VALUES (
  $1, $2,$3,$4,$5,$6,$7,$8,$9,$10
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUser :many
SELECT * FROM users 
ORDER BY id 
LIMIT $1 OFFSET $2;

-- name: ListDescUser :many
SELECT * FROM users 
ORDER BY id  DESC
LIMIT $1 OFFSET $2;


-- name: UpdateUser :one
UPDATE users SET 
name=$2,
mobile=$3,
email=$4,
address=$5,
hashed_password=$6,
last_updated_by=$7,
updated_at=$8
WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
UPDATE users SET 
is_active=false
WHERE id = $1 RETURNING *;
