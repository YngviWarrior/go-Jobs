 DELIMITER $

DROP PROCEDURE IF EXISTS `searchBestPriceForEndGame`;
CREATE PROCEDURE `searchBestPriceForEndGame`(
    IN `in_game_id` BIGINT
)

BEGIN
    START TRANSACTION;
        
        SET @S = 1;
        LOOP
            SELECT `id`, `hash_id`, `id_game`, `id_usuario`, `id_choice`, `id_balance`, `bet_amount_dolar`, `amount_win_dolar`, `price_amount_selected`, 
                `status_received_win_payment`, `id_trader_follower`, `bot_use_status`, `date_register`, `bonus_trader_percent_from_tax_bet_win`, 
                `bonus_indication_percent_from_tax_bet_win`, `status_received_refund_payment`, `refund`, `id_balance` 
            INTO @b_has, @b_id, @b_hash_id, @b_id_game, @b_id_usuario, @b_id_choice, @b_id_balance, @b_bet_amount_dolar, @b_amount_win_dolar, @b_price_amount_selected,
                @b_status_received_win_payment, @b_id_trader_follower, @b_bot_use_status, @b_date_register, @b_bonus_trader_percent_from_tax_bet_win,
                @b_bonus_indication_percent_from_tax_bet_win, @b_status_received_refund_payment, @b_refund, @b_id_balance
            FROM binary_option_game_bet
            WHERE id_game = `in_game_id`
            LIMIT @S, 1;

            LEAVE
                -- DROP TEMPORARY TABLE IF EXISTS `temp_search_best_price_for_end_game`;
                CREATE TEMPORARY TABLE IF NOT EXISTS `temp_search_best_price_for_end_game` (
                "id" BIGINT NOT NULL PRIMARY_KEY AUTO_INCREMENT,
                "id_game" BIGINT NOT NULL,
                "price" DECIMAL(60,8) NOT NULL,
                "total_win_dolar" DECIMAL(60,8) DEFAULT 0,
                "total_lose_dolar" DECIMAL(60,8) DEFAULT 0,
                "total_dif_win_lose" DECIMAL(60,8) DEFAULT 0,
                "total_equal_dolar" DECIMAL(60,8) DEFAULT 0,
                "list_players_win" BIGINT DEFAULT 0,
                "list_players_lose" BIGINT DEFAULT 0,
                "list_players_equal" BIGINT DEFAULT 0,
                "total_lose_dolar_trader_bot" DECIMAL(60,8) DEFAULT 0,
                "total_win_dolar_trader_bot" DECIMAL(60,8) DEFAULT 0
            );

            @price := SELECT `close` FROM `candle_1m` ORDER BY `mts` DESC LIMIT 1;
            INSERT INTO `temp_search_best_price_for_end_game`(`id_game`, `price`) VALUE (@price, `in_game_id`);

            IF @b_has = 0 THEN
            ELSE

            END IF;
            SET @S = @S + 1;
        END LOOP

        

    COMMIT;
END $

DELIMITER ;
