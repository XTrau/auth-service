CREATE TABLE blocked_tokens(
    token text NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_blocked_tokens_token ON blocked_tokens USING HASH (token);
CREATE INDEX idx_blocked_tokens_expires_at ON blocked_tokens USING BTREE (expires_at);
