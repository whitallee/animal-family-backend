CREATE TABLE IF NOT EXISTS `animalUser` (
  `animalId` INT UNSIGNED NOT NULL,
  `userId` INT UNSIGNED NOT NULL,
  
  PRIMARY KEY (`animalId`, `userId`),
  FOREIGN KEY (`animalId`) REFERENCES animals(`animalId`),
  FOREIGN KEY (`userId`) REFERENCES users(`userId`)
);