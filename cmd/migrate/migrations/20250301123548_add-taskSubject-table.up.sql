CREATE TABLE IF NOT EXISTS "taskSubject" (
    "taskId" INTEGER NOT NULL,
    "animalId" INTEGER NOT NULL,
    "enclosureId" INTEGER NOT NULL,
    
    PRIMARY KEY ("taskId", "animalId", "enclosureId"),
    FOREIGN KEY ("taskId") REFERENCES tasks("taskId"),
    FOREIGN KEY ("animalId") REFERENCES animals("animalId"),
    FOREIGN KEY ("enclosureId") REFERENCES enclosures("enclosureId")
);