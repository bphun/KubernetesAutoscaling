from flask import Flask, request
import numpy as np
app = Flask(__name__)
app.config["DEBUG"] = False

@app.route("/", methods=["POST"])
def home():
    if not request.json or not request.json["arr"]:
        return {"message": "No input array", "status": 400}

    inputArr = request.json["arr"]

    return {"message": "success", "status": 200, "data": np.cumsum(inputArr).tolist()}


def cumSum(arr):
    for i in range(1, len(arr)):
        arr[i] += arr[i - 1]

    return arr

if __name__ == "__main__":
    app.run(host='0.0.0.0')