## install air

```bash
 go install github.com/air-verse/air@latest
 air
```

this is just helpful when working on a server, it rebuilds whenever you change code

- add tmp/ to gitignore

if you want to change any of the default behavior, you can configure the air.toml file in the top level dir

example

```toml
[build]
cmd = "go build -o tmp/main ."
bin = "tmp/main"
full_bin = "tmp/main"
include_ext = ["go"]
exclude_dir = ["vendor", "tmp"]
[log]
time = false
```

for now i will just leave the defaults and see how it goes. what happens if I change the readme?
nah, doesn't reload, good
