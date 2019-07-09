import MediaIds from "./media_ids";
import Episode from "./episode";

type TvSeason = {
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
};

type TvShow = {
    ID: number,
    CreatedAt: Date,
    DeletedAt: Date,
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
