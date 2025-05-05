CREATE TABLE currency (
                          id INT AUTO_INCREMENT PRIMARY KEY,
                          type VARCHAR(16) NOT NULL,
                          symbol VARCHAR(8) NOT NULL
);

INSERT INTO currency (type, symbol) VALUES
                                        ('USD', '$'),
                                        ('Euro', '€');