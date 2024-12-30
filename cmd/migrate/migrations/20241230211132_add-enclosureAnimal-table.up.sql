CREATE TABLE IF NOT EXISTS `enclosureAnimal` (
  `enclosureId` INT UNSIGNED NOT NULL,
  `animalID` INT UNSIGNED NOT NULL,
  
  PRIMARY KEY (`enclosureId`, `animalId`),
  FOREIGN KEY (`enclosureId`) REFERENCES enclosures(`enclosureId`),
  FOREIGN KEY (`animalId`) REFERENCES animals(`animalId`)
);