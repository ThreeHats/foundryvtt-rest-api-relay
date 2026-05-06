#!/bin/sh
set -e

# Intel/AMD GPU: join the render group matching the host's /dev/dri GID.
# This eliminates the need for group_add in docker-compose — just expose
# devices: /dev/dri:/dev/dri and GPU mode is auto-configured.
if [ -d /dev/dri ]; then
    RENDER_GID=$(stat -c %g /dev/dri/renderD128 2>/dev/null \
        || stat -c %g /dev/dri/renderD129 2>/dev/null \
        || stat -c %g /dev/dri/card0 2>/dev/null \
        || echo "")
    if [ -n "$RENDER_GID" ] && [ "$RENDER_GID" != "0" ]; then
        groupadd -g "$RENDER_GID" hostrender 2>/dev/null || true
        usermod -aG hostrender appuser 2>/dev/null || true
    fi
fi

# NVIDIA GPU: most setups use 666 root:root (no group needed), but some
# distributions configure 660 root:video. Handle defensively.
if [ -e /dev/nvidiactl ]; then
    NVIDIA_GID=$(stat -c %g /dev/nvidiactl 2>/dev/null || echo "")
    if [ -n "$NVIDIA_GID" ] && [ "$NVIDIA_GID" != "0" ]; then
        groupadd -g "$NVIDIA_GID" hostnvidia 2>/dev/null || true
        usermod -aG hostnvidia appuser 2>/dev/null || true
    fi
fi

# Ensure appuser owns the data directory so it can read secrets and write the DB
# regardless of the host UID that created the mounted volume.
chown -R appuser:appuser /app/data 2>/dev/null || true

# Remove stale Chrome profile singleton locks left by a previous container run.
rm -f /app/data/chrome-profile/SingletonLock \
      /app/data/chrome-profile/SingletonCookie \
      /app/data/chrome-profile/SingletonSocket 2>/dev/null || true

# Make workspace directories that the test runner writes to writable by appuser.
# The workspace is owned by the host user (UID 1000) while appuser is UID 1001.
# chmod a+w on the directory allows appuser to create new files inside it.
for dir in /workspace/tests /workspace/docs/examples /workspace/test-results; do
    if [ -d "$dir" ]; then
        chmod -R a+w "$dir" 2>/dev/null || true
    fi
done
if [ -d /workspace/tests ]; then
    touch /workspace/tests/.global-vars.json 2>/dev/null || true
    chmod 666 /workspace/tests/.global-vars.json 2>/dev/null || true
fi

exec gosu appuser "$@"
