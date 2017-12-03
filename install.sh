#!/bin/bash

echo "--------------- FLEMZERD INSTALL ---------------"

RUN=/var/run/flemzerd
ETC=/etc/flemzerd
BIN=/usr/bin

USER=flemzer
GROUP=flemzer

# Copy exec file
printf "Copying flemzerd binary"
cp flemzerd $BIN/flemzerd
chmod a+x $BIN/flemzerd
printf "\t\t\tdone\n"

printf "Creating flemzer user"
id -u $USER 2>&1 /dev/null
if [ $? -eq 0 ]; then
    useradd -M $USER
    printf "\t\t\tdone\n"
else
    printf "\t\tskipping\n"
fi

# Ensure file system is prepared for flemzerd
printf "Creating folder hierarchy"
mkdir -p $RUN
mkdir -p $ETC
chown $USER:$GROUP $ETC
chown $USER:$GROUP $RUN
mkdir -p /var/lib/flemzerd/Library/Shows
mkdir -p /var/lib/flemzerd/Library/Movies
printf "\t\tdone\n"

printf "Copying default configuration file"
if [ ! -f $ETC/flemzerd.yml ]; then
    cp install/flemzerd.yml $ETC
    printf "\tdone\n"
else
    printf "\tskipping\n"
fi

printf "Creating systemd unit"
# Create systemd unit
cp install/flemzerd.service /etc/systemd/system/
chmod 0644 /etc/systemd/system/flemzerd.service
printf "\t\t\tdone\n"

# Reload systemd units
systemctl daemon-reload

exit 0
