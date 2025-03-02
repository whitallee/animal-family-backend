CREATE TABLE IF NOT EXISTS `taskSubject` (
  `taskId` INT UNSIGNED NOT NULL,
  `animalId` INT UNSIGNED NOT NULL,
  `enclosureId` INT UNSIGNED NOT NULL,
  
  PRIMARY KEY (`taskId`, `animalId`, `enclosureId`),
  FOREIGN KEY (`taskId`) REFERENCES tasks(`taskId`),
  FOREIGN KEY (`animalId`) REFERENCES animals(`animalId`),
  FOREIGN KEY (`enclosureId`) REFERENCES enclosures(`enclosureId`)
);