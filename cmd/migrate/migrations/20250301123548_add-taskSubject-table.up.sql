CREATE TABLE IF NOT EXISTS "taskSubject" (
    "taskSubjectId" SERIAL PRIMARY KEY,
    "taskId" INTEGER NOT NULL,
    "animalId" INTEGER,
    "enclosureId" INTEGER,
    
    FOREIGN KEY ("taskId") REFERENCES tasks("taskId"),
    FOREIGN KEY ("animalId") REFERENCES animals("animalId"),
    FOREIGN KEY ("enclosureId") REFERENCES enclosures("enclosureId")
);