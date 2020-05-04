import MediaIds from "./media_ids";
import Episode from "./episode";

export type TvSeason = {
    ID: number,
    CreatedAt: Date,
	AirDate      :Date,
	EpisodeCount :number,
	SeasonNumber :number,
	PosterPath   :string,
	TvShowID     :number,
};

export type SeasonDetails = {
	Info        :TvSeason,
	EpisodeList :Episode[],
    LoadError   :boolean,
	LoadPending :boolean,
};

type TvShow = {
    ID: number,
    CreatedAt: Date,
    DeletedAt: Date | null,
	MediaIds         :MediaIds,
	Banner           :string,
	Poster           :string,
	FirstAired       :Date,
	Overview         :string,
	Title            :string,
	OriginalTitle    :string,
	CustomTitle      :string,
	DisplayTitle     :string,
	Status           :number,
	NumberOfEpisodes :number,
	NumberOfSeasons  :number,
	Seasons          :TvSeason[],
	UseDefaultTitle  :boolean,
	IsAnime          :boolean,
};

export default TvShow;
