from fx import fx
from flask import Flask, request,jsonify
app = Flask(__name__)

@app.route('/', methods=['POST'])
def do_it():
    request.get_json(force=True)
    j = request.json
    return jsonify(fx(j))
