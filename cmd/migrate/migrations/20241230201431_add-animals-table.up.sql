CREATE TABLE IF NOT EXISTS `animals` (
  `animalId` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `animalName` VARCHAR(255) NOT NULL,
  `image` VARCHAR(255) NOT NULL,
  `notes` VARCHAR(1500) NOT NULL,
  `speciesID` INT UNSIGNED NOT NULL,
  
  PRIMARY KEY (`animalId`),
  FOREIGN KEY (`speciesId`) REFERENCES species(`speciesId`)
);