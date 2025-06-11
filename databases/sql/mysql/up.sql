USE `dbdev`;

DROP TABLE IF EXISTS `notifications`;
DROP TABLE IF EXISTS `tasks`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `role` enum('manager','technician') NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
LOCK TABLES `users` WRITE;
INSERT INTO `users` VALUES (1,'John Smith','john.smith@company.com','$2a$10$dummyhash1','manager','2025-06-06 18:29:04','2025-06-06 18:29:04'),(2,'Sarah Johnson','sarah.j@company.com','$2a$10$dummyhash2','technician','2025-06-06 18:29:04','2025-06-06 18:29:04'),(3,'Mike Wilson','mike.w@company.com','$2a$10$dummyhash3','technician','2025-06-06 18:29:04','2025-06-06 18:29:04');
UNLOCK TABLES;

CREATE TABLE `tasks` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `technician_id` bigint NOT NULL,
  `title` varchar(255) NOT NULL,
  `summary` text NOT NULL,
  `performed_at` timestamp NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `technician_id` (`technician_id`),
  CONSTRAINT `tasks_ibfk_1` FOREIGN KEY (`technician_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `notifications` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL,
  `message` text NOT NULL,
  `is_read` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `task_id` (`task_id`),
  CONSTRAINT `notifications_ibfk_1` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Insert some tasks for technicians
INSERT INTO `tasks` (`technician_id`, `title`, `summary`, `performed_at`) VALUES
(2, 'Server Maintenance', 'Regular server maintenance and updates', '2024-03-20 14:30:00'),
(2, 'Network Configuration', 'Updated network settings for new office', '2024-03-21 09:15:00'),
(3, 'Hardware Installation', 'Installed new workstations in Marketing', '2024-03-22 11:45:00'),
(3, 'Software Update', 'Updated security software across all systems', '2024-03-23 16:20:00');

-- Insert notifications based on the tasks
INSERT INTO `notifications` (`task_id`, `message`, `is_read`) VALUES
(1, 'The tech Sarah Johnson performed the task on 2024-03-20 14:30:00', 0),
(2, 'The tech Sarah Johnson performed the task on 2024-03-21 09:15:00', 0),
(3, 'The tech Mike Wilson performed the task on 2024-03-22 11:45:00', 0),
(4, 'The tech Mike Wilson performed the task on 2024-03-23 16:20:00', 0); 