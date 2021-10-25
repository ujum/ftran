# The cli tool for moving files into directories by file extensions

The UI app based on cross platform GUI toolkit: [ftranUI](https://github.com/ujum/ftranUI ) 

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

- moving without saving directory structure into extension directory result:

```
result
│
├── EXT_HTML
│   └── main.html
├── EXT_MD
│   └── README.md
└── EXT_TXT
    ├── text1.txt
    └── text2.txt
```


- moving with saving directory structure into extension directory result:

```
result
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