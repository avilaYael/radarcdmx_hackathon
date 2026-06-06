

-- user selects:
-- name: FetchUserByUuid :many
SELECT `uuid`,`name`,`lastname`,`email`,`password`,`status`,`updated_at`,`created_by`,`updated_by`,`created_at`
FROM `user`
WHERE 
    `uuid` = ? ;

        
-- name: FetchUserByEmail :many
SELECT `uuid`,`name`,`lastname`,`email`,`password`,`status`,`updated_at`,`created_by`,`updated_by`,`created_at`
FROM `user`
WHERE 
    `email` = ? 
LIMIT ?, ?;
        
-- name: FetchUserByUuidForUpdate :many
SELECT `uuid`,`name`,`lastname`,`email`,`password`,`status`,`updated_at`,`created_by`,`updated_by`,`created_at`
FROM `user`
WHERE 
    `uuid` = ? 
FOR UPDATE;
        




-- establecimiento selects:
-- name: FetchEstablecimientoByUuid :many
SELECT `uuid`,`id_denue`,`clee`,`nombre`,`razon_social`,`per_ocu`,`codigo_actividad`,`nombre_actividad`,`uso_de_suelo`,`clave_catastral`,`contacto`,`ubicacion`,`fecha_alta`,`created_at`,`updated_at`
FROM `establecimiento`
WHERE 
    `uuid` = ? ;

        
-- name: FetchEstablecimientoByUuidForUpdate :many
SELECT `uuid`,`id_denue`,`clee`,`nombre`,`razon_social`,`per_ocu`,`codigo_actividad`,`nombre_actividad`,`uso_de_suelo`,`clave_catastral`,`contacto`,`ubicacion`,`fecha_alta`,`created_at`,`updated_at`
FROM `establecimiento`
WHERE 
    `uuid` = ? 
FOR UPDATE;
        


