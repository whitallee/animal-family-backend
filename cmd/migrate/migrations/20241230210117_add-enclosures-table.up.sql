CREATE TABLE IF NOT EXISTS `enclosures` (
  `enclosureId` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `enclosureName` VARCHAR(255) NOT NULL,
  `image` VARCHAR(255) NOT NULL,
  `notes` VARCHAR(1500) NOT NULL,
  `habitatId` INT UNSIGNED NOT NULL,
  
  PRIMARY KEY (`enclosureId`),
  FOREIGN KEY (`habitatId`) REFERENCES habitats(`habitatId`)
);