-- Seed data for habitats
INSERT INTO "habitats" ("habitatName", "habitatDesc", "image", "humidity", "dayTempRange", "nightTempRange") VALUES
('No Habitat', 'Default habitat for unassigned species.', 'no_habitat.jpg', 'N/A', 'N/A', 'N/A'),
('Desert', 'Arid environment with high temperatures and low humidity. Features sandy terrain and sparse vegetation.', 'desert.jpg', '20-30%', '85-95°F', '65-75°F'),
('Tropical Rainforest', 'Humid environment with dense vegetation and high rainfall. Features lush greenery and diverse plant life.', 'rainforest.jpg', '70-90%', '75-85°F', '65-75°F'),
('Grassland', 'Open environment with moderate temperatures and seasonal rainfall. Features vast plains of grasses and scattered trees.', 'grassland.jpg', '40-60%', '70-80°F', '60-70°F'),
('Temperate Forest', 'Moderate climate with seasonal changes. Features mixed vegetation and moderate humidity.', 'temperate_forest.jpg', '50-70%', '65-75°F', '55-65°F'),
('Aquatic', 'Water-based environment requiring specific water parameters and filtration.', 'aquatic.jpg', 'N/A', '75-82°F', '72-78°F'),
('Semi-Arid', 'Moderately dry environment with some humidity. Features mixed terrain and vegetation.', 'semi_arid.jpg', '30-50%', '75-85°F', '65-75°F'),
('Indoor', 'Climate-controlled indoor environment suitable for domestic pets.', 'indoor.jpg', '40-60%', '68-75°F', '65-72°F');

-- Seed data for species
INSERT INTO "species" ("comName", "sciName", "image", "speciesDesc", "habitatId", "baskTemp", "diet", "sociality", "lifespan", "size", "weight", "conservationStatus", "extraCare") VALUES
('No Species', 'N/A', 'no_species.jpg', 'Default species for unassigned entries.', 1, 'N/A', 'N/A', 'N/A', 'N/A', 'N/A', 'N/A', 'N/A', 'N/A'),
('Bearded Dragon', 'Pogona vitticeps', 'bearded_dragon.jpg', 'A popular reptile pet known for its docile nature and distinctive spiky scales.', 2, '95-105°F', 'Omnivorous diet including insects and vegetables', 'Generally solitary but can be kept in pairs', '10-15 years', '16-24 inches', '10-18 ounces', 'Least Concern', 'Requires UVB lighting and proper temperature gradient'),
('Green Tree Python', 'Morelia viridis', 'green_tree_python.jpg', 'A stunning arboreal snake with vibrant green coloration.', 3, '85-90°F', 'Carnivorous diet of small mammals and birds', 'Solitary', '15-20 years', '4-6 feet', '2-4 pounds', 'Least Concern', 'Needs high humidity and vertical climbing space'),
('Leopard Gecko', 'Eublepharis macularius', 'leopard_gecko.jpg', 'A small, nocturnal lizard with distinctive spotted pattern.', 2, '88-92°F', 'Insectivorous diet', 'Solitary', '15-20 years', '7-10 inches', '2-3 ounces', 'Least Concern', 'No UVB required, but needs heat source'),
('Ferret', 'Mustela putorius furo', 'ferret.jpg', 'Playful and intelligent mustelid with high energy levels.', 8, 'N/A', 'Carnivorous diet of high-quality ferret food', 'Social, should be kept in pairs or groups', '6-10 years', '18-24 inches', '1-4 pounds', 'Domesticated', 'Requires daily exercise and mental stimulation'),
('White''s Dumpy Tree Frog', 'Litoria caerulea', 'whites_tree_frog.jpg', 'Large, docile tree frog with distinctive green coloration.', 3, 'N/A', 'Insectivorous diet', 'Can be kept in groups', '15-20 years', '3-4 inches', '2-4 ounces', 'Least Concern', 'Needs high humidity and vertical space'),
('Ball Python', 'Python regius', 'ball_python.jpg', 'Docile snake known for its defensive balling behavior.', 7, '88-92°F', 'Carnivorous diet of rodents', 'Solitary', '20-30 years', '3-5 feet', '2-5 pounds', 'Least Concern', 'Requires proper humidity and temperature gradient'),
('Veiled Chameleon', 'Chamaeleo calyptratus', 'veiled_chameleon.jpg', 'Colorful arboreal lizard with distinctive casque.', 7, '85-95°F', 'Insectivorous diet with occasional plant matter', 'Solitary', '5-8 years', '14-24 inches', '3-6 ounces', 'Least Concern', 'Needs UVB lighting and misting system'),
('African Fat-Tailed Gecko', 'Hemitheconyx caudicinctus', 'fat_tailed_gecko.jpg', 'Nocturnal gecko with distinctive fat tail.', 7, '88-92°F', 'Insectivorous diet', 'Solitary', '15-20 years', '7-9 inches', '2-3 ounces', 'Least Concern', 'No UVB required, but needs heat source'),
('Betta Fish', 'Betta splendens', 'betta_fish.jpg', 'Colorful freshwater fish known for its flowing fins.', 6, 'N/A', 'Carnivorous diet of pellets and live/frozen food', 'Solitary', '3-5 years', '2-3 inches', '0.1-0.2 ounces', 'Least Concern', 'Requires clean water and proper filtration'),
('Blue Tiger Polar Parrot Cichlid', 'Hybrid', 'polar_parrot_cichlid.jpg', 'Colorful hybrid cichlid with distinctive blue coloration.', 6, 'N/A', 'Omnivorous diet of pellets and live/frozen food', 'Can be kept in groups', '8-10 years', '6-8 inches', '8-12 ounces', 'Hybrid', 'Requires large tank and proper water parameters'),
('Red-Eared Slider Turtle', 'Trachemys scripta elegans', 'red_eared_slider.jpg', 'Popular aquatic turtle with distinctive red ear markings.', 6, '85-95°F', 'Omnivorous diet of pellets, vegetables, and protein', 'Can be kept in groups', '20-30 years', '8-12 inches', '1-2 pounds', 'Least Concern', 'Needs basking area and UVB lighting'),
('Plecostomus', 'Hypostomus plecostomus', 'plecostomus.jpg', 'Popular algae-eating catfish.', 6, 'N/A', 'Omnivorous diet including algae and sinking pellets', 'Can be kept in groups', '10-15 years', '12-24 inches', '1-2 pounds', 'Least Concern', 'Requires driftwood and hiding places'),
('Brazilian Rainbow Boa', 'Epicrates cenchria', 'rainbow_boa.jpg', 'Beautiful snake with iridescent scales.', 3, '80-85°F', 'Carnivorous diet of rodents', 'Solitary', '20-25 years', '5-7 feet', '3-5 pounds', 'Least Concern', 'Needs high humidity and proper temperature gradient'),
('Pacman Frog', 'Ceratophrys ornata', 'pacman_frog.jpg', 'Round, colorful frog with large mouth.', 3, 'N/A', 'Carnivorous diet of insects and small rodents', 'Solitary', '10-15 years', '4-6 inches', '4-8 ounces', 'Least Concern', 'Needs high humidity and shallow water dish'),
('Crested Gecko', 'Correlophus ciliatus', 'crested_gecko.jpg', 'Arboreal gecko with distinctive crest and eyelashes.', 3, '72-80°F', 'Omnivorous diet of prepared food and insects', 'Can be kept in pairs', '15-20 years', '7-9 inches', '1-2 ounces', 'Least Concern', 'No heat source required, but needs misting'),
('Chiweenie', 'Canis lupus familiaris', 'chiweenie.jpg', 'Small mixed breed dog combining Chihuahua and Dachshund traits.', 8, 'N/A', 'High-quality dog food', 'Social', '12-16 years', '8-12 inches', '5-12 pounds', 'Domesticated', 'Regular exercise and dental care needed'),
('Elephant Ear Betta', 'Betta splendens', 'elephant_ear_betta.jpg', 'Betta fish with distinctive large, flowing fins.', 6, 'N/A', 'Carnivorous diet of pellets and live/frozen food', 'Solitary', '3-5 years', '2-3 inches', '0.1-0.2 ounces', 'Least Concern', 'Requires clean water and proper filtration');

-- Seed data for users
INSERT INTO "users" ("firstName", "lastName", "email", "password", "createdAt") VALUES
('Admin', 'User', 'admin@animalfamily.com', '$2a$10$2X6wPck6z4qzer3M4l1FROXjGqLKkCpRjowLWQciIFWiAVnw5dDk', NOW()),
('Whit', 'Allee', 'whitallee@gmail.com', '$2a$10$2X6wPck6z4qzer3M4l1FROXjGqLKkCpRjowLWQciIFWiAVnw5dDk', NOW());

-- Seed data for enclosures
INSERT INTO "enclosures" ("enclosureName", "habitatId", "image", "notes") VALUES
('Ferret Cage', 8, 'ferret_cage.jpg', '2-tier ferret cage with lots ofhammocks'),
('Turtle Tank', 6, 'turtle_tank.jpg', '55 gallon aquatic tank with basking area'),
('Bedroom Tank', 6, 'bedroom_tank.jpg', '20 gallon fish tank'),
('Franny''s Tank', 3, 'franny_tank.jpg', '40 gallon tank'),
('Kiwi''s Tank', 7, 'kiwi_tank.jpg', '20 gallon semi-arid setup for African Fat-Tailed Gecko'),
('Chococat''s Tank', 7, 'chococat_tank.jpg', '40 gallon Ball Python enclosure with hides'),
('Jellybean''s Tank', 7, 'jellybean_tank.jpg', '1''4"x1''4"x2''6" arboreal setup for Veiled Chameleon'),
('Rosalina''s Tank', 3, 'rosalina_tank.jpg', '20 gallon humid setup for Pacman Frog'),
('Guava''s Tank', 3, 'guava_tank.jpg', 'Large Exo Terra arboreal setup for Crested Gecko'),
('Dumpy Tank', 3, 'dumpy_tank.jpg', 'Humid setup for White''s Dumpy Tree Frogs'),
('Betta Sorority Tank', 6, 'betta_sorority.jpg', 'Community tank for female bettas');

-- Seed data for enclosure ownership
INSERT INTO "enclosureUser" ("enclosureId", "userId") VALUES
(1, 2), (2, 2), (3, 2), (4, 2), (5, 2), (6, 2), (7, 2), (8, 2), (9, 2), (10, 2), (11, 2);

-- Seed data for animals
INSERT INTO "animals" ("animalName", "speciesId", "enclosureId", "image", "gender", "dob", "personalityDesc", "dietDesc", "routineDesc", "extraNotes") VALUES
('Blueberry', 17, NULL, 'blueberry.jpg', 'Female', '2020-01-01', 'Sweet and playful', 'High-quality dog food twice daily', 'Daily walks and playtime', 'Loves belly rubs'),
('Eevee', 5, 1, 'eevee.jpg', 'Female', '2021-06-15', 'Energetic and curious', 'Ferret kibble and treats', 'Daily playtime and cage cleaning', 'Loves tunnels'),
('Cinnamaroll', 5, 1, 'cinnamaroll.jpg', 'Female', '2021-06-15', 'Playful and social', 'Ferret kibble and treats', 'Daily playtime and cage cleaning', 'Loves hammocks'),
('Strawberry Milk', 5, 1, 'strawberry_milk.jpg', 'Female', '2021-06-15', 'Adventurous and bold', 'Ferret kibble and treats', 'Daily playtime and cage cleaning', 'Loves exploring'),
('Winston', 5, 1, 'winston.jpg', 'Male', '2021-06-15', 'Gentle and friendly', 'Ferret kibble and treats', 'Daily playtime and cage cleaning', 'Loves cuddles'),
('Wendy', 12, 2, 'wendy.jpg', 'Female', '2019-05-20', 'Active and curious', 'Turtle pellets and vegetables', 'Daily feeding and tank maintenance', 'Loves basking'),
('Iggy', 12, 2, 'iggy.jpg', 'Male', '2019-05-20', 'Shy but friendly', 'Turtle pellets and vegetables', 'Daily feeding and tank maintenance', 'Loves swimming'),
('Abe DeCatfish', 13, 3, 'abe.jpg', 'Male', '2022-03-10', 'Peaceful and nocturnal', 'Algae wafers and sinking pellets', 'Weekly water changes', 'Loves hiding spots'),
('Blue Tiger Polar Parrot Cichlids x4', 11, 3, 'btpc.jpg', 'Mixed', '2022-03-10', 'Active and social', 'Cichlid pellets and frozen food', 'Weekly water changes', 'Group of 4 fish'),
('Francesca', 14, 4, 'francesca.jpg', 'Female', '2020-08-15', 'Calm and gentle', 'Frozen/thawed rodents', 'Weekly feeding and enclosure maintenance', 'Loves climbing'),
('Kiwi', 9, 5, 'kiwi.jpg', 'Female', '2021-04-20', 'Shy but curious', 'Crickets and mealworms', 'Weekly feeding and spot cleaning', 'Loves warm hides'),
('Chococat', 7, 6, 'chococat.jpg', 'Male', '2020-11-30', 'Docile and calm', 'Frozen/thawed rodents', 'Weekly feeding and enclosure maintenance', 'Loves tight spaces'),
('Jellybean', 8, 7, 'jellybean.jpg', 'Male', '2021-07-15', 'Active and alert', 'Crickets and dubia roaches', 'Daily misting and feeding', 'Loves climbing'),
('Princess Rosalina', 15, 8, 'rosalina.jpg', 'Female', '2021-09-10', 'Bold and voracious', 'Crickets and pinky mice', 'Weekly feeding and substrate change', 'Loves burrowing'),
('Guava', 16, 9, 'guava.jpg', 'Female', '2021-10-05', 'Active at night', 'Crested gecko diet and insects', 'Daily misting and weekly feeding', 'Loves jumping'),
('Kuromi', 6, 10, 'kuromi.jpg', 'Female', '2022-01-15', 'Active and vocal', 'Crickets and dubia roaches', 'Daily misting and feeding', 'Loves climbing'),
('Kerropi II', 6, 10, 'kerropi.jpg', 'Male', '2022-01-15', 'Shy but active', 'Crickets and dubia roaches', 'Daily misting and feeding', 'Loves hiding'),
('Rainy', 10, 11, 'rainy.jpg', 'Female', '2022-06-01', 'Peaceful and graceful', 'Betta pellets and frozen food', 'Weekly water changes', 'Loves swimming'),
('Misty', 10, 11, 'misty.jpg', 'Female', '2022-06-01', 'Active and curious', 'Betta pellets and frozen food', 'Weekly water changes', 'Loves exploring');

-- Seed data for animal ownership
INSERT INTO "animalUser" ("animalId", "userId") VALUES
(1, 2), (2, 2), (3, 2), (4, 2), (5, 2), (6, 2), (7, 2), (8, 2), (9, 2), (10, 2), (11, 2), (12, 2), (13, 2), (14, 2), (15, 2), (16, 2), (17, 2), (18, 2), (19, 2);
