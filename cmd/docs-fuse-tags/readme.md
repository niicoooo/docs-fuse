# Usage
```console
foo@bar:~$ ./docs-fuse-tags -a http://127.0.0.1:8100 -u admin -p admin -d mnt &
Usage: ./docs-fuse-tags -h to display help
Connecting: admin@http://127.0.0.1:8100
Ctrl-C or umount to stop
foo@bar:~$ tree mnt
mnt
├── docs
│   ├── doc1
│   │   ├── data.json
│   │   └── files.json
│   ├── doc1(2)
│   │   ├── data.json
│   │   └── files.json
│   ├── doc1(3)
│   │   ├── data.json
│   │   └── files.json
│   └── Item 1
│       ├── data.json
│       ├── example(2).txt
│       ├── example.txt
│       └── files.json
└── Tag1
    ├── docs
    │   ├── doc1
    │   │   ├── data.json
    │   │   └── files.json
    │   ├── doc1(2)
    │   │   ├── data.json
    │   │   └── files.json
    │   └── Item 1
    │       ├── data.json
    │       ├── example(2).txt
    │       ├── example.txt
    │       └── files.json
    └── Tag2
        └── docs
            └── Item 1
                ├── data.json
                ├── example(2).txt
                ├── example.txt
                └── files.json

13 directories, 22 files
foo@bar:~$ umount mnt
```
