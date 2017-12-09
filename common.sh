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
BLUE="\033[0;34m"
GREEN="\033[92m"
GRAY="\033[0;37m"
PURPLE="\033[0;35m"
RED="\033[0;31m"
RESET="\033[0m"

# Helper functions
function print_done {
    printf "${GREEN}done${RESET}\n"
}

function print_skipping {
    printf "${GRAY}skip${RESET}\n"
}

# Setup vars
RUN=/var/run/flemzerd
ETC=/etc/flemzerd
BIN=/usr/bin

USER=flemzer
GROUP=flemzer


function copy_binary {
    # Copy exec file
    printf -- "- Copying flemzerd binary\t\t"
    cp flemzerd $BIN/flemzerd
    chmod a+x $BIN/flemzerd
    print_done
}

function create_user {
    printf -- "- Creating flemzer user\t\t\t"
    id -u $USER > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        useradd -M $USER
        print_done
    else
        print_skipping
    fi
}

function create_folder_structure {
    printf -- "- Creating folder hierarchy\t\t"
    mkdir -p $RUN
    mkdir -p $ETC
    chown $USER:$GROUP $ETC
    chown $USER:$GROUP $RUN
    mkdir -p /var/lib/flemzerd/Library/Shows
    mkdir -p /var/lib/flemzerd/Library/Movies
    print_done
}

function copy_config_files {
    printf -- "- Copying default configuration file\t"
    if [ ! -f $ETC/flemzerd.yml ]; then
        cp install/flemzerd.yml $ETC
        print_done
    else
        print_skipping
    fi
}

function create_systemd_unit {
    printf -- "- Creating systemd unit\t\t\t"
    # Create systemd unit
    cp install/flemzerd.service /etc/systemd/system/
    chmod 0644 /etc/systemd/system/flemzerd.service
    print_done
}

function reload_systemd_units {
    printf -- "- Reloading systemd units\t\t"
    systemctl daemon-reload
    print_done
}

function start_flemzerd {
    printf -- "- Starting flemzerd service\t\t"
    systemctl start flemzerd
    print_done
}

function stop_flemzerd {
    printf -- "- Stopping flemzerd service\t\t"
    systemctl stop flemzerd
    print_done
}
