from collections import namedtuple
from itertools import groupby
#import pyodbc
#cnxn = pyodbc.connect('DRIVER={SQL Server};SERVER=localhost;DATABASE=wumbo;UID=root;PWD=b8c328c4f2')
#db = cnxn.cursor()
"""This pure Python MySQL client provides a DB-API to a MySQL database by talking directly to the server via the binary client/server protocol."""
import mysql.connector

cnx = mysql.connector.connect(user='root', password='b8c328c4f2',
                              host='127.0.0.1',
                              database='wumbo')
cur = cnx.cursor()
# "chunk" our input file, delimited by blank lines
filename=raw_input("enter S00E00:").upper()
for i in range(1,11):
	filename = "S02E"+"{0:0=2d}".format(i)
	with open(filename+".srt") as f:
	    res = [list(g) for b,g in groupby(f, lambda x: bool(x.strip())) if b]
	subs = []
	for sub in res:
	    if len(sub) >= 3: # not strictly necessary, but better safe than sorry
	        sub = [x.strip() for x in sub]
		start_end = sub[1]
		content = " ".join(sub[2::]).replace('"', '\'')
	        #number, start_end,*content = sub # py3 syntax
	        start, end = start_end.split(' --> ')
	        subs.append('("%s","%s","%s","%s")'% (filename,start,end,content))
	subs
	query= 'INSERT INTO `data` (`episode`, `startTime`, `stopTime`, `text`) VALUES '+','.join(subs)+';'
	cur.execute(query)
	for response in cur:
	    print(response)
cur.close()
cnx.close()
