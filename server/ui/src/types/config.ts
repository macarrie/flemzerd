type Config = {
    Interface :{
        Enabled :boolean,
        Port    :number,
        Auth: {
            Username :string,
        },
    },
    Providers     :any,
    Indexers      :any,
    Notifiers     :any,
    Downloaders   :any,
    Watchlists    :any,
    MediaCenters  :any,
    Notifications :{
        Enabled                :boolean,
        NotifyNewEpisode       :boolean,
        NotifyNewMovie         :boolean,
        NotifyDownloadStart    :boolean,
        NotifyDownloadComplete :boolean,
        NotifyFailure          :boolean,
    },
    System :{
        CheckInterval                :number,
        TorrentDownloadAttemptsLimit :number,
        TrackShows                   :boolean,
        TrackMovies                  :boolean,
        AutomaticShowDownload        :boolean,
        AutomaticMovieDownload       :boolean,
        ShowDownloadDelay            :number,
        MovieDownloadDelay           :number,
        PreferredMediaQuality        :string,
        ExcludedReleaseTypes         :string,
        StrictTorrentCheck           :boolean,
    },
    Library :{
        ShowPath      :string,
        MoviePath     :string,
        CustomTmpPath :string,
    }
    Version :string,
};

export default Config;
