<?php

namespace App\Models;

use Jenssegers\Mongodb\Eloquent\Model as Eloquent;

class Game extends Eloquent
{
    protected $connection = 'mongodb'; // Use MongoDB connection
    protected $collection = 'game'; // Specify the collection name

    protected $fillable = [
        'title',
        'genre',
        'platform',
        'developer',
        'publisher',
        'status',
        'played_time',
        'time_to_beat',
        'ranking',
        'rating',
    ];

    // You can add additional methods or relationships here as needed
}
