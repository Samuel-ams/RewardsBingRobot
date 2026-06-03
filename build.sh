#!/bin/bash
EDGE_PATH="${EDGE_PATH:-C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe}"
USER_EDGE_DIR="${USER_EDGE_DIR:-AppData\\Local\\Microsoft\\Edge\\User Data}"

go build \
  -ldflags="-s -w -H=windowsgui \
    -X 'rewardsAutomation/internal/config.BuildEdgePath=${EDGE_PATH}' \
    -X 'rewardsAutomation/internal/config.BuildUserEdgeDir=${USER_EDGE_DIR}'" \
  -o bin/rewardsRobot.exe ./cmd/
