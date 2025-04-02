## CLI Commands Documentation

### `cat`
**Description**: Concatenate and print file contents or stdin  
**Usage**: `cat [FILE]...`  
**Arguments**:
- `FILE`: One or more files to display (optional, reads from stdin if omitted)

**Behavior**:
- With files: Outputs contents of all specified files concatenated
- Without files: Outputs contents from stdin (if provided)
- Returns error if no input source is available

**Examples**:
```bash
> cat file.txt
> echo "text" | cat
```

---

### `echo`
**Description**: Print arguments to stdout  
**Usage**: `echo [STRING]...`  
**Arguments**:
- `STRING`: Text to output (supports quoted strings and `\n` escapes)

**Features**:
- Removes outer quotes (`"` or `'`) if present
- Converts `\n` to newlines in single-quoted strings
- Joins multiple arguments with spaces

**Examples**:
```bash
> echo hello world
> echo '"quoted"' → quoted
> echo -e 'line1\nline2' → line1[newline]line2
```

---

### `exit`
**Description**: Terminate the shell  
**Usage**: `exit [STATUS]`  
**Arguments**:
- `STATUS`: Exit code (default: 0)

**Behavior**:
- Accepts numeric exit code (0-255)
- Immediately terminates process with specified code

**Examples**:
```bash
> exit
> exit 1
```

---

### `pwd`
**Description**: Print working directory  
**Usage**: `pwd`  
**Output**: Absolute path of current directory

**Error Cases**:
- Fails if directory cannot be determined

**Example**:
```bash
> pwd
/home/user/project
```

---

### `wc`
**Description**: Count lines, words, and characters  
**Usage**: `wc [FILE]`  
**Arguments**:
- `FILE`: File to analyze (reads from stdin if omitted)

**Output Format**: `lines words bytes`

**Behavior**:
- Counts:
  - Lines (`\n` separated)
  - Words (whitespace-separated)
  - Bytes (raw length)
- Requires either file or stdin input

**Examples**:
```bash
> wc file.txt
3 15 102
> echo "hello world" | wc
1 2 12
```

---

### `grep`
**Description**: Pattern search in text  
**Usage**: `grep [OPTIONS] PATTERN`  
**Arguments**:
- `PATTERN`: Regular expression to search

**Options**:
| Flag | Description                          |
|------|--------------------------------------|
| `-i` | Case-insensitive search              |
| `-w` | Match whole words only               |
| `-A N` | Print N lines after each match    |

**Behavior**:
- Processes stdin only
- Supports PCRE regex syntax
- Handles overlapping `-A` contexts
- Returns error for invalid regex

**Examples**:
```bash
> grep "error" log.txt
> grep -i -A 1 "warning" < input.txt
> echo "test" | grep -w "test"
```

---

### Notes
1. All commands:
   - Accept input via stdin when piped
   - Return `*bytes.Buffer` with output
   - Propagate errors with context

2. Common patterns:
   - Empty args → use stdin (where applicable)
   - File operations → relative to current `pwd`
   - String parsing → handles basic quoting

3. Error handling:
   - Commands return descriptive errors
   - Exit codes follow Unix conventions (where applicable)