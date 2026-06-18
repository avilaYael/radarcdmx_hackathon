CREATE TABLE IF NOT EXISTS `user` (
    `uuid` CHAR(36) NOT NULL,
    `name` VARCHAR(255),
    `lastname` VARCHAR(255),
    `email` VARCHAR(512) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `status` INT NOT NULL,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `created_by` CHAR(36) NOT NULL,
    `updated_by` CHAR(36) NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`uuid`),
    UNIQUE INDEX `unique_email` (`email`)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS `establecimiento` (
    `uuid` CHAR(36) NOT NULL,
    `id_denue` INT,
    `clee` VARCHAR(255),
    `nombre` VARCHAR(255),
    `razon_social` VARCHAR(255),
    `per_ocu` VARCHAR(255),
    `codigo_actividad` INT,
    `nombre_actividad` VARCHAR(255),
    `uso_de_suelo` VARCHAR(255),
    `clave_catastral` VARCHAR(255),
    `contacto` JSON,
    `ubicacion` JSON,
    `fecha_alta` DATE,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`uuid`)
) ENGINE = InnoDB;

