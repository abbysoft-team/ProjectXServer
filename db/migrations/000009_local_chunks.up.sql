ALTER TABLE chunks
ADD COLUMN number integer NOT NULL DEFAULT 0,
DROP CONSTRAINT IF EXISTS chunks_x_y_key;

CREATE UNIQUE INDEX chunks_id ON chunks (x, y, number);
