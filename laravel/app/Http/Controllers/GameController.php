<?php

namespace App\Http\Controllers;

use App\Models\Game;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class GameController extends Controller
{
    public function badges()
    {
        // Logic to get badges for games
        return response()->json(['message' => 'Badges data']);
    }

    public function createGame(Request $request)
    {
        // Validate the incoming request
        $validator = Validator::make($request->all(), [
            'title' => 'required|string|max:255',
            // Add other validation rules as necessary
        ]);

        if ($validator->fails()) {
            return response()->json($validator->errors(), 400);
        }

        // Create a new game
        $game = Game::create($request->all());

        return response()->json($game, 201);
    }

    public function deleteGame(Request $request)
    {
        // Assuming you are passing the game ID in the request
        $gameId = $request->input('id');

        $game = Game::find($gameId);

        if (!$game) {
            return response()->json(['message' => 'Game not found'], 404);
        }

        $game->delete();

        return response()->json(['message' => 'Game deleted successfully']);
    }

    public function getGames()
    {
        $games = Game::all();
        return response()->json($games);
    }

    public function startGame(Request $request)
    {
        // Logic to start a game, assuming you may want to update the game status
        $gameId = $request->input('id');
        $game = Game::find($gameId);

        if (!$game) {
            return response()->json(['message' => 'Game not found'], 404);
        }

        // Update game status to "started"
        $game->status = 'started'; // Adjust based on your status field
        $game->save();

        return response()->json($game);
    }

    public function stopGame(Request $request)
    {
        // Logic to stop a game
        $gameId = $request->input('id');
        $game = Game::find($gameId);

        if (!$game) {
            return response()->json(['message' => 'Game not found'], 404);
        }

        // Update game status to "stopped"
        $game->status = 'stopped'; // Adjust based on your status field
        $game->save();

        return response()->json($game);
    }

    public function stopwatch(Request $request)
    {
        // Logic for the stopwatch, return the time elapsed, etc.
        return response()->json(['message' => 'Stopwatch data']);
    }

    public function updateGame(Request $request)
    {
        // Validate the incoming request
        $validator = Validator::make($request->all(), [
            'id' => 'required|exists:games,_id', // Assuming the ID is the MongoDB _id
            'title' => 'sometimes|required|string|max:255',
            // Add other validation rules as necessary
        ]);

        if ($validator->fails()) {
            return response()->json($validator->errors(), 400);
        }

        $game = Game::find($request->input('id'));

        if (!$game) {
            return response()->json(['message' => 'Game not found'], 404);
        }

        // Update game attributes
        $game->update($request->all());

        return response()->json($game);
    }
}