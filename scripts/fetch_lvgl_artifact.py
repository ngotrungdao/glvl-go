#!/usr/bin/env python3
"""Download a build artifact from the ngotrungdao/lvgl "Build and Package Artifact"
workflow and unpack it into the local lvgl-c/ folder.

Auth: needs a GitHub token (artifact downloads are not anonymous).
Looked up in this order: --token, $GITHUB_TOKEN, $GH_TOKEN, `gh auth token`.
"""

from __future__ import annotations

import argparse
import json
import os
import shutil
import subprocess
import sys
import tarfile
import tempfile
import urllib.error
import urllib.parse
import urllib.request
import zipfile
from pathlib import Path

API = "https://api.github.com"
DEFAULT_REPO = "ngotrungdao/lvgl"
DEFAULT_WORKFLOW = "Build and Package Artifact"
KNOWN_ARTIFACTS = [
    "lvgl-build-fedora",
    "lvgl-build-linux-arm64",
    "lvgl-build-macos",
    "lvgl-build-windows",
]


def find_token(explicit: str | None) -> str:
    if explicit:
        return explicit
    for var in ("GITHUB_TOKEN", "GH_TOKEN"):
        if os.environ.get(var):
            return os.environ[var]
    gh = shutil.which("gh")
    if gh:
        try:
            out = subprocess.run(
                [gh, "auth", "token"], capture_output=True, text=True, check=True
            )
            token = out.stdout.strip()
            if token:
                return token
        except subprocess.CalledProcessError:
            pass
    sys.exit(
        "No GitHub token found. Set GITHUB_TOKEN, run `gh auth login`, or pass --token."
    )


class StripAuthOnRedirect(urllib.request.HTTPRedirectHandler):
    """Artifact downloads redirect to a pre-signed storage URL that rejects our token."""

    def redirect_request(self, req, fp, code, msg, headers, newurl):
        new = super().redirect_request(req, fp, code, msg, headers, newurl)
        if new is not None:
            for header in ("Authorization", "authorization"):
                new.headers.pop(header, None)
                new.unredirected_hdrs.pop(header, None)
        return new


def build_request(url: str, token: str | None) -> urllib.request.Request:
    req = urllib.request.Request(url)
    if token:
        req.add_header("Authorization", f"Bearer {token}")
    req.add_header("Accept", "application/vnd.github+json")
    req.add_header("X-GitHub-Api-Version", "2022-11-28")
    req.add_header("User-Agent", "fetch-lvgl-artifact")
    return req


def request(url: str, token: str | None, opener=None):
    req = build_request(url, token)
    try:
        return opener.open(req) if opener else urllib.request.urlopen(req)
    except urllib.error.HTTPError as exc:
        body = exc.read().decode("utf-8", "replace")[:500]
        sys.exit(f"GitHub API error {exc.code} for {url}\n{body}")


def get_json(url: str, token: str) -> dict:
    with request(url, token) as resp:
        return json.load(resp)


def latest_run_with_artifacts(repo: str, workflow: str, token: str) -> tuple[dict, list[dict]]:
    """Return the newest successful run of `workflow` that still has artifacts."""
    query = urllib.parse.urlencode({"status": "success", "per_page": 50})
    runs = get_json(f"{API}/repos/{repo}/actions/runs?{query}", token).get(
        "workflow_runs", []
    )
    candidates = [r for r in runs if r.get("name") == workflow] or runs
    if not candidates:
        sys.exit(f"No successful workflow runs found in {repo}.")
    for run in candidates:
        artifacts = [
            a
            for a in get_json(run["artifacts_url"], token).get("artifacts", [])
            if not a.get("expired")
        ]
        if artifacts:
            return run, artifacts
    sys.exit("No unexpired artifacts found in recent successful runs.")


def choose(artifacts: list[dict], preselected: str | None) -> dict:
    by_name = {a["name"]: a for a in artifacts}
    if preselected:
        if preselected not in by_name:
            sys.exit(
                f"Artifact '{preselected}' not in this run. Available: "
                + ", ".join(sorted(by_name))
            )
        return by_name[preselected]

    ordered = [by_name[n] for n in KNOWN_ARTIFACTS if n in by_name]
    ordered += [a for a in artifacts if a["name"] not in KNOWN_ARTIFACTS]

    print("\nAvailable artifacts:")
    for i, a in enumerate(ordered, 1):
        print(f"  {i}) {a['name']}  ({a['size_in_bytes'] / 1_048_576:.1f} MiB)")
    while True:
        answer = input("Select artifact [1-%d]: " % len(ordered)).strip()
        if answer.isdigit() and 1 <= int(answer) <= len(ordered):
            return ordered[int(answer) - 1]
        if answer in by_name:
            return by_name[answer]
        print("Invalid choice.")


def download(artifact: dict, token: str, dest: Path) -> None:
    print(f"Downloading {artifact['name']} ...")
    opener = urllib.request.build_opener(StripAuthOnRedirect)
    with request(artifact["archive_download_url"], token, opener) as resp:
        with dest.open("wb") as fh:
            shutil.copyfileobj(resp, fh)


def safe_extract(archive: Path, target: Path) -> None:
    target.mkdir(parents=True, exist_ok=True)
    if zipfile.is_zipfile(archive):
        with zipfile.ZipFile(archive) as zf:
            for member in zf.namelist():
                out = (target / member).resolve()
                if not str(out).startswith(str(target.resolve())):
                    sys.exit(f"Refusing unsafe path in archive: {member}")
            zf.extractall(target)
        return
    if tarfile.is_tarfile(archive):
        with tarfile.open(archive) as tf:
            for member in tf.getmembers():
                out = (target / member.name).resolve()
                if not str(out).startswith(str(target.resolve())):
                    sys.exit(f"Refusing unsafe path in archive: {member.name}")
            tf.extractall(target)
        return
    sys.exit(f"Unsupported archive format: {archive.name}")


def unpack_nested(staging: Path) -> None:
    """CI artifacts often wrap a single zip/tarball; unpack it in place."""
    entries = list(staging.iterdir())
    if len(entries) != 1 or not entries[0].is_file():
        return
    inner = entries[0]
    if inner.suffix in {".zip", ".gz", ".tgz", ".xz", ".bz2"} or inner.name.endswith(
        ".tar"
    ):
        print(f"Unpacking nested archive {inner.name} ...")
        moved = inner.parent / ("__nested__" + inner.name)
        inner.rename(moved)
        safe_extract(moved, staging)
        moved.unlink()


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--repo", default=DEFAULT_REPO)
    parser.add_argument("--workflow", default=DEFAULT_WORKFLOW)
    parser.add_argument("--artifact", help="skip the prompt and pick this artifact")
    parser.add_argument("--token", help="GitHub token")
    parser.add_argument(
        "--dest",
        default=str(Path(__file__).resolve().parent / "lvgl-c"),
        help="destination folder (default: ./lvgl-c)",
    )
    parser.add_argument(
        "--clean",
        action="store_true",
        help="delete the destination folder before extracting",
    )
    args = parser.parse_args()

    token = find_token(args.token)
    run, artifacts = latest_run_with_artifacts(args.repo, args.workflow, token)
    print(
        f"Latest run: #{run['run_number']} ({run['head_branch']} @ "
        f"{run['head_sha'][:8]}) {run['created_at']}\n{run['html_url']}"
    )

    artifact = choose(artifacts, args.artifact)
    dest = Path(args.dest).resolve()

    with tempfile.TemporaryDirectory() as tmp:
        tmp_path = Path(tmp)
        zip_path = tmp_path / f"{artifact['name']}.zip"
        download(artifact, token, zip_path)

        staging = tmp_path / "staging"
        safe_extract(zip_path, staging)
        unpack_nested(staging)

        if args.clean and dest.exists():
            shutil.rmtree(dest)
        dest.mkdir(parents=True, exist_ok=True)
        for item in staging.iterdir():
            target = dest / item.name
            if target.exists():
                shutil.rmtree(target) if target.is_dir() else target.unlink()
            shutil.move(str(item), str(target))

    print(f"Done: {artifact['name']} extracted into {dest}")


if __name__ == "__main__":
    main()
