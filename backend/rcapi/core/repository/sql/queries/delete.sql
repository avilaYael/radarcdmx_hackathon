-- name: DeleteUser :execresult
DELETE FROM `user`
WHERE
`uuid` = ?;

-- name: DeleteEstablecimiento :execresult
DELETE FROM `establecimiento`
WHERE
`uuid` = ?;

