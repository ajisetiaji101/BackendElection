CREATE TABLE public.kandidats (
    kandidat_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gubernur_name VARCHAR(255) NOT NULL,
    wakil_gubernur_name VARCHAR(255) NOT NULL,
    partai VARCHAR(255),
    visi TEXT,
    foto TEXT,
    foto_partai TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
