import sys
sys.path.append('./helpers')
import re
import restHelper
import dockerHelper
import fileHelper
import databaseHelper
import json

dockerHelper.startDocker()
databaseHelper.connect()

testDataPath = "./testData/"

def test_HappyPath():
    data = fileHelper.openFile(testDataPath + "ports.json")
    res = restHelper.post("localhost", 8080, "/v1/sendports", data)
    assert res['success'] == True

    # check in the database if the ports exist
    jsonData = json.loads(data)
    for primary_unloc in jsonData:
        res = databaseHelper.runQueryWithResult(
        'SELECT primary_unloc FROM ports WHERE primary_unloc=\'' + primary_unloc +'\'' +
        ' and deleted_at is NULL')
        assert res[0] == primary_unloc

def test_codeTooLong():
    data = fileHelper.openFile(testDataPath + "codeTooLong.json")
    res = restHelper.post("localhost", 8080, "/v1/sendports", data)
    assert res['success'] == False
    assert re.search("ERROR: value too long for type",res['error'])

    # check that noen of the ports were written to the database
    jsonData = json.loads(data)
    for primary_unloc in jsonData:
        res = databaseHelper.runQueryWithResult(
        'SELECT primary_unloc FROM ports WHERE primary_unloc=\'' + primary_unloc +'\'' +
        ' and deleted_at is NULL')
        assert res == None

# def test_badJSON():
#     data = fileHelper.openFile(testDataPath + "badJSON.json")
#     res = restHelper.post("localhost", 8080, "/v1/sendports", data)
#     assert res['success'] == False
#     assert re.search("invalid character ",res['error'])
