CREATE TABLE IF NOT EXISTS `taskUser` (
  `taskId` INT UNSIGNED NOT NULL,
  `userId` INT UNSIGNED NOT NULL,
  
  PRIMARY KEY (`taskId`, `userId`),
  FOREIGN KEY (`taskId`) REFERENCES tasks(`taskId`),
  FOREIGN KEY (`userId`) REFERENCES users(`userId`)
);