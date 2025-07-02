## [1.0.1](https://github.com/agkunz/docker-psa/compare/v1.0.0...v1.0.1) (2025-07-02)


### Bug Fixes

* updated readme ([87c590c](https://github.com/agkunz/docker-psa/commit/87c590cdc8d3d05c8a2bb3ba78bde82993365941))

# 1.0.0 (2025-07-02)


### Features

* f1rst p0st ([b9358b5](https://github.com/agkunz/docker-psa/commit/b9358b53e1a9aa4af24b0e74319083535b956853))

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Features
- Enhanced container output formatting with color-coded status indicators
- Multiple verbosity levels (-v, -vv) for different levels of detail
- Watch mode (-w) that continuously updates container status
- Support for filtering containers by name or image using regex
- Human-readable time and duration formatting
- Color-coded output using lipgloss styling

### Documentation
- Added comprehensive README with usage examples
- Created contributing guidelines with conventional commit format
- Set up automated changelog generation

### CI/CD
- Configured GitLab CI/CD pipeline with semantic release
- Added GitHub Actions workflow for cross-platform compatibility
- Automated multi-platform binary builds (Linux, macOS, Windows)
- Integrated semantic versioning and release automation
