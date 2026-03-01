set shell := ["nu.exe", "-c"]

EXE_DIR := "build"
MAIN_ENTRY := "cmd\\app\\main.go"

default:
    go build -o {{EXE_DIR}}\main.exe {{MAIN_ENTRY}}
