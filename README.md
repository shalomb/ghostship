Ghostship is like [Starship](https://starship.rs/) but a simpler prompt manager for bash.

![image](https://github.com/user-attachments/assets/88c5b78b-9ca0-4c79-95bf-6b706909d75e)

It is an implementation in go with some fundamental philosophical differences.

- Colours indicate significant shell state change.
- Features are costly. Only the core will remain in go.
- Speed is secondary. Ghostship is _un_blazingly fast!

## SSH Compatibility

Ghostship is designed to work reliably in SSH environments. Git operations are automatically timeout-protected (2 seconds) to prevent hangs during:

- Slow network connections
- Unreachable git remotes
- Authentication delays
- Network-mounted filesystems

This ensures that `source <(ghostship init bash)` works smoothly even in challenging network conditions, and Ctrl-C remains responsive.
