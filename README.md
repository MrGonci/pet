# Pattern Enumeration Tool
The **Pattern Enumeration Tool** is designed to find and enumerate common patterns that may contain juicy information, such as credentials, sensitive URLs, or other valuable data. 

### Features
- **Custom Pattern:** Create y our patterns in *config.json*.
- **Multiplatform**: Runs seamlessly on Linux, macOS, and Windows.
- **Verbose**: Detailed mode for pattern matches.
- **File Output:**: Save results to a file.
- **Colored Output**: Uses colored text to highlight findings

### Requirements
- **Go**: Version 1.16 or higher is required to compile the tool.

### Steps
#### 1. Clone the repository:
```bash
git clone https://github.com/MrGonci/pet.git
cd pet
```

#### 2. Compiling the program
##### Windows (x64)
```bash
GOOS=windows GOARCH=amd64 go build -o pet.exe
```

##### For linux (x64)
```bash
GOOS=linux GOARCH=amd64 go build -o pet
```

### Usage
### Help
```bash
pet -h
```

##### Options
**-d:** Specify the directory to search (default: current directory).\
**-v:** Enable verbose mode to display detailed results.\
**-o:** Redirect the output to a file (e.g., -o results.txt).

### Example Commands
```bash
./pet -d .
./pet -d /var/www/admin_panel -v
./pet -d /var/www/admin_panel -o out.txt -v
```

### Output
![Output Image](/imgs/out.png)
