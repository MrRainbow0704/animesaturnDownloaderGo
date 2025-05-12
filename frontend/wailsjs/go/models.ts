export namespace helper {
	
	export class AnimeInfo {
	    EpisodeCount: number;
	    Tags: string[];
	    Studio: string;
	    Status: string;
	    Plot: string;
	    FirstEpisode: number;
	    LastEpisode: number;
	
	    static createFrom(source: any = {}) {
	        return new AnimeInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.EpisodeCount = source["EpisodeCount"];
	        this.Tags = source["Tags"];
	        this.Studio = source["Studio"];
	        this.Status = source["Status"];
	        this.Plot = source["Plot"];
	        this.FirstEpisode = source["FirstEpisode"];
	        this.LastEpisode = source["LastEpisode"];
	    }
	}
	export class Anime {
	    Info: AnimeInfo;
	    Title: string;
	    Url: string;
	    Poster: string;
	
	    static createFrom(source: any = {}) {
	        return new Anime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Info = this.convertValues(source["Info"], AnimeInfo);
	        this.Title = source["Title"];
	        this.Url = source["Url"];
	        this.Poster = source["Poster"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

