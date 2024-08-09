CREATE TABLE IF NOT EXISTS flats (
    id SERIAL PRIMARY KEY,
    house_id INT NOT NULL,
    flat_number INT NOT NULL,
    price INT NOT NULL,
    rooms INT NOT NULL,
    status TEXT NOT NULL,
    FOREIGN KEY (house_id) REFERENCES houses (id),
    UNIQUE (house_id, flat_number)
);