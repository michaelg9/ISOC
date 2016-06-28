-- MySQL dump 10.13  Distrib 5.5.44, for debian-linux-gnu (armv7l)
--
-- Host: localhost    Database: mobile_data
-- ------------------------------------------------------
-- Server version	5.5.44-0+deb8u1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `BatteryStatus`
--

DROP TABLE IF EXISTS `BatteryStatus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `BatteryStatus` (
  `id` int(11) NOT NULL DEFAULT '0',
  `batteryPercentage` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `BatteryStatus_ibfk_1` FOREIGN KEY (`id`) REFERENCES `Data` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `BatteryStatus`
--

LOCK TABLES `BatteryStatus` WRITE;
/*!40000 ALTER TABLE `BatteryStatus` DISABLE KEYS */;
INSERT INTO `BatteryStatus` VALUES (1,70),(2,69),(3,68),(4,67),(6,41),(7,40);
/*!40000 ALTER TABLE `BatteryStatus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Data`
--

DROP TABLE IF EXISTS `Data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Data` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `device` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `device` (`device`),
  CONSTRAINT `Data_ibfk_1` FOREIGN KEY (`device`) REFERENCES `Device` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Data`
--

LOCK TABLES `Data` WRITE;
/*!40000 ALTER TABLE `Data` DISABLE KEYS */;
INSERT INTO `Data` VALUES (1,'2016-05-31 11:48:48',1),(2,'2016-05-31 11:55:36',1),(3,'2016-05-31 11:57:15',1),(4,'2016-05-31 11:58:30',1),(5,'2016-05-30 19:20:34',1),(6,'2016-05-30 19:20:34',1),(7,'2016-05-30 19:21:55',1);
/*!40000 ALTER TABLE `Data` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Device`
--

DROP TABLE IF EXISTS `Device`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Device` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `manufacturer` varchar(50) DEFAULT NULL,
  `modelName` varchar(50) DEFAULT NULL,
  `osVersion` varchar(50) DEFAULT NULL,
  `user` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user` (`user`),
  CONSTRAINT `Device_ibfk_1` FOREIGN KEY (`user`) REFERENCES `User` (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Device`
--

LOCK TABLES `Device` WRITE;
/*!40000 ALTER TABLE `Device` DISABLE KEYS */;
INSERT INTO `Device` VALUES (1,'Motorola','Moto X (2nd Generation)','Android 5.0',1);
/*!40000 ALTER TABLE `Device` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `User`
--

DROP TABLE IF EXISTS `User`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `User` (
  `uid` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(20) NOT NULL,
  `passwordHash` char(64) NOT NULL,
  `apiKey` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`uid`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `apiKey` (`apiKey`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `User`
--

LOCK TABLES `User` WRITE;
/*!40000 ALTER TABLE `User` DISABLE KEYS */;
INSERT INTO `User` VALUES (1,'user@usermail.com','$2a$10$539nT.CNbxpyyqrL9mro3OQEKuAjhTD3UjEa8JYPbZMZEM/HizvxK','37e72ff927f511e688adb827ebf7e157');
/*!40000 ALTER TABLE `User` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-06-02 16:20:38
