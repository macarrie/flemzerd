import Torrent from "./torrent";

type DownloadingItem = {
    ID: number,
    CreatedAt: Date,
    Pending: boolean,
    Downloading: boolean,
    Downloaded: boolean,
    TorrentList: Torrent[],
    CurrentTorrent: Torrent,
    CurrentDownloaderId: string,
    DownloadFailed: boolean,
    TorrentsNotFound: boolean,
};

export default DownloadingItem;
