from json import dumps
from util import convert
from flask import Blueprint, jsonify, request
from pymongo import MongoClient
from bson.objectid import ObjectId
from db import Connection

game = Blueprint('game', __name__)
db = Connection('ditto')

@game.route('/api/game', methods=['GET'])
def get_games():
    games = db.game.find({})
    return list(games)

@game.route('/api/game/<string:id>', methods=['GET'])
def get_game(id):
    try:
        game_id = ObjectId(id)
    except:
        return jsonify({'msg': 'Invalid game ID'}), 400

    game = db.game.find_one({'_id': game_id})
    if game:
        game = convert(game)
        return jsonify(game)
    return jsonify({'msg': 'Game not found'}), 404