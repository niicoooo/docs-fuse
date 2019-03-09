# Usage
```console
foo@bar:~$ ./docs-fuse -a http://127.0.0.1:8100 -u admin -p admin -d mnt &
Usage: ./docs-fuse -h to display help
Connecting: admin@http://127.0.0.1:8100
Ctrl-C or umount to stop
foo@bar:~$ tree mnt
mnt
+-- doc1
¦   +-- data.json
¦   +-- files.json
+-- doc1 (2)
¦   +-- data.json
¦   +-- files.json
+-- doc1 (3)
¦   +-- data.json
¦   +-- files.json
+-- Item 1
    +-- data.json
    +-- example(2).txt
    +-- example.txt
    +-- files.json

4 directories, 10 files
foo@bar:~$ cat mnt/Item\ 1/files.json | jq .
{
  "files": [
    {
      "id": "79b30f96-7539-4afa-a3e8-2a4a6ce3e812",
      "processing": false,
      "name": "example.txt",
      "version": 0,
      "mimetype": "text/plain",
      "document_id": "903ba9d9-31df-46b5-9d54-30695d6f095a",
      "create_date": 1551622448451,
      "size": 1
    },
    {
      "id": "44adfb82-aec8-419c-84ee-c758b1c8c7f9",
      "processing": false,
      "name": "example.txt",
      "version": 0,
      "mimetype": "text/plain",
      "document_id": "903ba9d9-31df-46b5-9d54-30695d6f095a",
      "create_date": 1551622448548,
      "size": 1
    }
  ]
}
foo@bar:~$ umount mnt
```
