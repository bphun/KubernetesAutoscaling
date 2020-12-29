from flask import Flask, request
import numpy as np
import time, json

app = Flask(__name__)
app.config["DEBUG"] = False


@app.route("/", methods=["POST"])
def home():
    response = None

    if not request.json or not request.json["arr"]:
        response = app.response_class(
            response=json.dumps({"message": "No input array"}),
            status=400,
            mimetype="application/json",
        )
    else:
        inputArr = request.json["arr"]
        cumSumMode = request.json["mode"]
        cumSumArr, executionTime = [], 0

        if cumSumMode == "np":
            cumSumArr, executionTime = npCumSum(inputArr)
        elif cumSumMode == "man":
            cumSumArr, executionTime = cumSum(inputArr)

        executionTime *= 1000

        app.logger.info("{} digit CumSum in {}ms".format(len(inputArr), executionTime))

        response = app.response_class(
            response=json.dumps(
                {
                    "message": "success",
                    "data": cumSumArr,
                    "exec_ms": executionTime,
                }
            ),
            status=200,
            mimetype="application/json",
        )

    return response


def npCumSum(arr):
    startTime = time.time()
    cumSumArr = np.cumsum(arr).tolist()
    endTime = time.time()

    return (cumSumArr, (endTime - startTime))


def cumSum(arr):
    startTime = time.time()
    for i in range(1, len(arr)):
        arr[i] += arr[i - 1]
    endTime = time.time()

    return (arr, (endTime - startTime))


if __name__ == "__main__":
    app.run(host="0.0.0.0")
