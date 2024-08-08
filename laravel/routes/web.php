<?php

use Illuminate\Support\Facades\Route;

Route::get('admin', [AdminController::class, 'index']);

Route::prefix('api')->group(function () {
    Route::get('csrf', [ApiController::class, 'getCsrf'])->name('get_csrf');
    
    Route::prefix('game')->group(function () {
        Route::get('badge', [GameController::class, 'badges'])->name('badges');
        Route::post('create', [GameController::class, 'createGame'])->name('create_game');
        Route::delete('delete', [GameController::class, 'deleteGame'])->name('delete_game');
        Route::get('pages', [GameController::class, 'getGames'])->name('get_games');
        Route::post('start', [GameController::class, 'startGame'])->name('start_game');
        Route::post('stop', [GameController::class, 'stopGame'])->name('stop_game');
        Route::get('stopwatch', [GameController::class, 'stopwatch'])->name('stopwatch');
        Route::put('update', [GameController::class, 'updateGame'])->name('update_game');
    });
});