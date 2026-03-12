-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fn_cascade_delete_user()
    RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM stress_samples WHERE user_id = OLD.id;
    DELETE FROM stress_scores WHERE user_id = OLD.id;

    DELETE FROM mood_entries WHERE user_id = OLD.id;

    DELETE FROM appointments WHERE veteran_id = OLD.id;
    DELETE FROM availability_slots WHERE provider_id = OLD.id;

    DELETE FROM support_invites WHERE sender_id = OLD.id OR receiver_id = OLD.id;
    DELETE FROM user_supports WHERE veteran_id = OLD.id OR support_id = OLD.id;

    DELETE FROM authorization_rules WHERE owner_id = OLD.id OR viewer_id = OLD.id;

    DELETE FROM user_addresses WHERE user_id = OLD.id;
    DELETE FROM user_devices WHERE user_id = OLD.id;

    DELETE FROM messages WHERE sender_id = OLD.id;
    DELETE FROM conversation_participants WHERE user_id = OLD.id;
    DELETE FROM conversations
    WHERE id NOT IN (SELECT DISTINCT conversation_id FROM conversation_participants);

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_cascade_delete_user
    BEFORE DELETE ON users
    FOR EACH ROW
EXECUTE FUNCTION fn_cascade_delete_user();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_cascade_delete_user ON users;
DROP FUNCTION IF EXISTS fn_cascade_delete_user();
-- +goose StatementEnd

