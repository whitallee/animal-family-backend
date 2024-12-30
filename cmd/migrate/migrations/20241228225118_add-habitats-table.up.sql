CREATE TABLE IF NOT EXISTS `habitats` (
  `habitatId` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `habitatName` VARCHAR(255) NOT NULL,
  `habitatDesc` VARCHAR(1500) NOT NULL,
  `image` VARCHAR(255) NOT NULL,
  `humidity` VARCHAR(255) NOT NULL,
  `dayTempRange` VARCHAR(255) NOT NULL,
  `nightTempRange` VARCHAR(255) NOT NULL,
  
  PRIMARY KEY (`habitatId`)
);