from fx import fx
from flask import Flask, request, jsonify
app = Flask(__name__)

@app.route('/', methods=['POST', 'GET'])
def handle():
    return fx(request)
