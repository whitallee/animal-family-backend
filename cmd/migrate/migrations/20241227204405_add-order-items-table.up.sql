CREATE TABLE IF NOT EXISTS `order_items` (
  `orderItemsId` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `orderId` INT UNSIGNED NOT NULL,
  `productId` INT UNSIGNED NOT NULL,
  `quantity` INT NOT NULL,
  `price` DECIMAL(10, 2) NOT NULL,
  
  PRIMARY KEY (`orderItemsId`),
  FOREIGN KEY (`orderId`) REFERENCES orders(`orderId`),
  FOREIGN KEY (`productId`) REFERENCES products(`productId`)
);