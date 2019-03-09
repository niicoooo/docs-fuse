# Usage
```console
foo@bar:~$ ./docs-fuse-raw -a http://127.0.0.1:8100 -u admin -p admin -d mnt &
Usage: ./docs-fuse-raw -h to display help
Connecting: admin@http://127.0.0.1:8100
Ctrl-C or umount to stop
foo@bar:~$ tree mnt
mnt
+-- 16fe70b8-fd93-41a1-872b-dfa89c84adfb doc1
¦   +-- data.json
¦   +-- files.json
+-- 1be48911-5f8c-4306-827c-565b5ea7cf0e doc1
¦   +-- data.json
¦   +-- files.json
+-- 7e1b9f80-662f-4bf2-9c38-0c18ed55779f doc1
¦   +-- data.json
¦   +-- files.json
+-- 903ba9d9-31df-46b5-9d54-30695d6f095a Item 1
    +-- 44adfb82-aec8-419c-84ee-c758b1c8c7f9 example 2.txt
    +-- 79b30f96-7539-4afa-a3e8-2a4a6ce3e812 example 1.txt
    +-- data.json
    +-- files.json

4 directories, 10 files
foo@bar:~$ cat mnt/903ba9d9-31df-46b5-9d54-30695d6f095a\ Item\ 1/files.json | jq .
{
  "files": [
    {
      "id": "79b30f96-7539-4afa-a3e8-2a4a6ce3e812",
      "processing": false,
      "name": "example 1.txt",
      "version": 0,
      "mimetype": "text/plain",
      "document_id": "903ba9d9-31df-46b5-9d54-30695d6f095a",
      "size": 1
    },
    {
      "id": "44adfb82-aec8-419c-84ee-c758b1c8c7f9",
      "processing": false,
      "name": "example 2.txt",
      "version": 0,
      "mimetype": "text/plain",
      "document_id": "903ba9d9-31df-46b5-9d54-30695d6f095a",
      "size": 1
    }
  ]
}
foo@bar:~$ umount mnt
```
