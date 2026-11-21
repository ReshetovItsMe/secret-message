#!/usr/bin/env bash
# Deploy only the three application images. Redis and host services are untouched.
set -euo pipefail

if [[ $# != 1 || ! "$1" =~ ^[0-9a-f]{40}$ ]]; then
    echo 'Expected one lowercase 40-character commit SHA' >&2
    exit 2
fi
candidate=$1
root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)"
active_compose="$root/docker-compose.prod.yml"
staged_compose="$root/staging/$candidate/docker-compose.prod.yml"
candidate_compose="${CANDIDATE_COMPOSE_FILE:-$active_compose}"
if [[ ! -f "$candidate_compose" ]]; then
    echo 'Candidate Compose file unavailable' >&2
    exit 1
fi
candidate_compose="$(cd -- "$(dirname -- "$candidate_compose")" && pwd -P)/$(basename -- "$candidate_compose")"
if [[ "$candidate_compose" != "$active_compose" && "$candidate_compose" != "$staged_compose" ]]; then
    echo 'Candidate Compose file must be the active file or the SHA-specific staging file' >&2
    exit 1
fi
settings="${DEPLOY_SETTINGS_FILE:-/etc/secret-message/deploy.conf}"
if [[ ! -f "$settings" || -L "$settings" ]]; then
    echo 'Protected host deployment settings unavailable' >&2
    exit 1
fi
# The owner provisions this trusted, restricted file outside the repository.
# It defines APP_HOST and RSA_PRIVATE_KEY_FILE; never enable shell tracing.
# shellcheck source=/dev/null
source "$settings"
if [[ ! "${APP_HOST:-}" =~ ^([a-z0-9][-a-z0-9]*\.)+[a-z]{2,}$ ||
      ! "${RSA_PRIVATE_KEY_FILE:-}" == /* || ! -f "$RSA_PRIVATE_KEY_FILE" ]]; then
    echo 'Host deployment settings are incomplete' >&2
    exit 1
fi

marker="$root/.deployed-sha"
previous=''
if [[ -f "$marker" ]]; then
    IFS= read -r previous < "$marker" || true
    if [[ ! "$previous" =~ ^[0-9a-f]{40}$ ]]; then
        echo 'Existing deployed revision marker is invalid; manual recovery required' >&2
        exit 1
    fi
fi

if command -v flock >/dev/null; then
    exec 9>"$root/.deploy.lock"
    flock -x 9
fi

free_kb="$(df -Pk "$root" | awk 'NR == 2 {print $4}')"
if [[ ! "$free_kb" =~ ^[0-9]+$ ]] || (( free_kb < 2097152 )); then
    echo 'Less than 2 GiB free: no image pull attempted' >&2
    exit 1
fi
if [[ -r /proc/meminfo ]]; then
    available_kb="$(awk '/^MemAvailable:/ {print $2}' /proc/meminfo)"
    if [[ ! "$available_kb" =~ ^[0-9]+$ ]] || (( available_kb < 131072 )); then
        echo 'Less than 128 MiB available RAM: rollout refused' >&2
        exit 1
    fi
fi

compose_file="$candidate_compose"
compose() {
    docker compose --project-name secret-message -f "$compose_file" "$@"
}
smoke() {
    python3 "$root/deploy/smoke.py" --base-url "https://$APP_HOST"
}
export RSA_PRIVATE_KEY_FILE
export IMAGE_TAG="$candidate"
compose config --quiet
if [[ -z "$(compose ps --status running -q redis)" ]]; then
    echo 'Redis must be bootstrapped separately and running before deployment' >&2
    exit 1
fi
compose pull frontend gateway secret-assistant
if ! compose up -d --no-deps frontend gateway secret-assistant || ! smoke; then
    if [[ -n "$previous" && "$previous" != "$candidate" ]]; then
        export IMAGE_TAG="$previous"
        compose_file="$root/staging/$previous/docker-compose.prod.yml"
        # Bootstrap releases may predate staged snapshots; the active file is
        # still the last successful config because promotion happens below.
        if [[ ! -f "$compose_file" ]]; then
            compose_file="$active_compose"
        fi
        if compose config --quiet && compose pull frontend gateway secret-assistant &&
           compose up -d --no-deps frontend gateway secret-assistant && smoke; then
            echo 'Candidate failed; previous app images restored' >&2
            exit 1
        fi
        echo 'Candidate and rollback unhealthy; manual recovery required' >&2
    else
        echo 'Initial or same-revision rollout failed; manual recovery required' >&2
    fi
    exit 2
fi

# Keep a versioned copy for recovery; do not alter the active config until
# all app containers and the public synthetic check have passed.
mkdir -p "$root/staging/$candidate"
if [[ "$candidate_compose" != "$staged_compose" ]]; then
    cp -- "$candidate_compose" "$staged_compose"
fi
pending_compose="$(mktemp "$root/.compose.XXXXXXXX")"
pending="$(mktemp "$root/.deployed-sha.XXXXXXXX")"
trap 'rm -f "$pending_compose" "$pending"' EXIT
install -m 0644 "$candidate_compose" "$pending_compose"
printf '%s
' "$candidate" > "$pending"
mv -f -- "$pending_compose" "$active_compose"
mv -f -- "$pending" "$marker"
echo 'Synthetic check passed; application revision recorded'
