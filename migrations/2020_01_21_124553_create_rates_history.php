<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateRatesHistory extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('rates_history', function (Blueprint $table) {
            $table->increments('id');
            $table->unsignedInteger('rate_id')->nullable(false);
            $table->decimal('value', 36, 18)->nullable(false);
            $table->char('provider', 40)->nullable(false);
            $table->dateTime('provider_time')->nullable(false);
            $table->timestamp('created_at')->nullable(true);

            $table->foreign('rate_id')->references('id')->on('rates')->onDelete('cascade');
            $table->index('provider_time');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('rates_history');
    }
}
