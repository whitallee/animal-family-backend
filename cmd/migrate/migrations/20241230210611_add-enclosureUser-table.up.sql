CREATE TABLE IF NOT EXISTS `enclosureUser` (
  `enclosureId` INT UNSIGNED NOT NULL,
  `userId` INT UNSIGNED NOT NULL,
  
  PRIMARY KEY (`enclosureId`, `userId`),
  FOREIGN KEY (`enclosureId`) REFERENCES enclosures(`enclosureId`),
  FOREIGN KEY (`userId`) REFERENCES users(`userId`)
);