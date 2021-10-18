# The tool for moving files into directories by file extensions

## Example

- before moving structure:

```
resources
│
├── res1
│   └── docs
│       └── README.md
├── res2
│   ├── text1.txt
│   └── text2.txt
└── main.html
```

- moving into same extension dir result:

```
resources
│
├── EXT_HTML
│   └── main.html
├── EXT_MD
│   └── README.md
└── EXT_TXT
    ├── text1.txt
    └── text2.txt
```


- moving into different extension directory result:

```
resources
│
├── EXT_HTML
│   └── main.html
├── res1
│   └── docs
│       └── EXT_MD
│           └── README.md
└── res2 
    └── EXT_TXT
        ├── text1.txt
        └── text2.txt
```