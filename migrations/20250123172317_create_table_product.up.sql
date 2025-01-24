CREATE TABLE products
(
    `id`           INT         NOT NULL AUTO_INCREMENT,
    `name`         VARCHAR(45) NOT NULL,
    `price`        INT         NOT NULL,
    `expired_date` DATETIME    NOT NULL,
    `product_type` ENUM('Sayuran', 'Protein', 'Buah', 'Snack') NOT NULL,
    PRIMARY KEY (`id`)
);