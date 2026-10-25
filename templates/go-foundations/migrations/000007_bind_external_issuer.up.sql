-- Existing external users must be explicitly bound to their trusted issuer before login.
ALTER TABLE users ADD COLUMN external_issuer TEXT NOT NULL DEFAULT '';
DROP INDEX idx_users_external_identity;
CREATE UNIQUE INDEX idx_users_external_identity
ON users(auth_provider, external_issuer, external_subject)
WHERE external_subject IS NOT NULL;
