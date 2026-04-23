CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Prerrequisitos mínimos para que las FK de services tengan sentido
CREATE TABLE users (
  id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
  email         VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role          VARCHAR(20)  NOT NULL,
  created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE providers (
  user_id    UUID          PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  nombre     VARCHAR(100)  NOT NULL,
  apellido   VARCHAR(100)  NOT NULL,
  telefono   VARCHAR(20)
);

-- Las 3 tablas del módulo catalog
CREATE TABLE categories (
  id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
  nombre     VARCHAR(100) NOT NULL,
  slug       VARCHAR(100) NOT NULL UNIQUE,
  parent_id  UUID         REFERENCES categories(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE services (
  id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
  provider_id UUID          NOT NULL REFERENCES providers(user_id) ON DELETE CASCADE,
  category_id UUID          NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
  titulo      VARCHAR(200)  NOT NULL,
  descripcion TEXT,
  precio_base NUMERIC(10,2) NOT NULL CHECK (precio_base >= 0),
  lat         NUMERIC(9,6),      -- latitud  (-90 a 90)
  lng         NUMERIC(9,6),      -- longitud (-180 a 180)
  radio_km    NUMERIC(5,2),
  is_active   BOOLEAN       NOT NULL DEFAULT true,
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TABLE portfolio_items (
  id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
  service_id  UUID         NOT NULL REFERENCES services(id) ON DELETE CASCADE,
  storage_url TEXT         NOT NULL,
  titulo      VARCHAR(200),
  orden       INT          NOT NULL DEFAULT 0
);

-- Seed para probar services (1 user + 1 provider)
INSERT INTO users (id, email, password_hash, role)
VALUES ('00000000-0000-0000-0000-000000000001', 'seed@test.com', 'x', 'provider');
INSERT INTO providers (user_id, nombre, apellido, telefono)
VALUES ('00000000-0000-0000-0000-000000000001', 'Juan', 'Pérez', '+57300');
