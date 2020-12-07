# compactdisk

Yet another tool to save disk space by hardlinking the same file content on the disk.

### how it works
Scans the given `target path` recursively and stores sha1 per each file.
Once an entry has an sha1 that was seen before, deletes that file and hardlinks from
the first instance.

- Uses multiple cores in a min(4, cpucount) thread pool (if available)
- Uses a concatened bunch of SHA1s for large files (each per 1GB segment of the file)

### how to run
```
go run compactdisk.go <target path>
```

### disclaimer
_Use it on your own risk_