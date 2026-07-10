#!/bin/bash
source .env

go build \
  -ldflags="-s -w -H=windowsgui \
    -X 'rewardsAutomation/internal/config.BuildUserEdgeDir=${USER_EDGE_DIR}' \
    -X 'rewardsAutomation/internal/config.BuildTmpUserEdgeDir=${TMP_USER_EDGE_DIR}'" \
  -o bin/rewardsRobot.exe ./cmd/
