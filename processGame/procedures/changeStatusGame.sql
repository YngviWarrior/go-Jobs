DELIMITER $

DROP PROCEDURE IF EXISTS `changeStatusGame`;
CREATE PROCEDURE `changeStatusGame`(
    IN `in_status_id` INT
)

BEGIN
    START TRANSACTION;
        -- SET SESSION group_concat_max_len = 18446744073709551615;
        DROP TEMPORARY TABLE IF EXISTS `temp_process_game`;
        CREATE TEMPORARY TABLE `temp_process_game` (
            `id` BIGINT,
            `id_moedas_pares` INT,
            `game_id_type_time` INT,
            `game_status` INT
        );

        SET @temp_status_id = `in_status_id` - 1;

        INSERT INTO `temp_process_game`
        SELECT `id`, `id_moedas_pares`, `game_id_type_time`, 2 as `game_status`
        FROM `binary_option_game`
        WHERE `game_id_status` = @temp_status_id
        AND `game_date_start` <= DATE_ADD(NOW(), INTERVAL 1 SECOND)
        FOR UPDATE;

        UPDATE `binary_option_game` b
        SET b.`game_id_status` = `in_status_id` 
        WHERE b.`id` IN (
            SELECT `id` FROM `temp_process_game`
        );

    -- O OBJETIVO É APENAS O PROCESSAMENTO FINISH

        IF `in_status_id` = 4 THEN
            SELECT NOW();
        END IF;
        
    COMMIT;
END $

DELIMITER ;

CALL changeStatusGame(2);
