from flask import Flask, jsonify
from flask_cors import CORS
from prometheus_flask_exporter import PrometheusMetrics
import requests
import os

app = Flask(__name__)
CORS(app)
PrometheusMetrics(app)  # exposes /metrics automatically

@app.route("/")
def health():
    return "The service is running", 200

@app.route('/<city>')
def hello(city):
    url = "https://weatherapi-com.p.rapidapi.com/current.json"
    querystring = {"q": city}
    headers = {
        'x-rapidapi-host': "weatherapi-com.p.rapidapi.com",
        'x-rapidapi-key': os.getenv("APIKEY")
    }
    response = requests.get(url, headers=headers, params=querystring, timeout=10)
    response.raise_for_status()
    return jsonify(response.json())
