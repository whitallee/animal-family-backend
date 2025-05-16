CREATE TABLE IF NOT EXISTS "taskUser" (
    "taskId" INTEGER NOT NULL,
    "userId" INTEGER NOT NULL,
    
    PRIMARY KEY ("taskId", "userId"),
    FOREIGN KEY ("taskId") REFERENCES tasks("taskId"),
    FOREIGN KEY ("userId") REFERENCES users("userId")
);