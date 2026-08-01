-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gateway_services (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    base_url VARCHAR(500) NOT NULL,
    base_path VARCHAR(200) NOT NULL UNIQUE,
    protocol VARCHAR(20) NOT NULL DEFAULT 'http',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    rate_limit_per_minute INT UNSIGNED NULL,
    health_status VARCHAR(20) NOT NULL DEFAULT 'unknown',
    health_checked_at DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    INDEX idx_gateway_services_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS gateway_routes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    service_id BIGINT UNSIGNED NOT NULL,
    method VARCHAR(10) NOT NULL,
    path_pattern VARCHAR(500) NOT NULL,
    permission_match_mode VARCHAR(10) NOT NULL DEFAULT 'any',
    permissions JSON NULL,
    public BOOLEAN NOT NULL DEFAULT FALSE,
    rate_limit_per_minute INT UNSIGNED NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    INDEX idx_gateway_routes_service_id (service_id),
    INDEX idx_gateway_routes_deleted_at (deleted_at),
    UNIQUE KEY uq_gateway_routes_service_method_path (service_id, method, path_pattern),
    FOREIGN KEY (service_id) REFERENCES gateway_services(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Singleton row (id=1) tracking Service/Route CUD count — lets RouteManager's periodic
-- refresh skip rebuilding when nothing changed. A monotonic counter (not a timestamp) so
-- comparisons never depend on clock sync across app instances.
CREATE TABLE IF NOT EXISTS gateway_cache_meta (
    id BIGINT UNSIGNED PRIMARY KEY,
    version BIGINT UNSIGNED NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO gateway_cache_meta (id, version) VALUES (1, 1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gateway_cache_meta;
DROP TABLE IF EXISTS gateway_routes;
DROP TABLE IF EXISTS gateway_services;
-- +goose StatementEnd
