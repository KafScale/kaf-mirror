# Changelog

All notable changes to this project are documented here.

## [1.1.3] - 2025-11-26
### Highlights
- First open-source release of `kaf-mirror`; source code, binaries, and container images are now published under the project license.
- Public GitHub Releases now include version-stamped CLI binaries and GHCR container images for easy adoption.

### Fixes & stability
- Resolved the dashboard crash when clusters are absent/null and improved Chart.js canvas reuse.
- Fixed login loops when stale or compromised tokens are present.
- Release workflow hardened to ensure GitHub Releases succeed with the correct permissions.

### Packaging & distribution
- Prebuilt binaries for Linux, macOS, and Windows (amd64) are published with each tag.
- Multi-arch container images (`linux/amd64`, `linux/arm64`) are pushed to GHCR with both versioned and `latest` tags.

### Upgrade notes
- No breaking configuration changes; rolling upgrade is safe. Restart `kaf-mirror` or pull the new GHCR image to adopt 1.1.3.

## [1.1.0] - 2025-08
- First GA feature-complete release. See `RELEASE_v1.1.0.md` for the full notes.
