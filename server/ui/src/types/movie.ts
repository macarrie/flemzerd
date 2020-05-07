import MediaIds from "./media_ids";
import DownloadingItem from "./downloading_item";

type Movie = {
    ID: number,
    CreatedAt: Date,
    DeletedAt: Date | null,
    MediaIds: MediaIds,
    Title: string,
    OriginalTitle: string,
    CustomTitle: string,
    DisplayTitle: string,
    Overview: string,
    Poster: string,
    Background: string,
    Date: Date,
    Notified: boolean,
    DownloadingItem: DownloadingItem,
    UseDefaultTitle: boolean,
};

export default Movie;
