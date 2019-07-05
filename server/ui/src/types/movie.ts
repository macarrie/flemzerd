import MediaIds from "./media_ids";
import DownloadingItem from "./downloading_item";

// TODO: Fix any
type Movie = {
    ID: number,
    CreatedAt: Date,
    MediaIds: MediaIds,
    Title: string,
    OriginalTitle: string,
    CustomTitle: string,
    DisplayTitle: string,
    Overview: string,
    Poster: string,
    Date: Date,
    Notified: boolean,
    DownloadingItem: DownloadingItem,
    UseDefaultTitle: boolean,
};

export default Movie;