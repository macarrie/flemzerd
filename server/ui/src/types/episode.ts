import MediaIds from "./media_ids";
import DownloadingItem from "./downloading_item";
import TvShow from "./tvshow";

type Episode = {
    ID                :number,
    CreatedAt         :Date,
    DeletedAt         :Date,
	MediaIds          :MediaIds,
	TvShow            :TvShow,
	AbsoluteNumber    :number,
	Number            :number,
	Season            :number,
	Title             :string,
	Date              :Date,
	Overview          :string,
	Notified          :boolean,
	DownloadingItem   :DownloadingItem,
};

export default Episode;
