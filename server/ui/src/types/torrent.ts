type Torrent = {
    ID: number,
    CreatedAt: Date,
    FailedTorrentID: number,
    TorrentId: string,
    Failed :boolean,
    Name: string,
    Link: string,
    DownloadDir: string,
    Seeders: number,
    PercentDone :number,
    TotalSize :number,
    RateUpload :number,
    RateDownload :number,
    ETA :Date,
};

export default Torrent;
