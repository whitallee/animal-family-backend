CREATE TABLE IF NOT EXISTS animals (
    "animalId" SERIAL PRIMARY KEY,
    "animalName" VARCHAR(255) NOT NULL,
    "image" VARCHAR(255) NOT NULL,
    "gender" VARCHAR(255) NOT NULL,
    "dob" DATE NOT NULL,
    "personalityDesc" VARCHAR(1500) NOT NULL,
    "dietDesc" VARCHAR(1500) NOT NULL,
    "routineDesc" VARCHAR(1500) NOT NULL,
    "extraNotes" VARCHAR(1500) NOT NULL,
    "speciesId" INTEGER NOT NULL,
    "enclosureId" INTEGER,
    
    FOREIGN KEY ("speciesId") REFERENCES species("speciesId"),
    FOREIGN KEY ("enclosureId") REFERENCES enclosures("enclosureId")
);