CREATE TABLE ram_spec (
                          id INT AUTO_INCREMENT PRIMARY KEY,
                          type VARCHAR(32) NOT NULL UNIQUE
);

INSERT INTO ram_spec (type) VALUES ('DDR3'), ('DDR4');