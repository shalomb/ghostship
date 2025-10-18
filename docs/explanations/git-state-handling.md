# Git State Handling in Ghostship

## Overview

This document explains how Ghostship handles various git repository states and edge cases. It provides a comprehensive analysis of the current implementation and identifies areas for improvement.

## Current Git State Detection Logic

### 1. Repository Detection (`isGitDirectory()`)
```go
func isGitDirectory() bool {
    stdout, err := _git([]string{
        "git", "rev-parse", "--is-inside-work-tree",
    }...)
    if err != nil {
        return false
    }
    return stdout == "true"
}
```

**Edge Cases Handled:**
- ✅ **Not a git repo**: Returns `false` on error
- ✅ **Outside work tree**: Returns `false` when not inside work tree
- ❌ **Bare repository**: Uses `--is-inside-work-tree` which returns `false` for bare repos

### 2. Bare Repository Detection (`isGitRepoBare()`)
```go
func isGitRepoBare() bool {
    stdout, err := _git([]string{
        "git", "rev-parse", "--is-bare-repository",
    }...)
    if err != nil {
        return false // TODO: this is an erroneous reason to return false
    }
    return stdout == "true"
}
```

**Issues:**
- ❌ **Error handling**: Returns `false` on error, but comment indicates this is wrong
- ❌ **Bare repo workflow**: Bare repos are detected but not properly handled in init

### 3. Branch Detection (`gitCurrentBranch()`)
```go
func gitCurrentBranch() (string, error) {
    stdout, err := _git([]string{
        "git", "branch", "--show-current",
    }...)
    if err != nil {
        return "", nil  // BUG: Should return error, not nil
    }
    return stdout, err
}
```

**Edge Cases:**
- ✅ **No commits yet**: Returns empty string (correct)
- ✅ **Detached HEAD**: Returns empty string (correct)
- ❌ **Error handling**: Returns `nil` instead of error on failure
- ❌ **New repo**: After `git init`, no branch exists until first commit

### 4. Commit Hash Detection (`gitRev()`)
```go
func gitRev() (string, string, error) {
    stdout, err := _git([]string{
        "git", "rev-parse", "HEAD",
    }...)
    
    // Add bounds checking to prevent slice bounds out of range panic
    if len(stdout) < 8 {
        return string(stdout), string(stdout), err
    }
    return string(stdout), string(stdout)[0:8], err
}
```

**Edge Cases:**
- ✅ **No commits**: Returns error (correct)
- ✅ **Short hash**: Now handles gracefully with bounds checking
- ✅ **Detached HEAD**: Works correctly
- ❌ **Error propagation**: Errors are ignored in init function

### 5. Dirty State Detection (`isGitRepoDirty()`)
```go
func isGitRepoDirty() bool {
    stdout, err := _git([]string{
        "git", "diff", "--no-ext-diff", "--quiet", "--exit-code",
    }...)
    if err != nil {
        return true  // Assumes dirty on error
    }
    return stdout != ""
}
```

**Edge Cases:**
- ✅ **Clean repo**: Returns `false` correctly
- ✅ **Staged changes**: Returns `false` (only checks working directory)
- ✅ **Unstaged changes**: Returns `true` correctly
- ❌ **Error handling**: Assumes dirty on any error (may be incorrect)

### 6. Status Detection (`gitStatus()`)
```go
func gitStatus() (string, gitStatusMask, error) {
    var status gitStatusMask
    
    stdout, err := _git([]string{
        "git", "status", "--porcelain",
    }...)
    if err != nil {
        return "", status, err
    }
    
    for _, line := range strings.Split(string(stdout), "\n") {
        r, _ := regexp.Compile(`^\s*(\S+)`)
        statusField := r.FindString(line)
        for _, xy := range statusField {
            switch string(xy) {
            case "?":
                status |= untracked
            case "A":
                status |= staged
            case "M":
                status |= modified
            case "U":
                status |= unmerged
            }
        }
    }
    // ... status symbol generation
}
```

**Edge Cases:**
- ✅ **Empty repo**: Handles correctly (no status)
- ✅ **Various file states**: Handles untracked, staged, modified, unmerged
- ❌ **Incomplete parsing**: Only handles first character of status codes
- ❌ **Missing states**: Doesn't handle deleted, renamed, etc. in parsing

## Critical Edge Cases Analysis

### 1. Fresh Git Repository (Just After `git init`)
**Current Behavior:**
- ✅ Detected as git directory
- ✅ Branch is empty string (correct)
- ✅ No commit hash (error handled)
- ✅ Not dirty (correct)
- ✅ No status symbols (correct)

**Status:** ✅ **HANDLED CORRECTLY**

### 2. Detached HEAD State
**Current Behavior:**
- ✅ Detected as git directory
- ✅ Branch is empty string (correct)
- ✅ Commit hash works (correct)
- ✅ Dirty detection works
- ✅ Status detection works

**Status:** ✅ **HANDLED CORRECTLY**

### 3. Bare Repository
**Current Behavior:**
- ❌ Not detected as git directory (uses `--is-inside-work-tree`)
- ❌ No processing occurs
- ❌ Would show no git status

**Status:** ❌ **NOT HANDLED**

### 4. Repository with No Commits
**Current Behavior:**
- ✅ Detected as git directory
- ✅ Branch is empty string (correct)
- ❌ `git rev-parse HEAD` fails but error is ignored
- ✅ Dirty detection works
- ✅ Status detection works

**Status:** ⚠️ **PARTIALLY HANDLED** (errors ignored)

### 5. Repository with Staged Changes Only
**Current Behavior:**
- ✅ Detected as git directory
- ✅ Branch detection works
- ✅ Commit hash works
- ❌ `isGitRepoDirty()` returns `false` (only checks working directory)
- ✅ Status detection shows staged files

**Status:** ⚠️ **PARTIALLY HANDLED** (dirty detection incomplete)

### 6. Repository with Merge Conflicts
**Current Behavior:**
- ✅ Detected as git directory
- ✅ Branch detection works
- ✅ Commit hash works
- ✅ Dirty detection works
- ✅ Status detection shows unmerged files
- ✅ Special styling for conflicts

**Status:** ✅ **HANDLED CORRECTLY**

## Issues Found

### 1. Error Handling Problems
- **`gitCurrentBranch()`**: Returns `nil` instead of error
- **`init()` function**: Ignores all errors from git functions
- **`isGitRepoBare()`**: Comment indicates wrong error handling

### 2. Incomplete Status Parsing
- Only parses first character of git status codes
- Missing handling for deleted, renamed, stashed files
- Doesn't distinguish between staged and unstaged modifications

### 3. Bare Repository Support
- Not detected as git directory
- No special handling for bare repos

### 4. Dirty Detection Logic
- Only checks working directory changes
- Doesn't consider staged changes as "dirty"

## Recommendations

### 1. Fix Error Handling
```go
// Fix gitCurrentBranch error handling
func gitCurrentBranch() (string, error) {
    stdout, err := _git([]string{
        "git", "branch", "--show-current",
    }...)
    if err != nil {
        return "", err  // Return actual error
    }
    return stdout, nil
}
```

### 2. Improve Status Parsing
```go
// Parse both characters of status codes
for _, xy := range statusField {
    if len(xy) >= 2 {
        // Parse staged state (first char)
        switch string(xy[0]) {
        case "A": status |= staged
        case "M": status |= staged
        // ... etc
        }
        // Parse working directory state (second char)
        switch string(xy[1]) {
        case "M": status |= modified
        case "D": status |= deleted
        // ... etc
        }
    }
}
```

### 3. Add Bare Repository Support
```go
func isGitDirectory() bool {
    // Check both work tree and bare repository
    stdout, err := _git([]string{
        "git", "rev-parse", "--is-inside-work-tree",
    }...)
    if err == nil && stdout == "true" {
        return true
    }
    
    // Check if it's a bare repository
    stdout, err = _git([]string{
        "git", "rev-parse", "--is-bare-repository",
    }...)
    return err == nil && stdout == "true"
}
```

### 4. Improve Dirty Detection
```go
func isGitRepoDirty() bool {
    // Check working directory
    _, err1 := _git([]string{
        "git", "diff", "--no-ext-diff", "--quiet", "--exit-code",
    }...)
    
    // Check staged changes
    _, err2 := _git([]string{
        "git", "diff", "--cached", "--no-ext-diff", "--quiet", "--exit-code",
    }...)
    
    return err1 != nil || err2 != nil
}
```

## Summary

The application handles most common git states reasonably well, but has several issues:

1. **Error handling is poor** - errors are often ignored
2. **Status parsing is incomplete** - misses many git states
3. **Bare repositories not supported**
4. **Dirty detection is incomplete** - doesn't consider staged changes

The recent fix for the slice bounds issue was necessary and correct, but there are deeper architectural issues that should be addressed for robust git state handling.