-- Seed data for habitats
INSERT INTO "habitats" ("habitatName", "habitatDesc", "image", "humidity", "dayTempRange", "nightTempRange") VALUES
('No Habitat', 'Default habitat for unassigned species.', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/habitats/no-habitat.png', 'N/A', 'N/A', 'N/A'),
('Desert', 'Arid environment with high temperatures and low humidity. Features sandy terrain and sparse vegetation.', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/habitats/desert.png', '20-30%', '85-95°F', '65-75°F'),
('Tropical Rainforest', 'Humid environment with dense vegetation and high rainfall. Features lush greenery and diverse plant life.', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/habitats/tropical-rain-forest.png', '70-90%', '75-85°F', '65-75°F'),
('Grassland', 'Open environment with moderate temperatures and seasonal rainfall. Features vast plains of grasses and scattered trees.', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/habitats/grassland.png', '40-60%', '70-80°F', '60-70°F'),
('Temperate Forest', 'Moderate climate with seasonal changes. Features mixed vegetation and moderate humidity.', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/habitats/temperate-forest.png', '50-70%', '65-75°F', '55-65°F'),
('Aquatic', 'Water-based environment requiring specific water parameters and filtration.', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/habitats/aquatic.png', 'N/A', '75-82°F', '72-78°F'),
('Semi-Arid', 'Moderately dry environment with some humidity. Features mixed terrain and vegetation.', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/habitats/semi-arid.png', '30-50%', '75-85°F', '65-75°F'),
('Indoor', 'Climate-controlled indoor environment suitable for domestic pets.', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/habitats/indoor.png', '40-60%', '68-75°F', '65-72°F');

-- Seed data for species
INSERT INTO "species" ("comName", "sciName", "image", "speciesDesc", "habitatId", "baskTemp", "diet", "sociality", "lifespan", "size", "weight", "conservationStatus", "extraCare") VALUES
('No Species', 'N/A', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/no-species.png', 'Default species for unassigned entries.', 1, 'N/A', 'N/A', 'N/A', 'N/A', 'N/A', 'N/A', 'N/A', 'N/A'),
('Bearded Dragon', 'Pogona vitticeps', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/bearded-dragon.png', 'A popular reptile pet known for its docile nature and distinctive spiky scales.', 2, '95-105°F', 'Omnivorous diet including insects and vegetables', 'Generally solitary but can be kept in pairs', '10-15 years', '16-24 inches', '10-18 ounces', 'Least Concern', 'Requires UVB lighting and proper temperature gradient'),
('Green Tree Python', 'Morelia viridis', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/green-tree-python.png', 'A stunning arboreal snake with vibrant green coloration.', 3, '85-90°F', 'Carnivorous diet of small mammals and birds', 'Solitary', '15-20 years', '4-6 feet', '2-4 pounds', 'Least Concern', 'Needs high humidity and vertical climbing space'),
('Leopard Gecko', 'Eublepharis macularius', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/leopard-gecko.png', 'A small, nocturnal lizard with distinctive spotted pattern.', 2, '88-92°F', 'Insectivorous diet', 'Solitary', '15-20 years', '7-10 inches', '2-3 ounces', 'Least Concern', 'No UVB required, but needs heat source'),
('Ferret', 'Mustela putorius furo', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/ferret.png', 'Playful and intelligent mustelid with high energy levels.', 8, 'N/A', 'Carnivorous diet of high-quality ferret food', 'Social, should be kept in pairs or groups', '6-10 years', '18-24 inches', '1-4 pounds', 'Domesticated', 'Requires daily exercise and mental stimulation'),
('White''s Dumpy Tree Frog', 'Litoria caerulea', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/whites-dumpy-tree-frog.png', 'Large, docile tree frog with distinctive green coloration.', 3, 'N/A', 'Insectivorous diet', 'Can be kept in groups', '15-20 years', '3-4 inches', '2-4 ounces', 'Least Concern', 'Needs high humidity and vertical space'),
('Ball Python', 'Python regius', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/ball-python.png', 'Docile snake known for its defensive balling behavior.', 7, '88-92°F', 'Carnivorous diet of rodents', 'Solitary', '20-30 years', '3-5 feet', '2-5 pounds', 'Least Concern', 'Requires proper humidity and temperature gradient'),
('Veiled Chameleon', 'Chamaeleo calyptratus', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/veiled-chameleon.png', 'Colorful arboreal lizard with distinctive casque.', 7, '85-95°F', 'Insectivorous diet with occasional plant matter', 'Solitary', '5-8 years', '14-24 inches', '3-6 ounces', 'Least Concern', 'Needs UVB lighting and misting system'),
('African Fat-Tailed Gecko', 'Hemitheconyx caudicinctus', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/african-fat-tailed-gecko.png', 'Nocturnal gecko with distinctive fat tail.', 7, '88-92°F', 'Insectivorous diet', 'Solitary', '15-20 years', '7-9 inches', '2-3 ounces', 'Least Concern', 'No UVB required, but needs heat source'),
('Betta Fish', 'Betta splendens', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/betta-fish.png', 'Colorful freshwater fish known for its flowing fins.', 6, 'N/A', 'Carnivorous diet of pellets and live/frozen food', 'Solitary', '3-5 years', '2-3 inches', '0.1-0.2 ounces', 'Least Concern', 'Requires clean water and proper filtration'),
('Blue Tiger Polar Parrot Cichlid', 'Hybrid', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/blue-tiger-polar-parrot-cichlid.png', 'Colorful hybrid cichlid with distinctive blue coloration.', 6, 'N/A', 'Omnivorous diet of pellets and live/frozen food', 'Can be kept in groups', '8-10 years', '6-8 inches', '8-12 ounces', 'Hybrid', 'Requires large tank and proper water parameters'),
('Red-Eared Slider Turtle', 'Trachemys scripta elegans', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/red-eared-slider-turtle.png', 'Popular aquatic turtle with distinctive red ear markings.', 6, '85-95°F', 'Omnivorous diet of pellets, vegetables, and protein', 'Can be kept in groups', '20-30 years', '8-12 inches', '1-2 pounds', 'Least Concern', 'Needs basking area and UVB lighting'),
('Plecostomus', 'Hypostomus plecostomus', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/plecostomus.png', 'Popular algae-eating catfish.', 6, 'N/A', 'Omnivorous diet including algae and sinking pellets', 'Can be kept in groups', '10-15 years', '12-24 inches', '1-2 pounds', 'Least Concern', 'Requires driftwood and hiding places'),
('Brazilian Rainbow Boa', 'Epicrates cenchria', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/brazilian-rainbow-boa.png', 'Beautiful snake with iridescent scales.', 3, '80-85°F', 'Carnivorous diet of rodents', 'Solitary', '20-25 years', '5-7 feet', '3-5 pounds', 'Least Concern', 'Needs high humidity and proper temperature gradient'),
('Pacman Frog', 'Ceratophrys ornata', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/pacman-frog.png', 'Round, colorful frog with large mouth.', 3, 'N/A', 'Carnivorous diet of insects and small rodents', 'Solitary', '10-15 years', '4-6 inches', '4-8 ounces', 'Least Concern', 'Needs high humidity and shallow water dish'),
('Crested Gecko', 'Correlophus ciliatus', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/crested-gecko.png', 'Arboreal gecko with distinctive crest and eyelashes.', 3, '72-80°F', 'Omnivorous diet of prepared food and insects', 'Can be kept in pairs', '15-20 years', '7-9 inches', '1-2 ounces', 'Least Concern', 'No heat source required, but needs misting'),
('Chiweenie', 'Canis lupus familiaris', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/chiweenie.png', 'Small mixed breed dog combining Chihuahua and Dachshund traits.', 8, 'N/A', 'High-quality dog food', 'Social', '12-16 years', '8-12 inches', '5-12 pounds', 'Domesticated', 'Regular exercise and dental care needed'),
('Cat', 'Felis catus', 'https://raw.githubusercontent.com/whitallee/brindle-assets/main/species/cat.png', 'Lovable animals with a wide range of personalities.', 8, 'N/A', 'High-quality cat food', 'Social', '15 years', '8-12 inches', '5-12 pounds', 'Domesticated', 'Regular exercise and dental care needed');

-- Seed data for users
INSERT INTO "users" ("firstName", "lastName", "email", "password", "createdAt") VALUES
('Whit', 'Allee', 'whitallee@gmail.com', '$2a$10$cWGN86Vlpqci8.MzxraFZuiC413QbD/5V0z7X8U6uG2CcPIPf7.H.', NOW());

-- Seed data for enclosures
INSERT INTO "enclosures" ("enclosureName", "habitatId", "image", "notes") VALUES
('Ferret Cage', 8, '', '2-tier ferret cage with lots of hammocks'),
('Turtle Tank', 6, '', '55 gallon aquatic tank with basking area'),
('Bedroom Tank', 6, '', '20 gallon fish tank'),
('Franny''s Tank', 3, '', '40 gallon tank'),
('Kiwi''s Tank', 7, '', '20 gallon semi-arid setup for African Fat-Tailed Gecko'),
('Chococat''s Tank', 7, '', '40 gallon Ball Python enclosure with hides'),
('Jellybean''s Tank', 7, '', '1''4"x1''4"x2''6" arboreal setup for Veiled Chameleon'),
('Rosalina''s Tank', 3, '', '20 gallon humid setup for Pacman Frog'),
('Guava''s Tank', 3, '', 'Large Exo Terra arboreal setup for Crested Gecko'),
('Dumpy Tank', 3, '', 'Humid setup for White''s Dumpy Tree Frogs'),
('Betta Sorority Tank', 6, '', 'Community tank for female bettas');

-- Seed data for enclosure ownership
INSERT INTO "enclosureUser" ("enclosureId", "userId") VALUES
(1, 1), (2, 1), (3, 1), (4, 1), (5, 1), (6, 1), (7, 1), (8, 1), (9, 1), (10, 1), (11, 1);

-- Seed data for animals
INSERT INTO "animals" ("animalName", "speciesId", "enclosureId", "image", "gender", "dob", "personalityDesc", "dietDesc", "routineDesc", "extraNotes") VALUES
('Blueberry', 17, NULL, '', 'Male', '2020-01-01', 'Sweet and playful', 'High-quality dog food twice daily', 'Daily walks and playtime', 'Loves belly rubs'),
('Queso', 18, NULL, '', 'Male', '2020-01-01', 'Sweet and playful', 'High-quality cat food twice daily', 'Daily playtime and litter box cleaning', 'Loves cuddles'),
('Eevee', 5, 1, '', 'Male', '2021-06-15', 'Energetic and curious', 'Ferret kibble and treats', 'Daily playtime and cage cleaning', 'Loves tunnels'),
('Cinnamaroll', 5, 1, '', 'Male', '2021-06-15', 'Playful and social', 'Ferret kibble and treats', 'Daily playtime and cage cleaning', 'Loves hammocks'),
('Strawberry Milk', 5, 1, '', 'Male', '2021-06-15', 'Adventurous and bold', 'Ferret kibble and treats', 'Daily playtime and cage cleaning', 'Loves exploring'),
('Winston', 5, 1, '', 'Male', '2021-06-15', 'Gentle and friendly', 'Ferret kibble and treats', 'Daily playtime and cage cleaning', 'Loves cuddles'),
('Wendy', 12, 2, '', 'Female', '2019-05-20', 'Active and curious', 'Turtle pellets and vegetables', 'Daily feeding and tank maintenance', 'Loves basking'),
('Iggy', 12, 2, '', 'Male', '2019-05-20', 'Shy but friendly', 'Turtle pellets and vegetables', 'Daily feeding and tank maintenance', 'Loves swimming'),
('Abe DeCatfish', 13, 3, '', 'Male', '2022-03-10', 'Peaceful and nocturnal', 'Algae wafers and sinking pellets', 'Weekly water changes', 'Loves hiding spots'),
('Blue Tiger Polar Parrot Cichlids x4', 11, 3, '', 'Mixed', '2022-03-10', 'Active and social', 'Cichlid pellets and frozen food', 'Weekly water changes', 'Group of 4 fish'),
('Francesca', 14, 4, '', 'Female', '2020-08-15', 'Calm and gentle', 'Frozen/thawed rodents', 'Weekly feeding and enclosure maintenance', 'Loves climbing'),
('Kiwi', 9, 5, '', 'Female', '2021-04-20', 'Shy but curious', 'Crickets and mealworms', 'Weekly feeding and spot cleaning', 'Loves warm hides'),
('Chococat', 7, 6, '', 'Male', '2020-11-30', 'Docile and calm', 'Frozen/thawed rodents', 'Weekly feeding and enclosure maintenance', 'Loves tight spaces'),
('Jellybean', 8, 7, '', 'Female', '2021-07-15', 'Active and alert', 'Crickets and dubia roaches', 'Daily misting and feeding', 'Loves climbing'),
('Princess Rosalina', 15, 8, '', 'Female', '2021-09-10', 'Bold and voracious', 'Crickets and pinky mice', 'Weekly feeding and substrate change', 'Loves burrowing'),
('Guava', 16, 9, '', 'Female', '2021-10-05', 'Active at night', 'Crested gecko diet and insects', 'Daily misting and weekly feeding', 'Loves jumping'),
('Kuromi', 6, 10, '', 'Female', '2022-01-15', 'Active', 'Crickets and dubia roaches', 'Daily misting and feeding', 'Loves climbing'),
('Kerropi II', 6, 10, '', 'Male', '2022-01-15', 'Shy but vocal', 'Crickets and dubia roaches', 'Daily misting and feeding', 'Loves hiding'),
('Rainy', 10, 11, '', 'Female', '2022-06-01', 'Peaceful and graceful', 'Betta pellets and frozen food', 'Weekly water changes', 'Loves swimming'),
('Misty', 10, 11, '', 'Female', '2022-06-01', 'Active and curious', 'Betta pellets and frozen food', 'Weekly water changes', 'Loves exploring'),
('Stormy', 10, 11, '', 'Female', '2022-06-01', 'Big and strong', 'Betta pellets and frozen food', 'Weekly water changes', 'Loves eating');

-- Seed data for animal ownership
INSERT INTO "animalUser" ("animalId", "userId") VALUES
(1, 1), (2, 1), (3, 1), (4, 1), (5, 1), (6, 1), (7, 1), (8, 1), (9, 1), (10, 1),
(11, 1), (12, 1), (13, 1), (14, 1), (15, 1), (16, 1), (17, 1), (18, 1), (19, 1),
(20, 1), (21, 1);

-- Seed data for tasks
INSERT INTO "tasks" ("taskName", "taskDesc", "complete", "lastCompleted", "repeatIntervHours") VALUES
('Feed Blueberry', 'Feed Blueberry high-quality dog food twice daily', FALSE, NOW(), 12),
('Brush teeth', 'Brush Blueberry''s teeth every 3 days', FALSE, NOW(), 24*3),
('Dremel nails', 'Dremel Blueberry''s nails every 30 days', FALSE, NOW(), 24*30),
('Give flea meds', 'Give Blueberry flea meds every 30 days', FALSE, NOW(), 24*30),
('Bathe Blueberry', 'Bathe Blueberry every 2 months', FALSE, NOW(), 24*30*2),
('Feed the ferrets', 'Feed the ferrets high-quality ferret food twice daily', FALSE, NOW(), 12),
('Ferret Potty Pads', 'Change the ferret potty pads every other day', FALSE, NOW(), 24*2),
('Ferret Litter', 'Clean the ferret litter boxes every week', FALSE, NOW(), 24*7),
('Feed Turtle', 'Feed Turtle turtle pellets and vegetables every 3 days', FALSE, NOW(), 24*3),
('Turtle Water Change', 'Do a water change every 2 weeks', FALSE, NOW(), 24*14),
('Turtle UVB Bulb', 'Change the UVB bulb every 6 months', FALSE, NOW(), 24*30*6),
('Feed Abe', 'Feed algae wafer once a week', FALSE, NOW(), 24*7),
('Feed cichlids', 'Feed the cichlids cichlid pellets or frozen food every 3 days', FALSE, NOW(), 24*3),
('Bedroom Tank Water Change', 'Do a water change every month', FALSE, NOW(), 24*30),
('Feed Franny', 'Feed Franny frozen/thawed rodents every 2 weeks', FALSE, NOW(), 24*14),
('Spotclean Franny''s tank', 'Spotclean Franny''s tank every 2 months', FALSE, NOW(), 24*30*2),
('Feed Kiwi', 'Feed Kiwi crickets or mealworms every 4 days', FALSE, NOW(), 24*4),
('Spotclean Kiwi''s tank', 'Spotclean Kiwi''s tank every 2 months', FALSE, NOW(), 24*30*2),
('Feed Chococat', 'Feed Chococat frozen/thawed rodents every 2 weeks', FALSE, NOW(), 24*14),
('Spotclean Chococat''s tank', 'Spotclean Chococat''s tank every 2 months', FALSE, NOW(), 24*30*2),
('Feed Jellybean', 'Feed Jellybean crickets and dubia roaches every 3 days', FALSE, NOW(), 24*3),
('Spotclean Jellybean''s tank', 'Spotclean Jellybean''s tank every 2 months', FALSE, NOW(), 24*30*2),
('Feed Princess Rosalina', 'Feed Princess Rosalina crickets or worms or dubias or pinky mice every 4 days', FALSE, NOW(), 24*4),
('Spotclean Rosalina''s tank', 'Spotclean Princess Rosalina''s tank every 2 months', FALSE, NOW(), 24*30*2),
('Feed Guava', 'Feed Guava Crested Gecko diet and insects', FALSE, NOW(), 24*3),
('Spotclean Guava''s tank', 'Spotclean Guava''s tank every 2 months', FALSE, NOW(), 24*30*2),
('Feed Kuromi and Kerropi', 'Feed Kuromi and Kerropi crickets or dubia roaches every 3 days', FALSE, NOW(), 24*3),
('Spotclean Kuromi and Kerropi''s tank', 'Spotclean Kuromi and Kerropi''s tank every 2 months', FALSE, NOW(), 24*30*2),
('Feed Bettas', 'Feed Rainy, Misty, and Stormy Betta pellets or frozen food every 3 days', FALSE, NOW(), 24*3),
('Betta Water Change', 'Do a water change every month', FALSE, NOW(), 24*30);

-- Seed data for taskSubject
INSERT INTO "taskSubject" ("taskId", "animalId", "enclosureId") VALUES
(1, 1, NULL), (2, 1, NULL), (3, 1, NULL), (4, 1, NULL), (5, 1, NULL),
(6, NULL, 1), (7, NULL, 1), (8, NULL, 1),
(9, NULL, 2), (10, NULL, 2), (11, NULL, 2),
(12, 9, NULL), (13, 10, NULL), (14, NULL, 3),
(15, 11, NULL), (16, NULL, 4),
(17, 12, NULL), (18, NULL, 5),
(19, 13, NULL), (20, NULL, 6),
(21, 14, NULL), (22, NULL, 7),
(23, 15, NULL), (24, NULL, 8),
(25, 16, NULL), (26, NULL, 9),
(27, NULL, 10), (28, NULL, 10),
(29, NULL, 11), (30, NULL, 11);

-- Seed data for task ownership
INSERT INTO "taskUser" ("taskId", "userId") VALUES
(1, 1), (2, 1), (3, 1), (4, 1), (5, 1), (6, 1), (7, 1), (8, 1), (9, 1), (10, 1),
(11, 1), (12, 1), (13, 1), (14, 1), (15, 1), (16, 1), (17, 1), (18, 1), (19, 1),
(20, 1), (21, 1), (22, 1), (23, 1), (24, 1), (25, 1), (26, 1), (27, 1), (28, 1),
(29, 1), (30, 1);