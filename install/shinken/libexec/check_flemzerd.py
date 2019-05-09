#!/usr/bin/env python2

import json
import optparse
import sys
import time
import urllib2

VERSION = "0.1"
 
OK = 0
WARNING = 1
CRITICAL = 2
UNKNOWN = 3

STATUS_STR = ["OK", "WARNING", "CRITICAL", "CRITICAL", "UNKNOWN"]

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
long_output = ""

def add_perfdata(name, value, min="", max="", warning="", critical=""):
    global perfdata
    perfdata.append("%s=%s;%s;%s;%s;%s" % (name, value, min, max, warning, critical))


def print_status(status, exit_label=""):
    label = exit_label
    color = GRAY

    if status == OK:
        if not label:
            label = STATUS_STR[OK]
        color = GREEN
    elif status == WARNING:
        if not label:
            label = STATUS_STR[WARNING]
        color = ORANGE
    elif status == CRITICAL:
        if not label:
            label = STATUS_STR[CRITICAL]
        color = RED
    else:
        if not label:
            label = STATUS_STR[UNKNOWN]
        color = GRAY

    printed_status = "<span style=\"color:%s;font-weight: bold;\">[%s]</span>" % (color, label)
    return printed_status


def exit(status, exit_label=""):
    global perfdata
    global output
    global long_output

    print "%s %s | %s" % (print_status(status, exit_label), output, " ".join(perfdata))
    if long_output != "":
        print long_output

    sys.exit(status)


def http_call(hostname, port, https, endpoint):
    global output

    if https == 1:
        host = "https://%s:%d" % (hostname, port)
    else:
        host = "http://%s:%d" % (hostname, port)

    url = "%s%s" % (host, endpoint)

    try:
        start = time.time()
        req = urllib2.urlopen(url)
        end = time.time()

        data = req.read()
    except urllib2.URLError as e:
        output += "Could not contact flemzerd: %s" % e
        exit(CRITICAL)

    return end - start, data


def get_stats(hostname, port, https):
    global output
    global long_output

    response_time, data = http_call(hostname, port, https, "/api/v1/stats")

    parsed_data = json.loads(data)
    flemzerd_stats = parsed_data["stats"]

    add_perfdata("response_time", response_time)

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

    exit_status = OK

    _, config_check_data = http_call(hostname, port, https, "/api/v1/config/check")
    config_check_json = json.loads(config_check_data)

    if len(config_check_json) > 0:
        exit_status = WARNING
        long_output += "<ul>"
        for config_err in config_check_json:
            if config_err["Status"] == CRITICAL:
                exit_status = CRITICAL
            if config_err["Key"] == "":
                long_output += "<li>%s %s</li>" % (print_status(config_err["Status"]), config_err["Message"])
            else:
                long_output += "<li>%s %s (%s: '%s')</li>" % (
                    print_status(config_err["Status"]), config_err["Message"], config_err["Key"], config_err["Value"])

        output += "%d configuration errors/warnings<br />" % len(config_check_json)
        long_output += "</ul>"

    if flemzerd_stats["Notifications"]["Unread"] > 0:
        output += "%d unread notifications<br />" % flemzerd_stats["Notifications"]["Unread"]


    if flemzerd_stats["Movies"]["Tracked"] == 0 and flemzerd_stats["Shows"]["Tracked"] == 0:
        output += "Flemzerd is running but no movies and shows are tracked. Configure watchlists to track medias.<br />"
        exit_status = WARNING

    exit(exit_status)

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
