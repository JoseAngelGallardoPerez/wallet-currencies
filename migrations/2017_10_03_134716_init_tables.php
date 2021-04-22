<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class InitTables extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
    }

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        // skip the migration if there are another migrations
        // It means this migration was already applied
        $migrations = DB::select('SELECT * FROM migrations LIMIT 1');
        if (!empty($migrations)) {
            return;
        }
        $oldMigrationTable = DB::select("SHOW TABLES LIKE 'schema_migrations'");
        if (!empty($oldMigrationTable)) {
            return;
        }

        DB::beginTransaction();

        try {
            app("db")->getPdo()->exec($this->getSql());
        } catch (\Throwable $e) {
            DB::rollBack();
            throw $e;
        }

        DB::commit();
    }

    private function getSql()
    {
        return <<<SQL
            CREATE TABLE `currencies` (
              `id` int(10) UNSIGNED NOT NULL,
              `code` varchar(64) DEFAULT NULL,
              `decimal_places` tinyint(3) UNSIGNED NOT NULL,
              `active` tinyint(1) NOT NULL DEFAULT '0',
              `created_at` timestamp NULL DEFAULT NULL,
              `updated_at` timestamp NULL DEFAULT NULL,
              `type` varchar(50) NOT NULL DEFAULT 'fiat',
              `feed` varchar(64) DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            INSERT INTO `currencies` (`id`, `code`, `decimal_places`, `active`, `created_at`, `updated_at`, `type`, `feed`) VALUES
            (1, 'AUD', 2, 1, '2018-07-18 15:33:38', '2018-11-11 10:40:44', 'fiat', 'ECB'),
            (2, 'BGN', 2, 1, '2018-07-18 15:33:38', '2018-11-11 10:40:44', 'fiat', 'ECB'),
            (3, 'BRL', 2, 1, '2018-07-18 15:33:38', '2018-11-11 10:40:44', 'fiat', 'ECB'),
            (4, 'CAD', 2, 1, '2018-07-18 15:33:38', '2018-11-11 10:40:44', 'fiat', 'ECB'),
            (5, 'CHF', 2, 1, '2018-07-18 15:33:38', '2018-11-11 10:40:44', 'fiat', 'ECB'),
            (6, 'CNY', 2, 1, '2018-07-18 15:33:38', '2018-11-11 10:40:44', 'fiat', 'ECB'),
            (7, 'CZK', 2, 1, '2018-07-18 15:33:38', '2018-11-11 10:40:44', 'fiat', 'ECB'),
            (8, 'DKK', 2, 1, '2018-07-18 15:33:38', '2018-11-07 09:32:36', 'fiat', 'ECB'),
            (9, 'EUR', 2, 1, '2018-07-18 15:33:38', '2018-11-07 09:32:36', 'fiat', 'ECB'),
            (10, 'GBP', 2, 1, '2018-07-18 15:33:38', '2018-11-12 06:09:17', 'fiat', 'ECB'),
            (11, 'HKD', 2, 1, '2018-07-18 15:33:38', '2018-11-12 06:09:17', 'fiat', 'ECB'),
            (12, 'HRK', 2, 1, '2018-07-18 15:33:38', '2018-11-12 06:09:17', 'fiat', 'ECB'),
            (13, 'HUF', 2, 1, '2018-07-18 15:33:38', '2018-11-12 06:09:17', 'fiat', 'ECB'),
            (14, 'ILS', 2, 1, '2018-07-18 15:33:38', '2018-11-07 09:32:36', 'fiat', 'ECB'),
            (15, 'INR', 2, 1, '2018-07-18 15:33:38', '2018-11-12 06:09:17', 'fiat', 'ECB'),
            (16, 'JPY', 0, 1, '2018-07-18 15:33:38', '2018-11-12 06:09:17', 'fiat', 'ECB'),
            (17, 'MXN', 2, 1, '2018-07-18 15:33:38', '2018-11-12 06:09:17', 'fiat', 'ECB'),
            (18, 'NOK', 2, 1, '2018-07-18 15:33:38', '2018-11-08 12:39:40', 'fiat', 'ECB'),
            (19, 'NZD', 2, 1, '2018-07-18 15:33:38', '2018-12-03 11:50:36', 'fiat', 'ECB'),
            (20, 'PLN', 2, 1, '2018-07-18 15:33:38', '2018-11-07 09:32:36', 'fiat', 'ECB'),
            (21, 'RON', 2, 1, '2018-07-18 15:33:38', '2019-02-15 06:44:51', 'fiat', 'ECB'),
            (22, 'RUB', 2, 1, '2018-07-18 15:33:38', '2018-12-03 11:50:36', 'fiat', 'ECB'),
            (23, 'SEK', 2, 1, '2018-07-18 15:33:38', '2018-11-26 09:06:32', 'fiat', 'ECB'),
            (24, 'SGD', 2, 1, '2018-07-18 15:33:38', '2019-02-08 06:21:38', 'fiat', 'ECB'),
            (25, 'THB', 2, 1, '2018-07-18 15:33:38', '2019-02-08 06:21:38', 'fiat', 'ECB'),
            (26, 'TRY', 2, 1, '2018-07-18 15:33:38', '2018-11-30 05:17:33', 'fiat', 'ECB'),
            (27, 'USD', 2, 1, '2018-07-18 15:33:38', '2018-11-28 05:32:53', 'fiat', 'ECB'),
            (28, 'ZAR', 2, 1, '2018-07-18 15:33:38', '2019-02-08 06:21:35', 'fiat', 'ECB'),
            (30, 'BTC', 8, 1, '2018-07-18 15:33:38', '2019-02-08 11:51:05', 'crypto', NULL);

            CREATE TABLE `rates` (
              `id` int(10) UNSIGNED NOT NULL,
              `currency_from_id` int(10) UNSIGNED NOT NULL,
              `currency_to_id` int(10) UNSIGNED NOT NULL,
              `value` decimal(36,18) NOT NULL DEFAULT '0.000000000000000000',
              `exchange_margin` decimal(36,18) NOT NULL DEFAULT '0.000000000000000000',
              `created_at` timestamp NULL DEFAULT NULL,
              `updated_at` timestamp NULL DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            CREATE TABLE `schema_migrations` (
              `version` bigint(20) NOT NULL,
              `dirty` tinyint(1) NOT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            INSERT INTO `schema_migrations` (`version`, `dirty`) VALUES
            (20190620163909, 0);

            CREATE TABLE `settings` (
              `id` int(10) UNSIGNED NOT NULL,
              `main_currency_id` int(10) UNSIGNED NOT NULL,
              `auto_updating_rates` tinyint(1) NOT NULL DEFAULT '0'
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            INSERT INTO `settings` (`id`, `main_currency_id`, `auto_updating_rates`) VALUES
            (1, 9, 0);


            ALTER TABLE `currencies`
              ADD PRIMARY KEY (`id`),
              ADD UNIQUE KEY `code_UNIQUE` (`code`);

            ALTER TABLE `rates`
              ADD PRIMARY KEY (`id`),
              ADD KEY `currency_from_id` (`currency_from_id`),
              ADD KEY `currency_to_id` (`currency_to_id`);

            ALTER TABLE `schema_migrations`
              ADD PRIMARY KEY (`version`);

            ALTER TABLE `settings`
              ADD PRIMARY KEY (`id`),
              ADD KEY `main_currency_id` (`main_currency_id`);


            ALTER TABLE `currencies`
              MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=34;

            ALTER TABLE `rates`
              MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;

            ALTER TABLE `settings`
              MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;


            ALTER TABLE `rates`
              ADD CONSTRAINT `rates_ibfk_1` FOREIGN KEY (`currency_from_id`) REFERENCES `currencies` (`id`),
              ADD CONSTRAINT `rates_ibfk_2` FOREIGN KEY (`currency_to_id`) REFERENCES `currencies` (`id`);

            ALTER TABLE `settings`
              ADD CONSTRAINT `settings_ibfk_1` FOREIGN KEY (`main_currency_id`) REFERENCES `currencies` (`id`);
SQL;
    }
}
