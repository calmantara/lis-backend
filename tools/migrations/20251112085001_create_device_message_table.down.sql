-- Drop indexes
DROP INDEX IF EXISTS idx_device_message_type_code;
DROP INDEX IF EXISTS idx_users_company_ididx_device_message_device_id;

-- Drop table
DROP TABLE IF EXISTS device_messages;
