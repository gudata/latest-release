#!/usr/bin/env python3

import sys
import requests
import re
import os


# Get the GitHub repository URL from the command-line argument
repo_url = sys.argv[1]


# Extract the repository owner and name from the URL
match = re.search(r"github.com/([^/]+)/([^/]+).*$", repo_url)
if not match:
    print("Invalid GitHub repository URL")
    sys.exit(1)
owner, repo = match.groups()

# Construct the releases page URL
url = f"https://github.com/{owner}/{repo}/releases"
print(f"I got that we need to fetch latest release from {url}")

# Send a GET request to the releases page and get the response HTML
# response = requests.get(url)
# html = response.text
with open('h', 'r') as f:
    lines = f.readlines()
html = "".join(lines)

# Find the link to the latest release tarball using a regex pattern
# pattern = r"https://github\.com/{owner}/{repo}/releases/.*\.tar\.gz"
pattern = r'nclude-fragment loading="lazy" src="([^"]+)"'

latest_release_links = re.findall(pattern, html)
print(latest_release_links[0])

# response = requests.get(url)
# html = response.text

with open('hh', 'r') as f:
    lines = f.readlines()
html = "".join(lines)



pattern = fr'a href="(/{owner}/{repo}/releases/download/[^"]+)" rel="nofollow"'

print(pattern)
match = re.search(pattern, html)

download_links = re.findall(pattern, html)
print(download_links)

for path in download_links:
    if ("x86_64" in path or "amd64" in path or "musl" in path) and "linux" in path:
        print(f"Matched string: {path}")
        download_link = f"https://github.com{path}"
        break
else:
    print("No match found")

print(download_link)
filename = os.path.basename(download_link)

# Download the latest release tarball
response = requests.get(download_link)
with open(filename, "wb") as f:
    f.write(response.content)

print(f"Downloaded {filename}")

