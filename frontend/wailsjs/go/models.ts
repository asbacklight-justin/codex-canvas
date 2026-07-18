export namespace main {

	export class ActionResult {
	    success: boolean;
	    message: string;
	    port?: number;

	    static createFrom(source: any = {}) {
	        return new ActionResult(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.port = source["port"];
	    }
	}
	export class ThemeConfig {
	    imagePath: string;
	    imageName: string;
	    overlay: number;
	    surfaceOpacity: number;
	    sidebarOpacity: number;
	    blur: number;
	    radius: number;
	    scale: number;
	    position: string;
	    accent: string;
	    active: boolean;
	    lastPort: number;

	    static createFrom(source: any = {}) {
	        return new ThemeConfig(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.imagePath = source["imagePath"];
	        this.imageName = source["imageName"];
	        this.overlay = source["overlay"];
	        this.surfaceOpacity = source["surfaceOpacity"];
	        this.sidebarOpacity = source["sidebarOpacity"];
	        this.blur = source["blur"];
	        this.radius = source["radius"];
	        this.scale = source["scale"];
	        this.position = source["position"];
	        this.accent = source["accent"];
	        this.active = source["active"];
	        this.lastPort = source["lastPort"];
	    }
	}
	export class AppStatus {
	    platform: string;
	    codexFound: boolean;
	    codexPath: string;
	    codexVersion: string;
	    supported: boolean;
	    active: boolean;
	    savedTheme: ThemeConfig;
	    previewUrl: string;
	    statusMessage: string;
	    compatibility: string;

	    static createFrom(source: any = {}) {
	        return new AppStatus(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.platform = source["platform"];
	        this.codexFound = source["codexFound"];
	        this.codexPath = source["codexPath"];
	        this.codexVersion = source["codexVersion"];
	        this.supported = source["supported"];
	        this.active = source["active"];
	        this.savedTheme = this.convertValues(source["savedTheme"], ThemeConfig);
	        this.previewUrl = source["previewUrl"];
	        this.statusMessage = source["statusMessage"];
	        this.compatibility = source["compatibility"];
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
	export class ImageSelection {
	    path: string;
	    name: string;
	    size: number;
	    previewUrl: string;
	    cancelled: boolean;

	    static createFrom(source: any = {}) {
	        return new ImageSelection(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.previewUrl = source["previewUrl"];
	        this.cancelled = source["cancelled"];
	    }
	}
	export class PresetBackground {
	    id: string;
	    previewUrl: string;

	    static createFrom(source: any = {}) {
	        return new PresetBackground(source);
	    }

	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.previewUrl = source["previewUrl"];
	    }
	}
}
