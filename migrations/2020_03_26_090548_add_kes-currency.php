<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class AddKesCurrency extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        DB::table('currencies')->insert([
            'code' => 'KES',
            'decimal_places' => 2,
            'active' => 1,
            'type' => 'fiat',
            'name' => 'Kenya Shilling',
            'logo_file_id' => 0,
            'coin_market_cap_id' => 0,
        ]);
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
    }
}
