CREATE TABLE hdd_spec (
                          id INT AUTO_INCREMENT PRIMARY KEY,
                          type VARCHAR(32) NOT NULL UNIQUE
);

INSERT INTO hdd_spec (type) VALUES ('SATA2'), ('SAS'), ('SSD');