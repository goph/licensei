#!/bin/sh
set -e

readonly VERSION=${1:-}
readonly GIT_ROOT="$(git rev-parse --show-toplevel)"

if test "$DISTRIBUTION" = "pro"; then
	echo "Using Pro distribution..."
	RELEASES_URL="https://github.com/goreleaser/goreleaser-pro/releases"
	FILE_BASENAME="goreleaser-pro"
else
	echo "Using the OSS distribution..."
	RELEASES_URL="https://github.com/goreleaser/goreleaser/releases"
	FILE_BASENAME="goreleaser"
fi

if [ -z "$VERSION" ]; then
    echo "Usage: $0 v{version_to_download}"
    exit 1
fi

if [ -x ${GIT_ROOT}/bin/goreleaser-${VERSION} ]; then
    ln -sf goreleaser-${VERSION} bin/goreleaser
    echo "GoReleaser ${VERSION} is already installed"
    exit 0
fi

test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"
export TAR_FILE="$TMPDIR/${FILE_BASENAME}_$(uname -s)_$(uname -m).tar.gz"

(
	cd "$TMPDIR"
	echo "Downloading GoReleaser $VERSION..."
	curl -sfLo "$TAR_FILE" \
		"$RELEASES_URL/download/$VERSION/${FILE_BASENAME}_$(uname -s)_$(uname -m).tar.gz"
	curl -sfLo "checksums.txt" "$RELEASES_URL/download/$VERSION/checksums.txt"
	curl -sfLo "checksums.txt.sig" "$RELEASES_URL/download/$VERSION/checksums.txt.sig"
)

tar -xf "$TAR_FILE" -C "$TMPDIR"

mkdir -p ${GIT_ROOT}/bin bin
mv ${TMPDIR}/goreleaser ${GIT_ROOT}/bin/goreleaser-${VERSION}
ln -sf goreleaser-${VERSION} bin/goreleaser

echo "GoReleaser ${VERSION} is installed"
