
 SET NAMES utf8mb4 ;

--
-- Table structure for table `changepasses`
--

DROP TABLE IF EXISTS `changepasses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `changepasses` (
  `id` varchar(36) NOT NULL,
  `userid` varchar(36) NOT NULL,
  `password` varchar(200) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `changepasses`
--

LOCK TABLES `changepasses` WRITE;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `roles` (
  `id` varchar(36) NOT NULL,
  `name` varchar(50) DEFAULT NULL,
  `status` smallint(6) DEFAULT NULL,
  `active` tinyint(4) DEFAULT NULL,
  `description` varchar(100) DEFAULT NULL,
  `created` timestamp NULL DEFAULT NULL,
  `updated` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
INSERT INTO `roles` VALUES ('24ff8948-e472-4c36-ac3f-eb88580a361b','USER',0,0,'General User',NULL,NULL),('da0512a2-857f-11eb-b0ad-02bcbb6fc696','ROOT',1,0,'Root/System Administrator',NULL,NULL),('f7899c8f-857f-11eb-b0ad-02bcbb6fc696','ADMIN',1,1,'Site Admin',NULL,NULL);
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `users` (
  `id` varchar(36) NOT NULL,
  `username` varchar(50) DEFAULT NULL,
  `email` varchar(150) DEFAULT NULL,
  `password` varchar(200) DEFAULT NULL,
  `status` smallint(6) DEFAULT NULL,
  `active` tinyint(4) DEFAULT NULL,
  `roles` mediumtext,
  `profile_id` varchar(100) DEFAULT NULL,
  `is_root` tinyint(4) DEFAULT NULL,
  `profile` json DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `default_role` varchar(10) NOT NULL,
  `updated` timestamp NULL DEFAULT NULL,
  `is_sys` tinyint(4) DEFAULT NULL,
  `is_lockout` tinyint(4) DEFAULT NULL,
  `is_pass_force_reset` tinyint(1) DEFAULT NULL,
  `lockout_start` timestamp NULL DEFAULT NULL,
  `failpasscount` smallint(6) DEFAULT NULL,
  `lastfailpass` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
INSERT INTO `users` VALUES ('5a1b489e-b978-11ea-b0ad-02bcbb6fc696','admin','admin@example.com','$2a$10$35lMU/xSah2EJrqn9iCEoeih7KRbZBubJUbKc73R.lwWx3g2Jxe4q',1,1,'ROOT,ADMIN,USER',NULL,1,'{\"avatar\": \"./images/robi_logo.png\", \"designation\": \"Administrator\", \"display_name\": \"Administrator\"}','2020-06-01 05:40:23','ADMIN','2021-04-24 02:15:01',0,0,0,NULL,0,'2021-04-24 03:44:26'),('9c186ba3-925d-11eb-b0ad-02bcbb6fc696','superadmin','root@example.com','$2a$10$TwmdbngtKJKSwWd1bBmTOutDxHFVNE.WtygdrAbx9BNYNQUygNUnO',1,1,'ROOT,ADMIN,USER',NULL,1,'{\"avatar\": \"./images/robi_logo.png\", \"designation\": \"Administrator\", \"display_name\": \"Administrator\"}','2020-06-01 05:40:23','ROOT','2022-12-30 14:43:54',1,0,0,NULL,0,NULL),('e1319bbc-1788-40f3-bafc-af59ed3cee70','user','user@example.com','$2a$10$NpKFsIMymKe/Dl/tWVhQoOsx5O6N9y8Vxl/I4v57UrVO1K3t/yamy',1,0,'USER','',0,'{\"avatar\": \"\", \"designation\": \"ADMIN\", \"display_name\": \"test@te.co\"}','2021-04-24 10:37:38','USER',NULL,0,0,0,NULL,0,NULL);
UNLOCK TABLES;

