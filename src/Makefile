readDB:
	go run ex00/readDB.go -f data/original_database.xml
	go run ex00/readDB.go -f data/stolen_database.json

compareDB:
	go run ex01/compareDB.go --old data/original_database.xml --new data/stolen_database.json

compareFS:
	go run ex02/compareFS.go --old data/snapshot1.txt --new data/snapshot2.txt