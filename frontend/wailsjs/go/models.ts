export namespace helper {
	
	export class Anime {
	    Title: string;
	    Url: string;
	    Poster: string;
	
	    static createFrom(source: any = {}) {
	        return new Anime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.Url = source["Url"];
	        this.Poster = source["Poster"];
	    }
	}
	export class AnimeInfo {
	    EpisodeCount: number;
	    Is18plus: boolean;
	    Tags: string[];
	    Studio: string;
	    Status: string;
	    Plot: string;
	    EpisodesList: string[];
	
	    static createFrom(source: any = {}) {
	        return new AnimeInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EpisodeCount = source["EpisodeCount"];
	        this.Is18plus = source["Is18plus"];
	        this.Tags = source["Tags"];
	        this.Studio = source["Studio"];
	        this.Status = source["Status"];
	        this.Plot = source["Plot"];
	        this.EpisodesList = source["EpisodesList"];
	    }
	}

}

