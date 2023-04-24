import psycopg2
from configparser import ConfigParser

conn = None

def readConfig():
    parser = ConfigParser()
    parser.read("./config/testing_config.toml")

    db = {}
    if parser.has_section('DBConfig'):
        params = parser.items('DBConfig')
        for param in params:
            db[param[0]] = param[1]
    else:
        raise Exception('Database config file not found')

    # print(db)
    return db

def connect():
    global conn
    try:
        params = readConfig()

        print('Connecting to the database')
        conn = psycopg2.connect(**params)
        print('Successfully connected to the database')
       
    except (Exception, psycopg2.DatabaseError) as error:
        print(error)

def runQuery(query):
    global conn
    cur = conn.cursor()
    
    cur.execute(query.encode('utf-8'))

    cur.close()

def runQueryWithResult(query):
    global conn
    cur = conn.cursor()
    
    cur.execute(query.encode('utf-8'))

    out = cur.fetchone()
    cur.close()

    return out
