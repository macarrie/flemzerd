#!/bin/bash

. common.sh

printf "$CHAR_CORNER_TOP_LEFT"
printf "%0.s$CHAR_HBAR" {1..46}
printf "$CHAR_CORNER_TOP_RIGHT\n"
printf "$CHAR_VBAR               ${GREEN}FLEMZERD UPDATER${RESET}               $CHAR_VBAR\n"
printf "$CHAR_CORNER_BOTTOM_LEFT"
printf "%0.s$CHAR_HBAR" {1..46}
printf "$CHAR_CORNER_BOTTOM_RIGHT\n"

if [ ! -f $BIN/flemzerd ]; then
    printf -- "${RED}Flemzerd is not installed. Use the install script instead${RESET}\n"
    exit 1
fi

stop_flemzerd

copy_binary
copy_server_files

start_flemzerd

exit 0
