
-- name: CreateBorrower :one
INSERT INTO borrowers (
  user_id,
  loan_id,
  created_by,
  last_updated_by,
  ip_from,
  user_agent
) VALUES (
  $1, $2,$3,$4,$5,$6
)
RETURNING *;

-- name: GetBorrower :one
SELECT * FROM borrowers
WHERE id = $1 LIMIT 1;

-- name: GetBorrowerByUserIdAndLoanId :one
SELECT * FROM borrowers
WHERE user_id = $1 AND loan_id=$2  LIMIT 1;

-- name: ListBorrower :many
SELECT * FROM borrowers 
ORDER BY id 
LIMIT $1 OFFSET $2;

-- name: ListDescBorrower :many
SELECT * FROM borrowers 
ORDER BY id  DESC
LIMIT $1 OFFSET $2;


-- name: UpdateBorrower :one
UPDATE borrowers SET 
user_id=$2,
loan_id=$3,
last_updated_by=$4,
updated_at=$5
WHERE id = $1 RETURNING *;

-- name: DeleteBorrower :exec
UPDATE borrowers SET 
is_active=false
WHERE id = $1 RETURNING *;
