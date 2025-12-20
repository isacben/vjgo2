# vj

Interactive TUI JSON viewer.

## Features

* Folding
* Syntax highlighting
* Vim navigation and motions
* Relative line numbers

## Usage

Run this command to load a JSON file:

```bash
vj file.json
```

Or pipe JSON data into vj:

```bash
echo '{"helo": "world"}' | vj
```

## Key Bindings

### Folding

`h` or `←` - fold JSON object or array<br>
`l` or `→` - unfold JSON object or array<br>

### Navigation

`j` or `↓` - move cursor down<br>
`k` or `↑` - move cursor up<br>
`5j` - move cursor 5 lines down from current position<br>
`5k` - move cursor 5 lines up from current position<br>
`g` - move cursor to the first line of the document<br>
`G` - move cursor to the last line of the document<br>

### Command Mode

`:` - Switch to commands mode
`:.` - Find path in JSON for example `:.users[0].email`
`:q` - Close vj
