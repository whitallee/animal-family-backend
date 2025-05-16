CREATE TABLE IF NOT EXISTS "enclosureUser" (
    "enclosureId" INTEGER NOT NULL,
    "userId" INTEGER NOT NULL,
    
    PRIMARY KEY ("enclosureId", "userId"),
    FOREIGN KEY ("enclosureId") REFERENCES enclosures("enclosureId"),
    FOREIGN KEY ("userId") REFERENCES users("userId")
);