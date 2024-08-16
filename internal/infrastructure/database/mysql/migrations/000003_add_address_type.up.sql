ALTER TABLE addresses
ADD COLUMN type VARCHAR(255) NOT NULL,
ADD CONSTRAINT unique_user_address_type UNIQUE (user_id, `type`);
;