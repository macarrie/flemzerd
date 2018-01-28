#!/bin/bash

# Display chars
CHAR_VBAR='│'
CHAR_HBAR='─'
CHAR_HBAR_LIGHT='─'
CHAR_CORNER_TOP_LEFT='┌'
CHAR_CORNER_TOP_RIGHT='┐'
CHAR_CORNER_BOTTOM_LEFT='└'
CHAR_CORNER_BOTTOM_RIGHT='┘'

# Color escape codes
BLUE="\033[0;94m"
GREEN="\033[92m"
GRAY="\033[0;37m"
PURPLE="\033[0;95m"
RED="\033[0;91m"
RESET="\033[0m"

# Helper functions
function print_done {
    printf "${GREEN}done${RESET}\n"
}

function print_skipping {
    printf "${GRAY}skip${RESET}\n"
}

function log_line {
    printf "%-50s" "$1"
}

# Setup vars
LIB=/var/lib/flemzerd
ETC=/etc/flemzerd
BIN=/usr/bin

USER=flemzer
GROUP=flemzer


function copy_binary {
    # Copy exec file
    log_line "- Copying flemzerd binary"
    cp flemzerd $BIN/flemzerd
    chmod a+x $BIN/flemzerd
    print_done
}

function create_user {
    log_line "- Creating flemzer user"
    id -u $USER > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        useradd -M $USER
        print_done
    else
        print_skipping
    fi
}

function create_folder_structure {
    log_line "- Creating folder hierarchy"
    mkdir -p $LIB
    mkdir -p $ETC
    chown $USER:$GROUP $ETC
    chown $USER:$GROUP $LIB
    mkdir -p /var/lib/flemzerd/Library/Shows
    mkdir -p /var/lib/flemzerd/Library/Movies
    print_done
}

function copy_config_files {
    log_line "- Copying default configuration file"
    if [ ! -f $ETC/flemzerd.yml ]; then
        cp install/flemzerd.yml $ETC
        print_done
    else
        print_skipping
    fi
}

function create_systemd_unit {
    log_line "- Creating systemd unit"
    # Create systemd unit
    cp install/flemzerd.service /etc/systemd/system/
    chmod 0644 /etc/systemd/system/flemzerd.service
    print_done
}

function reload_systemd_units {
    log_line "- Reloading systemd units"
    systemctl daemon-reload
    print_done
}

function start_flemzerd {
    log_line "- Starting flemzerd service"
    systemctl start flemzerd
    print_done
}

function stop_flemzerd {
    log_line "- Stopping flemzerd service"
    systemctl stop flemzerd
    print_done
}
