type Torrent = {
    ID: number,
    CreatedAt: Date,
    FailedTorrentID: number,
    CurrentTorrentID: number,
    TorrentId: string,
    Name: string,
    Link: string,
    DownloadDir: string,
    Seeders: number,
};

export default Torrent;
