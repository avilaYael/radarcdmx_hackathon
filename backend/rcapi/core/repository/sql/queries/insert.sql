-- name: InsertUser :execresult
INSERT INTO `user`
(`uuid`,`name`,`lastname`,`email`,`password`,`status`,`updated_at`,`created_by`,`updated_by`,`created_at`)
VALUES
(?,?,?,?,?,?,?,?,?,?);

-- name: InsertEstablecimiento :execresult
INSERT INTO `establecimiento`
(`uuid`,`id_denue`,`clee`,`nombre`,`razon_social`,`per_ocu`,`codigo_actividad`,`nombre_actividad`,`uso_de_suelo`,`clave_catastral`,`contacto`,`ubicacion`,`fecha_alta`,`created_at`,`updated_at`)
VALUES
(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);

