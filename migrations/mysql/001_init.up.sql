create table `employees`
(
    id     BIGINT UNSIGNED auto_increment
        primary key,
    name   VARCHAR(255)      not null comment 'ПІП працівника з 1С',
    card   VARCHAR(20)       not null comment 'Код перепустки',
    active TINYINT DEFAULT 1 null comment 'Чи активний працівник'
) charset = utf8;

CREATE TABLE `doctors`
(
    `id`          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `employee_id` BIGINT UNSIGNED UNIQUE NOT NULL,
    CONSTRAINT fk_doctors_employees FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`)
) charset = utf8;

CREATE TABLE `health_checks`
(
    `id`          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `doctor_id`   BIGINT UNSIGNED,
    `employee_id` BIGINT UNSIGNED,
    `when`        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `until`       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `conclusion`  TINYINT NOT NULL,
    CONSTRAINT fk_health_checks_doctors FOREIGN KEY (`doctor_id`) REFERENCES `doctors` (`id`),
    CONSTRAINT fk_health_checks_employees FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE
) charset = utf8;

CREATE TABLE `groups`
(
    `id`   BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL COMMENT 'Ім''я групи доступу'
    -- `departmentId` VARCHAR(9) default '000000000' NOT NULL COMMENT 'Id підрозділу з 1С'
) charset = utf8;

CREATE TABLE `members`
(
    `employee_id` BIGINT UNSIGNED,
    `group_id`    BIGINT UNSIGNED,
    PRIMARY KEY (`employee_id`, `group_id`),
    CONSTRAINT fk_members_employees FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`) ON DELETE CASCADE,
    CONSTRAINT fk_members_groups FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`)
) charset = utf8;

CREATE TABLE `readers`
(
    `id`         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'Номер зчитувача',
    `ip_address` VARCHAR(39)
) charset = utf8;

CREATE TABLE `access_nodes`
(
    `id`              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `parent_id`       BIGINT UNSIGNED,
    `name`            VARCHAR(255),
    `entrance_reader` BIGINT UNSIGNED NOT NULL,
    `exit_reader`     BIGINT UNSIGNED,
    CONSTRAINT fk_access_nodes_parent FOREIGN KEY (`parent_id`) REFERENCES `access_nodes` (`id`),
    CONSTRAINT fk_access_nodes_entrance FOREIGN KEY (`entrance_reader`) REFERENCES `readers` (`id`),
    CONSTRAINT fk_access_nodes_exit FOREIGN KEY (`exit_reader`) REFERENCES `readers` (`id`)
) charset = utf8;

CREATE TABLE `permissions`
(
    `group_id`       BIGINT UNSIGNED,
    `node_id`        BIGINT UNSIGNED,
    `health_check`   TINYINT NOT NULL,
    `sanitary_check` TINYINT NOT NULL,
    PRIMARY KEY (`group_id`, `node_id`),
    CONSTRAINT fk_permissions_groups FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`) ON DELETE CASCADE,
    CONSTRAINT fk_permissions_access_nodes FOREIGN KEY (`node_id`) REFERENCES `access_nodes` (`id`)
) charset = utf8;

CREATE TABLE `access_log`
(
    `id`          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `employee_id` BIGINT UNSIGNED,
    `node_id`     BIGINT UNSIGNED,
    `access`      TINYINT NOT NULL,
    `exited`      TINYINT NOT NULL,
    `access_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_access_log_employees FOREIGN KEY (`employee_id`) REFERENCES `employees` (`id`),
    CONSTRAINT fk_access_log_access_nodes FOREIGN KEY (`node_id`) REFERENCES `access_nodes` (`id`)
) charset = utf8;
