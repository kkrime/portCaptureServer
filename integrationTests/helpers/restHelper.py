import json
import requests
import subprocess

# REST helpers
def post(host, port, url, data):
    response = requests.post("http://"+ host + ':' + str(port) + url, data=data.encode('utf-8'))
    jsonResponse = response.json()
    # printJson(jsonResponse)
    return jsonResponse

# Debugging
def printJson(data):
    formattedJson = json.dumps(data, indent=2)
    print(highlight(formattedJson, JsonLexer(), TerminalFormatter())) 
