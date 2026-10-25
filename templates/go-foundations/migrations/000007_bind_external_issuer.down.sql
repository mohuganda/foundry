-- Fails safely if issuer-specific accounts would collide under the old key.
DROP INDEX idx_users_external_identity;
CREATE UNIQUE INDEX idx_users_external_identity
ON users(auth_provider, external_subject)
WHERE external_subject IS NOT NULL;
ALTER TABLE users DROP COLUMN external_issuer;
