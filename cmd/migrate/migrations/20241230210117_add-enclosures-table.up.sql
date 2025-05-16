CREATE TABLE IF NOT EXISTS enclosures (
    "enclosureId" SERIAL PRIMARY KEY,
    "enclosureName" VARCHAR(255) NOT NULL,
    "image" VARCHAR(255) NOT NULL,
    "notes" VARCHAR(1500) NOT NULL,
    "habitatId" INTEGER NOT NULL,
    
    FOREIGN KEY ("habitatId") REFERENCES habitats("habitatId")
);