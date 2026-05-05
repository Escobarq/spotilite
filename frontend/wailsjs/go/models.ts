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

}

