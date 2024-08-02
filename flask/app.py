from flask import Flask, jsonify, request
from db import Connection
from bson.objectid import ObjectId
from bson.json_util import dumps
from game import game

app = Flask(__name__)

app.register_blueprint(game, url__prefix='/api/game')

if __name__ == '__main__':
    app.run(debug=True)