#!/bin/bash
source .env

go build \
  -ldflags="-s -w -H=windowsgui \
    -X 'rewardsAutomation/internal/config.BuildEdgePath=${EDGE_PATH}' \
    -X 'rewardsAutomation/internal/config.BuildUserDataDir=${USER_DATA_DIR}' \
    -X 'rewardsAutomation/internal/config.BuildTmpUserDataDir=${TMP_USER_DATA_DIR}'" \
  -o bin/rewardsRobot.exe ./cmd/
