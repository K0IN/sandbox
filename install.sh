#!/bin/sh

# Define the GitHub repository.
REPO="K0IN/sandbox"

# Find the architecture of the current system.
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm*)
        ARCH="arm"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Define the base directory for temporary extraction.
TMP_DIR=$(mktemp -d)

# Use the GitHub API to get the download URL for the tarball of the latest release.
TARBALL_URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest |
    grep "browser_download_url.*${ARCH}\.tar\.gz" |
    cut -d '"' -f 4 | head -n 1)

# Check if the tarball URL wasn't found.
if [ -z "$TARBALL_URL" ]; then
    echo "Error: Unable to retrieve the latest release tarball for architecture: $ARCH"
    exit 1
fi

# Download the tar.gz file.
echo "Downloading sandbox for $ARCH from $TARBALL_URL..."
curl -sSL -o "${TMP_DIR}/sandbox.tar.gz" "$TARBALL_URL"

# Check if the download was successful.
if [ $? -ne 0 ] || [ ! -f "${TMP_DIR}/sandbox.tar.gz" ]; then
    echo "Error: Download failed."
    rm -rf "$TMP_DIR"
    exit 1
fi

# Extract the 'sandbox' binary from the tar.gz into the temporary directory.
tar -xzf "${TMP_DIR}/sandbox.tar.gz" -C "$TMP_DIR"

# Assuming the binary is named 'sandbox' and is located in the root of the tar directory.
chmod +x "${TMP_DIR}/sandbox"

echo "Installing sandbox to /usr/local/bin."
# Move the binary to a location in the user's PATH.

if [ "$(id -u)" -eq 0 ]; then
    mv "${TMP_DIR}/sandbox" /usr/local/bin/sandbox
    chown root:root /usr/local/bin/sandbox
    chmod u+s /usr/local/bin/sandbox
else
    sudo mv "${TMP_DIR}/sandbox" /usr/local/bin/sandbox
    sudo chown root:root /usr/local/bin/sandbox
    sudo chmod u+s /usr/local/bin/sandbox
fi

# Cleanup the temporary directory.
rm -rf "$TMP_DIR"

# Verify that installation was successful.
if command -v sandbox >/dev/null 2>&1; then
    echo "Installation successful. You can now run 'sandbox' from the command line."
    echo "$ sandbox --help"
    echo "warning: sandbox is a setuid binary."
    exit 0
else
    echo "Error: sandbox was not installed correctly."
    exit 1
fi
