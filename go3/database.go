package main

import (
    "os"
    "path/filepath"
)

var datadir = filepath.Join(os.Getenv("HOME"), "ripple")
