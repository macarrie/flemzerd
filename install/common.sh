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
ETC=/etc/flemzerd BIN=/usr/bin 
USER=flemzer
GROUP=flemzer


function copy_binary {
    # Copy exec file
    log_line "- Copying flemzerd binary"
    cp dependencies/vidocq $BIN/vidocq
    cp bin/flemzerd $BIN/flemzerd
    chmod a+x $BIN/vidocq
    chmod a+x $BIN/flemzerd
    print_done
}

function copy_server_files {
    log_line "- Copying UI files"
    mkdir -p $LIB/server/ui
    cp -r ui/* $LIB/server/ui
    print_done
}

function copy_certs {
    log_line "- Copying SSL certs"
    mkdir -p $LIB/certs
    cp -r certs/* $LIB/certs
    chown -R $USER:$GROUP $LIB/certs
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
    mkdir -p $LIB/db
    mkdir -p $LIB/library/shows
    mkdir -p $LIB/library/movies
    mkdir -p $LIB/cache/shows
    mkdir -p $LIB/cache/movies
    mkdir -p $LIB/certs
    mkdir -p $LIB/tmp
    mkdir -p $ETC
    chown -R $USER:$GROUP $ETC
    chown -R $USER:$GROUP $LIB
    chmod 0755 $LIB/tmp

    print_done
}

function copy_config_files {
    log_line "- Copying default configuration file"
    if [ ! -f $ETC/flemzerd.toml ]; then
        cp flemzerd.toml $ETC
        print_done
    else
        print_skipping
    fi
}

function backup_db {
    log_line "- Creating DB backup"
    cp $LIB/db/flemzer.db $LIB/db/flemzer.db.$(date +%s).bkp
    print_done
}

function copy_dependencies_config_files {
    log_line "- Copying dependencies configuration files"
    cp flemzerd-docker.toml $ETC
    cp -r transmission $ETC/transmission
    cp -r jackett $ETC/jackett
    cp -r flemzerd-docker.toml $ETC
    print_done
}

function create_systemd_unit {
    log_line "- Creating systemd unit"
    # Create systemd unit
    cp flemzerd.service /etc/systemd/system/
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
