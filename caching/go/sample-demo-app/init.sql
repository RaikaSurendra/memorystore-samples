-- Copyright 2025 Google LLC
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    price DOUBLE PRECISION NOT NULL
);

-- Insert sample data
INSERT INTO items (name, description, price) VALUES 
('Vintage Typewriter', 'A classic 1950s typewriter in excellent condition.', 150.00),
('Leather Journal', 'Handcrafted leather journal with 200 pages.', 45.50),
('Ceramic Mug', 'Artisanal ceramic mug with a unique glaze.', 25.00),
('Fountain Pen', 'Luxury fountain pen with gold nib.', 120.00),
('Desk Lamp', 'Modern LED desk lamp with adjustable brightness.', 55.99);
