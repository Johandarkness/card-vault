
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE INDEX IF NOT EXISTS idx_cards_user_id ON cards(user_id);
-- CREATE INDEX IF NOT EXISTS idx_cards_created_at ON cards(created_at);
-- CREATE INDEX IF NOT EXISTS idx_cards_is_active ON cards(is_active);

-- Create audit table for PCI compliance
CREATE TABLE IF NOT EXISTS card_audit (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    card_id UUID,
    user_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    details JSONB
);

CREATE INDEX IF NOT EXISTS idx_audit_card_id ON card_audit(card_id);
CREATE INDEX IF NOT EXISTS idx_audit_user_id ON card_audit(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON card_audit(timestamp);