import os
import subprocess
import hashlib
import tarfile
import configparser
import lzma
from datetime import datetime

# ---- Config ----
SOURCE_FILE = "main.go"
OUTPUT_BASE = os.path.join(os.getcwd(), "build")
TARGETS = [
    ("windows", "amd64"),
    ("windows", "arm64"),
    ("linux",   "amd64"),
    ("linux",   "arm64"),
    ("linux",   "arm"),  # armv6 is supported by GOARM=6
]
GOARM = {
    "arm": "6"  # target armv6 for SBCs
}

CFG = configparser.ConfigParser()
CFG.read("build.config.ini")
GPG_KEY_ID = CFG["gpg"]["key_id"] 

# ---- Logging ----
def log(level, msg):
    ts = datetime.now().isoformat(timespec="seconds")
    print(f"[{ts}] [{level.upper():5}] {msg}")

def run_cmd(cmd, cwd=None):
    subprocess.run(cmd, check=True, cwd=cwd)

def sha256sum(filepath):
    sha256 = hashlib.sha256()
    with open(filepath, "rb") as f:
        for chunk in iter(lambda: f.read(4096), b""):
            sha256.update(chunk)
    return sha256.hexdigest()

# ---- Build ----
def build(goos, goarch):
    ext = ".exe" if goos == "windows" else ""
    name = f"main-{goos}-{goarch}{ext}"

    target_dir = f"{goos}-{goarch}"
    out_dir = os.path.join(OUTPUT_BASE, target_dir)
    os.makedirs(out_dir, exist_ok=True)
    binary_path = os.path.join(out_dir, name)

    env = os.environ.copy()
    env["GOOS"] = goos
    env["GOARCH"] = goarch
    if goarch == "arm":
        env["GOARM"] = GOARM["arm"]

    log("info", f"Building {goos}/{goarch}")
    try:
        subprocess.run(
            ["go", "build", "-o", binary_path, SOURCE_FILE],
            check=True,
            cwd=os.getcwd(),
            env=env,
        )
        log("info", f"Built binary → {binary_path}")
    except subprocess.CalledProcessError:
        log("error", f"Failed building {goos}/{goarch}")
        return

    compress_and_sign(binary_path, out_dir)


# ---- Compress + Sign + Checksum ----
def compress_and_sign(binary_path, out_dir):
    name = os.path.basename(binary_path)
    xz_path = os.path.join(out_dir, f"{name}.xz")

    # Compress to .xz
    with open(binary_path, "rb") as f_in:
        with lzma.open(xz_path, "wb") as f_out:
            f_out.write(f_in.read())
    log("info", f"Compressed → {xz_path}")

    # sha256
    sha_path = xz_path + ".sha256"
    with open(sha_path, "w") as f:
        f.write(sha256sum(xz_path) + "  " + os.path.basename(xz_path) + "\n")
    log("info", f"SHA256 → {sha_path}")

    # GPG sign
    asc_path = xz_path + ".asc"
    try:
        cmd = ["gpg", "--detach-sign", "--armor"]
        if GPG_KEY_ID:
            cmd.extend(["--local-user", GPG_KEY_ID])
        cmd.extend(["--output", asc_path, xz_path])
        run_cmd(cmd)
        log("info", f"Signed → {asc_path}")
    except subprocess.CalledProcessError:
        log("warn", "GPG sign failed")

# ---- Main ----
if __name__ == "__main__":
    log("info", "Start cross-build...")
    if not os.path.exists(SOURCE_FILE):
        log("error", f"{SOURCE_FILE} not found.")
        exit(1)
    for goos, goarch in TARGETS:
        build(goos, goarch)
    log("info", "🎉 All builds complete.")
