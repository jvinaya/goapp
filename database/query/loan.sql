
-- name: CreateLoan :one
INSERT INTO loans (
  amount,
  amount_need_to_pay,
  term,
  approval_status,
  repayment_status,
  created_by,
  last_updated_by,
  ip_from,
  user_agent
) VALUES (
  $1, $2,$3,$4,$5,$6,$7,$8,$9
)
RETURNING *;

-- name: GetLoan :one
SELECT * FROM loans
WHERE id = $1 LIMIT 1;

-- name: ListLoan :many
SELECT * FROM loans 
ORDER BY id 
LIMIT $1 OFFSET $2;

-- name: ListDescLoan :many
SELECT * FROM loans 
ORDER BY id  DESC
LIMIT $1 OFFSET $2;


-- name: UpdateLoan :one
UPDATE loans SET 
amount=$2,
amount_need_to_pay=$3,
term=$4,
approval_status=$5,
repayment_status=$6,
last_updated_by=$7,
updated_at=$8
WHERE id = $1 RETURNING *;

-- name: UpdateLoanStatus :one
UPDATE loans SET 
approval_status=$2,
last_updated_by=$3,
updated_at=$4
WHERE id = $1 RETURNING *;

-- name: DeleteLoan :exec
UPDATE loans SET 
is_active=false
WHERE id = $1 RETURNING *;
