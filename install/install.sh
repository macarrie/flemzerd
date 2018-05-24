#!/bin/bash

. common.sh

printf "$CHAR_CORNER_TOP_LEFT"
printf "%0.s$CHAR_HBAR" {1..48}
printf "$CHAR_CORNER_TOP_RIGHT\n"
printf "$CHAR_VBAR               ${GREEN}FLEMZERD INSTALLER${RESET}               $CHAR_VBAR\n"
printf "$CHAR_CORNER_BOTTOM_LEFT"
printf "%0.s$CHAR_HBAR" {1..48}
printf "$CHAR_CORNER_BOTTOM_RIGHT\n"


create_user
create_folder_structure

copy_binary
copy_server_files
copy_config_files

create_systemd_unit
reload_systemd_units

#start_flemzerd

exit 0
