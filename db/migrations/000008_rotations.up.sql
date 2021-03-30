ALTER TABLE towns
ADD COLUMN rotation float4 NOT NULL DEFAULT 0;

ALTER TABLE town_buildings
ADD COLUMN rotation float4 NOT NULL DEFAULT 0;
