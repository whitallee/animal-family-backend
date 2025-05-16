CREATE TABLE IF NOT EXISTS "animalUser" (
    "animalId" INTEGER NOT NULL,
    "userId" INTEGER NOT NULL,
    
    PRIMARY KEY ("animalId", "userId"),
    FOREIGN KEY ("animalId") REFERENCES animals("animalId"),
    FOREIGN KEY ("userId") REFERENCES users("userId")
);