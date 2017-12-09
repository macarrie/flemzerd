#!/bin/bash

. common.sh

printf "$CHAR_CORNER_TOP_LEFT"
printf "%0.s$CHAR_HBAR" {1..48}
printf "$CHAR_CORNER_TOP_RIGHT\n"
printf "$CHAR_VBAR               ${GREEN}FLEMZERD INSTALLER${RESET}               $CHAR_VBAR\n"
printf "$CHAR_CORNER_BOTTOM_LEFT"
printf "%0.s$CHAR_HBAR" {1..48}
printf "$CHAR_CORNER_BOTTOM_RIGHT\n"

# Copy flemzerd executable
copy_binary
# Create flemzerd user
create_user
# Ensure file system is prepared for flemzerd
create_folder_structure
# Copy default configuration files
copy_config_files
# Create systemd unit
create_systemd_unit
# Reload systemd units configuration
reload_systemd_units
# Start flemzerd daemon
start_flemzerd

exit 0
