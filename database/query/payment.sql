
-- name: CreatePayment :one
INSERT INTO payments (
  loan_id,
  user_id,
  amount,
  created_by,
  last_updated_by,
  ip_from,
  user_agent
) VALUES (
  $1, $2,$3,$4,$5,$6,$7
)
RETURNING *;

-- name: GetPayment :one
SELECT * FROM payments
WHERE id = $1 LIMIT 1;

-- name: ListPayment :many
SELECT * FROM payments 
ORDER BY id 
LIMIT $1 OFFSET $2;

-- name: ListDescPayment :many
SELECT * FROM payments 
ORDER BY id  DESC
LIMIT $1 OFFSET $2;


-- name: UpdatePayment :one
UPDATE payments SET 
loan_id=$2,
user_id=$3,
amount=$4,
last_updated_by=$5,
updated_at=$6
WHERE id = $1 RETURNING *;


