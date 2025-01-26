CREATE TABLE public.pemilih (
    pemilih_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nik VARCHAR(16) UNIQUE NOT NULL,
    nama VARCHAR(255) NOT NULL,
    tanggal_lahir DATE NOT NULL,
    alamat VARCHAR(255),
    sidik_jari TEXT,
    status_aktif BOOLEAN DEFAULT TRUE,
    has_vote BOOLEAN DEFAULT FALSE,
    provinsi VARCHAR(255),
    kota VARCHAR(255),
    kecamatan VARCHAR(255),
    tps VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
