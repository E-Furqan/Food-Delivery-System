CREATE TABLE IF NOT EXISTS Configuration (
    ConfigId INT PRIMARY KEY,
    ClientId VARCHAR(255) NOT NULL,
    ClientSecret VARCHAR(255) NOT NULL,
    TokenUri VARCHAR(255) NOT NULL,
    AuthUri VARCHAR(255) NOT NULL,
    RedirectUris TEXT NOT NULL,
    AuthProviderCertUrl VARCHAR(255) NOT NULL
);