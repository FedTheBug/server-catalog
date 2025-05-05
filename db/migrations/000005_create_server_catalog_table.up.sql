CREATE TABLE server_catalog (
                                id INT AUTO_INCREMENT PRIMARY KEY,
                                model VARCHAR(128) NOT NULL,
                                ram_size INT NOT NULL,
                                ram_type INT NOT NULL,
                                hdd_size INT NOT NULL,
                                hdd_count INT NOT NULL,
                                hdd_type INT NOT NULL,
                                location VARCHAR(128) NOT NULL,
                                price decimal(20,2) unsigned NOT NULL,
                                currency INT NOT NULL,
                                FOREIGN KEY (ram_type) REFERENCES ram_spec(id),
                                FOREIGN KEY (hdd_type) REFERENCES hdd_spec(id),
                                FOREIGN KEY (currency) REFERENCES currency(id)
);