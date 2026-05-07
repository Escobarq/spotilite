export namespace api {
	
	export class Settings {
	    language: string;
	    backgroundMode: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.backgroundMode = source["backgroundMode"];
	    }
	}
	export class SpotXSettings {
	    adBlock: boolean;
	    sectionBlock: boolean;
	    premiumSpoof: boolean;
	    experiments: boolean;
	    lyricsTheme: string;
	    trackHistory: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SpotXSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.adBlock = source["adBlock"];
	        this.sectionBlock = source["sectionBlock"];
	        this.premiumSpoof = source["premiumSpoof"];
	        this.experiments = source["experiments"];
	        this.lyricsTheme = source["lyricsTheme"];
	        this.trackHistory = source["trackHistory"];
	    }
	}

}

