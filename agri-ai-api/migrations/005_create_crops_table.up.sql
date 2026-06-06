CREATE TABLE IF NOT EXISTS crops (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    min_temp NUMERIC(5,2) NOT NULL,
    max_temp NUMERIC(5,2) NOT NULL,
    min_precipitation NUMERIC(6,2) NOT NULL,
    max_precipitation NUMERIC(6,2) NOT NULL,
    ideal_season VARCHAR(50) NOT NULL
);

-- Seed Data (Carga Inicial)
INSERT INTO crops (name, min_temp, max_temp, min_precipitation, max_precipitation, ideal_season) VALUES
('Milho', 21.00, 32.00, 500.00, 800.00, 'summer'),
('Soja', 20.00, 30.00, 450.00, 700.00, 'summer'),
('Trigo', 15.00, 24.00, 350.00, 600.00, 'winter'),
('Café', 18.00, 28.00, 1200.00, 2000.00, 'spring'),
('Feijão', 18.00, 30.00, 300.00, 400.00, 'autumn')
ON CONFLICT (name) DO NOTHING;
