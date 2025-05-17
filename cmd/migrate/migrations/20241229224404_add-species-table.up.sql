CREATE TABLE IF NOT EXISTS "species" (
    "speciesId" SERIAL PRIMARY KEY,
    "comName" VARCHAR(255) NOT NULL,
    "sciName" VARCHAR(255) NOT NULL,
    "image" VARCHAR(255) NOT NULL,
    "speciesDesc" VARCHAR(1500) NOT NULL,
    "habitatId" INTEGER NOT NULL,
    "baskTemp" VARCHAR(255) NOT NULL,
    "diet" VARCHAR(1500) NOT NULL,
    "sociality" VARCHAR(1500) NOT NULL,
    "lifespan" VARCHAR(255) NOT NULL,
    "size" VARCHAR(255) NOT NULL,
    "weight" VARCHAR(255) NOT NULL,
    "conservationStatus" VARCHAR(255) NOT NULL,
    "extraCare" VARCHAR(1500) NOT NULL,
    
    FOREIGN KEY ("habitatId") REFERENCES habitats("habitatId")
);