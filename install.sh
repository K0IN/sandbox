#!/bin/sh

# This is the repository where `sandbox` is hosted
REPO="K0IN/sandbox"

# Get the latest release from GitHub api
URL=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep "browser_download_url" | cut -d '"' -f 4)

# Check if the URL is non-empty
if [ -z "$URL" ]; then
    echo "Error: Unable to retrieve the latest release."
    exit 1
fi

# Download the latest binary release of `sandbox`
echo "Downloading sandbox from $URL..."
curl -L -o sandbox "$URL"

# Ensure the binary is executable
chmod +x sandbox

# Move the binary to a location in the user's PATH (requires sudo for /usr/local/bin)
echo "Installing sandbox to /usr/local/bin - you might be prompted for your password."
sudo mv sandbox /usr/local/bin/

# Verify that installation was successful
if [ -x "$(command -v sandbox)" ]; then
    echo "Installation successful. You can now run 'sandbox' from the command line."
else
    echo "Error: sandbox was not installed correctly."
    exit 1
fi