#!/usr/bin/env python2

import json
import optparse
import os
import sys
import time
import urllib2

VERSION = "0.1"
 
OK = 0
WARNING = 1
CRITICAL = 2
UNKNOWN = 3

GREEN = '#2A9A3D'
RED = '#FF0000'
ORANGE = '#f57700'
GRAY = '#f57700'

parser = optparse.OptionParser("%prog [options]", version="%prog " + VERSION)
parser.add_option('-H', '--hostname', dest="hostname", help='Hostname to connect to')
parser.add_option('-p', '--port', dest="port", type="int", default=8400, help='Flemzerd port (default: 8400)')
parser.add_option('-S', '--https', dest="https", type="int", default=0,  help='Use SSL')

perfdata = []
output = ""

def add_perfdata(name, value, min="", max="", warning="", critical=""):
    global perfdata
    perfdata.append("%s=%s;%s;%s;%s;%s" % (name, value, min, max, warning, critical))

def exit(status, exit_label=""):
    global perfdata
    global output

    label = exit_label
    color = GRAY

    if status == OK:
        if not label:
            label = "OK"
        color = GREEN
    elif status == WARNING:
        if not label:
            label = "WARNING"
        color = ORANGE
    elif status == CRITICAL:
        if not label:
            label = "CRITICAL"
        color = RED
    else:
        if not label:
            label = "UNKNOWN"
        color = GRAY

    print "<span style=\"color:%s;font-weight: bold;\">[%s]</span> %s | %s" % (color, label, output, " ".join(perfdata))
    sys.exit(status)


def get_stats(hostname, port, https):
    global output

    if https == 1:
        host = "https://%s:%d" % (hostname, port)
    else: 
        host = "http://%s:%d" % (hostname, port)

    url = "%s%s" % (host, "/api/v1/stats")

    try:
        start = time.time()
        req = urllib2.urlopen(url)
        end = time.time()

        data = req.read()
    except urllib2.URLError as e:
        output += "Could not contact flemzerd: %s" % e
        exit(CRITICAL)

    parsed_data = json.loads(data)
    flemzerd_stats = parsed_data["stats"]

    add_perfdata("response_time", end - start)

    add_perfdata("movies_tracked",          flemzerd_stats["Movies"]["Tracked"])
    add_perfdata("movies_downloading",      flemzerd_stats["Movies"]["Downloading"])
    add_perfdata("movies_downloaded",       flemzerd_stats["Movies"]["Downloaded"])
    add_perfdata("movies_removed",          flemzerd_stats["Movies"]["Removed"])

    add_perfdata("shows_tracked",           flemzerd_stats["Shows"]["Tracked"])
    add_perfdata("shows_removed",           flemzerd_stats["Shows"]["Removed"])
    add_perfdata("episodes_downloading",    flemzerd_stats["Episodes"]["Downloading"])
    add_perfdata("episodes_downloaded",     flemzerd_stats["Episodes"]["Downloaded"])

    add_perfdata("notifications_read",      flemzerd_stats["Notifications"]["Read"])
    add_perfdata("notifications_unread",    flemzerd_stats["Notifications"]["Unread"])

    add_perfdata("runtime_gomaxprocs",      flemzerd_stats["Runtime"]["GoMaxProcs"])
    add_perfdata("runtime_numcpu",          flemzerd_stats["Runtime"]["NumCPU"])
    add_perfdata("runtime_goroutines",      flemzerd_stats["Runtime"]["GoRoutines"])

    if flemzerd_stats["Notifications"]["Unread"] > 0:
        output += "%d unread notifications" % flemzerd_stats["Notifications"]["Unread"]

    exit_status = OK

    if flemzerd_stats["Movies"]["Tracked"] == 0 and flemzerd_stats["Shows"]["Tracked"] == 0:
        output += "Flemzerd is running but no movies and shows are tracked. Configure watchlists to track medias."
        exit_stats = WARNING

    exit(OK)

if __name__ == '__main__':
    # Ok first job : parse args
    opts, args = parser.parse_args()
    if args:
        parser.error("Does not accept any argument.")

    port = opts.port
    hostname = opts.hostname
    if not hostname:
        # print "<span style=\"color:#A9A9A9;font-weight: bold;\">[ERROR]</span> Hostname parameter (-H) is mandatory"
        output = "Hostname parameter (-H) is mandatory"
        exit(CRITICAL, "ERROR")

    get_stats(hostname, port, opts.https)
