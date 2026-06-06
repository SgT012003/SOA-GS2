-- Index para buscar o usuário rapidamente pelo e-mail
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Index para o cache de clima (buscas frequentes por lat/lon/data)
CREATE INDEX IF NOT EXISTS idx_weather_cache_lat_lon_date ON weather_cache(latitude, longitude, query_date);
