CREATE TABLE IF NOT EXISTS `species` (
  `speciesId` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `comName` VARCHAR(255) NOT NULL,
  `sciName` VARCHAR(255) NOT NULL,
  `speciesDesc` VARCHAR(1500) NOT NULL,
  `image` VARCHAR(255) NOT NULL,
  `habitatId` INT UNSIGNED NOT NULL,
  `baskTemp` VARCHAR(255) NOT NULL,
  `diet` VARCHAR(1500) NOT NULL,
  `sociality` VARCHAR(1500) NOT NULL,
  `extraCare` VARCHAR(1500) NOT NULL,
  
  PRIMARY KEY (`speciesId`),
  FOREIGN KEY (`habitatId`) REFERENCES habitats(`habitatId`)
);