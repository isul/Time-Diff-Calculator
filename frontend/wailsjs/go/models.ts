export namespace backend {
	
	export class FormatLine {
	    key: string;
	    label: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new FormatLine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.label = source["label"];
	        this.value = source["value"];
	    }
	}
	export class CalculateResponse {
	    ok: boolean;
	    error?: string;
	    warning?: string;
	    results?: FormatLine[];
	    locale: string;
	
	    static createFrom(source: any = {}) {
	        return new CalculateResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.error = source["error"];
	        this.warning = source["warning"];
	        this.results = this.convertValues(source["results"], FormatLine);
	        this.locale = source["locale"];
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
	
	export class Settings {
	    formats: Record<string, boolean>;
	    customFormat: string;
	    language: string;
	    zeroPadding?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.formats = source["formats"];
	        this.customFormat = source["customFormat"];
	        this.language = source["language"];
	        this.zeroPadding = source["zeroPadding"];
	    }
	}
	export class ValidateResponse {
	    ready: boolean;
	    level: string;
	    message?: string;
	
	    static createFrom(source: any = {}) {
	        return new ValidateResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ready = source["ready"];
	        this.level = source["level"];
	        this.message = source["message"];
	    }
	}

}

