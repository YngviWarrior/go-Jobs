DELIMITER $

DROP PROCEDURE IF EXISTS `setStatusFinishGame`;
CREATE PROCEDURE `setStatusFinishGame`(
    IN `in_game_id` BIGINT
)

BEGIN
    START TRANSACTION;
        SET @status_id = 4;

        -- SELECT `id`, `id_moedas_pares`, `game_id_type_time`, 4 as `game_status`
        SELECT `id`, `id_moedas_pares`, `game_id_status`, `game_id_type_time`, `game_date_start`, `game_date_process`, `game_date_finish`, 
            `game_profit_percent`, `game_bet_amount_dolar_total`, `bonus_trader_total_amount_dolar`, `bonus_indication_total_amount_dolar`, 
            `game_win_amount_percent`, `game_win_amount_dolar`, `game_lose_amount_percent`, `game_lose_amount_dolar`, `game_equal_amount_percent`, 
            `game_equal_amount_dolar`, `game_price_amount_selected_finish`, `company_tax_amount_dolar_from_game_tax`, `game_calculated`, 
            `partition_filter`, `total_win_dolar_trader_bot`, `total_lose_dolar_trader_bot`, `id_trader_bot_action`
        INTO
            @g_id, @g_id_moedas_pares, @g_game_id_status, @g_game_id_type_time, @g_game_date_start, @g_game_date_process, @g_game_date_finish, 
            @g_game_profit_percent, @g_game_bet_amount_dolar_total, @g_bonus_trader_total_amount_dolar, @g_bonus_indication_total_amount_dolar, 
            @g_game_win_amount_percent, @g_game_win_amount_dolar, @g_game_lose_amount_percent, @g_game_lose_amount_dolar, @g_game_equal_amount_percent, 
            @g_game_equal_amount_dolar, @g_game_price_amount_selected_finish, @g_company_tax_amount_dolar_from_game_tax, @g_game_calculated, 
            @g_partition_filter, @g_total_win_dolar_trader_bot, @g_total_lose_dolar_trader_bot, @g_id_trader_bot_action
        FROM `binary_option_game`
        WHERE `id` = `in_game_id`
        FOR UPDATE;

        UPDATE `binary_option_game`
        SET `game_id_status` = `in_status_id` 
        WHERE `id` = `in_game_id`;
        
       IF @g_id IS NOT NULL THEN
            CALL searchBestPriceForEndGame(@g_id);
       END IF;

    COMMIT;
END $

DELIMITER ;

CALL setStatusFinishGame(20008);
