export default class Const {
    static NOTIFICATION_NEW_EPISODE = 0;
    static NOTIFICATION_NEW_MOVIE = 1;
    static NOTIFICATION_DOWNLOAD_START = 2;
    static NOTIFICATION_DOWNLOAD_SUCCESS = 3;
    static NOTIFICATION_DOWNLOAD_FAILURE = 4;
    static NOTIFICATION_TEXT = 5;
    static NOTIFICATION_NO_TORRENT = 6;

    static OK = 0;
    static WARNING = 1;
    static CRITICAL = 2;
    static UNKNOWN = 3;
    static INFO = 4;

    static NOTIFICATIONS_REFRESH = 10000;
    static DATA_REFRESH = 10000;

    static TVSHOW_RETURNING = 1;
    static TVSHOW_PLANNED = 2;
    static TVSHOW_ENDED = 3;
    static TVSHOW_UNKNOWN = 4;

    static DISPLAY_MINIATURES = 0;
    static DISPLAY_LIST = 1;
}
